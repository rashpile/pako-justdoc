# Technical Specification: Channel Operations

- **Functional Specification:** `context/spec/003-channel-operations/functional-spec.md`
- **Status:** Approved
- **Author(s):** Claude

---

## 1. High-Level Technical Approach

This feature adds two new read-only endpoints to list channels and documents. The implementation follows the existing layered architecture:

1. **Storage Layer**: Add `ListDocuments()` and `ListChannels()` methods to the `Storage` interface and implement in `BoltStorage`
2. **Handler Layer**: Add `ListDocuments` and `ListChannels` handlers following existing patterns
3. **Router Layer**: Register new routes for `GET /` and `GET /{channel}/`
4. **OpenAPI**: Update the OpenAPI spec to document new endpoints

No database schema changes required — bbolt's bucket structure already supports iteration.

---

## 2. Proposed Solution & Implementation Plan

### 2.1 Storage Interface Changes

**File:** `internal/storage/storage.go`

Add new method signatures to the `Storage` interface:

```go
// ListDocuments returns all document names in a channel (sorted alphabetically)
// Returns ErrNotFound if channel doesn't exist
ListDocuments(channel string) ([]string, error)

// ListChannels returns all channels with document counts (sorted alphabetically)
ListChannels() ([]ChannelInfo, error)
```

**File:** `internal/storage/storage.go` (new struct)

```go
// ChannelInfo represents a channel with its document count
type ChannelInfo struct {
    Name          string `json:"name"`
    DocumentCount int    `json:"document_count"`
}
```

### 2.2 bbolt Implementation

**File:** `internal/storage/bbolt.go`

**ListDocuments implementation:**
```go
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
```

**ListChannels implementation:**
```go
func (s *BoltStorage) ListChannels() ([]ChannelInfo, error) {
    var channels []ChannelInfo
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
```

### 2.3 Handler Implementation

**File:** `internal/api/handler.go`

**ListDocuments handler:**
```go
// ListDocuments handles GET /{channel}/
func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
    channel := r.PathValue("channel")

    if !model.IsValidName(channel) {
        writeError(w, http.StatusBadRequest, model.ErrCodeInvalidName, "Invalid channel name")
        return
    }

    docs, err := h.storage.ListDocuments(channel)
    if err == storage.ErrNotFound {
        writeError(w, http.StatusNotFound, model.ErrCodeNotFound, "Channel not found")
        return
    }
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(docs)
}
```

**ListChannels handler:**
```go
// ListChannels handles GET /
func (h *Handler) ListChannels(w http.ResponseWriter, r *http.Request) {
    channels, err := h.storage.ListChannels()
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(channels)
}
```

### 2.4 Router Changes

**File:** `internal/api/router.go`

```go
func NewRouter(h *Handler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /openapi.json", OpenAPI)
    mux.HandleFunc("GET /{channel}/{document}", h.GetDocument)
    mux.HandleFunc("POST /{channel}/{document}", h.PostDocument)
    mux.HandleFunc("GET /{channel}/", h.ListDocuments)  // NEW
    mux.HandleFunc("GET /", h.ListChannels)              // NEW
    return mux
}
```

### 2.5 OpenAPI Updates

**File:** `internal/api/openapi.go`

Add two new path definitions:
- `GET /` — List all channels
- `GET /{channel}/` — List documents in channel

---

## 3. Impact and Risk Analysis

### System Dependencies
- **Storage layer**: New methods added to interface; all implementations must be updated
- **Existing endpoints**: No changes to existing `GET/POST /{channel}/{document}` behavior

### Potential Risks & Mitigations

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Performance with many documents/channels | Low | bbolt iteration is efficient; no pagination needed for MVP |
| Route conflict with `/openapi.json` | Low | `/openapi.json` registered first, takes precedence |
| Empty array vs null for empty results | Low | Always return `[]`, never `null` |

---

## 4. Testing Strategy

### Unit Tests

**Storage tests** (`internal/storage/bbolt_test.go`):
- `TestListDocuments_ReturnsAlphabeticalOrder`
- `TestListDocuments_ChannelNotFound_ReturnsError`
- `TestListDocuments_EmptyChannel_ReturnsEmptyArray`
- `TestListChannels_ReturnsAlphabeticalOrder`
- `TestListChannels_NoChannels_ReturnsEmptyArray`
- `TestListChannels_IncludesDocumentCount`

**Handler tests** (`internal/api/handler_test.go`):
- `TestListDocuments_Success`
- `TestListDocuments_ChannelNotFound_Returns404`
- `TestListDocuments_InvalidChannelName_Returns400`
- `TestListChannels_Success`
- `TestListChannels_Empty_ReturnsEmptyArray`

### Integration Tests
- Verify `/openapi.json` still works after adding `GET /` route
- End-to-end flow: create documents, list them, verify response