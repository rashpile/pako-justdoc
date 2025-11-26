# Functional Specification: Project Infrastructure

- **Roadmap Item:** Project Infrastructure (Phase 1)
- **Status:** Approved
- **Author:** Poe

---

## 1. Overview and Rationale (The "Why")

**Purpose:** Establish the foundational development environment and project infrastructure for JustDoc, enabling efficient feature development in subsequent phases.

**Problem:** Without proper project infrastructure, developers cannot build, test, or run the application consistently. A well-structured foundation ensures all team members can contribute effectively from day one.

**Desired Outcome:** A fully initialized Go project with version control, build tooling, and containerization that allows a developer to clone the repo and run `make build` within minutes.

**Success Criteria:**
- New developer can clone, build, and run the project in under 5 minutes
- All Makefile commands execute successfully
- Docker image builds and runs correctly

---

## 2. Functional Requirements (The "What")

### 2.1 Initialize Git Repository

- **As a** developer, **I want** a Git repository initialized with proper ignore patterns, **so that** I can track changes without committing build artifacts.
  - **Acceptance Criteria:**
    - [ ] Repository is initialized with `git init`
    - [ ] `.gitignore` includes: Go binaries, `*.exe`, `*.db` (bbolt data), `/bin/`, `/.idea/`, `/.vscode/`, `.DS_Store`

### 2.2 Project Structure

- **As a** developer, **I want** a clear directory layout, **so that** I know where to place code.
  - **Acceptance Criteria:**
    - [ ] Directory structure matches architecture spec:
      ```
      pako-justdoc/
      ├── cmd/justdoc/main.go
      ├── internal/api/
      ├── internal/storage/
      ├── internal/model/
      ├── Makefile
      ├── Dockerfile
      ├── go.mod
      └── README.md
      ```
    - [ ] `main.go` contains minimal placeholder (prints "JustDoc starting..." or similar)
    - [ ] `README.md` contains project name and brief description

### 2.3 Go Module

- **As a** developer, **I want** Go modules initialized, **so that** dependencies are managed properly.
  - **Acceptance Criteria:**
    - [ ] `go.mod` initialized with module path `github.com/rashpile/pako-justdoc`
    - [ ] Go version set to `1.25`
    - [ ] `go.sum` generated after adding dependencies

### 2.4 Makefile

- **As a** developer, **I want** common commands in a Makefile, **so that** I can build/test/run with simple commands.
  - **Acceptance Criteria:**
    - [ ] `make build` — Compiles binary to `bin/justdoc`
    - [ ] `make run` — Builds and runs the server (default port: `8080`, configurable via `PORT` env var)
    - [ ] `make test` — Runs all tests with coverage report
    - [ ] `make lint` — Runs golangci-lint
    - [ ] `make clean` — Removes build artifacts (`bin/`, `*.db`)
    - [ ] `make docker` — Builds Docker image tagged as `justdoc:latest`

### 2.5 Docker Image

- **As a** developer, **I want** a minimal Docker image, **so that** the service can be deployed in containers efficiently.
  - **Acceptance Criteria:**
    - [ ] Multi-stage Dockerfile (builder + runtime)
    - [ ] Final image uses `scratch` or `distroless` base (minimal size)
    - [ ] Binary is statically compiled (`CGO_ENABLED=0`)
    - [ ] Image exposes port `8080`
    - [ ] Image can be built with `make docker`
    - [ ] Container runs successfully with `docker run -p 8080:8080 justdoc:latest`

### 2.6 Linter Configuration

- **As a** developer, **I want** golangci-lint configured, **so that** code quality is enforced.
  - **Acceptance Criteria:**
    - [ ] `.golangci.yml` config file present
    - [ ] `make lint` runs without errors on initial codebase

---

## 3. Scope and Boundaries

### In-Scope

- Git repository initialization
- Project directory structure
- `go.mod` with module path and Go version
- Makefile with build, run, test, lint, clean, docker commands
- Minimal Dockerfile (multi-stage, scratch base)
- `.gitignore` for Go projects
- `.golangci.yml` linter configuration
- Placeholder `main.go`
- Basic `README.md`

### Out-of-Scope

*(Separate roadmap items to be addressed in Phase 2 and beyond)*

- Document Operations (Store/Update, Retrieve)
- Channel Operations (List Documents, List Channels)
- API Documentation
- Quick Start Guide
- Delete Document
- Authentication
- Rate Limiting