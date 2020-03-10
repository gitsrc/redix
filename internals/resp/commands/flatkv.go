package commands

import (
	"strconv"
	"time"

	"github.com/alash3al/redix/internals/db"
	"github.com/alash3al/redix/internals/resp"
)

func init() {
	prefix := usrPrefix

	resp.Handlers["set"] = func(c *resp.Context) {
		args := c.Args()

		if len(args) < 2 {
			c.Conn().WriteError(incorrectArgsCount)
			return
		}

		entry := db.Entry{
			Key:   append(prefix, args[0]...),
			Value: args[1],
		}

		if len(args) > 2 {
			parsedDur, err := time.ParseDuration(string(args[2]))
			if err != nil {
				c.Conn().WriteError(err.Error())
				return
			}

			entry.TTL = parsedDur
		}

		c.DB().Put(&entry)
		c.Conn().WriteString("OK")
	}

	resp.Handlers["get"] = func(c *resp.Context) {
		args := c.Args()

		if len(args) < 1 {
			c.Conn().WriteError(incorrectArgsCount)
			return
		}

		k := append(prefix, args[0]...)
		val, err := c.DB().Get(k)

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		if val == nil {
			c.Conn().WriteNull()
			return
		}

		c.Conn().WriteBulk(val)
	}

	resp.Handlers["del"] = func(c *resp.Context) {
		args := c.Args()

		if len(args) < 1 {
			c.Conn().WriteInt(0)
			return
		}

		entries := []*db.Entry{}
		for _, k := range args {
			entry := db.Entry{Key: append(prefix, k...)}
			entries = append(entries, &entry)
		}

		c.DB().Batch(entries)
		c.Conn().WriteInt(len(entries))
	}

	resp.Handlers["incr"] = func(c *resp.Context) {
		args := c.Args()

		if len(args) < 1 {
			c.Conn().WriteError(incorrectArgsCount)
			return
		}

		entry := db.Entry{
			Key: append(prefix, args[0]...),
		}

		delta := float64(0)

		if len(args) > 1 {
			delta, _ = strconv.ParseFloat(string(args[1]), 64)
		}

		if len(args) > 2 {
			parsedDur, err := time.ParseDuration(string(args[2]))
			if err != nil {
				c.Conn().WriteError(err.Error())
				return
			}
			entry.TTL = parsedDur
		}

		newVal, err := c.DB().Incr(entry.Key, delta, entry.TTL)
		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		if float64(int64(*newVal)) == *newVal {
			c.Conn().WriteInt64(int64(*newVal))
		} else {
			c.Conn().WriteBulkString(strconv.FormatFloat(*newVal, 'f', -1, 64))
		}
	}
}