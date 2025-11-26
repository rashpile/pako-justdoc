package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rashpile/pako-justdoc/internal/model"
	"github.com/rashpile/pako-justdoc/internal/storage"
)

// MaxBodySize is the maximum allowed request body size (10MB)
const MaxBodySize = 10 * 1024 * 1024

// Handler handles HTTP requests for the document API
type Handler struct {
	storage storage.Storage
}

// NewHandler creates a new Handler with the given storage
func NewHandler(s storage.Storage) *Handler {
	return &Handler{storage: s}
}

// PostDocument handles POST /{channel}/{document}
func (h *Handler) PostDocument(w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")
	document := r.PathValue("document")

	// Validate names
	if !model.IsValidName(channel) || !model.IsValidName(document) {
		writeError(w, http.StatusBadRequest, model.ErrCodeInvalidName, "Invalid channel or document name")
		return
	}

	// Limit body size
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

	// Read body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		if err.Error() == "http: request body too large" {
			writeError(w, http.StatusRequestEntityTooLarge, model.ErrCodePayloadTooLarge, "Request body exceeds 10MB limit")
			return
		}
		writeError(w, http.StatusBadRequest, model.ErrCodeInvalidJSON, "Failed to read request body")
		return
	}

	// Validate JSON
	if !json.Valid(data) {
		writeError(w, http.StatusBadRequest, model.ErrCodeInvalidJSON, "Invalid JSON body")
		return
	}

	// Store document
	created, err := h.storage.PutDocument(channel, document, data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to store document")
		return
	}

	// Build response
	status := "updated"
	statusCode := http.StatusOK
	if created {
		status = "created"
		statusCode = http.StatusCreated
	}

	writeJSON(w, statusCode, model.SuccessResponse{
		Status:   status,
		Channel:  channel,
		Document: document,
	})
}

// GetDocument handles GET /{channel}/{document}
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")
	document := r.PathValue("document")

	// Validate names
	if !model.IsValidName(channel) || !model.IsValidName(document) {
		writeError(w, http.StatusBadRequest, model.ErrCodeInvalidName, "Invalid channel or document name")
		return
	}

	data, err := h.storage.GetDocument(channel, document)
	if err == storage.ErrNotFound {
		writeError(w, http.StatusNotFound, model.ErrCodeNotFound, "Document not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func writeError(w http.ResponseWriter, statusCode int, errCode, message string) {
	writeJSON(w, statusCode, model.ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}
