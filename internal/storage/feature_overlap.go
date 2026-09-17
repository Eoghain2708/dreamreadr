package storage

import (
	"fmt"
	"strings"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

func (sr *SQLRepository) FindOverlappingFeatures(values []string) ([]dream.Dream, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("no values provided")
	}

	// Normalise and deduplicate search values.
	normalised := make([]string, 0, len(values))
	seen := make(map[string]struct{})

	for _, v := range values {
		v = normaliseVal(v)

		if v == "" {
			continue
		}

		if _, exists := seen[v]; exists {
			continue
		}

		seen[v] = struct{}{}
		normalised = append(normalised, v)
	}

	if len(normalised) == 0 {
		return nil, fmt.Errorf("no valid values provided")
	}

	requested := make([]string, len(normalised))
	args := make([]any, 0, len(normalised)*2)

	for i, value := range normalised {
		requested[i] = "(?, ?)"
		args = append(args, i, value)
	}

	fmt.Println(requested)
	fmt.Println(args)

	query := fmt.Sprintf(`
		WITH features AS (
			SELECT dream_id, symbol AS value
			FROM dream_symbols

			UNION ALL

			SELECT dream_id, location AS value
			FROM dream_locations

			UNION ALL

			SELECT dream_id, theme AS value
			FROM dream_themes

			UNION ALL

			SELECT dream_id, emotion AS value
			FROM dream_emotions

			UNION ALL

			SELECT dream_id, person AS value
			FROM dream_people
		),

		requested(term_id, value) AS (
			VALUES %s
		),

		matches AS (
			SELECT DISTINCT
				f.dream_id,
				r.term_id
			FROM features f
			JOIN requested r
				ON f.value LIKE '%%' || r.value || '%%'
		)

		SELECT dream_id
		FROM matches
		GROUP BY dream_id
		HAVING COUNT(DISTINCT term_id) = ?
	`, strings.Join(requested, ", "))

	args = append(args, len(normalised))

	rows, err := sr.db.Query(query, args...)
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

		d, err := sr.GetDream(id)
		if err != nil {
			return nil, err
		}

		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func normaliseVal(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func (sr *SQLRepository) FindCommonCoOccs(limit int) ([]dream.FeatureCooccurence, error) {
	query := `
		WITH features AS (
			SELECT dream_id, symbol AS feature
			FROM dream_symbols
			UNION
			SELECT dream_id, location AS feature
			FROM dream_locations
			UNION
			SELECT dream_id, emotion AS feature
			FROM dream_emotions
			UNION
			SELECT dream_id, person AS feature
			FROM dream_people
			UNION
			SELECT dream_id, theme AS feature
			FROM dream_themes
		)

		SELECT
		a.feature as feature_a,
		b.feature as feature_b,
		COUNT (DISTINCT a.dream_id) AS count
		FROM features a
		JOIN features b
			ON a.dream_id = b.dream_id
			AND a.feature < b.feature
		GROUP BY
			a.feature,
			b.feature
		ORDER BY 
			count DESC
		LIMIT ?
	`

	rows, err := sr.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []dream.FeatureCooccurence

	for rows.Next() {
		var fc dream.FeatureCooccurence
		if err := rows.Scan(&fc.FeatureA, &fc.FeatureB, &fc.Count); err != nil {
			return nil, err
		}

		res = append(res, fc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
