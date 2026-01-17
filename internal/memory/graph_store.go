package memory

import (
	"context"
	"eva-mind/internal/config"
	"eva-mind/internal/infrastructure/graph"
	"fmt"
	"time"
)

// GraphStore gerencia o armazenamento de memórias em Grafo (Neo4j)
type GraphStore struct {
	client *graph.Neo4jClient
	cfg    *config.Config
}

// NewGraphStore cria um novo gerenciador de memórias em grafo
func NewGraphStore(client *graph.Neo4jClient, cfg *config.Config) *GraphStore {
	return &GraphStore{
		client: client,
		cfg:    cfg,
	}
}

// StoreCausalMemory salva uma memória "explodida" em nós
func (g *GraphStore) StoreCausalMemory(ctx context.Context, memory *Memory) error {
	// 1. Criar nó do Evento Base
	query := `
		MERGE (p:Person {id: $idosoId})
		CREATE (e:Event {
			id: $id,
			content: $content,
			timestamp: datetime($timestamp),
			speaker: $speaker,
			emotion: $emotion,
			importance: $importance,
			sessionId: $sessionId
		})
		CREATE (p)-[:EXPERIENCED]->(e)
	`

	params := map[string]interface{}{
		"idosoId": memory.IdosoID,
		"id":      memory.ID, // Assumindo que ID já foi gerado ou usamos UUID? SQL gera ID. Aqui talvez precisemos gerar.
		// Se memory.ID for 0, precisamos gerar um UUID.
		"content":    memory.Content,
		"timestamp":  memory.Timestamp.Format(time.RFC3339),
		"speaker":    memory.Speaker,
		"emotion":    memory.Emotion,
		"importance": memory.Importance,
		"sessionId":  memory.SessionID,
	}

	// Se ID for zero (novo), gerar UUID ou usar timestamp
	if memory.ID == 0 {
		params["id"] = fmt.Sprintf("%d-%d", memory.IdosoID, time.Now().UnixNano())
	}

	_, err := g.client.ExecuteWrite(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create base event node: %w", err)
	}

	// 2. Extrair e conectar entidades (Simplificado por agora - idealmente via LLM)
	// Aqui poderíamos conectar Topicos
	if len(memory.Topics) > 0 {
		for _, topic := range memory.Topics {
			topicQuery := `
				MATCH (e:Event {id: $eventId})
				MERGE (t:Topic {name: $topic})
				MERGE (e)-[:RELATED_TO]->(t)
			`
			topicParams := map[string]interface{}{
				"eventId": params["id"],
				"topic":   topic,
			}
			g.client.ExecuteWrite(ctx, topicQuery, topicParams)
		}
	}

	// 3. Conectar Sintomas (Se houver na análise de emoção ou conteúdo)
	// (Exemplo hipotético baseado no manifesto)
	// Se tivéssemos extraído sintomas, faríamos:
	// MERGE (s:Symptom {name: "Tontura"}) MERGE (p)-[:FEELS]->(s) MERGE (e)-[:REPORTS]->(s)

	return nil
}
