package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type documentRepo interface {
	Create(ctx context.Context, doc *model.Document) error
	List(ctx context.Context) ([]model.Document, error)
}

type DocumentService struct {
	repo documentRepo
}

func NewDocumentService(repo documentRepo) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) CreateDocument(ctx context.Context, fileName string, fileSize int64, mimeType string) (*model.Document, error) {
	doc := &model.Document{
		ID:         uuid.New().String(),
		FileName:   fileName,
		FileSize:   fileSize,
		MimeType:   mimeType,
		Status:     "uploaded",
		UploadedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	return doc, nil
}

func (s *DocumentService) ListDocuments(ctx context.Context) ([]model.Document, error) {
	docs, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return docs, nil
}