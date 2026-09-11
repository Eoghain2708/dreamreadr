package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

type LocalLLMEmbedder struct {
	Client  *http.Client
	BaseURL string
	Model   string
}

func NewLocalLLMEmbedder(url, model string) (*LocalLLMEmbedder, error) {
	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("url is blank")
	}

	if strings.TrimSpace(model) == "" {
		return nil, fmt.Errorf("model is blank")
	}

	return &LocalLLMEmbedder{Client: &http.Client{}, BaseURL: url, Model: model}, nil
}

type localLLMEmbedderResponse struct {
	Data []embeddingResponse `json:"data"`
}

type embeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type embeddingRequest struct {
	Input string `json:"input"`
}

func (e *LocalLLMEmbedder) Embed(ctx context.Context, d dream.Dream) (*dream.DreamEmbedding, error) {
	reqBody := embeddingRequest{
		Input: d.RawText,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling dream text, %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, e.BaseURL+"/v1/embeddings", bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create http request, %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request, %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid request: %v - %v\n FOR REQUEST: %v", resp.StatusCode, resp.Status, strings.TrimSpace(string(body)))
	}

	var result localLLMEmbedderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 || len(result.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("embedding contained no data")
	}

	embedding := result.Data[0].Embedding

	return &dream.DreamEmbedding{
		DreamID:    d.ID,
		Model:      e.Model,
		Dimensions: len(embedding),
		Embedding:  embedding,
	}, nil

}
