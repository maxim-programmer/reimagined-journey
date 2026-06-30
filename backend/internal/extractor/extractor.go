package extractor

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

func ExtractPages(filePath, mimeType string) ([]model.PageText, error) {
	switch mimeType {
	case "application/pdf":
		return extractPDFPages(filePath)
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return extractDOCXPages(filePath)
	default:
		return nil, fmt.Errorf("unsupported mime type: %s", mimeType)
	}
}

func ExtractText(filePath, mimeType string) (string, error) {
	pages, err := ExtractPages(filePath, mimeType)
	if err != nil {
		return "", err
	}
	var parts []string
	for _, p := range pages {
		if p.Text != "" {
			parts = append(parts, p.Text)
		}
	}
	return strings.Join(parts, "\n"), nil
}

func extractPDFPages(filePath string) ([]model.PageText, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	totalPages := r.NumPage()
	var pages []model.PageText
	for i := 1; i <= totalPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		text = strings.TrimSpace(text)
		if text != "" {
			pages = append(pages, model.PageText{
				PageNumber: i,
				Text:       text,
			})
		}
	}
	return pages, nil
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

func extractDOCXPages(filePath string) ([]model.PageText, error) {
	zr, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("open docx zip: %w", err)
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
		return nil, fmt.Errorf("word/document.xml not found in docx")
	}

	rc, err := docFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open document.xml: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read document.xml: %w", err)
	}

	var doc docxDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse document.xml: %w", err)
	}

	var sb strings.Builder
	for _, para := range doc.Body.Paragraphs {
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
	}

	fullText := sb.String()
	if fullText == "" {
		return nil, nil
	}

	return []model.PageText{
		{PageNumber: 1, Text: fullText},
	}, nil
}
