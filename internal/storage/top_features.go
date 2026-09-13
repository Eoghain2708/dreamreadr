package storage

import (
	"fmt"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

type feature struct {
	name  string
	table string
}

var (
	symbols   = feature{"symbol", "dream_symbols"}
	locations = feature{"location", "dream_locations"}
	themes    = feature{"theme", "dream_themes"}
	emotions  = feature{"emotion", "dream_emotions"}
	people    = feature{"people", "dream_people"}
)

func (sr *SQLRepository) getFeatureCount(f feature, limit int) ([]dream.FeatureCount, error) {
	query := fmt.Sprintf(`
    SELECT %s, COUNT(*)
    FROM %s
    GROUP BY %s
    ORDER BY COUNT(*) DESC
    LIMIT ?
	`, f.name, f.table, f.name)

	rows, err := sr.db.Query(query, limit)

	if err != nil {
		return nil, err
	}

	var res []dream.FeatureCount

	for rows.Next() {
		var fc dream.FeatureCount
		if err = rows.Scan(&fc.Value, &fc.Count); err != nil {
			return nil, err
		}

		res = append(res, fc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (sr *SQLRepository) GetTopLocations(limit int) ([]dream.FeatureCount, error) {
	return sr.getFeatureCount(locations, limit)
}

func (sr *SQLRepository) GetTopSymbols(limit int) ([]dream.FeatureCount, error) {
	return sr.getFeatureCount(symbols, limit)
}

func (sr *SQLRepository) GetTopThemes(limit int) ([]dream.FeatureCount, error) {
	return sr.getFeatureCount(themes, limit)
}

func (sr *SQLRepository) GetTopPeople(limit int) ([]dream.FeatureCount, error) {
	return sr.getFeatureCount(people, limit)
}

func (sr *SQLRepository) GetTopEmotions(limit int) ([]dream.FeatureCount, error) {
	return sr.getFeatureCount(emotions, limit)
}
