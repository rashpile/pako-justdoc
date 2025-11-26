# Tasks: Project Infrastructure

---

## Slice 1: Minimal Runnable Go Project
*Goal: `go build` works and binary runs*

- [x] **Slice 1: Create minimal buildable Go project**
  - [x] Create directory structure: `cmd/justdoc/`, `internal/api/`, `internal/storage/`, `internal/model/`
  - [x] Create `go.mod` with module `github.com/rashpile/pako-justdoc` and Go 1.25
  - [x] Create placeholder `cmd/justdoc/main.go` that prints startup message
  - [x] Add `.gitkeep` files to empty `internal/` subdirectories
  - [x] **Verify:** `go build ./cmd/justdoc` succeeds and binary runs

---

## Slice 2: Makefile with Build & Run
*Goal: `make build` and `make run` work*

- [x] **Slice 2: Add Makefile with core commands**
  - [x] Create `Makefile` with `build`, `run`, `clean` targets
  - [x] **Verify:** `make build` creates `bin/justdoc`
  - [x] **Verify:** `make run` starts the application
  - [x] **Verify:** `make clean` removes `bin/`

---

## Slice 3: Git Repository & Ignore Patterns
*Goal: Project is version-controlled*

- [x] **Slice 3: Initialize Git with proper .gitignore**
  - [x] Create `.gitignore` (binaries, `*.db`, IDE, OS files)
  - [x] Run `git init`
  - [x] **Verify:** `bin/` and `*.db` files are ignored

---

## Slice 4: Linter Configuration
*Goal: `make lint` passes*

- [x] **Slice 4: Configure golangci-lint**
  - [x] Create `.golangci.yml` with v2 config
  - [x] Add `lint` target to Makefile
  - [x] **Verify:** `make lint` runs without errors

---

## Slice 5: Test Infrastructure
*Goal: `make test` works*

- [x] **Slice 5: Add test command**
  - [x] Add `test` target to Makefile with coverage
  - [x] **Verify:** `make test` runs (passes with no tests)

---

## Slice 6: Docker Image
*Goal: `make docker` builds and container runs*

- [x] **Slice 6: Create minimal Docker image**
  - [x] Create multi-stage `Dockerfile` (builder + scratch)
  - [x] Add `docker` target to Makefile
  - [x] **Verify:** `make docker` builds `justdoc:latest`
  - [x] **Verify:** `docker run -p 8080:8080 justdoc:latest` starts

---

## Slice 7: Documentation & Dependencies
*Goal: Project is ready for Phase 2*

- [x] **Slice 7: Finalize project setup**
  - [x] Create `README.md` with quick start and commands
  - [x] Add bbolt dependency: `go get go.etcd.io/bbolt@v1.3.7`
  - [x] **Verify:** `go.sum` is generated
  - [x] **Verify:** All `make` commands work