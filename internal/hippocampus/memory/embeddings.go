package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// gemini-embedding-001 é o novo modelo recomendado (substitui text-embedding-004)
	// Suporta 100+ idiomas, 2048 tokens, dimensões: 3072, 1536 ou 768
	geminiEmbeddingEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent"
	expectedDimension       = 3072 // Máxima qualidade - Qdrant recriado com 3072
)

// EmbeddingService gera embeddings usando Gemini API
type EmbeddingService struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewEmbeddingService cria um novo serviço de embeddings
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// EmbeddingRequest representa o payload da API
type embeddingRequest struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
	OutputDimensionality int `json:"outputDimensionality,omitempty"` // Para gemini-embedding-001: 768, 1536 ou 3072
}

// EmbeddingResponse representa a resposta da API
type embeddingResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

// GenerateEmbedding gera um vetor de embedding para o texto fornecido
func (e *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Truncar texto se muito longo (limite da API: ~2048 tokens)
	if len(text) > 8000 {
		text = text[:8000]
	}

	// Construir request
	reqBody := embeddingRequest{}
	reqBody.Content.Parts = []struct {
		Text string `json:"text"`
	}{
		{Text: text},
	}
	reqBody.OutputDimensionality = expectedDimension // 3072 dimensões - máxima qualidade semântica

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar request: %w", err)
	}

	// Fazer request HTTP
	url := fmt.Sprintf("%s?key=%s", geminiEmbeddingEndpoint, e.APIKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	if len(result.Embedding.Values) == 0 {
		return nil, fmt.Errorf("embedding vazio retornado pela API")
	}

	// ✅ VALIDAÇÃO CRÍTICA DE DIMENSÃO
	actualDim := len(result.Embedding.Values)
	if actualDim != expectedDimension {
		return nil, fmt.Errorf(
			"❌ DIMENSION MISMATCH DETECTED!\n"+
				"   Expected: %d (Postgres schema)\n"+
				"   Got: %d (Gemini API)\n"+
				"   This will cause ALL searches to fail!\n"+
				"   Run migration: migrations/004_fix_embedding_dimension.sql",
			expectedDimension,
			actualDim,
		)
	}

	log.Printf("✅ [EMBEDDING] Generated %d dimensions (validated)", actualDim)
	return result.Embedding.Values, nil
}

// GenerateBatch gera embeddings para múltiplos textos
// Útil para re-embedding em massa
func (e *EmbeddingService) GenerateBatch(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		emb, err := e.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("erro no texto %d: %w", i, err)
		}
		embeddings[i] = emb

		// Rate limiting: 10 requests/second
		if i < len(texts)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return embeddings, nil
}
