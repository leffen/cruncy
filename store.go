package cruncy

import (
  "github.com/boltdb/bolt"
  "errors"
  "time"
  "sync"
)

type Store struct {
  db    *bolt.DB
  mutex sync.Mutex
}

var (
  ErrNotFound = errors.New("ression: key not found")
  ErrBadValue = errors.New("ression: bad value")
)

func Open(path string) (*Store, error) {
  opts := &bolt.Options{
    Timeout: 50 * time.Millisecond,
  }
  db, err := bolt.Open(path, 0640, opts)

  if  err != nil {
    return nil, err
  }

  return &Store{db: db}, nil
}

func (store *Store) CreateBucket(bucket string) error {
  return store.db.Update(func(tx *bolt.Tx) error {
      _, err := tx.CreateBucketIfNotExists([]byte(bucket))
      return err
  })
}

func (store *Store) Close() error {
  return store.db.Close()
}

func (store *Store) Put(bucket string, key string, value string) error {
  store.mutex.Lock()
  defer store.mutex.Unlock()

  return store.db.Update(func(tx *bolt.Tx) error {
    return tx.Bucket([]byte(bucket)).Put([]byte(key), ([]byte(value)))
  })
}

func (store *Store) Get(bucket, key string, value *string) error {
  store.mutex.Lock()
  defer store.mutex.Unlock()

  return store.db.View(func(tx *bolt.Tx) error {
    c := tx.Bucket([]byte(bucket)).Cursor()
    if k, v := c.Seek([]byte(key)); k == nil || string(k) != key {
      return ErrNotFound
    } else if value == nil {
      return nil
    } else {
      *value = string(v[:])
      return nil
    }
  })
}

// Delete the entry with the given key. If no such key is present in the store,
// it returns ErrNotFound.
//
//	store.Delete("key42")
func (store *Store) Delete(bucket string,key string) error {
  store.mutex.Lock()
  defer store.mutex.Unlock()

  return store.db.Update(func(tx *bolt.Tx) error {
    c := tx.Bucket([]byte(bucket)).Cursor()
    if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
      return ErrNotFound
    } else {
      return c.Delete()
    }
  })
}

