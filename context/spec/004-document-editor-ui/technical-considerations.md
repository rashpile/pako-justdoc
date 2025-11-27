# Technical Specification: Document Editor UI

- **Functional Specification:** `context/spec/004-document-editor-ui/functional-spec.md`
- **Status:** Approved
- **Author(s):** Claude

---

## 1. High-Level Technical Approach

This feature adds a lightweight web UI for editing JSON documents. The implementation follows the existing architecture principles:

1. **UI Handler**: Add `GET /_/edit/{channel}/{document}` route serving an HTML editor page
2. **Static File Serving**: Add `GET /_/static/` route for CSS/JS assets
3. **Embedded Assets**: Use Go's `embed` package to bundle assets into the binary
4. **Custom JS Highlighting**: ~50 lines of vanilla JS for JSON syntax highlighting

No external dependencies — all assets are embedded in the single binary.

**Note:** Routes use `/_/` prefix to avoid conflicts with channel names (e.g., a channel named "static" or "edit").

---

## 2. Proposed Solution & Implementation Plan

### 2.1 Project Structure Changes

```
justdoc/
├── internal/
│   ├── api/
│   │   ├── handler.go
│   │   ├── router.go
│   │   ├── ui.go              # NEW - UI handler + static serving
│   │   └── static/            # NEW - embedded static files
│   │       ├── editor.html    # HTML template
│   │       ├── editor.css     # Styles (~2KB)
│   │       └── editor.js      # Logic + highlighting (~3KB)
```

### 2.2 Router Changes

**File:** `internal/api/router.go`

```go
func NewRouter(h *Handler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /openapi.json", OpenAPI)
    mux.HandleFunc("GET /_/static/", ServeStatic)                      // NEW
    mux.HandleFunc("GET /_/edit/{channel}/{document}", h.EditorUI)     // NEW
    mux.HandleFunc("GET /", h.ListChannels)
    mux.HandleFunc("GET /{channel}/", h.ListDocuments)
    mux.HandleFunc("GET /{channel}/{document}", h.GetDocument)
    mux.HandleFunc("POST /{channel}/{document}", h.PostDocument)
    return mux
}
```

### 2.3 UI Handler Implementation

**File:** `internal/api/ui.go`

```go
package api

import (
    "embed"
    "html/template"
    "io/fs"
    "net/http"

    "github.com/rashpile/pako-justdoc/internal/model"
)

//go:embed static/*
var staticFiles embed.FS

var editorTemplate *template.Template

func init() {
    editorTemplate = template.Must(template.ParseFS(staticFiles, "static/editor.html"))
}

// ServeStatic serves embedded static files
func ServeStatic(w http.ResponseWriter, r *http.Request) {
    // Strip /_/static/ prefix and serve from embedded fs
    sub, _ := fs.Sub(staticFiles, "static")
    http.StripPrefix("/_/static/", http.FileServer(http.FS(sub))).ServeHTTP(w, r)
}

// EditorUI serves the document editor HTML page
func (h *Handler) EditorUI(w http.ResponseWriter, r *http.Request) {
    channel := r.PathValue("channel")
    document := r.PathValue("document")

    if !model.IsValidName(channel) || !model.IsValidName(document) {
        http.Error(w, "Invalid channel or document name", http.StatusBadRequest)
        return
    }

    data := struct {
        Channel  string
        Document string
    }{
        Channel:  channel,
        Document: document,
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    _ = editorTemplate.Execute(w, data)
}
```

### 2.4 HTML Template

**File:** `internal/api/static/editor.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Channel}} / {{.Document}} - JustDoc</title>
    <link rel="stylesheet" href="/_/static/editor.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>{{.Channel}} / {{.Document}}</h1>
        </header>
        <main>
            <div class="editor-wrapper">
                <pre id="highlight" aria-hidden="true"></pre>
                <textarea id="editor" spellcheck="false"></textarea>
            </div>
        </main>
        <footer>
            <button id="save-btn">Save</button>
            <span id="status"></span>
        </footer>
    </div>
    <script>
        window.CHANNEL = "{{.Channel}}";
        window.DOCUMENT = "{{.Document}}";
    </script>
    <script src="/_/static/editor.js"></script>
</body>
</html>
```

### 2.5 CSS Styles

**File:** `internal/api/static/editor.css` (~2KB)

