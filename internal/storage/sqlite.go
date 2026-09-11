package storage

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

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

func (sr *SQLRepository) SaveDreamEmbedding(de dream.DreamEmbedding) error {
	if err := sr.DeleteDreamEmbedding(de.DreamID); err != nil {
		return err
	}

	_, err := sr.db.Exec(`
		INSERT INTO dream_embeddings (
			dream_id, model, dimensions, embedding
		) VALUES (?, ?, ?, ?)
		 ON CONFLICT(dream_id) DO UPDATE SET 
		 	model = excluded.model,
			dimensions = excluded.dimensions,
			embedding = excluded.embedding
	`, de.DreamID, de.Model, len(de.Embedding), floatsToBytes(de.Embedding))

	return err
}

func (sr *SQLRepository) DeleteDreamEmbedding(dreamID string) error {
	_, err := sr.db.Exec(`
		DELETE FROM dream_embeddings
		WHERE dream_id = ?
	`, dreamID)

	return err
}

func (sr *SQLRepository) GetDreamEmbedding(dreamID string) (*dream.DreamEmbedding, error) {
	var de dream.DreamEmbedding
	var embeddingBytes []byte
	err := sr.db.QueryRow(`
		SELECT model, dimensions, embedding
		FROM dream_embeddings
		WHERE dream_id = ?
	`, dreamID).Scan(&de.Model, &de.Dimensions, &embeddingBytes)

	if err != nil {
		return nil, err
	}

	embedding, err := bytesToFloats(embeddingBytes)
	if err != nil {
		return nil, err
	}
	de.DreamID = dreamID
	de.Embedding = embedding

	return &de, nil
}

func floatsToBytes(values []float32) []byte {
	buf := make([]byte, len(values)*4)

	for i, value := range values {
		bits := math.Float32bits(value)
		binary.LittleEndian.PutUint32(buf[i*4:], bits)
	}

	return buf
}

func bytesToFloats(data []byte) ([]float32, error) {
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("invalid embedding byte length")
	}

	values := make([]float32, len(data)/4)

	for i := range values {
		bits := binary.LittleEndian.Uint32(data[i*4:])
		values[i] = math.Float32frombits(bits)
	}

	return values, nil
}

func (sr *SQLRepository) ListDreamEmbeddings() ([]dream.DreamEmbedding, error) {
	rows, err := sr.db.Query(`
	SELECT dream_id, model, dimensions, embedding
	FROM dream_embeddings
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []dream.DreamEmbedding

	for rows.Next() {
		var de dream.DreamEmbedding
		var embeddingBytes []byte
		err := rows.Scan(&de.DreamID, &de.Model, &de.Dimensions, &embeddingBytes)
		if err != nil {
			return nil, err
		}
		de.Embedding, err = bytesToFloats(embeddingBytes)
		if err != nil {
			return nil, err
		}

		result = append(result, de)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
