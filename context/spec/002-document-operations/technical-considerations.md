# Technical Specification: Document Operations

- **Functional Specification:** `context/spec/002-document-operations/functional-spec.md`
- **Status:** Approved
- **Author(s):** Poe

---

## 1. High-Level Technical Approach

Implement the Document Operations API using a layered architecture:

1. **Storage Layer (`internal/storage/`)**: bbolt wrapper with channels as buckets, documents as key-value pairs
2. **API Layer (`internal/api/`)**: HTTP handlers using standard `net/http` with Go 1.22+ routing
3. **Model Layer (`internal/model/`)**: Request/response types, validation, and error definitions
4. **Main (`cmd/justdoc/`)**: Wire components together, start HTTP server

**Key Principles:**
- Idiomatic Go patterns
- Interface-based design for testability
- Standard library HTTP (no frameworks)
- Clean separation of concerns

---

## 2. Proposed Solution & Implementation Plan (The "How")

### 2.1 Project Structure (Updated)

```
pako-justdoc/
├── cmd/justdoc/
│   └── main.go              # Entry point, wiring, server startup
├── internal/
│   ├── api/
│   │   ├── handler.go       # HTTP handlers (GetDocument, PostDocument)
│   │   ├── router.go        # Route setup
│   │   └── middleware.go    # Request size limiting, logging
│   ├── storage/
│   │   ├── storage.go       # Storage interface
│   │   └── bbolt.go         # bbolt implementation
│   └── model/
│       ├── response.go      # API response types
│       ├── errors.go        # Error types and codes
│       └── validation.go    # Name validation
├── justdoc.db               # Database file (gitignored)
└── ...
```

### 2.2 Data Model (bbolt)

```
bbolt database (justdoc.db)
└── Bucket: "<channel-name>"
    ├── Key: "<document-name>" → Value: <raw JSON bytes>
    ├── Key: "<document-name>" → Value: <raw JSON bytes>
    └── ...
```

- Each channel is a bbolt bucket
- Documents are key-value pairs within the bucket
- Values stored as raw JSON bytes (no transformation)

### 2.3 Storage Layer

**internal/storage/storage.go:**
```go
package storage

// Storage defines the document storage interface
type Storage interface {
    // GetDocument retrieves a document from a channel
    // Returns ErrNotFound if channel or document doesn't exist
    GetDocument(channel, document string) ([]byte, error)

    // PutDocument stores a document in a channel
    // Creates the channel if it doesn't exist
    // Returns created=true if document was new, false if updated
    PutDocument(channel, document string, data []byte) (created bool, err error)

    // Close closes the storage connection
    Close() error
}

// Sentinel errors
var (
    ErrNotFound = errors.New("not found")
)
```

**internal/storage/bbolt.go:**
```go
package storage

import (
    "go.etcd.io/bbolt"
)

type BoltStorage struct {
    db *bbolt.DB
}

func NewBoltStorage(path string) (*BoltStorage, error) {
    db, err := bbolt.Open(path, 0600, nil)
    if err != nil {
        return nil, err
    }
    return &BoltStorage{db: db}, nil
}

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
        data = make([]byte, len(v))
        copy(data, v)
        return nil
    })
    return data, err
}

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

func (s *BoltStorage) Close() error {
    return s.db.Close()
}
```

### 2.4 Model Layer

**internal/model/response.go:**
```go
package model

// SuccessResponse for POST operations
type SuccessResponse struct {
    Status   string `json:"status"`   // "created" or "updated"
    Channel  string `json:"channel"`
    Document string `json:"document"`
}

// ErrorResponse for all error cases
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
```

**internal/model/errors.go:**
```go
package model

const (
    ErrCodeInvalidJSON     = "invalid_json"
    ErrCodeInvalidName     = "invalid_name"
    ErrCodeNotFound        = "not_found"
    ErrCodePayloadTooLarge = "payload_too_large"
)
```

**internal/model/validation.go:**
```go
package model

import "regexp"

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

func IsValidName(name string) bool {
    return validName.MatchString(name)
}
```

### 2.5 API Layer

