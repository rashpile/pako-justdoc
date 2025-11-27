package api

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/rashpile/pako-justdoc/internal/model"
)

//go:embed static/*
var staticFiles embed.FS

var editorTemplate *template.Template

func init() {
	editorTemplate = template.Must(template.ParseFS(staticFiles, "static/editor.html"))
}

// ServeStatic serves embedded static files
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Strip /_/static/ prefix and serve from embedded fs
	sub, _ := fs.Sub(staticFiles, "static")
	http.StripPrefix("/_/static/", http.FileServer(http.FS(sub))).ServeHTTP(w, r)
}

// EditorUI serves the document editor HTML page
func (h *Handler) EditorUI(w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")
	document := r.PathValue("document")

	if !model.IsValidName(channel) || !model.IsValidName(document) {
		http.Error(w, "Invalid channel or document name", http.StatusBadRequest)
		return
	}

	data := struct {
		Channel  string
		Document string
	}{
		Channel:  channel,
		Document: document,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = editorTemplate.Execute(w, data)
}
