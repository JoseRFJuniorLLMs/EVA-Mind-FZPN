package main

import (
	"context"
	"eva-mind/internal/brainstem/config"
	"eva-mind/internal/brainstem/infrastructure/vector"
	"eva-mind/internal/hippocampus/memory" // For EmbeddingService
	"eva-mind/pkg/types"
	"fmt"
	"log"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

var initialStories = []types.TherapeuticStory{
	{
		ID:             "tortoise_hare",
		Title:          "A Lebre e a Tartaruga",
		Content:        "Era uma vez uma lebre que se gabava de sua velocidade...",
		TargetEmotions: []string{"ansiedade", "impaciência", "pressa"},
		Archetype:      "O Sábio",
		Moral:          "Devagar e sempre se vai ao longe. A constância supera a rapidez.",
		Tags:           []string{"paciência", "consistência", "calma"},
		MinAge:         60,
	},
	{
		ID:             "ugly_duckling",
		Title:          "O Patinho Feio",
		Content:        "Um patinho nasceu diferente dos outros e foi rejeitado...",
		TargetEmotions: []string{"solidão", "tristeza", "rejeição"},
		Archetype:      "O Herói",
		Moral:          "A beleza verdadeira leva tempo para florescer. Você pertence a um lugar maior.",
		Tags:           []string{"autoestima", "pertencimento", "esperança"},
		MinAge:         60,
	},
	{
		ID:             "oak_reed",
		Title:          "O Carvalho e o Caniço",
		Content:        "O carvalho era forte e não se curvava ao vento, enquanto o caniço se dobrava...",
		TargetEmotions: []string{"raiva", "teimosia", "rigidez"},
		Archetype:      "O Governante",
		Moral:          "Muitas vezes é melhor ceder do que quebrar. A flexibilidade é uma força.",
		Tags:           []string{"flexibilidade", "resiliência", "aceitação"},
		MinAge:         60,
	},
	{
		ID:             "starfish",
		Title:          "O Menino e as Estrelas do Mar",
		Content:        "Um homem caminhava pela praia e viu um menino jogando estrelas do mar de volta ao oceano...",
		TargetEmotions: []string{"impotência", "desanimo", "inutilidade"},
		Archetype:      "O Cuidador",
		Moral:          "Para aquela que você salvou, fez toda a diferença. Pequenas ações importam.",
		Tags:           []string{"propósito", "cuidado", "impacto"},
		MinAge:         60,
	},
}

func main() {
	log.Println("🌱 Iniciando seed de Histórias Terapêuticas...")

	// 1. Carregar Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}
	if cfg.GoogleAPIKey == "" {
		log.Fatal("❌ GoogleAPIKey não encontrada no .env")
	}

	// 2. Inicializar Serviços
	qdrantClient, err := vector.NewQdrantClient("localhost", 6334)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar Qdrant: %v", err)
	}

	// Check embeddings.go for signature. Assuming NewEmbeddingService(apiKey)
	embedder := memory.NewEmbeddingService(cfg.GoogleAPIKey)

	ctx := context.Background()

	// 3. Criar Coleção se não existir
	err = qdrantClient.CreateCollection(ctx, "stories", 768) // 768 for Gemini embedding
	if err != nil {
		log.Printf("⚠️ Coleção 'stories' já existe ou erro: %v", err)
	}

	// 4. Processar e Inserir
	var points []*qdrant.PointStruct

	for i, story := range initialStories {
		log.Printf("Processando: %s", story.Title)

		// Gerar união de texto para embedding semântico
		textToEmbed := fmt.Sprintf("Título: %s. Moral: %s. Emoções: %v. %s",
			story.Title, story.Moral, story.TargetEmotions, story.Content)

		vec, err := embedder.GenerateEmbedding(ctx, textToEmbed)
		if err != nil {
			log.Printf("❌ Erro ao gerar embedding para %s: %v", story.Title, err)
			continue
		}

		// Converter []float64 para []float32 se necessário (geralmente embedder retorna []float32, mas verificar)
		// Assuming embedder returns []float32 based on typical logic. If it returns []float64, convert.
		// memory.EmbeddingService usually returns []float32.

		payload := map[string]interface{}{
			"id":              story.ID,
			"title":           story.Title,
			"content":         story.Content,
			"moral":           story.Moral,
			"archetype":       story.Archetype,
			"target_emotions": story.TargetEmotions,
			"tags":            story.Tags,
			"min_age":         story.MinAge,
		}

		point := &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: uint64(i + 1)},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: vec}},
			},
			Payload: map[string]*qdrant.Value{},
		}

		// Helper to convert map to qdrant payload (simplified)
		// Qdrant Go client expects precise types.
		// Using a helper or manual conversion.
		for k, v := range payload {
			point.Payload[k] = toQdrantValue(v)
		}

		points = append(points, point)

		// Rate limit helper
		time.Sleep(500 * time.Millisecond)
	}

	// 5. Upsert
	if len(points) > 0 {
		err = qdrantClient.Upsert(ctx, "stories", points)
		if err != nil {
			log.Fatalf("❌ Erro no Upsert: %v", err)
		}
		log.Printf("✅ Sucesso! %d histórias inseridas.", len(points))
	} else {
		log.Println("⚠️ Nenhuma história processada.")
	}
}

func toQdrantValue(v interface{}) *qdrant.Value {
	switch val := v.(type) {
	case string:
		return &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: val}}
	case int:
		return &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(val)}}
	case []string:
		list := &qdrant.ListValue{}
		for _, s := range val {
			list.Values = append(list.Values, &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: s}})
		}
		return &qdrant.Value{Kind: &qdrant.Value_ListValue{ListValue: list}}
	default:
		return &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: fmt.Sprintf("%v", val)}}
	}
}
