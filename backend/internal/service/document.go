package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/cache"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/chunker"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/elastic"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/extractor"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type documentRepo interface {
	Create(ctx context.Context, doc *model.Document) error
	List(ctx context.Context, userID string) ([]model.Document, error)
	GetByID(ctx context.Context, id, userID string) (*model.Document, error)
	Delete(ctx context.Context, id, userID string) error
}

type chunkRepo interface {
	CreateBatch(ctx context.Context, chunks []model.Chunk) error
	ListByDocument(ctx context.Context, documentID string) ([]model.Chunk, error)
}

type DocumentService struct {
	repo        documentRepo
	chunkRepo   chunkRepo
	esClient    *elastic.Client
	cache       *cache.RedisCache
	historyRepo searchHistoryRepo
	uploadDir   string
}

func NewDocumentService(repo documentRepo, chunkRepo chunkRepo, esClient *elastic.Client, redisCache *cache.RedisCache, historyRepo searchHistoryRepo, uploadDir string) *DocumentService {
	return &DocumentService{
		repo:        repo,
		chunkRepo:   chunkRepo,
		esClient:    esClient,
		cache:       redisCache,
		historyRepo: historyRepo,
		uploadDir:   uploadDir,
	}
}

func (s *DocumentService) CreateDocument(ctx context.Context, fileName string, fileSize int64, mimeType, filePath, userID string) (*model.Document, error) {
	pages, err := extractor.ExtractPages(filePath, mimeType)
	if err != nil {
		log.Printf("text extraction warning for %s: %v", fileName, err)
	}

	fullText, _ := extractor.ExtractText(filePath, mimeType)

	doc := &model.Document{
		ID:            uuid.New().String(),
		UserID:        userID,
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
				UserID:     userID,
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

func (s *DocumentService) ListDocuments(ctx context.Context, userID string) ([]model.Document, error) {
	docs, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return docs, nil
}

func (s *DocumentService) GetDocument(ctx context.Context, id, userID string) (*model.Document, error) {
	doc, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	return doc, nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, id, userID string) error {
	doc, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("get document: %w", err)
	}
	if doc == nil {
		return fmt.Errorf("document not found")
	}

	if err := s.esClient.DeleteByDocument(ctx, id); err != nil {
		log.Printf("elasticsearch delete warning for document %s: %v", id, err)
	}

	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	ext := filepath.Ext(doc.FileName)
	filePath := filepath.Join(s.uploadDir, doc.ID+ext)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("file removal warning for document %s: %v", id, err)
	}

	return nil
}

func (s *DocumentService) Search(ctx context.Context, query, userID string) ([]elastic.SearchHit, error) {
	if query == "" {
		return []elastic.SearchHit{}, nil
	}

	cacheKey := cache.SearchKey(query, userID)

	var cached []elastic.SearchHit
	hit, err := s.cache.Get(ctx, cacheKey, &cached)
	if err != nil {
		log.Printf("cache get warning for query %q: %v", query, err)
	}

	var hits []elastic.SearchHit
	if hit {
		log.Printf("cache hit for query %q", query)
		hits = cached
	} else {
		hits, err = s.esClient.Search(ctx, query, userID)
		if err != nil {
			return nil, fmt.Errorf("search: %w", err)
		}
		if setErr := s.cache.Set(ctx, cacheKey, hits, cache.SearchTTL); setErr != nil {
			log.Printf("cache set warning for query %q: %v", query, setErr)
		}
	}

	if userID != "" && s.historyRepo != nil {
		if addErr := s.historyRepo.Add(ctx, userID, query); addErr != nil {
			log.Printf("history add warning for user %s: %v", userID, addErr)
		}
	}

	return hits, nil
}
