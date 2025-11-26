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
