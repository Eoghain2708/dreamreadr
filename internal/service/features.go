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
