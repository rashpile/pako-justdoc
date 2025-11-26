# Tasks: Document Operations

---

## Slice 1: Model Layer (Foundation)
*Goal: Types and validation ready for use*

- [x] **Slice 1: Create model layer with types and validation**
  - [x] Create `internal/model/response.go` with `SuccessResponse` and `ErrorResponse` types
  - [x] Create `internal/model/errors.go` with error code constants
  - [x] Create `internal/model/validation.go` with `IsValidName()` function
  - [x] Create `internal/model/validation_test.go` with tests for valid/invalid names
  - [x] **Verify:** `make test` passes, validation works correctly

---

## Slice 2: Storage Layer (bbolt)
*Goal: Can store and retrieve documents programmatically*

- [x] **Slice 2: Implement bbolt storage layer**
  - [x] Create `internal/storage/storage.go` with `Storage` interface and `ErrNotFound`
  - [x] Create `internal/storage/bbolt.go` with `BoltStorage` implementation
  - [x] Create `internal/storage/bbolt_test.go` with unit tests
  - [x] **Verify:** `make test` passes, storage operations work

---

## Slice 3: POST Endpoint (Store Document)
*Goal: Can store documents via HTTP POST*

- [x] **Slice 3: Implement POST /{channel}/{document} endpoint**
  - [x] Create `internal/api/handler.go` with `Handler` struct and `PostDocument` method
  - [x] Create `internal/api/router.go` with route setup
  - [x] Update `cmd/justdoc/main.go` to wire storage and API
  - [x] **Verify:** `make run` starts server, `curl -X POST` stores document with 201 response

---

## Slice 4: GET Endpoint (Retrieve Document)
*Goal: Can retrieve stored documents via HTTP GET*

- [x] **Slice 4: Implement GET /{channel}/{document} endpoint**
  - [x] Add `GetDocument` method to `internal/api/handler.go`
  - [x] Add GET route to `internal/api/router.go`
  - [x] **Verify:** `curl GET` returns stored document with 200 response

---

## Slice 5: Error Handling
*Goal: All error cases return proper responses*

- [x] **Slice 5: Implement error handling for all cases**
  - [x] Invalid JSON → 400 Bad Request
  - [x] Invalid name → 400 Bad Request
  - [x] Document not found → 404 Not Found
  - [x] Payload too large → 413 Payload Too Large
  - [x] **Verify:** All error cases return correct status codes and JSON error format

---

## Slice 6: API Tests
*Goal: HTTP handlers fully tested*

- [x] **Slice 6: Add API handler tests**
  - [x] Create `internal/api/handler_test.go`
  - [x] Test POST create (201), POST update (200)
  - [x] Test GET success (200), GET not found (404)
  - [x] Test all error cases (400, 413)
  - [x] **Verify:** `make test` passes with 80%+ coverage on handlers

---

## Slice 7: Integration & Lint
*Goal: Full end-to-end verification*

- [x] **Slice 7: Integration testing and final verification**
  - [x] Run `make lint` and fix any issues
  - [x] Run `make test` with coverage
  - [x] Test full flow: POST → GET → POST (update) → GET
  - [x] Test with `make docker` and verify container works
  - [x] **Verify:** All make commands pass, Docker container serves API
