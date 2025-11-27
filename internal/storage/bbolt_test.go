package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBoltStorage(t *testing.T) {
	// Create temp directory for test database
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	t.Run("GetDocument_NotFound_NoChannel", func(t *testing.T) {
		_, err := storage.GetDocument("nonexistent", "doc")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("PutDocument_Create", func(t *testing.T) {
		data := []byte(`{"key": "value"}`)
		created, err := storage.PutDocument("channel1", "doc1", data)
		if err != nil {
			t.Fatalf("PutDocument failed: %v", err)
		}
		if !created {
			t.Error("Expected created=true for new document")
		}
	})

	t.Run("GetDocument_Success", func(t *testing.T) {
		data, err := storage.GetDocument("channel1", "doc1")
		if err != nil {
			t.Fatalf("GetDocument failed: %v", err)
		}
		expected := `{"key": "value"}`
		if string(data) != expected {
			t.Errorf("Got %q, want %q", string(data), expected)
		}
	})

	t.Run("PutDocument_Update", func(t *testing.T) {
		data := []byte(`{"key": "updated"}`)
		created, err := storage.PutDocument("channel1", "doc1", data)
		if err != nil {
			t.Fatalf("PutDocument failed: %v", err)
		}
		if created {
			t.Error("Expected created=false for existing document")
		}

		// Verify the update
		retrieved, err := storage.GetDocument("channel1", "doc1")
		if err != nil {
			t.Fatalf("GetDocument failed: %v", err)
		}
		if string(retrieved) != string(data) {
			t.Errorf("Got %q, want %q", string(retrieved), string(data))
		}
	})

	t.Run("GetDocument_NotFound_NoDocument", func(t *testing.T) {
		_, err := storage.GetDocument("channel1", "nonexistent")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("PutDocument_AutoCreateChannel", func(t *testing.T) {
		data := []byte(`{"auto": "created"}`)
		created, err := storage.PutDocument("newchannel", "doc1", data)
		if err != nil {
			t.Fatalf("PutDocument failed: %v", err)
		}
		if !created {
			t.Error("Expected created=true for new document in new channel")
		}

		// Verify retrieval
		retrieved, err := storage.GetDocument("newchannel", "doc1")
		if err != nil {
			t.Fatalf("GetDocument failed: %v", err)
		}
		if string(retrieved) != string(data) {
			t.Errorf("Got %q, want %q", string(retrieved), string(data))
		}
	})
}

func TestBoltStorage_MultipleDocuments(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// Store multiple documents in same channel
	docs := map[string]string{
		"doc1": `{"id": 1}`,
		"doc2": `{"id": 2}`,
		"doc3": `{"id": 3}`,
	}

	for name, content := range docs {
		_, err := storage.PutDocument("multi", name, []byte(content))
		if err != nil {
			t.Fatalf("Failed to store %s: %v", name, err)
		}
	}

	// Verify all documents
	for name, expected := range docs {
		data, err := storage.GetDocument("multi", name)
		if err != nil {
			t.Fatalf("Failed to get %s: %v", name, err)
		}
		if string(data) != expected {
			t.Errorf("Doc %s: got %q, want %q", name, string(data), expected)
		}
	}
}

func TestListDocuments_ReturnsAlphabeticalOrder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// Store documents in non-alphabetical order
	docs := []string{"zebra", "alpha", "middle", "beta"}
	for _, name := range docs {
		_, err := storage.PutDocument("testchan", name, []byte(`{}`))
		if err != nil {
			t.Fatalf("Failed to store %s: %v", name, err)
		}
	}

	// List documents
	result, err := storage.ListDocuments("testchan")
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}

	// Verify alphabetical order
	expected := []string{"alpha", "beta", "middle", "zebra"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d documents, got %d", len(expected), len(result))
	}
	for i, name := range expected {
		if result[i] != name {
			t.Errorf("Position %d: expected %q, got %q", i, name, result[i])
		}
	}
}

func TestListDocuments_ChannelNotFound_ReturnsError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// Try to list documents in non-existent channel
	_, err = storage.ListDocuments("nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestListDocuments_EmptyChannel(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// Store a document and then list (bbolt creates bucket with first doc)
	_, err = storage.PutDocument("emptychan", "doc1", []byte(`{}`))
	if err != nil {
		t.Fatalf("Failed to store document: %v", err)
	}

	// List documents - should have one document
	result, err := storage.ListDocuments("emptychan")
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}

	if len(result) != 1 || result[0] != "doc1" {
		t.Errorf("Expected [doc1], got %v", result)
	}
}

func TestListChannels_ReturnsAlphabeticalOrder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// Create channels in non-alphabetical order
	channels := []string{"zebra", "alpha", "middle", "beta"}
	for _, channel := range channels {
		_, err := storage.PutDocument(channel, "doc1", []byte(`{}`))
		if err != nil {
			t.Fatalf("Failed to create channel %s: %v", channel, err)
		}
	}

	// List channels
	result, err := storage.ListChannels()
	if err != nil {
		t.Fatalf("ListChannels failed: %v", err)
	}

	// Verify alphabetical order
	expected := []string{"alpha", "beta", "middle", "zebra"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d channels, got %d", len(expected), len(result))
	}
	for i, name := range expected {
		if result[i].Name != name {
			t.Errorf("Position %d: expected %q, got %q", i, name, result[i].Name)
		}
	}
}

func TestListChannels_NoChannels_ReturnsEmptySlice(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// List channels on empty database
	result, err := storage.ListChannels()
	if err != nil {
		t.Fatalf("ListChannels failed: %v", err)
	}

	// Verify empty slice (not nil)
	if result == nil {
		t.Error("Expected empty slice, got nil")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d channels", len(result))
	}
}

func TestListChannels_IncludesDocumentCount(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "justdoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	storage, err := NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("Failed to close storage: %v", err)
		}
	}()

	// Create channel1 with 3 documents
	for i := 1; i <= 3; i++ {
		_, err := storage.PutDocument("channel1", "doc"+string(rune('0'+i)), []byte(`{}`))
		if err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}
	}

	// Create channel2 with 1 document
	_, err = storage.PutDocument("channel2", "doc1", []byte(`{}`))
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	// List channels
	result, err := storage.ListChannels()
	if err != nil {
		t.Fatalf("ListChannels failed: %v", err)
	}

	// Verify document counts
	if len(result) != 2 {
		t.Fatalf("Expected 2 channels, got %d", len(result))
	}

	// Find channel1 and channel2 in results
	var ch1, ch2 *ChannelInfo
	for i := range result {
		if result[i].Name == "channel1" {
			ch1 = &result[i]
		}
		if result[i].Name == "channel2" {
			ch2 = &result[i]
		}
	}

	if ch1 == nil {
		t.Fatal("channel1 not found in results")
	}
	if ch1.DocumentCount != 3 {
		t.Errorf("channel1: expected document_count=3, got %d", ch1.DocumentCount)
	}

	if ch2 == nil {
		t.Fatal("channel2 not found in results")
	}
	if ch2.DocumentCount != 1 {
		t.Errorf("channel2: expected document_count=1, got %d", ch2.DocumentCount)
	}
}
