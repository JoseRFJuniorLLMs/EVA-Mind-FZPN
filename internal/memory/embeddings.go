package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	geminiEmbeddingEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/text-embedding-004:embedContent"
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

	return result.Embedding.Values, nil
}

// GenerateBatch gera embeddings para múltiplos textos (otimização futura)
func (e *EmbeddingService) GenerateBatch(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		emb, err := e.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("erro no texto %d: %w", i, err)
		}
		embeddings[i] = emb
	}

	return embeddings, nil
}
