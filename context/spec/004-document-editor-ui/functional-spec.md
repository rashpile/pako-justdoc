# Functional Specification: Document Editor UI

- **Roadmap Item:** Document Editor UI — Simple web UI for document management
- **Status:** Approved
- **Author:** Claude

---

## 1. Overview and Rationale (The "Why")

### Context
JustDoc provides a REST API for storing and retrieving JSON documents. Currently, developers must use tools like curl, Postman, or write code to interact with their data.

### Problem
When prototyping or debugging, developers often want to quickly view or edit a JSON document without writing code or switching to API tools. This slows down the development workflow.

### Desired Outcome
Developers can view and edit any JSON document directly in the browser by navigating to a simple URL. The UI is lightweight, fast-loading, and requires no setup.

### Success Metrics
- Page loads in under 500ms
- Zero external CDN dependencies (all assets served locally)
- Works in all modern browsers without JavaScript frameworks

---

## 2. Functional Requirements (The "What")

### 2.1 Access the Editor

**Endpoint:** `GET /_/edit/<channel>/<document>`

**As a** developer, **I want to** open a document in my browser, **so that** I can view and edit it without using API tools.

**Acceptance Criteria:**
- [x] When I navigate to `/_/edit/<channel>/<document>`, I see a web page with a JSON editor
- [x] If the document exists, the editor is pre-filled with the document's JSON content
- [x] If the document does not exist, the editor shows an empty text area (ready for creating a new document)
- [x] The page title shows the channel and document name (e.g., "myapp / settings - JustDoc")
- [x] The URL path (channel/document) is displayed on the page for clarity

---

### 2.2 Edit and Save Document

**As a** developer, **I want to** edit JSON and save it, **so that** I can update my document without using the API directly.

**Acceptance Criteria:**
- [x] The editor is a text area where I can type/edit JSON
- [x] Basic JSON syntax highlighting is applied (keywords, strings, numbers in different colors)
- [x] A "Save" button is visible
- [x] When I click "Save" with valid JSON:
  - The document is saved via `POST /<channel>/<document>`
  - A success message appears (e.g., "Saved successfully")
  - I remain on the same page
- [x] When I click "Save" with invalid JSON:
  - The save is prevented
  - An error message appears (e.g., "Invalid JSON: [error details]")
  - The error location is indicated if possible

---

### 2.3 Visual Design

**As a** developer, **I want** the UI to be clean and fast, **so that** it doesn't distract from my work.

**Acceptance Criteria:**
- [x] Light color scheme (white/light gray background)
- [x] Minimalist design with no unnecessary elements
- [x] Monospace font for the editor
- [x] Page contains only: title/path display, editor area, save button, status message area
- [x] No navigation menus, sidebars, or footer
- [x] Total page size under 50KB (HTML + CSS + JS combined)

---

### 2.4 JSON Syntax Highlighting

**Acceptance Criteria:**
- [x] JSON keys are highlighted in one color
- [x] Strings are highlighted in another color
- [x] Numbers and booleans are highlighted in another color
- [x] Uses a lightweight approach (no heavy libraries like CodeMirror or Monaco)
- [x] Highlighting updates as the user types (or on save attempt)

---

## 3. Scope and Boundaries

### In-Scope
- Single-document editor at `/_/edit/<channel>/<document>`
- View existing document content
- Create new document (by navigating to non-existent document URL)
- Edit and save document
- Basic JSON syntax highlighting
- Inline success/error messages
- Vanilla HTML/CSS/JS implementation

### Out-of-Scope
- Channel/document browsing UI (Phase 4 or later)
- Document deletion from UI
- Multiple document tabs
- Dark mode
- Undo/redo beyond browser native
- Keyboard shortcuts (except Ctrl+S/Cmd+S for save)
- Mobile-optimized layout
- Authentication (Phase 4)
- Real-time collaboration
