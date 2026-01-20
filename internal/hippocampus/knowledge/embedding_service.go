package knowledge

import (
	"context"
	"encoding/json"
	"eva-mind/internal/brainstem/config"
	"eva-mind/internal/brainstem/infrastructure/vector" // Import wrapper
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

// EmbeddingService rastreia cadeias de significantes via embeddings semânticos
// Implementa "L'instance de la lettre" - a letra que circula e se repete
type EmbeddingService struct {
	cfg            *config.Config
	qdrantClient   *vector.QdrantClient // Use wrapper
	httpClient     *http.Client
	collectionName string
}

// SignifierChain representa uma cadeia de significantes detectada
type SignifierChain struct {
	CoreSignifier   string    `json:"core_signifier"`   // Significante nuclear
	RelatedWords    []string  `json:"related_words"`    // Palavras semanticamente próximas
	EmotionalCharge float64   `json:"emotional_charge"` // Carga afetiva (0.0-1.0)
	Frequency       int       `json:"frequency"`        // Vezes que apareceu
	LastOccurrence  time.Time `json:"last_occurrence"`
	Contexts        []string  `json:"contexts"` // Frases onde apareceu
}

// NewEmbeddingService cria serviço de embeddings
func NewEmbeddingService(cfg *config.Config, qdrantClient *vector.QdrantClient) (*EmbeddingService, error) {
	svc := &EmbeddingService{
		cfg:            cfg,
		qdrantClient:   qdrantClient,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		collectionName: "signifier_chains",
	}

	// Criar coleção se não existir
	// Check qdrantClient is not nil to avoid panic during testing if passed nil
	if qdrantClient != nil {
		if err := svc.ensureCollection(context.Background()); err != nil {
			log.Printf("⚠️ Warning: Could not create Qdrant collection: %v", err)
		}
	}

	return svc, nil
}

// ensureCollection garante que a coleção existe
func (e *EmbeddingService) ensureCollection(ctx context.Context) error {
	// wrapper provides GetCollectionInfo
	_, err := e.qdrantClient.GetCollectionInfo(ctx, e.collectionName)
	if err == nil {
		return nil // Já existe
	}

	// Criar coleção (Gemini embeddings são 768 dimensões)
	// wrapper provides CreateCollection
	err = e.qdrantClient.CreateCollection(ctx, e.collectionName, 768)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	log.Printf("✅ Created Qdrant collection: %s", e.collectionName)
	return nil
}

// GenerateEmbedding gera embedding usando Gemini API
func (e *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/text-embedding-004:embedContent?key=%s",
		e.cfg.GoogleAPIKey,
	)

	payload := map[string]interface{}{
		"model": "models/text-embedding-004",
		"content": map[string]interface{}{
			"parts": []map[string]string{
				{"text": text},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embedding.Values, nil
}

// TrackSignifierChain rastreia cadeia de significantes
func (e *EmbeddingService) TrackSignifierChain(ctx context.Context, idosoID int64, text string, emotionalCharge float64) error {
	if e.qdrantClient == nil {
		return nil
	}

	// Gerar embedding do texto
	embedding, err := e.GenerateEmbedding(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Extrair palavras-chave (simples split por enquanto, ideal seria NLP)
	keywords := strings.Fields(text)

	// Criar ponto no Qdrant
	pointID := uint64(time.Now().UnixNano())

	payload := map[string]interface{}{
		"idoso_id":         idosoID,
		"text":             text,
		"keywords":         keywords,
		"emotional_charge": emotionalCharge,
		"timestamp":        time.Now().Unix(),
	}

	// Use wrapper's CreatePoint and Upsert
	point := vector.CreatePoint(pointID, embedding, payload)

	err = e.qdrantClient.Upsert(ctx, e.collectionName, []*qdrant.PointStruct{point})
	if err != nil {
		return fmt.Errorf("failed to upsert point: %w", err)
	}

	log.Printf("🔗 [QDRANT] Signifier chain tracked: idoso=%d", idosoID)
	return nil
}

// FindRelatedSignifiers busca significantes semanticamente relacionados
func (e *EmbeddingService) FindRelatedSignifiers(ctx context.Context, idosoID int64, text string, limit int) ([]SignifierChain, error) {
	if e.qdrantClient == nil {
		return nil, nil
	}

	// Gerar embedding da consulta
	embedding, err := e.GenerateEmbedding(ctx, text)
	if err != nil {
		return nil, err
	}

	// Buscar pontos similares usando wrapper
	// Filter just for user_id
	// Need to manually create filter for wrapper Search, or use SearchWithScore which has user filter built-in
	// Let's use Search with explicit filter as wrapper's SearchWithScore enforces a min score which we can set to 0.7 maybe?

	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: "idoso_id",
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Integer{Integer: idosoID},
						},
					},
				},
			},
		},
	}

	searchResult, err := e.qdrantClient.Search(ctx, e.collectionName, embedding, uint64(limit), filter)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	var chains []SignifierChain
	for _, point := range searchResult {
		chain := SignifierChain{}

		if kw, ok := point.Payload["keywords"]; ok {
			if list, ok := kw.GetKind().(*qdrant.Value_ListValue); ok {
				for _, v := range list.ListValue.Values {
					if s, ok := v.GetKind().(*qdrant.Value_StringValue); ok {
						chain.RelatedWords = append(chain.RelatedWords, s.StringValue)
					}
				}
			}
		}

		if len(chain.RelatedWords) > 0 {
			chain.CoreSignifier = chain.RelatedWords[0]
		}

		if charge, ok := point.Payload["emotional_charge"]; ok {
			if dbl, ok := charge.GetKind().(*qdrant.Value_DoubleValue); ok {
				chain.EmotionalCharge = dbl.DoubleValue
			}
		}

		if txt, ok := point.Payload["text"]; ok {
			if str, ok := txt.GetKind().(*qdrant.Value_StringValue); ok {
				chain.Contexts = append(chain.Contexts, str.StringValue)
			}
		}

		if ts, ok := point.Payload["timestamp"]; ok {
			if intVal, ok := ts.GetKind().(*qdrant.Value_IntegerValue); ok {
				chain.LastOccurrence = time.Unix(intVal.IntegerValue, 0)
			}
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

// GetSemanticContext monta contexto para o prompt usando similaridade semântica
func (e *EmbeddingService) GetSemanticContext(ctx context.Context, idosoID int64, currentText string) string {
	chains, err := e.FindRelatedSignifiers(ctx, idosoID, currentText, 5)
	if err != nil {
		log.Printf("⚠️ Error finding related signifiers: %v", err)
		return ""
	}

	if len(chains) == 0 {
		return ""
	}

	context := "\n🔗 CADEIA DE SIGNIFICANTES (Análise Semântica):\n\n"
	context += "O sistema detectou que palavras/temas similares já apareceram antes:\n"

	for i, chain := range chains {
		context += fmt.Sprintf("%d. '%s' (carga emocional: %.2f)\n",
			i+1, chain.CoreSignifier, chain.EmotionalCharge)

		if len(chain.Contexts) > 0 {
			context += fmt.Sprintf("   Contexto anterior: \"%s\"\n",
				truncateText(chain.Contexts[0], 100))
		}

		context += fmt.Sprintf("   Última vez: %s\n\n",
			chain.LastOccurrence.Format("02/01/2006 15:04"))
	}

	context += "→ Use essas informações para fazer conexões entre o que o paciente disse antes e agora.\n"
	context += "→ Se houver repetição de temas, isso pode indicar um nó sintomático.\n"

	return context
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
