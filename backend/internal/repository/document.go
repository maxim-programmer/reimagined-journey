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
		`INSERT INTO documents (id, file_name, file_size, mime_type, status, extracted_text, uploaded_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		doc.ID, doc.FileName, doc.FileSize, doc.MimeType, doc.Status, doc.ExtractedText, doc.UploadedAt,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) List(ctx context.Context) ([]model.Document, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, file_name, file_size, mime_type, status, extracted_text, uploaded_at
		 FROM documents ORDER BY uploaded_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query documents: %w", err)
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		var d model.Document
		if err := rows.Scan(&d.ID, &d.FileName, &d.FileSize, &d.MimeType, &d.Status, &d.ExtractedText, &d.UploadedAt); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *DocumentRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	var d model.Document
	err := r.db.QueryRow(ctx,
		`SELECT id, file_name, file_size, mime_type, status, extracted_text, uploaded_at
		 FROM documents WHERE id = $1`,
		id,
	).Scan(&d.ID, &d.FileName, &d.FileSize, &d.MimeType, &d.Status, &d.ExtractedText, &d.UploadedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get document by id: %w", err)
	}
	return &d, nil
}