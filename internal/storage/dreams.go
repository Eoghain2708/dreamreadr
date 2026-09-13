package storage

import "github.com/Eoghain2708/dreamreadr/internal/dream"

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
