package cruncy

import (
	"errors"
	"fmt"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Store is db abstarction on top of bolt db
type Store struct {
	db    *bolt.DB
	mutex sync.Mutex
}

var (
	// ErrNotFound error key not found
	ErrNotFound = errors.New("store: key not found")
	// ErrBadValue error bad value
	ErrBadValue = errors.New("store: bad value")
)

// Open a database file
func Open(path string) (*Store, error) {
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}
	db, err := bolt.Open(path, 0640, opts)

	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// CreateBucket creates a buck
func (store *Store) CreateBucket(bucket string) error {
	if store == nil {
		return fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return fmt.Errorf("Missing store db object")
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

// Close the database
func (store *Store) Close() error {
	if store == nil {
		return fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return fmt.Errorf("Missing store db object")
	}

	return store.db.Close()
}

// Put a key/value into a given bucket
func (store *Store) Put(bucket string, key string, value string) error {
	if store == nil {
		return fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return fmt.Errorf("Missing store db object")
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Put([]byte(key), ([]byte(value)))
	})
}

// Get a key/value from a given bucket
func (store *Store) Get(bucket, key string, value *string) error {
	if value == nil {
		return fmt.Errorf("store.Get requires pointer to a string. Param must be given")
	}
	if store == nil {
		return fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return fmt.Errorf("Missing store db object")
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		c := b.Cursor()
		k, v := c.Seek([]byte(key))
		if k == nil || string(k) != key {
			return ErrNotFound
		}
		*value = string(v[:])
		return nil

	})
}

// Delete the entry with the given key. If no such key is present in the store,
// it returns ErrNotFound.
//
//	store.Delete("key42")
func (store *Store) Delete(bucket string, key string) error {
	if store == nil {
		return fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return fmt.Errorf("Missing store db object")
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		c := b.Cursor()
		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		}
		return c.Delete()

	})
}

// ForEach iterates over a given bucket
func (store *Store) ForEach(bucket string, fn func(k, v []byte) error) error {
	if store == nil {
		return fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return fmt.Errorf("Missing store db object")
	}

	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("No bucket found")
		}
		return b.ForEach(fn)
	})
	return nil
}

// ListBuckets lists all buckets
func (store *Store) ListBuckets() ([]string, error) {
	if store == nil {
		return nil, fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return nil, fmt.Errorf("Missing store db object")
	}

	rc := []string{}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	err := store.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			rc = append(rc, string(name))
			return nil
		})
	})

	return rc, err
}

// ListBucket lists all buckets
func (store *Store) ListBucket(bucket string, filter func(k, v string) (bool, error)) ([]string, error) {
	if store == nil {
		return nil, fmt.Errorf("Missing store object")
	}

	if store.db == nil {
		return nil, fmt.Errorf("Missing store db object")
	}

	rc := []string{}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	err := store.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			add := true
			var err error
			if filter != nil {
				add, err = filter(string(k), string(v))
				if err != nil {
					return err
				}
			}
			if add {
				rc = append(rc, string(k))
			}
		}
		return nil
	})

	return rc, err
}

// GetDB returns pointer to underlaying DB.
func (store *Store) GetDB() *bolt.DB {
	return store.db
}
