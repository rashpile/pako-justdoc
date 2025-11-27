package api

import "net/http"

// NewRouter creates a new HTTP router with the document API routes
func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.json", OpenAPI)
	mux.HandleFunc("GET /_/static/", ServeStatic)
	mux.HandleFunc("GET /_/edit/{channel}/{document}", h.EditorUI)
	mux.HandleFunc("GET /", h.ListChannels)
	mux.HandleFunc("GET /{channel}/", h.ListDocuments)
	mux.HandleFunc("GET /{channel}/{document}", h.GetDocument)
	mux.HandleFunc("POST /{channel}/{document}", h.PostDocument)
	return mux
}
