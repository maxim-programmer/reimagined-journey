package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const indexName = "document_chunks"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type ChunkDocument struct {
	ChunkID    string `json:"chunk_id"`
	DocumentID string `json:"document_id"`
	UserID     string `json:"user_id"`
	FileName   string `json:"file_name"`
	PageNumber int    `json:"page_number"`
	Text       string `json:"text"`
}

type SearchHit struct {
	ChunkID  string  `json:"chunk_id"`
	FileName string  `json:"file_name"`
	Page     int     `json:"page"`
	Text     string  `json:"text"`
	Score    float64 `json:"score"`
}

func (c *Client) EnsureIndex(ctx context.Context) error {
	exists, err := c.indexExists(ctx)
	if err != nil {
		return fmt.Errorf("check index existence: %w", err)
	}
	if exists {
		log.Printf("elasticsearch index %q already exists, skipping creation", indexName)
		return nil
	}

	body := map[string]any{
		"settings": map[string]any{
			"analysis": map[string]any{
				"filter": map[string]any{
					"russian_stop": map[string]any{
						"type":      "stop",
						"stopwords": "_russian_",
					},
					"russian_stemmer": map[string]any{
						"type":     "stemmer",
						"language": "russian",
					},
				},
				"analyzer": map[string]any{
					"russian": map[string]any{
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"russian_stop",
							"russian_stemmer",
						},
					},
				},
			},
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"chunk_id": map[string]any{
					"type": "keyword",
				},
				"document_id": map[string]any{
					"type": "keyword",
				},
				"user_id": map[string]any{
					"type": "keyword",
				},
				"file_name": map[string]any{
					"type": "text",
					"fields": map[string]any{
						"keyword": map[string]any{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"page_number": map[string]any{
					"type": "integer",
				},
				"text": map[string]any{
					"type":     "text",
					"analyzer": "russian",
					"fields": map[string]any{
						"standard": map[string]any{
							"type":     "text",
							"analyzer": "standard",
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal index body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/%s", c.baseURL, indexName),
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create index request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errBody map[string]any
		json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("create index failed with status %d: %v", resp.StatusCode, errBody)
	}

	log.Printf("elasticsearch index %q created successfully", indexName)
	return nil
}

func (c *Client) IndexChunk(ctx context.Context, doc ChunkDocument) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal chunk document: %w", err)
	}

	url := fmt.Sprintf("%s/%s/_doc/%s", c.baseURL, indexName, doc.ChunkID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create index request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("index chunk request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errBody map[string]any
		json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("index chunk failed with status %d: %v", resp.StatusCode, errBody)
	}

	return nil
}

func (c *Client) Search(ctx context.Context, query, userID string) ([]SearchHit, error) {
	body := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": map[string]any{
					"multi_match": map[string]any{
						"query":  query,
						"fields": []string{"text", "text.standard", "file_name"},
						"type":   "best_fields",
					},
				},
				"filter": map[string]any{
					"term": map[string]any{
						"user_id": userID,
					},
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal search body: %w", err)
	}

	url := fmt.Sprintf("%s/%s/_search", c.baseURL, indexName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody map[string]any
		json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("search failed with status %d: %v", resp.StatusCode, errBody)
	}

	var esResp struct {
		Hits struct {
			Hits []struct {
				Score  float64       `json:"_score"`
				Source ChunkDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&esResp); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	hits := make([]SearchHit, 0, len(esResp.Hits.Hits))
	for _, h := range esResp.Hits.Hits {
		hits = append(hits, SearchHit{
			ChunkID:  h.Source.ChunkID,
			FileName: h.Source.FileName,
			Page:     h.Source.PageNumber,
			Text:     h.Source.Text,
			Score:    h.Score,
		})
	}

	return hits, nil
}

func (c *Client) DeleteByDocument(ctx context.Context, documentID string) error {
	body := map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"document_id": documentID,
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal delete by query body: %w", err)
	}

	url := fmt.Sprintf("%s/%s/_delete_by_query", c.baseURL, indexName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create delete by query request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete by query request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody map[string]any
		json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("delete by query failed with status %d: %v", resp.StatusCode, errBody)
	}

	return nil
}

func (c *Client) indexExists(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead,
		fmt.Sprintf("%s/%s", c.baseURL, indexName),
		nil,
	)
	if err != nil {
		return false, fmt.Errorf("create head request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("head request: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}