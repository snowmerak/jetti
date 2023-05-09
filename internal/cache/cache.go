package cache

import (
	"encoding/binary"
	"github.com/dgraph-io/badger/v4"
	"time"
)

type Cache struct {
	db *badger.DB
}

func NewCache(path string) (*Cache, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &Cache{
		db: db,
	}, nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}

func (c *Cache) Get(key string) (int64, bool) {
	var value int64
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = int64(binary.BigEndian.Uint64(val))
			return nil
		})
	})
	if err != nil {
		return 0, false
	}
	return value, true
}

func (c *Cache) Set(key string, value int64) error {
	return c.db.Update(func(txn *badger.Txn) error {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(value))
		e := badger.NewEntry([]byte(key), buf).WithTTL(time.Hour * 24 * 30)
		return txn.SetEntry(e)
	})
}

func (c *Cache) Delete(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (c *Cache) CompareAndSet(key string, value int64) error {
	return c.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		var oldValue int64
		err = item.Value(func(val []byte) error {
			oldValue = int64(binary.BigEndian.Uint64(val))
			return nil
		})
		if err != nil {
			return err
		}
		if oldValue >= value {
			return nil
		}
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(value))
		return txn.Set([]byte(key), buf)
	})
}

func (c *Cache) NeedToUpdate(key string, value int64) bool {
	var oldValue int64
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			oldValue = int64(binary.BigEndian.Uint64(val))
			return nil
		})
	})
	if err != nil {
		return false
	}
	return oldValue < value
}
