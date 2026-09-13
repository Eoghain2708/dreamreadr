package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/google/uuid"
)

type DreamService struct {
	dreams     DreamRepository
	analyses   AnalysisRepository
	embeddings EmbeddingRepository
	features   FeatureRepository

	analyser dream.DreamAnalyser
	embedder dream.DreamEmbedder
}

func NewDreamService(dreams DreamRepository,
	analyses AnalysisRepository,
	embeddings EmbeddingRepository,
	features FeatureRepository,
	analyser dream.DreamAnalyser,
	embedder dream.DreamEmbedder) *DreamService {
	return &DreamService{dreams, analyses, embeddings, features, analyser, embedder}
}

func (s *DreamService) CreateDream(title, rawText string) (*dream.Dream, error) {
	d, err := dream.NewDream(generateID(), title, rawText, time.Now())
	if err != nil {
		return nil, err
	}

	if err := s.dreams.CreateDream(*d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *DreamService) ListDreams() ([]dream.Dream, error) {
	return s.dreams.ListDreams()
}

func (s *DreamService) DeleteDream(id string) error {
	return s.dreams.DeleteDream(id)
}

func (s *DreamService) GetDream(id string) (dream.Dream, error) {
	return s.dreams.GetDream(id)
}

func generateID() string {
	return uuid.NewString()
}

func (s *DreamService) AnalyseDream(ctx context.Context, id string) (*dream.DreamAnalysis, error) {
	d, err := s.dreams.GetDream(id)
	if err != nil {
		return nil, err
	}

	analysis, err := s.analyser.Analyse(ctx, d)
	if err != nil {
		return nil, err
	}

	err = s.analyses.SaveDreamAnalysis(*analysis)
	if err != nil {
		return nil, err
	}

	embedding, err := s.embedder.Embed(ctx, d)
	if err != nil {
		return nil, err
	}

	fmt.Println("saving embedding for dream")

	err = s.embeddings.SaveDreamEmbedding(*embedding)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}

func (s *DreamService) GetDreamAnalysis(dreamID string) (*dream.DreamAnalysis, error) {
	analysis, err := s.analyses.GetDreamAnalysis(dreamID)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}
