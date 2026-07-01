package model

import "time"

type Document struct {
	ID            string    `json:"id"`
	UserID        string    `json:"-"`
	FileName      string    `json:"file_name"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	Status        string    `json:"status"`
	ExtractedText string    `json:"extracted_text"`
	UploadedAt    time.Time `json:"uploaded_at"`
}

type PageText struct {
	PageNumber int
	Text       string
}

type Chunk struct {
	ID         int64     `json:"id"`
	DocumentID string    `json:"document_id"`
	ChunkIndex int       `json:"chunk_index"`
	PageNumber int       `json:"page_number"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
