# Functional Specification: Channel Operations

- **Roadmap Item:** Channel Operations — List Documents in Channel, List All Channels
- **Status:** Approved
- **Author:** Claude

---

## 1. Overview and Rationale (The "Why")

### Context
JustDoc currently allows developers to store and retrieve JSON documents via `POST` and `GET` requests to `/<channel>/<document>`. However, developers cannot discover what data they've stored without knowing exact document names.

### Problem
A developer building a prototype may store multiple documents across several channels over time. Without a way to list their stored data, they must remember or track document names externally, which defeats the "zero-config" simplicity goal.

### Desired Outcome
Developers can browse their stored data by:
1. Listing all channels to see how their data is organized
2. Listing all documents within a specific channel to find what they need

### Success Metrics
- API response time under 200ms for listing operations
- Response format is intuitive and requires no documentation to understand

---

## 2. Functional Requirements (The "What")

### 2.1 List Documents in Channel

**Endpoint:** `GET /<channel>/`

**As a** developer, **I want to** list all documents in a channel, **so that** I can discover what data I've stored without remembering exact document names.

**Acceptance Criteria:**
- [ ] When I send `GET /my-channel/`, I receive a JSON array of document names stored in that channel
- [ ] Response format: `["document1", "document2", "document3"]`
- [ ] Response Content-Type is `application/json`
- [ ] When the channel does not exist (no documents have ever been stored in it), I receive a `404 Not Found` response
- [ ] All documents in the channel are returned (no pagination)
- [ ] Document names are returned in alphabetical order

**Example Response (Success):**
```json
["config", "user-preferences", "session-data"]
```

**Example Response (Channel Not Found):**
```
HTTP 404 Not Found
```

---

### 2.2 List All Channels

**Endpoint:** `GET /`

**As a** developer, **I want to** list all channels, **so that** I can see an overview of how my data is organized.

**Acceptance Criteria:**
- [ ] When I send `GET /`, I receive a JSON array of channel objects
- [ ] Each channel object includes: `name` (string) and `document_count` (number)
- [ ] Response Content-Type is `application/json`
- [ ] When no channels exist, I receive an empty array `[]`
- [ ] Channels are returned in alphabetical order by name
- [ ] The existing `GET /openapi.json` endpoint continues to work and is not affected

**Example Response (With Channels):**
```json
[
  {"name": "app-config", "document_count": 3},
  {"name": "user-data", "document_count": 12}
]
```

**Example Response (No Channels):**
```json
[]
```

---

## 3. Scope and Boundaries

### In-Scope
- `GET /<channel>/` endpoint returning array of document names
- `GET /` endpoint returning array of channel objects with document counts
- 404 response for non-existent channels
- Alphabetical sorting for both endpoints

### Out-of-Scope
- Pagination or filtering
- Document metadata (size, created date, modified date)
- Channel metadata beyond document count
- Search functionality
- DELETE operations (Phase 4)
- Authentication (Phase 4)
- Document Editor UI (Phase 3)