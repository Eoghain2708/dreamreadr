package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"google.golang.org/genai"
)

type GeminiEmbedder struct {
	Client *genai.Client
	Model  string
}

func NewGeminiEmbedder(model string) (*GeminiEmbedder, error) {
	if strings.TrimSpace(model) == "" {
		return nil, fmt.Errorf("model is blank")
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})

	if err != nil {
		return nil, err
	}

	return &GeminiEmbedder{Client: client, Model: model}, nil
}

func (e *GeminiEmbedder) Embed(ctx context.Context, d dream.Dream) (*dream.DreamEmbedding, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(d.RawText, genai.RoleUser),
	}

	res, err := e.Client.Models.EmbedContent(ctx, e.Model, contents, nil)
	if err != nil {
		return nil, err
	}

	if len(res.Embeddings) == 0 {
		return nil, fmt.Errorf("gemini returned no embeddings")
	}

	embedding := res.Embeddings[0].Values

	return &dream.DreamEmbedding{
		DreamID:    d.ID,
		Model:      e.Model,
		Dimensions: len(embedding),
		Embedding:  embedding,
	}, nil

}
