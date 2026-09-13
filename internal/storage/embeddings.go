package storage

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

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
