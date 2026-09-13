package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0},
			b:        []float32{1, 0},
			expected: 1,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0},
			b:        []float32{0, 1},
			expected: 0,
		},
		{
			name:     "same direction",
			a:        []float32{1, 1},
			b:        []float32{1, 1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cosineSimilarity(tt.a, tt.b)
			if err != nil {
				t.Fatal(err)
			}

			if result != tt.expected {
				t.Errorf("expected %f, got %f for %s", tt.expected, result, tt.name)
			}
		})
	}
}

type fakeDreamRepository struct {
	dreams     map[string]dream.Dream
	embeddings map[string]dream.DreamEmbedding
}

func (f *fakeDreamRepository) GetDream(id string) (dream.Dream, error) {
	d, ok := f.dreams[id]
	if !ok {
		return dream.Dream{}, fmt.Errorf("dream not found")
	}

	return d, nil
}

func (f *fakeDreamRepository) GetDreamEmbedding(id string) (*dream.DreamEmbedding, error) {
	e, ok := f.embeddings[id]
	if !ok {
		return nil, fmt.Errorf("embedding not found")
	}

	return &e, nil
}

func (f *fakeDreamRepository) ListDreamEmbeddings() ([]dream.DreamEmbedding, error) {
	result := make([]dream.DreamEmbedding, 0, len(f.embeddings))

	for _, e := range f.embeddings {
		result = append(result, e)
	}

	return result, nil
}

func (f *fakeDreamRepository) DeleteDreamEmbedding(dreamID string) error {
	return nil
}

func (f *fakeDreamRepository) CreateDream(d dream.Dream) error {
	f.dreams[d.ID] = d
	return nil
}

func (f *fakeDreamRepository) ListDreams() ([]dream.Dream, error) {
	result := make([]dream.Dream, 0, len(f.dreams))

	for _, d := range f.dreams {
		result = append(result, d)
	}

	return result, nil
}

func (f *fakeDreamRepository) DeleteDream(id string) error {
	delete(f.dreams, id)
	return nil
}

func (f *fakeDreamRepository) SaveDreamAnalysis(analysis dream.DreamAnalysis) error {
	return nil
}

func (f *fakeDreamRepository) GetDreamAnalysis(dreamID string) (*dream.DreamAnalysis, error) {
	return nil, nil
}

func (f *fakeDreamRepository) SaveDreamEmbedding(embedding dream.DreamEmbedding) error {
	f.embeddings[embedding.DreamID] = embedding
	return nil
}

func TestFindSimilarDreams(t *testing.T) {
	repo := &fakeDreamRepository{
		dreams: map[string]dream.Dream{
			"a": {ID: "a", Title: "Original dream"},
			"b": {ID: "b", Title: "Very similar"},
			"c": {ID: "c", Title: "Somewhat similar"},
			"d": {ID: "d", Title: "Completely different"},
		},
		embeddings: map[string]dream.DreamEmbedding{
			"a": {
				DreamID:   "a",
				Embedding: []float32{1, 0},
			},
			"b": {
				DreamID:   "b",
				Embedding: []float32{1, 0},
			},
			"c": {
				DreamID:   "c",
				Embedding: []float32{0.8, 0.6},
			},
			"d": {
				DreamID:   "d",
				Embedding: []float32{-1, 0},
			},
		},
	}

	service := &DreamService{dreams: repo, embeddings: repo}

	results, err := service.FindSimilarDreams(
		context.Background(),
		"a",
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Dream.ID != "b" {
		t.Errorf("expected first result to be b, got %s", results[0].Dream.ID)
	}

	if results[1].Dream.ID != "c" {
		t.Errorf("expected second result to be c, got %s", results[1].Dream.ID)
	}
}
