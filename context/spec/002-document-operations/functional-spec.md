# Functional Specification: Document Operations

- **Roadmap Item:** Document Operations (Phase 2)
- **Status:** Approved
- **Author:** Poe

---

## 1. Overview and Rationale (The "Why")

**Purpose:** Enable developers to store and retrieve JSON documents via simple HTTP API endpoints, providing a zero-config persistence layer for frontend applications.

**Problem:** Frontend developers building prototypes or small applications need a way to persist JSON data without the overhead of setting up a database and backend infrastructure.

**Desired Outcome:** A simple REST API where developers can store any valid JSON document with a single POST request and retrieve it with a GET request.

**Success Criteria:**
- Developers can store and retrieve their first document within 5 minutes
- API response time under 200ms for typical operations
- Documents up to 10MB supported

---

## 2. Functional Requirements (The "What")

### 2.1 Store/Update Document

- **As a** developer, **I want to** store a JSON document via POST request, **so that** I can persist my application data.

**Endpoint:** `POST /<channel>/<document>`

**Request:**
- Body must be valid JSON
- Maximum body size: 10MB
- `Content-Type: application/json`

**Behavior:**
- If the channel doesn't exist, create it automatically
- If the document doesn't exist, create it (return `201 Created`)
- If the document exists, update it (return `200 OK`)

**Response (Success):**
```json
{
  "status": "created",
  "channel": "<channel>",
  "document": "<document>"
}
```
or
```json
{
  "status": "updated",
  "channel": "<channel>",
  "document": "<document>"
}
```

**Acceptance Criteria:**
- [ ] `POST /myapp/settings` with valid JSON body stores the document
- [ ] Returns `201 Created` when document is new
- [ ] Returns `200 OK` when document is updated
- [ ] Channel is created automatically if it doesn't exist
- [ ] Invalid JSON body returns `400 Bad Request`
- [ ] Body exceeding 10MB returns `413 Payload Too Large`
- [ ] Invalid channel/document name returns `400 Bad Request`

---

### 2.2 Retrieve Document

- **As a** developer, **I want to** retrieve a stored JSON document via GET request, **so that** I can read my persisted data.

**Endpoint:** `GET /<channel>/<document>`

**Response (Success):**
- Status: `200 OK`
- Body: The raw JSON document
- Header: `Content-Type: application/json`

**Acceptance Criteria:**
- [ ] `GET /myapp/settings` returns the stored JSON document
- [ ] Returns `200 OK` with `Content-Type: application/json`
- [ ] Non-existent document returns `404 Not Found`
- [ ] Non-existent channel returns `404 Not Found`

---

### 2.3 Naming Conventions

**Channel and Document Names:**
- Allowed characters: alphanumeric (`a-z`, `A-Z`, `0-9`), hyphens (`-`), underscores (`_`)
- Maximum length: 128 characters
- Case-sensitive

**Acceptance Criteria:**
- [ ] `my-channel`, `my_channel`, `MyChannel123` are valid names
- [ ] Names with spaces, special characters return `400 Bad Request`
- [ ] Names exceeding 128 characters return `400 Bad Request`

---

### 2.4 Error Responses

All error responses follow a consistent format:

```json
{
  "error": "<error_code>",
  "message": "<human_readable_message>"
}
```

| Scenario | Status Code | Error Code |
|----------|-------------|------------|
| Invalid JSON body | `400` | `invalid_json` |
| Invalid name format | `400` | `invalid_name` |
| Document not found | `404` | `not_found` |
| Channel not found | `404` | `not_found` |
| Payload too large | `413` | `payload_too_large` |

---

## 3. Scope and Boundaries

### In-Scope

- `POST /<channel>/<document>` to store/update JSON
- `GET /<channel>/<document>` to retrieve JSON
- Automatic channel creation
- JSON validation
- 10MB document size limit
- Name validation (alphanumeric, hyphens, underscores, max 128 chars)
- Consistent error response format

### Out-of-Scope

*(Separate roadmap items)*

- Channel Operations (List Documents, List Channels)
- Developer Experience (API Documentation, Quick Start Guide)
- Delete Document
- Authentication
- Rate Limiting