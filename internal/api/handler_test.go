package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rashpile/pako-justdoc/internal/model"
	"github.com/rashpile/pako-justdoc/internal/storage"
)

// setupTestHandler creates a handler with a temporary bbolt database
func setupTestHandler(t *testing.T) (*Handler, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "justdoc-api-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := storage.NewBoltStorage(dbPath)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	handler := NewHandler(store)
	cleanup := func() {
		_ = store.Close()
		_ = os.RemoveAll(tmpDir)
	}

	return handler, cleanup
}

func TestPostDocument_Create(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{"theme": "dark"}`
	req := httptest.NewRequest(http.MethodPost, "/myapp/settings", strings.NewReader(body))
	req.SetPathValue("channel", "myapp")
	req.SetPathValue("document", "settings")

	w := httptest.NewRecorder()
	handler.PostDocument(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp model.SuccessResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Status != "created" {
		t.Errorf("Expected status 'created', got %q", resp.Status)
	}
	if resp.Channel != "myapp" {
		t.Errorf("Expected channel 'myapp', got %q", resp.Channel)
	}
	if resp.Document != "settings" {
		t.Errorf("Expected document 'settings', got %q", resp.Document)
	}
}

func TestPostDocument_Update(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// First create the document
	body1 := `{"theme": "dark"}`
	req1 := httptest.NewRequest(http.MethodPost, "/myapp/settings", strings.NewReader(body1))
	req1.SetPathValue("channel", "myapp")
	req1.SetPathValue("document", "settings")
	w1 := httptest.NewRecorder()
	handler.PostDocument(w1, req1)

	// Now update it
	body2 := `{"theme": "light"}`
	req2 := httptest.NewRequest(http.MethodPost, "/myapp/settings", strings.NewReader(body2))
	req2.SetPathValue("channel", "myapp")
	req2.SetPathValue("document", "settings")
	w2 := httptest.NewRecorder()
	handler.PostDocument(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w2.Code)
	}

	var resp model.SuccessResponse
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Status != "updated" {
		t.Errorf("Expected status 'updated', got %q", resp.Status)
	}
}

func TestPostDocument_InvalidJSON(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/test/doc", strings.NewReader(body))
	req.SetPathValue("channel", "test")
	req.SetPathValue("document", "doc")

	w := httptest.NewRecorder()
	handler.PostDocument(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != model.ErrCodeInvalidJSON {
		t.Errorf("Expected error code %q, got %q", model.ErrCodeInvalidJSON, resp.Error)
	}
}

func TestPostDocument_InvalidName(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name     string
		channel  string
		document string
	}{
		{"invalid channel with @", "test@invalid", "doc"},
		{"invalid document with space", "test", "doc space"},
		{"invalid channel with dot", "test.channel", "doc"},
		{"empty channel", "", "doc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := `{"test": "data"}`
			// Use a valid URL path but set invalid path values
			req := httptest.NewRequest(http.MethodPost, "/test/doc", strings.NewReader(body))
			req.SetPathValue("channel", tt.channel)
			req.SetPathValue("document", tt.document)

			w := httptest.NewRecorder()
			handler.PostDocument(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}

			var resp model.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if resp.Error != model.ErrCodeInvalidName {
				t.Errorf("Expected error code %q, got %q", model.ErrCodeInvalidName, resp.Error)
			}
		})
	}
}

func TestPostDocument_PayloadTooLarge(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create a body larger than 10MB
	largeData := make([]byte, 11*1024*1024)
	for i := range largeData {
		largeData[i] = 'a'
	}
	body := `{"data": "` + string(largeData) + `"}`

	req := httptest.NewRequest(http.MethodPost, "/test/large", bytes.NewReader([]byte(body)))
	req.SetPathValue("channel", "test")
	req.SetPathValue("document", "large")

	w := httptest.NewRecorder()
	handler.PostDocument(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status %d, got %d", http.StatusRequestEntityTooLarge, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != model.ErrCodePayloadTooLarge {
		t.Errorf("Expected error code %q, got %q", model.ErrCodePayloadTooLarge, resp.Error)
	}
}

func TestGetDocument_Success(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// First store a document
	body := `{"theme": "dark"}`
	postReq := httptest.NewRequest(http.MethodPost, "/myapp/settings", strings.NewReader(body))
	postReq.SetPathValue("channel", "myapp")
	postReq.SetPathValue("document", "settings")
	postW := httptest.NewRecorder()
	handler.PostDocument(postW, postReq)

	// Now retrieve it
	getReq := httptest.NewRequest(http.MethodGet, "/myapp/settings", nil)
	getReq.SetPathValue("channel", "myapp")
	getReq.SetPathValue("document", "settings")
	getW := httptest.NewRecorder()
	handler.GetDocument(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, getW.Code)
	}

	if getW.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %q", getW.Header().Get("Content-Type"))
	}

	if getW.Body.String() != body {
		t.Errorf("Expected body %q, got %q", body, getW.Body.String())
	}
}

func TestGetDocument_NotFound(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent/doc", nil)
	req.SetPathValue("channel", "nonexistent")
	req.SetPathValue("document", "doc")

	w := httptest.NewRecorder()
	handler.GetDocument(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != model.ErrCodeNotFound {
		t.Errorf("Expected error code %q, got %q", model.ErrCodeNotFound, resp.Error)
	}
}

func TestGetDocument_InvalidName(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/test@invalid/doc", nil)
	req.SetPathValue("channel", "test@invalid")
	req.SetPathValue("document", "doc")

	w := httptest.NewRecorder()
	handler.GetDocument(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != model.ErrCodeInvalidName {
		t.Errorf("Expected error code %q, got %q", model.ErrCodeInvalidName, resp.Error)
	}
}

func TestListDocuments_Success(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// First store some documents
	docs := []string{"doc1", "doc3", "doc2"} // Store in non-alphabetical order
	for _, doc := range docs {
		body := `{"test": "data"}`
		req := httptest.NewRequest(http.MethodPost, "/mychannel/"+doc, strings.NewReader(body))
		req.SetPathValue("channel", "mychannel")
		req.SetPathValue("document", doc)
		w := httptest.NewRecorder()
		handler.PostDocument(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create document %s: status %d", doc, w.Code)
		}
	}

	// Now list documents
	req := httptest.NewRequest(http.MethodGet, "/mychannel/", nil)
	req.SetPathValue("channel", "mychannel")
	w := httptest.NewRecorder()
	handler.ListDocuments(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %q", w.Header().Get("Content-Type"))
	}

	var result []string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify alphabetical order
	expected := []string{"doc1", "doc2", "doc3"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d documents, got %d", len(expected), len(result))
	}
	for i, name := range expected {
		if result[i] != name {
			t.Errorf("Position %d: expected %q, got %q", i, name, result[i])
		}
	}
}

func TestListDocuments_ChannelNotFound_Returns404(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent/", nil)
	req.SetPathValue("channel", "nonexistent")

	w := httptest.NewRecorder()
	handler.ListDocuments(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != model.ErrCodeNotFound {
		t.Errorf("Expected error code %q, got %q", model.ErrCodeNotFound, resp.Error)
	}
}

func TestListDocuments_InvalidChannelName_Returns400(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/invalid@channel/", nil)
	req.SetPathValue("channel", "invalid@channel")

	w := httptest.NewRecorder()
	handler.ListDocuments(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != model.ErrCodeInvalidName {
		t.Errorf("Expected error code %q, got %q", model.ErrCodeInvalidName, resp.Error)
	}
}

func TestListChannels_Success(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create some channels with documents
	channels := []struct {
		name     string
		docCount int
	}{
		{"zebra", 2},
		{"alpha", 1},
		{"beta", 3},
	}

	for _, ch := range channels {
		for i := 1; i <= ch.docCount; i++ {
			body := `{"test": "data"}`
			req := httptest.NewRequest(http.MethodPost, "/"+ch.name+"/doc", strings.NewReader(body))
			req.SetPathValue("channel", ch.name)
			req.SetPathValue("document", "doc"+string(rune('0'+i)))
			w := httptest.NewRecorder()
			handler.PostDocument(w, req)
			if w.Code != http.StatusCreated && w.Code != http.StatusOK {
				t.Fatalf("Failed to create document in %s: status %d", ch.name, w.Code)
			}
		}
	}

	// List channels
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ListChannels(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %q", w.Header().Get("Content-Type"))
	}

	var result []storage.ChannelInfo
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify alphabetical order and document counts
	if len(result) != 3 {
		t.Fatalf("Expected 3 channels, got %d", len(result))
	}

	expected := []storage.ChannelInfo{
		{Name: "alpha", DocumentCount: 1},
		{Name: "beta", DocumentCount: 3},
		{Name: "zebra", DocumentCount: 2},
	}

	for i, exp := range expected {
		if result[i].Name != exp.Name {
			t.Errorf("Position %d: expected name %q, got %q", i, exp.Name, result[i].Name)
		}
		if result[i].DocumentCount != exp.DocumentCount {
			t.Errorf("Channel %q: expected document_count=%d, got %d", exp.Name, exp.DocumentCount, result[i].DocumentCount)
		}
	}
}

func TestListChannels_Empty_ReturnsEmptyArray(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// List channels on empty database
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ListChannels(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %q", w.Header().Get("Content-Type"))
	}

	var result []storage.ChannelInfo
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify empty array (not null)
	if result == nil {
		t.Error("Expected empty array, got null")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty array, got %d channels", len(result))
	}

	// Also verify the raw JSON is "[]" not "null"
	body := strings.TrimSpace(w.Body.String())
	if body != "[]" {
		t.Errorf("Expected JSON '[]', got %q", body)
	}
}
