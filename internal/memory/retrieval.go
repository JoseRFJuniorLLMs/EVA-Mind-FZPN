package memory

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// RetrievalService busca memórias por similaridade semântica
type RetrievalService struct {
	db       *sql.DB
	embedder *EmbeddingService
}

// NewRetrievalService cria um novo serviço de busca
func NewRetrievalService(db *sql.DB, embedder *EmbeddingService) *RetrievalService {
	return &RetrievalService{
		db:       db,
		embedder: embedder,
	}
}

// SearchResult representa um resultado de busca com score de similaridade
type SearchResult struct {
	Memory     *Memory
	Similarity float64 // 0.0 (nada similar) a 1.0 (idêntico)
}

// Retrieve busca as K memórias mais relevantes para uma query
func (r *RetrievalService) Retrieve(ctx context.Context, idosoID int64, query string, k int) ([]*SearchResult, error) {
	// 1. Gerar embedding da query
	queryEmbedding, err := r.embedder.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar embedding: %w", err)
	}

	// 2. Buscar usando função SQL search_similar_memories
	sqlQuery := `
		SELECT * FROM search_similar_memories(
			$1,  -- idoso_id
			$2,  -- query_embedding
			$3,  -- limit
			$4   -- min_similarity
		)
	`

	rows, err := r.db.QueryContext(
		ctx,
		sqlQuery,
		idosoID,
		vectorToPostgres(queryEmbedding),
		k,
		0.5, // Similaridade mínima de 50%
	)
	if err != nil {
		return nil, fmt.Errorf("erro na busca SQL: %w", err)
	}
	defer rows.Close()

	var results []*SearchResult

	for rows.Next() {
		var (
			memoryID        int64
			content         string
			speaker         string
			memoryTimestamp string
			emotion         sql.NullString
			importance      float64
			topics          string
			similarity      float64
		)

		err := rows.Scan(
			&memoryID,
			&content,
			&speaker,
			&memoryTimestamp,
			&emotion,
			&importance,
			&topics,
			&similarity,
		)
		if err != nil {
			log.Printf("⚠️ Erro ao scanear resultado: %v", err)
			continue
		}

		memory := &Memory{
			ID:         memoryID,
			IdosoID:    idosoID,
			Speaker:    speaker,
			Content:    content,
			Importance: importance,
			Topics:     parsePostgresArray(topics),
		}

		if emotion.Valid {
			memory.Emotion = emotion.String
		}

		results = append(results, &SearchResult{
			Memory:     memory,
			Similarity: similarity,
		})
	}

	return results, rows.Err()
}

// RetrieveRecent busca memórias recentes (últimos N dias) sem usar embedding
// Útil para contexto temporal imediato
func (r *RetrievalService) RetrieveRecent(ctx context.Context, idosoID int64, days int, limit int) ([]*Memory, error) {
	query := `
		SELECT id, idoso_id, timestamp, speaker, content, emotion, 
		       importance, topics, session_id
		FROM episodic_memories
		WHERE idoso_id = $1
		  AND timestamp > NOW() - INTERVAL '1 day' * $2
		ORDER BY importance DESC, timestamp DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, idosoID, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// RetrieveHybrid combina busca semântica + temporal
// Retorna memórias relevantes E recentes
func (r *RetrievalService) RetrieveHybrid(ctx context.Context, idosoID int64, query string, k int) ([]*SearchResult, error) {
	// Buscar memórias semânticas
	semantic, err := r.Retrieve(ctx, idosoID, query, k)
	if err != nil {
		return nil, err
	}

	// Buscar memórias recentes (últimos 3 dias)
	recent, err := r.RetrieveRecent(ctx, idosoID, 3, k/2)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar memórias recentes: %v", err)
		return semantic, nil // Retorna apenas semânticas
	}

	// Mesclar e deduplicar
	seen := make(map[int64]bool)
	var combined []*SearchResult

	// Adicionar semânticas primeiro
	for _, res := range semantic {
		if !seen[res.Memory.ID] {
			combined = append(combined, res)
			seen[res.Memory.ID] = true
		}
	}

	// Adicionar recentes (se não duplicadas)
	for _, mem := range recent {
		if !seen[mem.ID] {
			combined = append(combined, &SearchResult{
				Memory:     mem,
				Similarity: 0.9, // Score artificial alto para recentes
			})
			seen[mem.ID] = true

			if len(combined) >= k {
				break
			}
		}
	}

	return combined, nil
}
