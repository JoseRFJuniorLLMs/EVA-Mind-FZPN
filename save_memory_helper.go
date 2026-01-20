// save_memory_helper.go
package main

import (
	"context"
	"eva-mind/internal/memory"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

// saveAsMemory salva uma transcrição como memória episódica
func (s *SignalingServer) saveAsMemory(idosoID int64, role, text string) {
	// Ignorar textos muito curtos
	if len(text) < 10 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Gerar embedding
	embedding, err := s.embeddingService.GenerateEmbedding(ctx, text)
	if err != nil {
		log.Printf("❌ [MEMORY] Erro ao gerar embedding: %v", err)
		return
	}

	// 2. Analisar metadados (emoção, importância, tópicos)
	// metadataAnalyzer está em s.metadataAnalyzer no main.go
	metadata, err := s.metadataAnalyzer.Analyze(ctx, text)
	if err != nil {
		log.Printf("⚠️ [MEMORY] Erro na análise (usando padrão): %v", err)
		metadata = &memory.Metadata{
			Emotion:    "neutro",
			Importance: 0.5,
			Topics:     []string{"geral"},
		}
	}

	// 3. Criar objeto Memory
	mem := &memory.Memory{
		IdosoID:    idosoID,
		Speaker:    role,
		Content:    text,
		Embedding:  embedding,
		Emotion:    metadata.Emotion,
		Importance: metadata.Importance,
		Topics:     metadata.Topics,
		Timestamp:  time.Now(),
	}

	// 4. Salvar no POSTGRES (CRÍTICO - deve bloquear)
	err = s.memoryStore.Store(ctx, mem)
	if err != nil {
		log.Printf("❌ [MEMORY] Erro ao salvar no Postgres: %v", err)
		return
	}

	log.Printf("✅ [POSTGRES] Memory saved: ID=%d, Speaker=%s", mem.ID, role)

	// ✅ 5. NOVO: UPSERT NO QDRANT (Assíncrono)
	if s.qdrantClient != nil {
		go func() {
			qctx, qcancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer qcancel()

			// Criar point para Qdrant
			point := &qdrant.PointStruct{
				Id: &qdrant.PointId{
					PointIdOptions: &qdrant.PointId_Num{Num: uint64(mem.ID)},
				},
				Vectors: &qdrant.Vectors{
					VectorsOptions: &qdrant.Vectors_Vector{
						Vector: &qdrant.Vector{Data: mem.Embedding},
					},
				},
				Payload: map[string]*qdrant.Value{
					"content": {
						Kind: &qdrant.Value_StringValue{StringValue: mem.Content},
					},
					"speaker": {
						Kind: &qdrant.Value_StringValue{StringValue: mem.Speaker},
					},
					"idoso_id": {
						Kind: &qdrant.Value_IntegerValue{IntegerValue: mem.IdosoID},
					},
					"timestamp": {
						Kind: &qdrant.Value_StringValue{
							StringValue: mem.Timestamp.Format(time.RFC3339),
						},
					},
					"emotion": {
						Kind: &qdrant.Value_StringValue{StringValue: mem.Emotion},
					},
					"importance": {
						Kind: &qdrant.Value_DoubleValue{DoubleValue: mem.Importance},
					},
					"topics": {
						Kind: &qdrant.Value_ListValue{
							ListValue: stringSliceToQdrantList(mem.Topics),
						},
					},
				},
			}

			// Upsert com retry (3 tentativas)
			for attempt := 1; attempt <= 3; attempt++ {
				err := s.qdrantClient.Upsert(qctx, "memories", []*qdrant.PointStruct{point})
				if err == nil {
					log.Printf("✅ [QDRANT] Memory %d indexed", mem.ID)
					break
				}

				if attempt < 3 {
					log.Printf("⚠️ [QDRANT] Attempt %d/3 failed: %v (retrying...)", attempt, err)
					time.Sleep(time.Second * time.Duration(attempt))
				} else {
					log.Printf("❌ [QDRANT] Failed after 3 attempts: %v", err)
				}
			}
		}()
	}

	// 6. Salvar no NEO4J (Grafo Causal) - Assíncrono
	if s.graphStore != nil {
		go func() {
			nctx, ncancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer ncancel()

			err := s.graphStore.StoreCausalMemory(nctx, mem)
			if err != nil {
				log.Printf("❌ [NEO4J] Erro ao salvar nó: %v", err)
			} else {
				log.Printf("✅ [NEO4J] Memory %d graphed", mem.ID)
			}

			// Rastrear significantes (Lacan)
			if s.signifierService != nil {
				s.signifierService.TrackSignifiers(nctx, mem.IdosoID, mem.Content)
			}
		}()
	}

	// 7. Atualizar estado de personalidade
	if s.personalityService != nil && role == "user" {
		go func() {
			pctx, pcancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer pcancel()

			err := s.personalityService.UpdateAfterConversation(
				pctx, idosoID, metadata.Emotion, metadata.Topics)
			if err != nil {
				log.Printf("⚠️ [PERSONALITY] Erro ao atualizar estado: %v", err)
			}
		}()
	}
}

// ✅ FIX 3: FDPN Hook Implementation
// handleUserTranscription processa transcrições do usuário em tempo real
func (s *SignalingServer) handleUserTranscription(client *PCMClient, text string) {
	if len(text) < 10 {
		return // Ignorar textos muito curtos
	}

	ctx := context.Background()
	userID := fmt.Sprintf("%d", client.IdosoID)

	// 🚀 ATIVAR FDPN SPREADING ACTIVATION
	go func() {
		start := time.Now()

		if s.fdpnEngine != nil {
			if err := s.fdpnEngine.StreamingPrime(ctx, userID, text); err != nil {
				log.Printf("⚠️ [FDPN] Prime error: %v", err)
			} else {
				elapsed := time.Since(start)
				log.Printf("✅ [FDPN] Primed in %dms (user=%s)", elapsed.Milliseconds(), userID)

				// Debug: mostrar keywords extraídas
				keywords := extractKeywords(text)
				if len(keywords) > 0 {
					log.Printf("🔍 [FDPN] Keywords: %v", keywords)
				}
			}
		}
	}()

	// Salvar memória (já existente)
	go s.saveAsMemory(client.IdosoID, "user", text)
}

// Helper para extrair keywords (rápido e simples)
func extractKeywords(text string) []string {
	stopwords := map[string]bool{
		"o": true, "a": true, "de": true, "que": true, "e": true,
		"do": true, "da": true, "em": true, "um": true, "para": true,
		"com": true, "não": true, "uma": true, "os": true, "no": true,
		"se": true, "na": true, "por": true, "mais": true, "as": true,
	}

	words := strings.Fields(strings.ToLower(text))
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
