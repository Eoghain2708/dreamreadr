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

func (sr *SQLRepository) findDreamsWithFeature(f feature, value string) ([]dream.Dream, error) {
	query := fmt.Sprintf(`
		SELECT dream_id from %s 
		WHERE %s = ?
	`, f.table, f.name)

	rows, err := sr.db.Query(query, value)
	if err != nil {
		return nil, err
	}

	var res []dream.Dream

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		dream, err := sr.GetDream(id)
		if err != nil {
			return nil, err
		}

		res = append(res, dream)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (sr *SQLRepository) FindDreamsWithEmotion(s string) ([]dream.Dream, error) {
	return sr.findDreamsWithFeature(emotions, s)
}

func (sr *SQLRepository) FindDreamsWithSymbol(s string) ([]dream.Dream, error) {
	return sr.findDreamsWithFeature(symbols, s)
}

func (sr *SQLRepository) FindDreamsWithPerson(s string) ([]dream.Dream, error) {
	return sr.findDreamsWithFeature(people, s)
}

func (sr *SQLRepository) FindDreamsWithTheme(s string) ([]dream.Dream, error) {
	return sr.findDreamsWithFeature(themes, s)
}

func (sr *SQLRepository) FindDreamsWithLocation(s string) ([]dream.Dream, error) {
	return sr.findDreamsWithFeature(locations, s)
}
