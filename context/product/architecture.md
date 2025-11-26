# System Architecture Overview: JustDoc

---

## 1. Application & Technology Stack

- **Language:** Go 1.25 (latest stable)
- **HTTP Router:** Standard library `net/http` (pure Go, no external dependencies)
- **Linter:** golangci-lint v2.6.2
- **Build Tool:** Standard `go build` + Makefile

---

## 2. Data & Persistence

- **Primary Database:** bbolt v1.3.7 (`go.etcd.io/bbolt`)
  - Pure Go, no CGO dependencies
  - Single-file embedded database
  - Key-value store with buckets (channels → documents)
  - Battle-tested (etcd, Consul)

---

## 3. Infrastructure & Deployment

- **Deployment:** Single statically-linked binary
- **Configuration:** Environment variables
- **Containerization:** Dockerfile (optional, for convenience)
- **External Dependencies:** None (self-contained)

---

## 4. Project Structure

```
justdoc/
├── cmd/justdoc/          # Application entry point
│   └── main.go
├── internal/
│   ├── api/              # HTTP handlers and routing
│   ├── storage/          # bbolt storage layer
│   └── model/            # Data types and structures
├── Makefile              # Build, run, test, clean commands
├── go.mod
├── go.sum
└── README.md
```

---

## 5. Design Principles

- **Simplicity:** Minimal dependencies, straightforward code
- **Idiomatic Go:** Follow Go conventions and best practices
- **Pure Go:** Standard library preferred, minimal external dependencies
- **Single Binary:** Easy deployment, no runtime dependencies
- **Fast Startup:** Sub-second initialization for development workflow