**internal/api/handler.go:**
```go
package api

import (
    "encoding/json"
    "io"
    "net/http"

    "github.com/rashpile/pako-justdoc/internal/model"
    "github.com/rashpile/pako-justdoc/internal/storage"
)

const MaxBodySize = 10 * 1024 * 1024 // 10MB

type Handler struct {
    storage storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
    return &Handler{storage: s}
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
    channel := r.PathValue("channel")
    document := r.PathValue("document")

    // Validate names
    if !model.IsValidName(channel) || !model.IsValidName(document) {
        writeError(w, http.StatusBadRequest, model.ErrCodeInvalidName, "Invalid channel or document name")
        return
    }

    data, err := h.storage.GetDocument(channel, document)
    if err == storage.ErrNotFound {
        writeError(w, http.StatusNotFound, model.ErrCodeNotFound, "Document not found")
        return
    }
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func (h *Handler) PostDocument(w http.ResponseWriter, r *http.Request) {
    channel := r.PathValue("channel")
    document := r.PathValue("document")

    // Validate names
    if !model.IsValidName(channel) || !model.IsValidName(document) {
        writeError(w, http.StatusBadRequest, model.ErrCodeInvalidName, "Invalid channel or document name")
        return
    }

    // Limit body size
    r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

    // Read body
    data, err := io.ReadAll(r.Body)
    if err != nil {
        if err.Error() == "http: request body too large" {
            writeError(w, http.StatusRequestEntityTooLarge, model.ErrCodePayloadTooLarge, "Request body exceeds 10MB limit")
            return
        }
        writeError(w, http.StatusBadRequest, model.ErrCodeInvalidJSON, "Failed to read request body")
        return
    }

    // Validate JSON
    if !json.Valid(data) {
        writeError(w, http.StatusBadRequest, model.ErrCodeInvalidJSON, "Invalid JSON body")
        return
    }

    // Store document
    created, err := h.storage.PutDocument(channel, document, data)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal_error", "Failed to store document")
        return
    }

    // Build response
    status := "updated"
    statusCode := http.StatusOK
    if created {
        status = "created"
        statusCode = http.StatusCreated
    }

    writeJSON(w, statusCode, model.SuccessResponse{
        Status:   status,
        Channel:  channel,
        Document: document,
    })
}

func writeError(w http.ResponseWriter, statusCode int, errCode, message string) {
    writeJSON(w, statusCode, model.ErrorResponse{
        Error:   errCode,
        Message: message,
    })
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(v)
}
```

**internal/api/router.go:**
```go
package api

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{channel}/{document}", h.GetDocument)
    mux.HandleFunc("POST /{channel}/{document}", h.PostDocument)
    return mux
}
```

### 2.6 Main Entry Point

**cmd/justdoc/main.go:**
```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/rashpile/pako-justdoc/internal/api"
    "github.com/rashpile/pako-justdoc/internal/storage"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "justdoc.db"
    }

    // Initialize storage
    store, err := storage.NewBoltStorage(dbPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer store.Close()

    // Initialize API
    handler := api.NewHandler(store)
    router := api.NewRouter(handler)

    // Graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh
        fmt.Println("\nShutting down...")
        store.Close()
        os.Exit(0)
    }()

    // Start server
    fmt.Printf("JustDoc starting on port %s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}
```

### 2.7 Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DB_PATH` | `justdoc.db` | Path to bbolt database file |

---

## 3. Impact and Risk Analysis

### System Dependencies

- **bbolt v1.3.7**: Already added as dependency in Phase 1
- **Go 1.22+**: Required for enhanced ServeMux routing patterns

### Potential Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Concurrent write conflicts | bbolt handles locking internally; single-writer, multiple-reader |
| Large documents (up to 10MB) | MaxBytesReader prevents memory exhaustion |
| Invalid JSON storage | Validate JSON before storing |
| Database corruption | bbolt is ACID-compliant, uses copy-on-write |

---

## 4. Testing Strategy

### Unit Tests

1. **Storage Layer (`internal/storage/bbolt_test.go`)**
   - Test GetDocument with existing/non-existing documents
   - Test PutDocument create vs update
   - Test channel auto-creation

2. **Validation (`internal/model/validation_test.go`)**
   - Test valid names (alphanumeric, hyphens, underscores)
   - Test invalid names (spaces, special chars, too long)

3. **Handlers (`internal/api/handler_test.go`)**
   - Test POST with valid JSON → 201/200
   - Test POST with invalid JSON → 400
   - Test POST with oversized body → 413
   - Test GET existing document → 200
   - Test GET non-existing document → 404
   - Test invalid names → 400

### Integration Tests

- End-to-end tests using `httptest.Server`
- Store document, retrieve it, verify content
- Update document, verify status changes from created to updated

### Test Coverage Target

- 80%+ code coverage
- All error paths tested