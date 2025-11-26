# Technical Specification: Project Infrastructure

- **Functional Specification:** `context/spec/001-project-infrastructure/functional-spec.md`
- **Status:** Approved
- **Author(s):** Poe

---

## 1. High-Level Technical Approach

This is a greenfield Go project setup establishing the foundation for JustDoc. The implementation involves:

1. Creating the idiomatic Go directory structure with placeholder files
2. Initializing Go modules with bbolt as the only external dependency
3. Writing build tooling (Makefile) with idiomatic Go build flags
4. Creating a minimal multi-stage Dockerfile using `scratch` base
5. Configuring golangci-lint v2 for idiomatic Go enforcement

**Key Principles:**
- Idiomatic Go patterns throughout
- Pure Go with minimal dependencies
- Single statically-linked binary
- Zero-config local development

---

## 2. Proposed Solution & Implementation Plan (The "How")

### 2.1 Project Structure

```
pako-justdoc/
├── cmd/justdoc/
│   └── main.go           # Application entry point
├── internal/
│   ├── api/              # HTTP handlers (placeholder)
│   │   └── .gitkeep
│   ├── storage/          # bbolt storage layer (placeholder)
│   │   └── .gitkeep
│   └── model/            # Data types (placeholder)
│       └── .gitkeep
├── bin/                  # Build output (gitignored)
├── .gitignore
├── .golangci.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

### 2.2 Go Module Configuration

**go.mod:**
```go
module github.com/rashpile/pako-justdoc

go 1.25

require go.etcd.io/bbolt v1.3.7
```

- Module path: `github.com/rashpile/pako-justdoc`
- Go version: `1.25`
- Dependencies: `go.etcd.io/bbolt v1.3.7` (added for Phase 2 readiness)

### 2.3 Makefile

```makefile
BINARY_NAME=justdoc
BUILD_DIR=bin
PORT?=8080

.PHONY: build run test lint clean docker

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/justdoc

run: build
	PORT=$(PORT) ./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR) *.db

docker:
	docker build -t justdoc:latest .
```

**Build flags:**
- `-race` in tests for race condition detection
- Default port `8080`, configurable via `PORT` env var

### 2.4 Dockerfile (Multi-stage, Scratch)

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o justdoc ./cmd/justdoc

# Runtime stage
FROM scratch
COPY --from=builder /app/justdoc /justdoc
EXPOSE 8080
ENTRYPOINT ["/justdoc"]
```

**Optimizations:**
- `CGO_ENABLED=0` for static linking
- `-ldflags="-s -w"` strips debug symbols (smaller binary)
- `scratch` base for minimal image size (~5-10MB expected)

### 2.5 golangci-lint Configuration

**.golangci.yml:**
```yaml
version: "2"
run:
  go: "1.25"
linters:
  default: standard
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
```

### 2.6 .gitignore

```gitignore
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Database
*.db

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Test
*.test
*.out
coverage.html
```

### 2.7 Placeholder main.go

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("JustDoc starting on port %s...\n", port)
}
```

### 2.8 README.md

```markdown
# JustDoc

A simple JSON document storage service.

## Quick Start

```bash
make build
make run
```

## Commands

- `make build` - Compile binary
- `make run` - Build and run server
- `make test` - Run tests with coverage
- `make lint` - Run linter
- `make clean` - Remove build artifacts
- `make docker` - Build Docker image
```

---

## 3. Impact and Risk Analysis

### System Dependencies

- **Go 1.25** must be installed for local development
- **golangci-lint v2.6.2** must be installed for `make lint`
- **Docker** required for `make docker`

### Potential Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Go 1.25 not yet widely available | Document installation steps in README |
| golangci-lint not installed | Makefile can check and warn, or provide install instructions |
| Docker build fails on non-Linux | Multi-platform build flags already included (`GOOS=linux`) |

---

## 4. Testing Strategy

### Phase 1 Testing

Since Phase 1 is infrastructure setup, testing is limited to:

1. **Build Verification:**
   - `make build` produces a binary in `bin/justdoc`
   - Binary executes and prints startup message

2. **Docker Verification:**
   - `make docker` builds successfully
   - `docker run -p 8080:8080 justdoc:latest` starts without error

3. **Lint Verification:**
   - `make lint` passes with zero errors

### Future Testing (Phase 2+)

- Unit tests in `*_test.go` files alongside source
- Table-driven tests (idiomatic Go)
- Integration tests for API endpoints
- Coverage target: 80%+