package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

type SimilarDream struct {
	Dream      dream.Dream
	Similarity float32
}

// track which dreams are most similar so as not to have to call GetDream() unnecessarily
type SimilarDreamTracker struct {
	DreamID    string
	Similarity float32
}

func (ds *DreamService) FindSimilarDreams(ctx context.Context, dreamID string, limit int) ([]SimilarDream, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit, must be over 0")
	}
	if strings.TrimSpace(dreamID) == "" {
		return nil, fmt.Errorf("invalid: dream ID is blank")
	}
	var similarDreams []SimilarDream

	embedding, err := ds.embeddings.GetDreamEmbedding(dreamID)
	if err != nil {
		return nil, err
	}

	otherEmbeddings, err := ds.embeddings.ListDreamEmbeddings()
	if err != nil {
		return nil, err
	}

	var rankedSimilarity []SimilarDreamTracker

	for _, other := range otherEmbeddings {
		if other.DreamID == dreamID {
			continue
		}

		cosine, err := cosineSimilarity(embedding.Embedding, other.Embedding)
		if err != nil {
			return nil, err
		}
		rankedSimilarity = append(rankedSimilarity,
			SimilarDreamTracker{DreamID: other.DreamID, Similarity: cosine})
	}

	sort.Slice(rankedSimilarity, func(i, j int) bool {
		return rankedSimilarity[i].Similarity > rankedSimilarity[j].Similarity
	})

	if len(rankedSimilarity) > limit {
		rankedSimilarity = rankedSimilarity[:limit]
	}

	for _, embedding := range rankedSimilarity {
		dream, err := ds.GetDream(embedding.DreamID)
		if err != nil {
			return nil, err
		}

		similarDreams = append(similarDreams, SimilarDream{dream, embedding.Similarity})
	}

	return similarDreams, nil
}

func cosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors are not the same length")
	}
	dot := 0.0
	for i := range a {
		dot += float64((a[i] * b[i]))
	}

	magnitudeA := 0.0
	for _, value := range a {
		magnitudeA += float64(value * value)
	}

	magnitudeA = math.Sqrt(magnitudeA)

	magnitudeB := 0.0
	for _, value := range b {
		magnitudeB += float64(value * value)
	}

	magnitudeB = math.Sqrt(magnitudeB)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0, fmt.Errorf("cannot calculate similarity for a zero vector")
	}

	return float32(dot / (magnitudeA * magnitudeB)), nil
}
