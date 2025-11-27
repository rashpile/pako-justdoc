package storage

import "errors"

// ErrNotFound is returned when a document or channel is not found
var ErrNotFound = errors.New("not found")

// ChannelInfo represents a channel with its document count
type ChannelInfo struct {
	Name          string `json:"name"`
	DocumentCount int    `json:"document_count"`
}

// Storage defines the document storage interface
type Storage interface {
	// GetDocument retrieves a document from a channel
	// Returns ErrNotFound if channel or document doesn't exist
	GetDocument(channel, document string) ([]byte, error)

	// PutDocument stores a document in a channel
	// Creates the channel if it doesn't exist
	// Returns created=true if document was new, false if updated
	PutDocument(channel, document string, data []byte) (created bool, err error)

	// ListDocuments returns all document names in a channel (sorted alphabetically)
	// Returns ErrNotFound if channel doesn't exist
	ListDocuments(channel string) ([]string, error)

	// ListChannels returns all channels with document counts (sorted alphabetically)
	ListChannels() ([]ChannelInfo, error)

	// Close closes the storage connection
	Close() error
}
