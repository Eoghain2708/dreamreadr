package service

import "github.com/Eoghain2708/dreamreadr/internal/dream"

func (ds *DreamService) FindCommonCoOccs(l int) ([]dream.FeatureCooccurence, error) {
	return ds.features.FindCommonCoOccs(l)
}

func (ds *DreamService) FindOverlappingFeatures(v []string) ([]dream.Dream, error) {
	return ds.features.FindOverlappingFeatures(v)
}
