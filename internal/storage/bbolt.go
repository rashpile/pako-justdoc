package storage

import (
	"sort"

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

// ListDocuments returns all document names in a channel (sorted alphabetically)
func (s *BoltStorage) ListDocuments(channel string) ([]string, error) {
	var docs []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(channel))
		if bucket == nil {
			return ErrNotFound
		}
		return bucket.ForEach(func(k, v []byte) error {
			docs = append(docs, string(k))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(docs)
	return docs, nil
}

// ListChannels returns all channels with document counts (sorted alphabetically)
func (s *BoltStorage) ListChannels() ([]ChannelInfo, error) {
	channels := make([]ChannelInfo, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			count := b.Stats().KeyN
			channels = append(channels, ChannelInfo{
				Name:          string(name),
				DocumentCount: count,
			})
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Name < channels[j].Name
	})
	return channels, nil
}

// Close closes the database connection
func (s *BoltStorage) Close() error {
	return s.db.Close()
}
