package driver

// Interface an interface describes a storage backend
type Interface interface {
	Open(string, map[string]interface{}) (Interface, error)
	Put([]byte, []byte) error
	Get([]byte) ([]byte, error)
	BulkPut([]KeyValue) error
	Delete([]byte) error
	Scan(ScanOpts)
	Close() error
}

// KeyValue represents a key - value pair
type KeyValue struct {
	Key   []byte
	Value []byte
}

// Registry a registry for available drivers
var Registry = map[string]Interface{}