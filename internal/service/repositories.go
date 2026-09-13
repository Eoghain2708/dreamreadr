package service

import "github.com/Eoghain2708/dreamreadr/internal/dream"

type DreamRepository interface {
	CreateDream(d dream.Dream) error
	ListDreams() ([]dream.Dream, error)
	GetDream(id string) (dream.Dream, error)
	DeleteDream(id string) error
}

type AnalysisRepository interface {
	SaveDreamAnalysis(da dream.DreamAnalysis) error
	GetDreamAnalysis(dreamID string) (*dream.DreamAnalysis, error)
	DeleteDreamAnalysis(dreamID string) error
}

type EmbeddingRepository interface {
	SaveDreamEmbedding(de dream.DreamEmbedding) error
	GetDreamEmbedding(dreamID string) (*dream.DreamEmbedding, error)
	DeleteDreamEmbedding(dreamID string) error
	ListDreamEmbeddings() ([]dream.DreamEmbedding, error)
}

type FeatureRepository interface {
	GetTopLocations(limit int) ([]dream.FeatureCount, error)
	GetTopEmotions(limit int) ([]dream.FeatureCount, error)
	GetTopPeople(limit int) ([]dream.FeatureCount, error)
	GetTopSymbols(limit int) ([]dream.FeatureCount, error)
	GetTopThemes(limit int) ([]dream.FeatureCount, error)

	FindDreamsWithEmotion(s string) ([]dream.Dream, error)
	FindDreamsWithLocation(s string) ([]dream.Dream, error)
	FindDreamsWithPerson(s string) ([]dream.Dream, error)
	FindDreamsWithSymbol(s string) ([]dream.Dream, error)
	FindDreamsWithTheme(s string) ([]dream.Dream, error)
}
