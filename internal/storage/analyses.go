package storage

import (
	"database/sql"
	"strings"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

func (sr *SQLRepository) SaveDreamAnalysis(da dream.DreamAnalysis) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := sr.DeleteDreamAnalysis(da.DreamID); err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO dream_analyses (
            dream_id, summary
        ) VALUES (?, ?)
    `, da.DreamID, da.Summary)

	if err != nil {
		return err
	}

	for _, theme := range uniqueStrings(da.Themes) {
		_, err = tx.Exec(`
            INSERT INTO dream_themes (
                dream_id, theme
            ) VALUES (?, ?)
        `, da.DreamID, theme)

		if err != nil {
			return err
		}
	}

	for _, emotion := range uniqueStrings(da.Emotions) {
		_, err = tx.Exec(`
            INSERT INTO dream_emotions (
                dream_id, emotion
            ) VALUES (?, ?)
        `, da.DreamID, emotion)

		if err != nil {
			return err
		}
	}

	for _, location := range uniqueStrings(da.Locations) {
		_, err = tx.Exec(`
            INSERT INTO dream_locations (
                dream_id, location
            ) VALUES (?, ?)
        `, da.DreamID, location)

		if err != nil {
			return err
		}
	}

	for _, person := range uniqueStrings(da.People) {
		_, err = tx.Exec(`
			INSERT INTO dream_people (
			dream_id, person
			) VALUES (?, ?)
		`, da.DreamID, person)

		if err != nil {
			return err
		}
	}

	for _, symbol := range uniqueStrings(da.Symbols) {
		_, err = tx.Exec(`
			INSERT INTO dream_symbols (
			dream_id, symbol
			) VALUES (?, ?)
		`, da.DreamID, symbol)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (sr *SQLRepository) GetDreamAnalysis(dreamID string) (*dream.DreamAnalysis, error) {
	var da dream.DreamAnalysis

	err := sr.db.QueryRow(`
	SELECT id, dream_id, summary
	FROM dream_analyses
	WHERE dream_id = ?
	`, dreamID).Scan(
		&da.ID,
		&da.DreamID,
		&da.Summary,
	)

	if err != nil {
		return nil, err
	}

	rows, err := sr.db.Query(`
		SELECT theme
		FROM dream_themes
		WHERE dream_id = ?
	`, dreamID)

	if err != nil {
		return nil, err
	}

	da.Themes, err = scanStrings(rows)
	if err != nil {
		return nil, err
	}

	rows, err = sr.db.Query(`
		SELECT emotion 
		FROM dream_emotions
		WHERE dream_id = ?
	`, dreamID)

	da.Emotions, err = scanStrings(rows)
	if err != nil {
		return nil, err
	}

	rows, err = sr.db.Query(`
		SELECT location
		FROM dream_locations
		WHERE dream_id = ?
	`, dreamID)

	da.Locations, err = scanStrings(rows)
	if err != nil {
		return nil, err
	}

	rows, err = sr.db.Query(`
		SELECT person
		FROM dream_people
		WHERE dream_id = ?
	`, dreamID)

	da.People, err = scanStrings(rows)
	if err != nil {
		return nil, err
	}

	rows, err = sr.db.Query(`
		SELECT symbol
		FROM dream_symbols
		WHERE dream_id = ?
	`, dreamID)

	da.Symbols, err = scanStrings(rows)
	if err != nil {
		return nil, err
	}

	return &da, nil
}

func (sr *SQLRepository) DeleteDreamAnalysis(dreamID string) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM dream_analyses WHERE dream_id = ?`, dreamID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM dream_themes WHERE dream_id = ?`, dreamID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM dream_emotions WHERE dream_id = ?`, dreamID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM dream_locations WHERE dream_id = ?`, dreamID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM dream_people WHERE dream_id = ?`, dreamID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM dream_symbols WHERE dream_id = ?`, dreamID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func scanStrings(rows *sql.Rows) ([]string, error) {
	defer rows.Close()
	var result []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			rows.Close()
			return nil, err
		}
		result = append(result, value)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))

	for _, word := range values {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		if _, exists := seen[word]; exists == true {
			continue
		}

		seen[word] = true
		result = append(result, word)
	}
	return result
}
