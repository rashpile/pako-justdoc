# Product Roadmap: JustDoc

_This roadmap outlines our strategic direction based on customer needs and business goals. It focuses on the "what" and "why," not the technical "how."_

---

### Phase 1 ✓

_Project foundation — completed._

- [x] **Project Infrastructure**
  - [x] Initialize Git Repository, Project Structure, Build Scripts, Makefile

---

### Phase 2 ✓

_Core functionality — completed._

- [x] **Document Operations**
  - [x] Store/Update Document: `POST /<channel>/<document>`
  - [x] Retrieve Document: `GET /<channel>/<document>`
  - [x] OpenAPI Documentation: `GET /openapi.json`

- [x] **Channel Operations**
  - [x] List Documents in Channel: `GET /<channel>/`
  - [x] List All Channels: `GET /`

---

### Phase 3 ✓

_User interface — completed._

- [x] **Document Editor UI**
  - [x] Edit documents: `GET /_/edit/<channel>/<document>`
  - [x] JSON editor with syntax highlighting
  - [x] Create, edit, and save documents via UI with minimalist design

---

### Phase 4

_Future considerations — features to evaluate based on user feedback._

- [ ] **Extended Functionality**
  - [ ] Delete Document: `DELETE /<channel>/<document>`
  - [ ] Authentication: API key-based secure access
  - [ ] Rate Limiting: Request throttling
