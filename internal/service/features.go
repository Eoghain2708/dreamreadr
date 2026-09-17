package service

import (
	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

func (ds *DreamService) GetTopLocations(limit int) ([]dream.FeatureCount, error) {
	return ds.features.GetTopLocations(limit)
}

func (ds *DreamService) GetTopSymbols(limit int) ([]dream.FeatureCount, error) {
	return ds.features.GetTopSymbols(limit)
}

func (ds *DreamService) GetTopThemes(limit int) ([]dream.FeatureCount, error) {
	return ds.features.GetTopThemes(limit)
}

func (ds *DreamService) GetTopPeople(limit int) ([]dream.FeatureCount, error) {
	return ds.features.GetTopPeople(limit)
}

func (ds *DreamService) GetTopEmotions(limit int) ([]dream.FeatureCount, error) {
	return ds.features.GetTopEmotions(limit)
}

func (ds *DreamService) FindDreamsWithEmotion(s string) ([]dream.Dream, error) {
	return ds.features.FindDreamsWithEmotion(s)
}

func (ds *DreamService) FindDreamsWithLocation(s string) ([]dream.Dream, error) {
	return ds.features.FindDreamsWithLocation(s)
}

func (ds *DreamService) FindDreamsWithPerson(s string) ([]dream.Dream, error) {
	return ds.features.FindDreamsWithPerson(s)
}

func (ds *DreamService) FindDreamsWithSymbol(s string) ([]dream.Dream, error) {
	return ds.features.FindDreamsWithSymbol(s)
}

func (ds *DreamService) FindDreamsWithTheme(s string) ([]dream.Dream, error) {
	return ds.features.FindDreamsWithTheme(s)
}

func (ds *DreamService) FindDreams(filters ...dream.FeatureFilter) ([]dream.Dream, error) {
	return ds.features.FindDreams(filters)
}