```css
* { box-sizing: border-box; margin: 0; padding: 0; }

body {
    font-family: system-ui, sans-serif;
    background: #f5f5f5;
    min-height: 100vh;
}

.container {
    max-width: 900px;
    margin: 0 auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}

header h1 {
    font-size: 1.2rem;
    font-weight: 500;
    color: #333;
    padding: 10px 0;
}

main { flex: 1; display: flex; flex-direction: column; }

.editor-wrapper {
    position: relative;
    flex: 1;
    min-height: 400px;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: #fff;
}

#editor, #highlight {
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    padding: 12px;
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 14px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-wrap: break-word;
    overflow: auto;
}

#editor {
    background: transparent;
    color: transparent;
    caret-color: #333;
    border: none;
    resize: none;
    z-index: 1;
}

#highlight {
    color: #333;
    pointer-events: none;
}

.json-key { color: #881391; }
.json-string { color: #1a1aa6; }
.json-number { color: #1c6b48; }
.json-boolean { color: #d73a49; }
.json-null { color: #6f42c1; }

footer {
    padding: 15px 0;
    display: flex;
    align-items: center;
    gap: 15px;
}

#save-btn {
    padding: 8px 24px;
    background: #0066cc;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
}

#save-btn:hover { background: #0052a3; }
#save-btn:disabled { background: #ccc; cursor: not-allowed; }

#status { font-size: 14px; }
#status.success { color: #1c6b48; }
#status.error { color: #d73a49; }
```

### 2.6 JavaScript Logic

**File:** `internal/api/static/editor.js` (~3KB)

```javascript
(function() {
    const editor = document.getElementById('editor');
    const highlight = document.getElementById('highlight');
    const saveBtn = document.getElementById('save-btn');
    const status = document.getElementById('status');
    const apiUrl = '/' + CHANNEL + '/' + DOCUMENT;

    // Syntax highlighting
    function highlightJSON(text) {
        if (!text) return '';
        return text
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/("(?:\\.|[^"\\])*")\s*:/g, '<span class="json-key">$1</span>:')
            .replace(/:(\s*)("(?:\\.|[^"\\])*")/g, ':$1<span class="json-string">$2</span>')
            .replace(/:\s*(-?\d+\.?\d*)/g, ': <span class="json-number">$1</span>')
            .replace(/:\s*(true|false)/g, ': <span class="json-boolean">$1</span>')
            .replace(/:\s*(null)/g, ': <span class="json-null">$1</span>');
    }

    function updateHighlight() {
        highlight.innerHTML = highlightJSON(editor.value) + '\n';
    }

    function syncScroll() {
        highlight.scrollTop = editor.scrollTop;
        highlight.scrollLeft = editor.scrollLeft;
    }

    // Load document
    async function loadDocument() {
        try {
            const res = await fetch(apiUrl);
            if (res.ok) {
                const data = await res.json();
                editor.value = JSON.stringify(data, null, 2);
            } else if (res.status === 404) {
                editor.value = '';
            }
            updateHighlight();
        } catch (e) {
            setStatus('Failed to load document', true);
        }
    }

    // Save document
    async function saveDocument() {
        const content = editor.value.trim();

        // Validate JSON
        try {
            if (content) JSON.parse(content);
        } catch (e) {
            setStatus('Invalid JSON: ' + e.message, true);
            return;
        }

        saveBtn.disabled = true;
        try {
            const res = await fetch(apiUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: content || '{}'
            });
            if (res.ok) {
                setStatus('Saved successfully', false);
            } else {
                const err = await res.json();
                setStatus('Error: ' + err.message, true);
            }
        } catch (e) {
            setStatus('Failed to save', true);
        }
        saveBtn.disabled = false;
    }

    function setStatus(msg, isError) {
        status.textContent = msg;
        status.className = isError ? 'error' : 'success';
        if (!isError) setTimeout(() => status.textContent = '', 3000);
    }

    // Event listeners
    editor.addEventListener('input', updateHighlight);
    editor.addEventListener('scroll', syncScroll);
    saveBtn.addEventListener('click', saveDocument);

    // Ctrl+S to save
    document.addEventListener('keydown', function(e) {
        if ((e.ctrlKey || e.metaKey) && e.key === 's') {
            e.preventDefault();
            saveDocument();
        }
    });

    // Initial load
    loadDocument();
})();
```

---

## 3. Impact and Risk Analysis

### System Dependencies
- **Existing API endpoints**: UI uses `GET` and `POST /{channel}/{document}` — no changes needed
- **Go embed**: Requires Go 1.16+ (we use 1.25)
- **Single binary**: Static assets are embedded, no external files required

### Potential Risks & Mitigations

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Route conflict with channel names | None | `/_/` prefix reserved for internal routes |
| Large JSON performance | Low | Browser handles rendering; no server-side processing |
| Embed increases binary size | Low | ~5KB total for HTML+CSS+JS |

---

## 4. Testing Strategy

### Handler Tests (`internal/api/ui_test.go`)
- `TestEditorUI_ReturnsHTML` — Verify Content-Type is `text/html`
- `TestEditorUI_IncludesChannelAndDocument` — Verify template variables are injected
- `TestEditorUI_InvalidName_Returns400` — Verify validation works
- `TestServeStatic_CSS` — Verify CSS file is served

### Manual Testing
- Load editor for existing document — verify content appears
- Load editor for non-existent document — verify blank editor
- Edit and save valid JSON — verify success message
- Edit and save invalid JSON — verify error message
- Verify syntax highlighting works for keys, strings, numbers, booleans
