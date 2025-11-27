# Tasks: Channel Operations

## Slice 1: List Documents in Channel (Full Vertical Slice)
_Delivers the complete `GET /<channel>/` endpoint._

- [x] **Slice 1: List Documents in Channel**
  - [x] Add `ListDocuments(channel string) ([]string, error)` method to `Storage` interface in `internal/storage/storage.go`
  - [x] Implement `ListDocuments` in `BoltStorage` (`internal/storage/bbolt.go`) with alphabetical sorting
  - [x] Add unit tests for `ListDocuments` in `internal/storage/bbolt_test.go`
  - [x] Add `ListDocuments` handler in `internal/api/handler.go` (validation, 404 handling, JSON response)
  - [x] Register `GET /{channel}/` route in `internal/api/router.go`
  - [x] Add handler tests for `ListDocuments` in `internal/api/handler_test.go`
  - [x] Manually verify: create documents, call `GET /test-channel/`, confirm JSON array response

---

## Slice 2: List All Channels (Full Vertical Slice)
_Delivers the complete `GET /` endpoint._

- [x] **Slice 2: List All Channels**
  - [x] Add `ChannelInfo` struct to `internal/storage/storage.go`
  - [x] Add `ListChannels() ([]ChannelInfo, error)` method to `Storage` interface
  - [x] Implement `ListChannels` in `BoltStorage` with alphabetical sorting and document count
  - [x] Add unit tests for `ListChannels` in `internal/storage/bbolt_test.go`
  - [x] Add `ListChannels` handler in `internal/api/handler.go`
  - [x] Register `GET /` route in `internal/api/router.go`
  - [x] Add handler tests for `ListChannels` in `internal/api/handler_test.go`
  - [x] Manually verify: `GET /` returns channel list, `GET /openapi.json` still works

---

## Slice 3: Update OpenAPI Documentation
_Documents the new endpoints in the API spec._

- [x] **Slice 3: Update OpenAPI Specification**
  - [x] Add `GET /` path definition to `internal/api/openapi.go`
  - [x] Add `GET /{channel}/` path definition to `internal/api/openapi.go`
  - [x] Verify `GET /openapi.json` returns updated spec with new endpoints