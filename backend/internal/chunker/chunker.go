package chunker

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