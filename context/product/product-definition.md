# Product Definition: JustDoc

- **Version:** 1.0
- **Status:** Proposed

---

## 1. The Big Picture (The "Why")

### 1.1. Project Vision & Purpose

To provide developers with a simple, reliable way to store and retrieve JSON documents via a clean API, without the overhead of setting up and managing a database.

### 1.2. Target Audience

Frontend developers building apps that need persistent data storage for prototypes, proof-of-concepts, and small applications.

### 1.3. User Personas

- **Persona 1: "Frontend Fiona"**
  - **Role:** Frontend developer building a React/Vue/Angular application.
  - **Goal:** Wants to persist app data (user preferences, form submissions, app state) without setting up a backend or database.
  - **Frustration:** Setting up a database and backend just to store some JSON feels like overkill for a POC or small project.

### 1.4. Success Metrics

- Developers can integrate and make their first API call within 5 minutes.
- Zero-config setup — no database knowledge required.
- API response time under 200ms for typical operations.

---

## 2. The Product Experience (The "What")

### 2.1. Core Features

- **Store/Update Document** — `POST /<channel-name>/<document-name>` to create or update a JSON document.
- **Retrieve Document** — `GET /<channel-name>/<document-name>` to fetch a stored JSON document.
- **List Documents** — `GET /<channel-name>/` to list all documents within a channel.
- **List Channels** — `GET /` to list all available channels.

### 2.2. User Journey

1. Developer discovers JustDoc and reads the simple API documentation.
2. Developer makes a POST request to `/<channel>/<document>` with JSON body to store data.
3. Developer retrieves the data anytime with a GET request to the same endpoint.
4. Developer can organize documents into channels and list them as needed.
5. Total time from discovery to first working integration: under 5 minutes.

---

## 3. Project Boundaries

### 3.1. What's In-Scope for this Version

- Store JSON documents via POST request.
- Retrieve JSON documents via GET request.
- Organize documents into channels.
- List all documents in a channel.
- List all channels.
- Simple RESTful API structure.

### 3.2. What's Out-of-Scope (Non-Goals)

- Authentication and authorization.
- User accounts and API keys.
- Full-text search within documents.
- Document versioning or history.
- Real-time subscriptions/webhooks.
- Rate limiting.
- Admin UI or dashboard.
- DELETE endpoint (for v1).
- Document validation or schema enforcement.
