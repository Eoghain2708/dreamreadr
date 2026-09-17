package storage

import (
	"fmt"
	"strings"

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
	people    = feature{"person", "dream_people"}
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

	defer rows.Close()

	var res []dream.FeatureCount

	for rows.Next() {
		var fc dream.FeatureCount
		if err = rows.Scan(&fc.Value, &fc.Count); err != nil {
			return nil, err
		}

		fc.Type = defineFeatureType(f)
		res = append(res, fc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func defineFeatureType(fc feature) dream.FeatureType {
	switch fc {
	case symbols:
		return dream.Symbol
	case emotions:
		return dream.Emotion
	case people:
		return dream.Person
	case themes:
		return dream.Theme
	case locations:
		return dream.Location
	default:
		return dream.Unknown
	}
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

	defer rows.Close()

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

func (sr *SQLRepository) FindDreams(filters []dream.FeatureFilter) ([]dream.Dream, error) {
	var conditions []string
	var args []any

	for _, f := range filters {
		conditions = append(conditions, "(feature_type = ? AND value LIKE ?)")
		args = append(args, strings.ToLower(string(f.Type)), strings.ToLower("%"+f.Value+"%"))
	}

	query := fmt.Sprintf(`
		WITH features AS (
			SELECT dream_id, 'symbol' AS feature_type, symbol as value
			FROM dream_symbols
			UNION ALL
			SELECT dream_id, 'theme', theme as value
			FROM dream_themes
			UNION ALL
			SELECT dream_id, 'person', person as value
			FROM dream_people
			UNION ALL
			SELECT dream_id, 'location', location as value
			FROM dream_locations
			UNION ALL
			SELECT dream_id, 'emotion', emotion as value
			FROM dream_emotions
		)
		SELECT d.id, d.title, d.raw_text
		FROM dreams d
		JOIN (
			SELECT dream_id
			FROM features
			WHERE %s
			GROUP BY dream_id
			HAVING COUNT(*) = ?
		) matching
		 ON matching.dream_id = d.id
	`, strings.Join(conditions, " OR "))

	args = append(args, len(filters))

	rows, err := sr.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []dream.Dream

	for rows.Next() {
		var d dream.Dream
		if err := rows.Scan(&d.ID, &d.Title, &d.RawText); err != nil {
			return nil, err
		}

		res = append(res, d)
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
