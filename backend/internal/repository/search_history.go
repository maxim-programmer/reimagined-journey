package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type SearchHistoryRepository struct {
	db *pgxpool.Pool
}

func NewSearchHistoryRepository(db *pgxpool.Pool) *SearchHistoryRepository {
	return &SearchHistoryRepository{db: db}
}

func (r *SearchHistoryRepository) Add(ctx context.Context, userID, query string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO search_history (user_id, query) VALUES ($1, $2)`,
		userID, query,
	)
	if err != nil {
		return fmt.Errorf("insert search history: %w", err)
	}
	return nil
}

func (r *SearchHistoryRepository) ListByUser(ctx context.Context, userID string, limit int) ([]model.SearchHistory, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, query, created_at
		 FROM search_history
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query search history: %w", err)
	}
	defer rows.Close()

	var items []model.SearchHistory
	for rows.Next() {
		var h model.SearchHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.Query, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan search history: %w", err)
		}
		items = append(items, h)
	}
	return items, rows.Err()
}

func (r *SearchHistoryRepository) DeleteByUser(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM search_history WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete search history: %w", err)
	}
	return nil
}