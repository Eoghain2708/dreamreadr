package dream

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Dream struct {
	ID        string
	Title     string
	CreatedAt time.Time
	RawText   string
}

func NewDream(id, title, rawText string, createdAt time.Time) (*Dream, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id is blank")
	}

	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title is blank")
	}

	if strings.TrimSpace(rawText) == "" {
		return nil, fmt.Errorf("content is blank")
	}

	if createdAt.IsZero() {
		return nil, fmt.Errorf("created at not provided")
	}

	return &Dream{
		ID: id, Title: title, RawText: rawText, CreatedAt: createdAt,
	}, nil
}

type DreamAnalysis struct {
	ID        int
	DreamID   string
	Summary   string
	Themes    []string
	Emotions  []string
	Locations []string
	People    []string
	Symbols   []string
}

type DreamRepository interface {
	CreateDream(d Dream) error
	ListDreams() ([]Dream, error)
	GetDream(id string) (Dream, error)
	DeleteDream(id string) error

	SaveDreamAnalysis(da DreamAnalysis) error
	GetDreamAnalysis(dreamID string) (*DreamAnalysis, error)
	DeleteDreamAnalysis(dreamID string) error

	SaveDreamEmbedding(de DreamEmbedding) error
	GetDreamEmbedding(dreamID string) (*DreamEmbedding, error)
	DeleteDreamEmbedding(dreamID string) error
	ListDreamEmbeddings() ([]DreamEmbedding, error)
}

type DreamAnalyser interface {
	Analyse(ctx context.Context, d Dream) (*DreamAnalysis, error)
}

type FakeAnalyser struct{}

func (a *FakeAnalyser) Analyse(ctx context.Context, d Dream) (*DreamAnalysis, error) {
	return &DreamAnalysis{
		DreamID:   d.ID,
		Summary:   "A test analysis",
		Themes:    []string{"testing", "romance"},
		Emotions:  []string{"misery"},
		Locations: []string{"Belfast", "Derry"},
		People:    []string{"Eamonn", "Amy"},
		Symbols:   []string{"Trophy", "Red Rose"},
	}, nil
}

type DreamEmbedding struct {
	DreamID    string
	Model      string
	Dimensions int
	Embedding  []float32
}

type DreamEmbedder interface {
	Embed(ctx context.Context, d Dream) (*DreamEmbedding, error)
}
