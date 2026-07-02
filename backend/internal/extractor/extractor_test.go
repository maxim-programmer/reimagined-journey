package extractor

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExtractText(t *testing.T) {
    t.Run("unsupported mime type returns error", func(t *testing.T) {
        _, err := ExtractText("test", "image/png")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "unsupported mime type")
    })

    t.Run("empty content returns error", func(t *testing.T) {
        _, err := ExtractText("", "application/pdf")
        assert.Error(t, err)
    })
}

func TestExtractPDFPages(t *testing.T) {
    t.Skip("integration test - requires fixture files")
}

func TestExtractDOCXPages(t *testing.T) {
    t.Skip("integration test - requires fixture files")
}
