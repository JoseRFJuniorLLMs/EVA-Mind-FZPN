package memory

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Memory representa uma memória episódica armazenada
type Memory struct {
	ID            int64     `json:"id"`
	IdosoID       int64     `json:"idoso_id"`
	Timestamp     time.Time `json:"timestamp"`
	Speaker       string    `json:"speaker"` // "user" ou "assistant"
	Content       string    `json:"content"`
	Embedding     []float32 `json:"-"` // Não serializar embedding (muito grande)
	Emotion       string    `json:"emotion"`
	Importance    float64   `json:"importance"`
	Topics        []string  `json:"topics"`
	SessionID     string    `json:"session_id,omitempty"`
	CallHistoryID *int64    `json:"call_history_id,omitempty"`
}

// MemoryStore gerencia o armazenamento de memórias
type MemoryStore struct {
	db *sql.DB
}

// NewMemoryStore cria um novo gerenciador de memórias
func NewMemoryStore(db *sql.DB) *MemoryStore {
	return &MemoryStore{db: db}
}

// Store salva uma nova memória no banco
func (m *MemoryStore) Store(ctx context.Context, memory *Memory) error {
	query := `
		INSERT INTO episodic_memories 
		(idoso_id, speaker, content, embedding, emotion, importance, topics, session_id, call_history_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, timestamp
	`

	embeddingStr := vectorToPostgres(memory.Embedding)

	err := m.db.QueryRowContext(
		ctx,
		query,
		memory.IdosoID,
		memory.Speaker,
		memory.Content,
		embeddingStr,
		memory.Emotion,
		memory.Importance,
		pqArray(memory.Topics),
		memory.SessionID,
		memory.CallHistoryID,
	).Scan(&memory.ID, &memory.Timestamp)

	return err
}

// GetByID recupera uma memória por ID
func (m *MemoryStore) GetByID(ctx context.Context, id int64) (*Memory, error) {
	query := `
		SELECT id, idoso_id, timestamp, speaker, content, emotion, 
		       importance, topics, session_id, call_history_id
		FROM episodic_memories
		WHERE id = $1
	`

	memory := &Memory{}
	var topics string

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&memory.ID,
		&memory.IdosoID,
		&memory.Timestamp,
		&memory.Speaker,
		&memory.Content,
		&memory.Emotion,
		&memory.Importance,
		&topics,
		&memory.SessionID,
		&memory.CallHistoryID,
	)

	if err != nil {
		return nil, err
	}

	// Parse topics array
	memory.Topics = parsePostgresArray(topics)

	return memory, nil
}

// GetRecent retorna as N memórias mais recentes de um idoso
func (m *MemoryStore) GetRecent(ctx context.Context, idosoID int64, limit int) ([]*Memory, error) {
	query := `
		SELECT id, idoso_id, timestamp, speaker, content, emotion, 
		       importance, topics, session_id
		FROM episodic_memories
		WHERE idoso_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := m.db.QueryContext(ctx, query, idosoID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return m.scanMemories(rows)
}

// DeleteOld remove memórias mais antigas que X dias (para LGPD/GDPR)
func (m *MemoryStore) DeleteOld(ctx context.Context, idosoID int64, olderThanDays int) (int64, error) {
	query := `
		DELETE FROM episodic_memories
		WHERE idoso_id = $1
		  AND timestamp < NOW() - INTERVAL '1 day' * $2
		  AND importance < 0.7  -- Preservar memórias importantes
	`

	result, err := m.db.ExecContext(ctx, query, idosoID, olderThanDays)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// scanMemories helper para converter rows em slice de Memory
func (m *MemoryStore) scanMemories(rows *sql.Rows) ([]*Memory, error) {
	var memories []*Memory

	for rows.Next() {
		memory := &Memory{}
		var topics string

		err := rows.Scan(
			&memory.ID,
			&memory.IdosoID,
			&memory.Timestamp,
			&memory.Speaker,
			&memory.Content,
			&memory.Emotion,
			&memory.Importance,
			&topics,
			&memory.SessionID,
		)

		if err != nil {
			return nil, err
		}

		memory.Topics = parsePostgresArray(topics)
		memories = append(memories, memory)
	}

	return memories, rows.Err()
}

// Helpers para conversão de tipos PostgreSQL

func vectorToPostgres(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}

	result := "["
	for i, v := range vec {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%f", v)
	}
	result += "]"

	return result
}

func pqArray(arr []string) string {
	if len(arr) == 0 {
		return "{}"
	}

	result := "{"
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("\"%s\"", s)
	}
	result += "}"

	return result
}

func parsePostgresArray(s string) []string {
	if s == "{}" || s == "" {
		return []string{}
	}

	// Remove {} e split por vírgula
	s = s[1 : len(s)-1]
	var result []string

	// Parse manual para lidar com aspas
	var current string
	inQuotes := false

	for _, c := range s {
		switch c {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes {
				if current != "" {
					result = append(result, current)
					current = ""
				}
			} else {
				current += string(c)
			}
		default:
			current += string(c)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}
