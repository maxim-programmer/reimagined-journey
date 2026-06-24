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
	"time"

	"github.com/google/uuid"
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

type documentRepo interface {
	Create(ctx context.Context, doc *model.Document) error
	List(ctx context.Context) ([]model.Document, error)
}

type DocumentHandler struct {
	repo      documentRepo
	uploadDir string
}

func NewDocumentHandler(repo documentRepo, uploadDir string) *DocumentHandler {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}
	return &DocumentHandler{repo: repo, uploadDir: uploadDir}
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

	id := uuid.NewString()
	destPath := filepath.Join(h.uploadDir, id+ext)

	dst, err := os.Create(destPath)
	if err != nil {
		log.Printf("create file error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("copy file error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	doc := &model.Document{
		ID:         id,
		FileName:   header.Filename,
		FileSize:   header.Size,
		MimeType:   mimeType,
		Status:     "uploaded",
		UploadedAt: time.Now().UTC(),
	}

	if err := h.repo.Create(r.Context(), doc); err != nil {
		log.Printf("db insert error: %v", err)
		os.Remove(destPath)
		writeError(w, http.StatusInternalServerError, "failed to save document metadata")
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	docs, err := h.repo.List(r.Context())
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