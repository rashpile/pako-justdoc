# JustDoc

> Zero-config JSON document storage for frontend developers

Store and retrieve JSON documents with simple HTTP requests. No database setup, no configuration files, no hassle.

## Features

- **Simple API** - Just `GET` and `POST` to store/retrieve JSON
- **Zero Config** - Works out of the box, no setup required
- **Channels** - Organize documents into logical groups
- **10MB Documents** - Store large JSON payloads
- **OpenAPI Spec** - Built-in API documentation at `/openapi.json`
- **Tiny Docker Image** - ~2MB multi-arch image (amd64/arm64)

## Quick Start

```bash
# Using Docker
docker run -p 8080:8080 ghcr.io/rashpile/pako-justdoc:latest

# Or build from source
make build && make run
```

## Usage

### Store a Document

```bash
curl -X POST http://localhost:8080/myapp/settings \
  -H "Content-Type: application/json" \
  -d '{"theme": "dark", "language": "en"}'
```

Response (201 Created):
```json
{"status": "created", "channel": "myapp", "document": "settings"}
```

### Retrieve a Document

```bash
curl http://localhost:8080/myapp/settings
```

Response (200 OK):
```json
{"theme": "dark", "language": "en"}
```

### Update a Document

```bash
curl -X POST http://localhost:8080/myapp/settings \
  -H "Content-Type: application/json" \
  -d '{"theme": "light", "language": "fr"}'
```

Response (200 OK):
```json
{"status": "updated", "channel": "myapp", "document": "settings"}
```

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/{channel}/{document}` | Retrieve a document |
| `POST` | `/{channel}/{document}` | Store or update a document |
| `GET` | `/openapi.json` | OpenAPI 3.0 specification |

### Naming Rules

- **Allowed characters**: `a-z`, `A-Z`, `0-9`, `-`, `_`
- **Max length**: 128 characters
- **Case-sensitive**: `MyApp` and `myapp` are different

### Error Responses

| Status | Error Code | Description |
|--------|------------|-------------|
| 400 | `invalid_json` | Request body is not valid JSON |
| 400 | `invalid_name` | Channel or document name is invalid |
| 404 | `not_found` | Document does not exist |
| 413 | `payload_too_large` | Request body exceeds 10MB |

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DB_PATH` | `justdoc.db` | Path to database file |

## Development

```bash
make build    # Compile binary
make run      # Build and run server
make test     # Run tests with coverage
make lint     # Run linter
make docker   # Build Docker image
make clean    # Remove build artifacts
```

## License

MIT
