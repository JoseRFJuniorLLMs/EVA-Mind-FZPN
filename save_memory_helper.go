package main

import (
	"context"
	"eva-mind/internal/memory"
	"log"
	"time"
)

// saveAsMemory salva uma transcrição como memória episódica (async)
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
	metadata, err := s.metadataAnalyzer.Analyze(ctx, text)
	if err != nil {
		log.Printf("⚠️ [MEMORY] Erro na análise (usando padrão): %v", err)
		// Fallback manual se necessário, mas Analyze já deve tratar isso
		metadata = &memory.Metadata{
			Emotion:    "neutro",
			Importance: 0.5,
			Topics:     []string{"geral"},
		}
	}

	// 3. Salvar no banco
	mem := &memory.Memory{
		IdosoID:    idosoID,
		Speaker:    role,
		Content:    text,
		Embedding:  embedding,
		Emotion:    metadata.Emotion,
		Importance: metadata.Importance,
		Topics:     metadata.Topics,
	}

	err = s.memoryStore.Store(ctx, mem)
	if err != nil {
		log.Printf("❌ [MEMORY] Erro ao salvar: %v", err)
		return
	}

	// log.Printf("🧠 [MEMORY] Salva: [%s] %s (importância: %.2f)", role, text[:50], metadata.Importance)
}
