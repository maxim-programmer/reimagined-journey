package extractor

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractText(filePath, mimeType string) (string, error) {
	switch mimeType {
	case "application/pdf":
		return extractPDF(filePath)
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return extractDOCX(filePath)
	default:
		return "", fmt.Errorf("unsupported mime type: %s", mimeType)
	}
}

func extractPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	totalPages := r.NumPage()
	for i := 1; i <= totalPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		if i < totalPages {
			buf.WriteRune('\n')
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

type docxDocument struct {
	Body docxBody `xml:"body"`
}

type docxBody struct {
	Paragraphs []docxParagraph `xml:"p"`
}

type docxParagraph struct {
	Runs []docxRun `xml:"r"`
}

type docxRun struct {
	Text string `xml:"t"`
}

func extractDOCX(filePath string) (string, error) {
	zr, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("open docx zip: %w", err)
	}
	defer zr.Close()

	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return "", fmt.Errorf("word/document.xml not found in docx")
	}

	rc, err := docFile.Open()
	if err != nil {
		return "", fmt.Errorf("open document.xml: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("read document.xml: %w", err)
	}

	var doc docxDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("parse document.xml: %w", err)
	}

	var sb strings.Builder
	for i, para := range doc.Body.Paragraphs {
		var line strings.Builder
		for _, run := range para.Runs {
			line.WriteString(run.Text)
		}
		text := strings.TrimSpace(line.String())
		if text != "" {
			if sb.Len() > 0 {
				sb.WriteRune('\n')
			}
			sb.WriteString(text)
		}
		_ = i
	}
	return sb.String(), nil
}