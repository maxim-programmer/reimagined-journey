package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/chunker"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/elastic"
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
	esClient  *elastic.Client
}

func NewDocumentService(repo documentRepo, chunkRepo chunkRepo, esClient *elastic.Client) *DocumentService {
	return &DocumentService{repo: repo, chunkRepo: chunkRepo, esClient: esClient}
}

func (s *DocumentService) CreateDocument(ctx context.Context, fileName string, fileSize int64, mimeType, filePath string) (*model.Document, error) {
	pages, err := extractor.ExtractPages(filePath, mimeType)
	if err != nil {
		log.Printf("text extraction warning for %s: %v", fileName, err)
	}

	fullText, _ := extractor.ExtractText(filePath, mimeType)

	doc := &model.Document{
		ID:            uuid.New().String(),
		FileName:      fileName,
		FileSize:      fileSize,
		MimeType:      mimeType,
		Status:        "uploaded",
		ExtractedText: fullText,
		UploadedAt:    time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	if len(pages) > 0 {
		chunks := chunker.SplitPages(pages)
		for i := range chunks {
			chunks[i].DocumentID = doc.ID
		}

		if err := s.chunkRepo.CreateBatch(ctx, chunks); err != nil {
			log.Printf("chunk creation warning for document %s: %v", doc.ID, err)
		}

		for _, chunk := range chunks {
			chunkID := fmt.Sprintf("%s-%d", doc.ID, chunk.ChunkIndex)
			esDoc := elastic.ChunkDocument{
				ChunkID:    chunkID,
				DocumentID: doc.ID,
				FileName:   fileName,
				PageNumber: chunk.PageNumber,
				Text:       chunk.Content,
			}
			if err := s.esClient.IndexChunk(ctx, esDoc); err != nil {
				log.Printf("elasticsearch index warning for chunk %s: %v", chunkID, err)
			}
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