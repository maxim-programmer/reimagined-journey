package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type ChunkRepository struct {
	db *pgxpool.Pool
}

func NewChunkRepository(db *pgxpool.Pool) *ChunkRepository {
	return &ChunkRepository{db: db}
}

func (r *ChunkRepository) CreateBatch(ctx context.Context, chunks []model.Chunk) error {
	if len(chunks) == 0 {
		return nil
	}

	rows := make([][]any, len(chunks))
	for i, c := range chunks {
		rows[i] = []any{c.DocumentID, c.ChunkIndex, c.PageNumber, c.Content}
	}

	_, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{"document_chunks"},
		[]string{"document_id", "chunk_index", "page_number", "content"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("copy chunks: %w", err)
	}

	return nil
}

func (r *ChunkRepository) ListByDocument(ctx context.Context, documentID string) ([]model.Chunk, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, document_id, chunk_index, page_number, content, created_at
		 FROM document_chunks
		 WHERE document_id = $1
		 ORDER BY chunk_index`,
		documentID,
	)
	if err != nil {
		return nil, fmt.Errorf("query chunks: %w", err)
	}
	defer rows.Close()

	var chunks []model.Chunk
	for rows.Next() {
		var c model.Chunk
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.ChunkIndex, &c.PageNumber, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan chunk: %w", err)
		}
		chunks = append(chunks, c)
	}
	return chunks, rows.Err()
}
