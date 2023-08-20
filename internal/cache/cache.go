package cache

import (
	"encoding/binary"
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/snowmerak/jetti/v2/gen/grpc/model/store"
	"github.com/snowmerak/jetti/v2/lib/model"
	"google.golang.org/protobuf/proto"
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

func (c *Cache) GetBytes(key string) ([]byte, bool) {
	var value []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = val
			return nil
		})
	})
	if err != nil {
		return nil, false
	}
	return value, true
}

func (c *Cache) SetBytes(key string, value []byte) error {
	return c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		return txn.SetEntry(e)
	})
}

func (c *Cache) SetInterface(key string, value model.InterfaceTransferObject) error {
	return c.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(key), data)
		return txn.SetEntry(e)
	})
}

func (c *Cache) GetInterface(key string) (model.InterfaceTransferObject, bool) {
	var value model.InterfaceTransferObject
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &value)
		})
	})
	if err != nil {
		return model.InterfaceTransferObject{}, false
	}
	return value, true
}

func (c *Cache) SetInterfaceNames(key string, value []string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		data, err := proto.Marshal(&store.StringList{
			Values: value,
		})
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(key), data)
		return txn.SetEntry(e)
	})
}

func (c *Cache) GetInterfaceNames(key string) ([]string, bool) {
	var value store.StringList
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &value)
		})
	})
	if err != nil {
		return nil, false
	}
	return value.Values, true
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
