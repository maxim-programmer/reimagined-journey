package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxim-programmer/reimagined-journey/backend/internal/elastic"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/middleware"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

const (
	maxFileSize = 20 << 20
)

var allowedExtensions = map[string]string{
	".pdf":  "application/pdf",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
}

var zipMagic = []byte{0x50, 0x4B, 0x03, 0x04}
var pdfMagic = []byte{0x25, 0x50, 0x44, 0x46}

type documentService interface {
	CreateDocument(ctx context.Context, fileName string, fileSize int64, mimeType, filePath, userID string) (*model.Document, error)
	ListDocuments(ctx context.Context, userID string) ([]model.Document, error)
	GetDocument(ctx context.Context, id, userID string) (*model.Document, error)
	DeleteDocument(ctx context.Context, id, userID string) error
	Search(ctx context.Context, query, userID string) ([]elastic.SearchHit, error)
}

type DocumentHandler struct {
	svc       documentService
	uploadDir string
}

func NewDocumentHandler(svc documentService, uploadDir string) *DocumentHandler {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}
	return &DocumentHandler{svc: svc, uploadDir: uploadDir}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func detectMIMEType(buf []byte, ext string) (string, bool) {
	switch ext {
	case ".pdf":
		if len(buf) >= 4 && bytes.Equal(buf[:4], pdfMagic) {
			return "application/pdf", true
		}
	case ".docx":
		if len(buf) >= 4 && bytes.Equal(buf[:4], zipMagic) {
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", true
		}
	}
	return "", false
}

func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or malformed")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "field 'file' is required")
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		writeError(w, http.StatusBadRequest, "file size exceeds 20 MB limit")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if _, ok := allowedExtensions[ext]; !ok {
		writeError(w, http.StatusBadRequest, "unsupported file type: only PDF and DOCX are allowed")
		return
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		writeError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	mimeType, ok := detectMIMEType(buf[:n], ext)
	if !ok {
		writeError(w, http.StatusBadRequest, "file content does not match its extension: only PDF and DOCX are allowed")
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process file")
		return
	}

	tmp, err := os.CreateTemp(h.uploadDir, "upload-*"+ext)
	if err != nil {
		log.Printf("create temp file error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	tmpPath := tmp.Name()

	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		log.Printf("copy file error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	tmp.Close()

	doc, err := h.svc.CreateDocument(r.Context(), header.Filename, header.Size, mimeType, tmpPath, userID)
	if err != nil {
		os.Remove(tmpPath)
		log.Printf("create document error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save document metadata")
		return
	}

	finalPath := filepath.Join(h.uploadDir, doc.ID+ext)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		log.Printf("rename file error: %v", err)
	}

	writeJSON(w, http.StatusCreated, doc)
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	docs, err := h.svc.ListDocuments(r.Context(), userID)
	if err != nil {
		log.Printf("list documents error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to retrieve documents")
		return
	}
	if docs == nil {
		docs = []model.Document{}
	}
	writeJSON(w, http.StatusOK, docs)
}

func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "document id is required")
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.svc.GetDocument(r.Context(), id, userID)
	if err != nil {
		log.Printf("get document error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to retrieve document")
		return
	}
	if doc == nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}

	writeJSON(w, http.StatusOK, doc)
}

func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "document id is required")
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.svc.GetDocument(r.Context(), id, userID)
	if err != nil {
		log.Printf("get document error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to retrieve document")
		return
	}
	if doc == nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if err := h.svc.DeleteDocument(r.Context(), id, userID); err != nil {
		log.Printf("delete document error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to delete document")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *DocumentHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	hits, err := h.svc.Search(r.Context(), query, userID)
	if err != nil {
		log.Printf("search error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to perform search")
		return
	}

	writeJSON(w, http.StatusOK, hits)
}
