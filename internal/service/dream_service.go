package service

import (
	"context"
	"time"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/google/uuid"
)

type DreamService struct {
	repo     dream.DreamRepository
	analyser dream.DreamAnalyser
}

func NewDreamService(repo dream.DreamRepository, analyser dream.DreamAnalyser) *DreamService {
	return &DreamService{repo: repo, analyser: analyser}
}

func (s *DreamService) CreateDream(title, rawText string) (*dream.Dream, error) {
	d, err := dream.NewDream(generateID(), title, rawText, time.Now())
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDream(*d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *DreamService) ListDreams() ([]dream.Dream, error) {
	return s.repo.ListDreams()
}

func (s *DreamService) DeleteDream(id string) error {
	return s.repo.DeleteDream(id)
}

func (s *DreamService) GetDream(id string) (dream.Dream, error) {
	return s.repo.GetDream(id)
}

func generateID() string {
	return uuid.NewString()
}

func (s *DreamService) AnalyseDream(ctx context.Context, id string) (*dream.DreamAnalysis, error) {
	d, err := s.repo.GetDream(id)
	if err != nil {
		return nil, err
	}

	analysis, err := s.analyser.Analyse(ctx, d)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveDreamAnalysis(*analysis)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}

func (s *DreamService) GetDreamAnalysis(dreamID string) (*dream.DreamAnalysis, error) {
	analysis, err := s.repo.GetDreamAnalysis(dreamID)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}
