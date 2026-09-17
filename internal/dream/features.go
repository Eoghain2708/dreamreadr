package dream

type FeatureType string

const (
	Location FeatureType = "location"
	Symbol   FeatureType = "symbol"
	Theme    FeatureType = "theme"
	Person   FeatureType = "person"
	Emotion  FeatureType = "emotion"
	Unknown  FeatureType = "unknown"
)

type FeatureCount struct {
	Type  FeatureType
	Value string
	Count int
}

type FeatureCooccurence struct {
	FeatureA string
	FeatureB string
	Count    int
}

type FeatureFilter struct {
	Type  FeatureType
	Value string
}
