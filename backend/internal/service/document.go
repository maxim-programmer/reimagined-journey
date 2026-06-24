package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/chunker"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/extractor"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type documentRepo interface {
	Create(ctx context.Context, doc *model.Document) error
	List(ctx context.Context) ([]model.Document, error)
}

type chunkRepo interface {
	CreateBatch(ctx context.Context, chunks []model.Chunk) error
	ListByDocument(ctx context.Context, documentID string) ([]model.Chunk, error)
}

type DocumentService struct {
	repo      documentRepo
	chunkRepo chunkRepo
}

func NewDocumentService(repo documentRepo, chunkRepo chunkRepo) *DocumentService {
	return &DocumentService{repo: repo, chunkRepo: chunkRepo}
}

func (s *DocumentService) CreateDocument(ctx context.Context, fileName string, fileSize int64, mimeType, filePath string) (*model.Document, error) {
	text, err := extractor.ExtractText(filePath, mimeType)
	if err != nil {
		log.Printf("text extraction warning for %s: %v", fileName, err)
	}

	doc := &model.Document{
		ID:            uuid.New().String(),
		FileName:      fileName,
		FileSize:      fileSize,
		MimeType:      mimeType,
		Status:        "uploaded",
		ExtractedText: text,
		UploadedAt:    time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	if text != "" {
		parts := chunker.Split(text)
		chunks := make([]model.Chunk, len(parts))
		for i, part := range parts {
			chunks[i] = model.Chunk{
				DocumentID: doc.ID,
				ChunkIndex: i,
				Content:    part,
			}
		}
		if err := s.chunkRepo.CreateBatch(ctx, chunks); err != nil {
			log.Printf("chunk creation warning for document %s: %v", doc.ID, err)
		}
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