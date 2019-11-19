package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/alash3al/redix/db/driver"
)

// Driver represents a driver
type Driver struct {
	db *leveldb.DB
}

// Open implements driver.Open
func (drv Driver) Open(dbname string, opts map[string]interface{}) (driver driver.Interface, err error) {
	db, err := leveldb.OpenFile(dbname, nil)

	if err != nil {
		return driver, err
	}

	return Driver{
		db: db,
	}, nil
}

// Put implements driver.Put
func (drv Driver) Put(k, v []byte) error {
	return drv.db.Put(k, v, nil)
}

// BulkPut perform multi put operation
func (drv Driver) BulkPut(pairs []driver.KeyValue) error {
	batch := new(leveldb.Batch)

	for _, pair := range pairs {
		batch.Put(pair.Key, pair.Value)
	}

	return drv.db.Write(batch, nil)
}

// Get implements driver.Get
func (drv Driver) Get(k []byte) ([]byte, error) {
	return drv.db.Get(k, nil)
}

// Delete implements driver.Delete
func (drv Driver) Delete(k []byte) error {
	return drv.db.Delete(k, nil)
}

// Close implements driver.Close
func (drv Driver) Close() error {
	return drv.db.Close()
}

// Scan implements driver.Scan
func (drv Driver) Scan(opts driver.ScanOpts) {
	if opts.Filter == nil {
		return
	}

	var iter iterator.Iterator
	var next func() bool

	if opts.Prefix != nil {
		iter = drv.db.NewIterator(util.BytesPrefix(opts.Prefix), nil)
	} else {
		iter = drv.db.NewIterator(nil, nil)
	}

	if opts.ReverseScan {
		next = iter.Prev
	} else {
		next = iter.Next
	}

	if opts.Offset != nil {
		iter.Seek(opts.Offset)
	}

	if opts.ReverseScan && opts.Offset == nil && opts.Prefix == nil {
		iter.Last()
	}

	if opts.Offset != nil && !opts.IncludeOffset {
		next()
	}

	defer iter.Release()
	for next() {
		if err := iter.Error(); err != nil {
			break
		}

		if !iter.Valid() {
			break
		}

		_k, _v := iter.Key(), iter.Value()

		if _k == nil {
			break
		}

		newK := make([]byte, len(_k))
		newV := make([]byte, len(_v))

		copy(newK, _k)
		copy(newV, _v)

		if !opts.Filter(newK, newV) {
			break
		}
	}
}