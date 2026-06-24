package chunker

import (
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

const (
	ChunkSize    = 1000
	ChunkOverlap = 100
)

func Split(text string) []string {
	runes := []rune(text)
	total := len(runes)

	if total == 0 {
		return nil
	}

	var chunks []string
	step := ChunkSize - ChunkOverlap

	for start := 0; start < total; start += step {
		end := start + ChunkSize
		if end > total {
			end = total
		}
		chunks = append(chunks, string(runes[start:end]))
		if end == total {
			break
		}
	}

	return chunks
}

func SplitPages(pages []model.PageText) []model.Chunk {
	var result []model.Chunk
	chunkIndex := 0
	step := ChunkSize - ChunkOverlap

	for _, page := range pages {
		runes := []rune(page.Text)
		total := len(runes)
		if total == 0 {
			continue
		}

		for start := 0; start < total; start += step {
			end := start + ChunkSize
			if end > total {
				end = total
			}
			result = append(result, model.Chunk{
				ChunkIndex: chunkIndex,
				PageNumber: page.PageNumber,
				Content:    string(runes[start:end]),
			})
			chunkIndex++
			if end == total {
				break
			}
		}
	}

	return result
}