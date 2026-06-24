package model

import "time"

type Document struct {
	ID         string    `json:"id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}