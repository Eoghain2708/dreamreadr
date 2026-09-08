package storage

import (
	"database/sql"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (sr *SQLRepository) CreateDream(d dream.Dream) error {
	_, err := sr.db.Exec(`
	INSERT INTO dreams (
		id, title, created_at, raw_text
	)
		VALUES (?, ?, ?, ?)
	`,
		d.ID, d.Title, d.CreatedAt, d.RawText,
	)
	return err
}

func (sr *SQLRepository) GetDream(id string) (dream.Dream, error) {
	var d dream.Dream
	err := sr.db.QueryRow(`
	SELECT 
		id, title, created_at, raw_text
	FROM dreams 
	WHERE id = ?
	`, id).Scan(
		&d.ID,
		&d.Title,
		&d.CreatedAt,
		&d.RawText,
	)

	if err != nil {
		return dream.Dream{}, err
	}

	return d, nil
}

func (sr *SQLRepository) ListDreams() ([]dream.Dream, error) {
	rows, err := sr.db.Query(`
		SELECT 
		id, title, created_at, raw_text
		FROM dreams
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dreams []dream.Dream

	for rows.Next() {
		var d dream.Dream
		err := rows.Scan(
			&d.ID, &d.Title, &d.CreatedAt, &d.RawText,
		)

		if err != nil {
			return nil, err
		}

		dreams = append(dreams, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dreams, nil
}

func (sr *SQLRepository) DeleteDream(id string) error {
	_, err := sr.db.Exec(`
	DELETE FROM dreams where id = ?
	`, id)

	return err
}

func (sr *SQLRepository) CreateDreamAnalysis(da dream.DreamAnalysis) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
        INSERT INTO dream_analyses (
            dream_id, summary
        ) VALUES (?, ?)
    `, da.DreamID, da.Summary)

	if err != nil {
		return err
	}

	for _, theme := range da.Themes {
		_, err = tx.Exec(`
            INSERT INTO dream_themes (
                dream_id, theme
            ) VALUES (?, ?)
        `, da.DreamID, theme)

		if err != nil {
			return err
		}
	}

	for _, emotion := range da.Emotions {
		_, err = tx.Exec(`
            INSERT INTO dream_emotions (
                dream_id, emotion
            ) VALUES (?, ?)
        `, da.DreamID, emotion)

		if err != nil {
			return err
		}
	}

	for _, location := range da.Locations {
		_, err = tx.Exec(`
            INSERT INTO dream_locations (
                dream_id, location
            ) VALUES (?, ?)
        `, da.DreamID, location)

		if err != nil {
			return err
		}
	}

	for _, person := range da.People {
		_, err = tx.Exec(`
			INSERT INTO dream_people (
			dream_id, person
			) VALUES (?, ?)
		`, da.DreamID, person)

		if err != nil {
			return err
		}
	}

	for _, symbol := range da.Symbols {
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
