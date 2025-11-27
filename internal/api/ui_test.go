package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEditorUI_ReturnsHTML(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/myapp/settings/ui", nil)
	req.SetPathValue("channel", "myapp")
	req.SetPathValue("document", "settings")

	w := httptest.NewRecorder()
	handler.EditorUI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		t.Errorf("Expected Content-Type to start with 'text/html', got %q", contentType)
	}
}

func TestEditorUI_IncludesChannelAndDocument(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/testchannel/testdoc/ui", nil)
	req.SetPathValue("channel", "testchannel")
	req.SetPathValue("document", "testdoc")

	w := httptest.NewRecorder()
	handler.EditorUI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()

	// Check for channel and document in title
	if !strings.Contains(body, "<title>testchannel / testdoc - JustDoc</title>") {
		t.Error("Expected title to contain 'testchannel / testdoc - JustDoc'")
	}

	// Check for channel and document in header
	if !strings.Contains(body, "<h1>testchannel / testdoc</h1>") {
		t.Error("Expected h1 to contain 'testchannel / testdoc'")
	}

	// Check for JavaScript variables
	if !strings.Contains(body, `window.CHANNEL = "testchannel"`) {
		t.Error("Expected JavaScript to contain window.CHANNEL = 'testchannel'")
	}
	if !strings.Contains(body, `window.DOCUMENT = "testdoc"`) {
		t.Error("Expected JavaScript to contain window.DOCUMENT = 'testdoc'")
	}
}

func TestEditorUI_InvalidName_Returns400(t *testing.T) {
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
		{"empty document", "channel", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test/doc/ui", nil)
			req.SetPathValue("channel", tt.channel)
			req.SetPathValue("document", tt.document)

			w := httptest.NewRecorder()
			handler.EditorUI(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}

			body := w.Body.String()
			if !strings.Contains(body, "Invalid channel or document name") {
				t.Errorf("Expected error message about invalid name, got %q", body)
			}
		})
	}
}

func TestServeStatic_CSS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_/static/editor.css", nil)
	w := httptest.NewRecorder()

	ServeStatic(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/css") {
		t.Errorf("Expected Content-Type to contain 'text/css', got %q", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, ".container") {
		t.Error("Expected CSS to contain '.container' class")
	}
	if !strings.Contains(body, "#editor") {
		t.Error("Expected CSS to contain '#editor' selector")
	}
}
