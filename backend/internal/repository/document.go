package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type DocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, doc *model.Document) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO documents (id, user_id, file_name, file_size, mime_type, status, extracted_text, uploaded_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		doc.ID, doc.UserID, doc.FileName, doc.FileSize, doc.MimeType, doc.Status, doc.ExtractedText, doc.UploadedAt,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) List(ctx context.Context, userID string) ([]model.Document, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, file_name, file_size, mime_type, status, extracted_text, uploaded_at
		 FROM documents WHERE user_id = $1 ORDER BY uploaded_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query documents: %w", err)
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		var d model.Document
		if err := rows.Scan(&d.ID, &d.UserID, &d.FileName, &d.FileSize, &d.MimeType, &d.Status, &d.ExtractedText, &d.UploadedAt); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *DocumentRepository) GetByID(ctx context.Context, id, userID string) (*model.Document, error) {
	var d model.Document
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, file_name, file_size, mime_type, status, extracted_text, uploaded_at
		 FROM documents WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&d.ID, &d.UserID, &d.FileName, &d.FileSize, &d.MimeType, &d.Status, &d.ExtractedText, &d.UploadedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get document by id: %w", err)
	}
	return &d, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id, userID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM documents WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
