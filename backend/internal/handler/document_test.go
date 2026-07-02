package handler

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDetectMIMEType(t *testing.T) {
    tests := []struct {
        name     string
        data     []byte
        ext      string
        expectedMIME string
        expectedValid bool
    }{
        {"valid PDF", []byte("%PDF-1.4"), ".pdf", "application/pdf", true},
        {"valid DOCX", []byte("PK\x03\x04"), ".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", true},
        {"fake PDF", []byte("MZ\x90\x00"), ".pdf", "", false},
        {"empty file", []byte{}, ".pdf", "", false},
        {"wrong extension", []byte("%PDF-1.4"), ".docx", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mimeType, valid := detectMIMEType(tt.data, tt.ext)
            assert.Equal(t, tt.expectedMIME, mimeType)
            assert.Equal(t, tt.expectedValid, valid)
        })
    }
}
