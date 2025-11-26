package storage

import (
	"go.etcd.io/bbolt"
)

// BoltStorage implements Storage using bbolt
type BoltStorage struct {
	db *bbolt.DB
}

// NewBoltStorage creates a new bbolt-backed storage
func NewBoltStorage(path string) (*BoltStorage, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltStorage{db: db}, nil
}

// GetDocument retrieves a document from a channel
func (s *BoltStorage) GetDocument(channel, document string) ([]byte, error) {
	var data []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(channel))
		if bucket == nil {
			return ErrNotFound
		}
		v := bucket.Get([]byte(document))
		if v == nil {
			return ErrNotFound
		}
		// Copy the data since bbolt values are only valid during the transaction
		data = make([]byte, len(v))
		copy(data, v)
		return nil
	})
	return data, err
}

// PutDocument stores a document in a channel
func (s *BoltStorage) PutDocument(channel, document string, data []byte) (bool, error) {
	var created bool
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(channel))
		if err != nil {
			return err
		}
		existing := bucket.Get([]byte(document))
		created = existing == nil
		return bucket.Put([]byte(document), data)
	})
	return created, err
}

// Close closes the database connection
func (s *BoltStorage) Close() error {
	return s.db.Close()
}
