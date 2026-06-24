package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS documents (
			id             TEXT PRIMARY KEY,
			file_name      TEXT        NOT NULL,
			file_size      BIGINT      NOT NULL,
			mime_type      TEXT        NOT NULL,
			status         TEXT        NOT NULL DEFAULT 'uploaded',
			extracted_text TEXT        NOT NULL DEFAULT '',
			uploaded_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create table documents: %w", err)
	}

	_, err = db.Exec(ctx, `
		ALTER TABLE documents ADD COLUMN IF NOT EXISTS extracted_text TEXT NOT NULL DEFAULT ''
	`)
	if err != nil {
		return fmt.Errorf("alter table documents: %w", err)
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS document_chunks (
			id           BIGSERIAL   PRIMARY KEY,
			document_id  TEXT        NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
			chunk_index  INT         NOT NULL,
			content      TEXT        NOT NULL,
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create table document_chunks: %w", err)
	}

	_, err = db.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_document_chunks_document_id ON document_chunks(document_id)
	`)
	if err != nil {
		return fmt.Errorf("create index document_chunks: %w", err)
	}

	return nil
}