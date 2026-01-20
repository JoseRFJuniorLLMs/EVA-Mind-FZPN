package brain

import (
	"context"
	"eva-mind/pkg/types"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pgvector/pgvector-go"
	"github.com/qdrant/go-client/qdrant"
)

// ProcessUserSpeech handles user transcription in real-time (FDPN Hook)
func (s *Service) ProcessUserSpeech(ctx context.Context, idosoID int64, text string) {
	if len(text) < 10 {
		return // Ignore short texts
	}

	userID := fmt.Sprintf("%d", idosoID)

	// Output log to track flow
	log.Printf("🗣️ [User Speech] Processing for user %s: %s", userID, text)

	// 🚀 ACTIVATE UNIFIED RETRIEVAL PRIMING (RSI + FDPN)
	if s.unifiedRetrieval != nil {
		go s.unifiedRetrieval.Prime(ctx, idosoID, text)
	}

	// Save memory (Fire and forget)
	go s.SaveEpisodicMemory(idosoID, "user", text)
}

// SaveEpisodicMemory saves memory to Postgres and Qdrant
func (s *Service) SaveEpisodicMemory(idosoID int64, role, content string) {
	ctx := context.Background()

	// 1. Generate Embedding
	embedding, err := s.embeddingService.GenerateEmbedding(ctx, content)
	if err != nil {
		log.Printf("⚠️ [MEMORY] Erro ao gerar embedding: %v", err)
		return
	}

	// 2. Save to Postgres
	query := `
		INSERT INTO episodic_memories (idoso_id, speaker, content, embedding, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`
	var memoryID int64
	err = s.db.QueryRow(query, idosoID, role, content, pgvector.NewVector(embedding)).Scan(&memoryID)
	if err != nil {
		log.Printf("❌ [POSTGRES] Erro ao salvar memória: %v", err)
		return
	}
	log.Printf("✅ [POSTGRES] Memory saved: ID=%d, Speaker=%s", memoryID, role)

	// 3. Upsert to Qdrant (Retry Logic)
	if s.qdrantClient != nil {
		go func() {
			metadata := types.MemoryMetadata{
				Emotion:    "neutral",
				Importance: 0.5,
				Topics:     extractKeywords(content),
			}

			// Tentar 3 vezes
			for attempt := 1; attempt <= 3; attempt++ {
				points := []*qdrant.PointStruct{
					{
						Id: &qdrant.PointId{
							PointIdOptions: &qdrant.PointId_Num{Num: uint64(memoryID)},
						},
						Vectors: &qdrant.Vectors{
							VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: embedding}},
						},
						Payload: map[string]*qdrant.Value{
							"idoso_id":   {Kind: &qdrant.Value_IntegerValue{IntegerValue: idosoID}},
							"role":       {Kind: &qdrant.Value_StringValue{StringValue: role}},
							"content":    {Kind: &qdrant.Value_StringValue{StringValue: content}},
							"created_at": {Kind: &qdrant.Value_StringValue{StringValue: time.Now().Format(time.RFC3339)}},
							"emotion":    {Kind: &qdrant.Value_StringValue{StringValue: metadata.Emotion}},
							"topics":     stringSliceToQdrantList(metadata.Topics),
						},
					},
				}

				if err := s.qdrantClient.Upsert(ctx, "memories", points); err != nil {
					log.Printf("⚠️ [QDRANT] Upsert falhou (tentativa %d): %v", attempt, err)
					time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
					continue
				}

				log.Printf("✅ [QDRANT] Memory %d indexed", memoryID)

				// 4. Update Personality State (Async)
				if role == "user" && s.personalityService != nil {
					go func() {
						pctx, pcancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer pcancel()
						s.personalityService.UpdateAfterConversation(pctx, idosoID, metadata.Emotion, metadata.Topics)
					}()
				}
				break
			}
		}()
	}
}

// Helper to extract keywords
func extractKeywords(text string) []string {
	stopwords := map[string]bool{
		"o": true, "a": true, "de": true, "que": true, "e": true,
		"do": true, "da": true, "em": true, "um": true, "para": true,
		"com": true, "não": true, "uma": true, "os": true, "no": true,
		"se": true, "na": true, "por": true, "mais": true, "as": true,
	}

	var keywords []string
	seen := make(map[string]bool)

	for _, w := range strings.Fields(strings.ToLower(text)) {
		w = strings.Trim(w, ".,!?;:'\"")
		if len(w) < 3 || stopwords[w] || seen[w] {
			continue
		}
		keywords = append(keywords, w)
		seen[w] = true
	}

	return keywords
}

func stringSliceToQdrantList(slice []string) *qdrant.Value {
	values := make([]*qdrant.Value, len(slice))
	for i, s := range slice {
		values[i] = &qdrant.Value{
			Kind: &qdrant.Value_StringValue{StringValue: s},
		}
	}
	return &qdrant.Value{
		Kind: &qdrant.Value_ListValue{
			ListValue: &qdrant.ListValue{Values: values},
		},
	}
}
