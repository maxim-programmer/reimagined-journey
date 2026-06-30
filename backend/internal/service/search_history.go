package service

import (
	"context"
	"fmt"

	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

const historyLimit = 50

type searchHistoryRepo interface {
	Add(ctx context.Context, userID, query string) error
	ListByUser(ctx context.Context, userID string, limit int) ([]model.SearchHistory, error)
	DeleteByUser(ctx context.Context, userID string) error
}

type SearchHistoryService struct {
	repo searchHistoryRepo
}

func NewSearchHistoryService(repo searchHistoryRepo) *SearchHistoryService {
	return &SearchHistoryService{repo: repo}
}

func (s *SearchHistoryService) Add(ctx context.Context, userID, query string) error {
	if err := s.repo.Add(ctx, userID, query); err != nil {
		return fmt.Errorf("add history: %w", err)
	}
	return nil
}

func (s *SearchHistoryService) ListHistory(ctx context.Context, userID string) ([]model.SearchHistory, error) {
	items, err := s.repo.ListByUser(ctx, userID, historyLimit)
	if err != nil {
		return nil, fmt.Errorf("list history: %w", err)
	}
	return items, nil
}

func (s *SearchHistoryService) ClearHistory(ctx context.Context, userID string) error {
	if err := s.repo.DeleteByUser(ctx, userID); err != nil {
		return fmt.Errorf("clear history: %w", err)
	}
	return nil
}
