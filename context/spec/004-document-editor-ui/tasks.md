# Tasks: Document Editor UI

## Slice 1: Basic Editor Page (Static HTML)
_Delivers a working editor page with no functionality — just the structure._

- [x] **Slice 1: Basic Editor Page**
  - [x] Create `internal/api/static/` directory
  - [x] Create `editor.html` template with basic structure (header, textarea, save button, status area)
  - [x] Create `editor.css` with minimalist styling (light theme, monospace font)
  - [x] Create `ui.go` with embedded static files and `EditorUI` handler
  - [x] Create `ServeStatic` handler for `/static/` route
  - [x] Register `GET /static/` and `GET /{channel}/{document}/ui` routes in router
  - [x] Add handler tests for `EditorUI` and `ServeStatic`
  - [ ] Manually verify: navigate to `/test/doc/ui`, see styled editor page

---

## Slice 2: Load and Display Document
_Adds JavaScript to fetch and display existing document content._

- [x] **Slice 2: Load Document Content**
  - [x] Create `editor.js` with `loadDocument()` function
  - [x] Fetch document via `GET /{channel}/{document}` on page load
  - [x] Display document content in textarea (pretty-printed JSON)
  - [x] Handle 404 (show empty textarea for new documents)
  - [x] Handle errors (show status message)
  - [ ] Manually verify: create document via API, open UI, see content

---

## Slice 3: Save Document
_Adds save functionality to persist changes._

- [x] **Slice 3: Save Document**
  - [x] Add `saveDocument()` function to `editor.js`
  - [x] Validate JSON before saving (show error if invalid)
  - [x] POST to `/{channel}/{document}` on save button click
  - [x] Show success/error status message
  - [x] Disable save button during request
  - [x] Add Ctrl+S / Cmd+S keyboard shortcut
  - [ ] Manually verify: edit and save document, verify via API

---

## Slice 4: JSON Syntax Highlighting
_Adds visual highlighting for JSON syntax._

- [x] **Slice 4: Syntax Highlighting**
  - [x] Add `highlightJSON()` function with regex-based highlighting
  - [x] Add `<pre id="highlight">` overlay element
  - [x] Sync scroll between textarea and highlight overlay
  - [x] Update highlighting on input
  - [x] Add CSS classes for keys, strings, numbers, booleans, null
  - [ ] Manually verify: type JSON, see colored syntax
