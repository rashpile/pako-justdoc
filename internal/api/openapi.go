package api

import "net/http"

// OpenAPISpec is the OpenAPI 3.0 specification for JustDoc API
const OpenAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "JustDoc API",
    "description": "Simple JSON document storage API for frontend developers",
    "version": "1.0.0",
    "contact": {
      "name": "JustDoc"
    },
    "license": {
      "name": "MIT"
    }
  },
  "servers": [
    {
      "url": "/",
      "description": "Current server"
    }
  ],
  "paths": {
    "/": {
      "get": {
        "summary": "List all channels",
        "description": "Returns a list of all channels with document counts",
        "operationId": "listChannels",
        "tags": ["Channels"],
        "responses": {
          "200": {
            "description": "List of channels retrieved successfully",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/ChannelInfo"
                  }
                },
                "example": [
                  {
                    "name": "app-config",
                    "document_count": 3
                  },
                  {
                    "name": "user-data",
                    "document_count": 12
                  }
                ]
              }
            }
          }
        }
      }
    },
    "/{channel}/": {
      "get": {
        "summary": "List documents in a channel",
        "description": "Returns a list of all document names in the specified channel",
        "operationId": "listDocuments",
        "tags": ["Channels"],
        "parameters": [
          {
            "name": "channel",
            "in": "path",
            "required": true,
            "description": "Channel name (alphanumeric, hyphens, underscores, max 128 chars)",
            "schema": {
              "type": "string",
              "pattern": "^[a-zA-Z0-9_-]{1,128}$"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "List of documents retrieved successfully",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                },
                "example": ["config", "settings", "user-preferences"]
              }
            }
          },
          "400": {
            "description": "Invalid channel name",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "invalid_name",
                  "message": "Invalid channel name"
                }
              }
            }
          },
          "404": {
            "description": "Channel not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "not_found",
                  "message": "Channel not found"
                }
              }
            }
          }
        }
      }
    },
    "/{channel}/{document}": {
      "get": {
        "summary": "Retrieve a document",
        "description": "Retrieves a stored JSON document from the specified channel",
        "operationId": "getDocument",
        "tags": ["Documents"],
        "parameters": [
          {
            "name": "channel",
            "in": "path",
            "required": true,
            "description": "Channel name (alphanumeric, hyphens, underscores, max 128 chars)",
            "schema": {
              "type": "string",
              "pattern": "^[a-zA-Z0-9_-]{1,128}$"
            }
          },
          {
            "name": "document",
            "in": "path",
            "required": true,
            "description": "Document name (alphanumeric, hyphens, underscores, max 128 chars)",
            "schema": {
              "type": "string",
              "pattern": "^[a-zA-Z0-9_-]{1,128}$"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Document retrieved successfully",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "additionalProperties": true
                }
              }
            }
          },
          "400": {
            "description": "Invalid channel or document name",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "invalid_name",
                  "message": "Invalid channel or document name"
                }
              }
            }
          },
          "404": {
            "description": "Document not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "not_found",
                  "message": "Document not found"
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Store or update a document",
        "description": "Stores a JSON document in the specified channel. Creates the channel if it doesn't exist. Returns 201 for new documents, 200 for updates.",
        "operationId": "postDocument",
        "tags": ["Documents"],
        "parameters": [
          {
            "name": "channel",
            "in": "path",
            "required": true,
            "description": "Channel name (alphanumeric, hyphens, underscores, max 128 chars)",
            "schema": {
              "type": "string",
              "pattern": "^[a-zA-Z0-9_-]{1,128}$"
            }
          },
          {
            "name": "document",
            "in": "path",
            "required": true,
            "description": "Document name (alphanumeric, hyphens, underscores, max 128 chars)",
            "schema": {
              "type": "string",
              "pattern": "^[a-zA-Z0-9_-]{1,128}$"
            }
          }
        ],
        "requestBody": {
          "required": true,
          "description": "JSON document to store (max 10MB)",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "additionalProperties": true
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Document updated successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SuccessResponse"
                },
                "example": {
                  "status": "updated",
                  "channel": "myapp",
                  "document": "settings"
                }
              }
            }
          },
          "201": {
            "description": "Document created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SuccessResponse"
                },
                "example": {
                  "status": "created",
                  "channel": "myapp",
                  "document": "settings"
                }
              }
            }
          },
          "400": {
            "description": "Invalid request (bad JSON or invalid name)",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "examples": {
                  "invalid_json": {
                    "value": {
                      "error": "invalid_json",
                      "message": "Invalid JSON body"
                    }
                  },
                  "invalid_name": {
                    "value": {
                      "error": "invalid_name",
                      "message": "Invalid channel or document name"
                    }
                  }
                }
              }
            }
          },
          "413": {
            "description": "Payload too large (exceeds 10MB)",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "payload_too_large",
                  "message": "Request body exceeds 10MB limit"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "ChannelInfo": {
        "type": "object",
        "required": ["name", "document_count"],
        "properties": {
          "name": {
            "type": "string",
            "description": "Channel name"
          },
          "document_count": {
            "type": "integer",
            "description": "Number of documents in the channel"
          }
        }
      },
      "SuccessResponse": {
        "type": "object",
        "required": ["status", "channel", "document"],
        "properties": {
          "status": {
            "type": "string",
            "enum": ["created", "updated"],
            "description": "Operation result"
          },
          "channel": {
            "type": "string",
            "description": "Channel name"
          },
          "document": {
            "type": "string",
            "description": "Document name"
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "required": ["error", "message"],
        "properties": {
          "error": {
            "type": "string",
            "enum": ["invalid_json", "invalid_name", "not_found", "payload_too_large"],
            "description": "Error code"
          },
          "message": {
            "type": "string",
            "description": "Human-readable error message"
          }
        }
      }
    }
  },
  "tags": [
    {
      "name": "Channels",
      "description": "Channel and document listing operations"
    },
    {
      "name": "Documents",
      "description": "Document storage operations"
    }
  ]
}`

// OpenAPI serves the OpenAPI specification
func OpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(OpenAPISpec))
}
