package lacan

import (
	"context"
	"database/sql"
	"encoding/json"
	"eva-mind/internal/brainstem/config"
	"eva-mind/internal/brainstem/infrastructure/graph"
	"eva-mind/internal/brainstem/infrastructure/vector"
	"eva-mind/internal/hippocampus/knowledge"
	"eva-mind/pkg/types"
	"fmt"
	"log"
	"strings"
)

// UnifiedRetrieval implementa "O Sinthoma" - a amarração dos registros RSI
// Real (trauma, corpo), Simbólico (linguagem, grafo), Imaginário (narrativa, memória)
// Integra TODOS os módulos lacanianos em um contexto coerente para o Gemini
type UnifiedRetrieval struct {
	// Módulos Lacanianos
	interpretation *InterpretationService
	embedding      *knowledge.EmbeddingService
	fdpn           *FDPNEngine
	zeta           *ZetaRouter

	// Infraestrutura
	db    *sql.DB
	neo4j *graph.Neo4jClient
	cfg   *config.Config
}

// UnifiedContext representa o contexto completo integrado
type UnifiedContext struct {
	// Identificação
	IdosoID   int64
	IdosoNome string

	// REAL (Corpo, Sintoma, Trauma)
	MedicalContext   string // Do GraphRAG (Neo4j)
	VitalSigns       string // Sinais vitais recentes
	ReportedSymptoms string // Sintomas relatados

	// SIMBÓLICO (Linguagem, Estrutura, Grafo)
	LacanianAnalysis *InterpretationResult // Análise lacaniana completa
	DemandGraph      string                // Grafo de demandas (FDPN)
	SignifierChains  string                // Cadeias de significantes (Qdrant)

	// IMAGINÁRIO (Narrativa, Memória, História)
	RecentMemories []string                  // Memórias episódicas recentes
	LifeStory      string                    // Narrativa de vida (se disponível)
	Patterns       []*types.RecurrentPattern // Padrões detectados

	// INTERVENÇÃO (Ética + Postura)
	EthicalStance *EthicalStance
	GurdjieffType int    // Tipo de atenção recomendado
	SystemPrompt  string // Prompt final integrado
}

// NewUnifiedRetrieval cria serviço de recuperação unificada
func NewUnifiedRetrieval(
	db *sql.DB,
	neo4j *graph.Neo4jClient,
	qdrant *vector.QdrantClient,
	cfg *config.Config,
) *UnifiedRetrieval {
	interpretation := NewInterpretationService(db, neo4j)

	embedding, err := knowledge.NewEmbeddingService(cfg, qdrant)
	if err != nil {
		log.Printf("⚠️ Warning: Embedding service initialization failed: %v", err)
	}

	fdpn := NewFDPNEngine(neo4j)
	zeta := NewZetaRouter(interpretation)

	return &UnifiedRetrieval{
		interpretation: interpretation,
		embedding:      embedding,
		fdpn:           fdpn,
		zeta:           zeta,
		db:             db,
		neo4j:          neo4j,
		cfg:            cfg,
	}
}

// BuildUnifiedContext constrói contexto completo integrando todos os módulos
func (u *UnifiedRetrieval) BuildUnifiedContext(
	ctx context.Context,
	idosoID int64,
	currentText string,
	previousText string,
) (*UnifiedContext, error) {

	unified := &UnifiedContext{
		IdosoID: idosoID,
	}

	// 1. ANÁLISE LACANIANA (Núcleo)
	lacanResult, err := u.interpretation.AnalyzeUtterance(ctx, idosoID, currentText, previousText)
	if err != nil {
		log.Printf("⚠️ Lacanian analysis failed: %v", err)
		// Continua mesmo com erro
	} else {
		unified.LacanianAnalysis = lacanResult
	}

	// 2. GRAFO DO DESEJO (A quem pede)
	if u.fdpn != nil {
		// Ajuste: usando LatentDesire do resultado Lacaniano
		var latent string
		if lacanResult != nil && lacanResult.DemandDesire != nil {
			latent = string(lacanResult.DemandDesire.LatentDesire)
		}
		addressee, _ := u.fdpn.AnalyzeDemandAddressee(ctx, idosoID, currentText, latent)
		unified.DemandGraph = u.fdpn.BuildGraphContext(ctx, idosoID)

		// Adicionar orientação do destinatário
		if addressee != ADDRESSEE_UNKNOWN {
			unified.DemandGraph += "\n" + GetClinicalGuidanceForAddressee(addressee)
		}
	}

	// 3. CADEIAS SEMÂNTICAS (Qdrant)
	if u.embedding != nil {
		unified.SignifierChains = u.embedding.GetSemanticContext(ctx, idosoID, currentText)
	}

	// 4. CONTEXTO MÉDICO (Neo4j GraphRAG)
	unified.MedicalContext = u.getMedicalContext(ctx, idosoID)

	// 5. MEMÓRIAS RECENTES (Postgres)
	unified.RecentMemories = u.getRecentMemories(ctx, idosoID, 5)

	// 6. POSTURA ÉTICA (Zeta Router)
	if lacanResult != nil {
		stance, _ := u.zeta.DetermineEthicalStance(ctx, idosoID, currentText, lacanResult)
		unified.EthicalStance = stance
		unified.GurdjieffType = u.zeta.DetermineGurdjieffType(ctx, idosoID, lacanResult)
	}

	// 7. CONSTRUIR PROMPT FINAL
	unified.SystemPrompt = u.buildIntegratedPrompt(unified)

	return unified, nil
}

// getMedicalContext recupera contexto médico do Neo4j
func (u *UnifiedRetrieval) getMedicalContext(ctx context.Context, idosoID int64) string {
	if u.neo4j == nil {
		return ""
	}

	query := `
		MATCH (p:Person {id: $idosoId})
		OPTIONAL MATCH (p)-[:HAS_CONDITION]->(c:Condition)
		OPTIONAL MATCH (p)-[:TAKES_MEDICATION]->(m:Medication)
		OPTIONAL MATCH (p)-[:EXPERIENCED]->(s:Symptom)
		WHERE s.timestamp > datetime() - duration('P7D')
		RETURN 
			p.name as name,
			collect(DISTINCT c.name) as conditions,
			collect(DISTINCT m.name) as medications,
			collect(DISTINCT s.description) as recent_symptoms
	`

	records, err := u.neo4j.ExecuteRead(ctx, query, map[string]interface{}{
		"idosoId": idosoID,
	})

	if err != nil || len(records) == 0 {
		return ""
	}

	record := records[0]
	name, _ := record.Get("name")
	conditions, _ := record.Get("conditions")
	medications, _ := record.Get("medications")
	symptoms, _ := record.Get("recent_symptoms")

	context := "\n🏥 CONTEXTO MÉDICO (GraphRAG):\n\n"

	if name != nil {
		context += fmt.Sprintf("Paciente: %s\n", name)
	}

	if conds, ok := conditions.([]interface{}); ok && len(conds) > 0 {
		context += "\nCondições conhecidas:\n"
		for _, c := range conds {
			context += fmt.Sprintf("- %s\n", c)
		}
	}

	if meds, ok := medications.([]interface{}); ok && len(meds) > 0 {
		context += "\nMedicamentos em uso:\n"
		for _, m := range meds {
			context += fmt.Sprintf("- %s\n", m)
		}
	}

	if symps, ok := symptoms.([]interface{}); ok && len(symps) > 0 {
		context += "\nSintomas recentes (última semana):\n"
		for _, s := range symps {
			context += fmt.Sprintf("- %s\n", s)
		}
	}

	return context
}

// getRecentMemories recupera memórias episódicas recentes
func (u *UnifiedRetrieval) getRecentMemories(ctx context.Context, idosoID int64, limit int) []string {
	query := `
		SELECT conteudo->'summary' as summary
		FROM analise_gemini
		WHERE idoso_id = $1 
		  AND tipo = 'AUDIO'
		  AND conteudo->'summary' IS NOT NULL
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := u.db.QueryContext(ctx, query, idosoID, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var memories []string
	for rows.Next() {
		var summary string
		if err := rows.Scan(&summary); err == nil {
			memories = append(memories, summary)
		}
	}

	return memories
}

// buildIntegratedPrompt constrói o prompt final integrando tudo
func (u *UnifiedRetrieval) buildIntegratedPrompt(unified *UnifiedContext) string {
	var builder strings.Builder

	// Cabeçalho
	builder.WriteString("═══════════════════════════════════════════════════════════\n")
	builder.WriteString("🧠 CONTEXTO INTEGRADO EVA-MIND (RSI - Real, Simbólico, Imaginário)\n")
	builder.WriteString("═══════════════════════════════════════════════════════════\n\n")

	// REAL (Corpo, Sintoma)
	if unified.MedicalContext != "" {
		builder.WriteString("▌REAL - CORPO E SINTOMA:\n")
		builder.WriteString(unified.MedicalContext)
		builder.WriteString("\n")
	}

	// SIMBÓLICO (Linguagem, Estrutura)
	builder.WriteString("▌SIMBÓLICO - ESTRUTURA E LINGUAGEM:\n\n")

	if unified.LacanianAnalysis != nil {
		builder.WriteString(unified.LacanianAnalysis.ClinicalGuidance)
		builder.WriteString("\n")
	}

	if unified.DemandGraph != "" {
		builder.WriteString(unified.DemandGraph)
		builder.WriteString("\n")
	}

	if unified.SignifierChains != "" {
		builder.WriteString(unified.SignifierChains)
		builder.WriteString("\n")
	}

	// IMAGINÁRIO (Narrativa, Memória)
	if len(unified.RecentMemories) > 0 {
		builder.WriteString("▌IMAGINÁRIO - NARRATIVA E MEMÓRIA:\n\n")
		builder.WriteString("Resumos de conversas recentes:\n")
		for i, mem := range unified.RecentMemories {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, mem))
		}
		builder.WriteString("\n")
	}

	// INTERVENÇÃO ÉTICA
	if unified.EthicalStance != nil {
		builder.WriteString(u.zeta.BuildEthicalPrompt(unified.EthicalStance))
		builder.WriteString("\n")
	}

	// Tipo de Atenção (Gurdjieff)
	var typeDirective string
	switch unified.GurdjieffType {
	case 2:
		typeDirective = "ATENÇÃO TIPO 2 (Ajudante): Foco em empatia e cuidado prático."
	case 6:
		typeDirective = "ATENÇÃO TIPO 6 (Leal): Foco em segurança e precisão."
	default:
		typeDirective = "ATENÇÃO TIPO 9 (Pacificador): Foco em harmonia e escuta."
	}
	builder.WriteString(fmt.Sprintf("🎯 %s\n\n", typeDirective))

	// Rodapé
	builder.WriteString("═══════════════════════════════════════════════════════════\n")
	builder.WriteString("⚠️ LEMBRE-SE: Você é EVA, não um modelo genérico.\n")
	builder.WriteString("Use este contexto como suas próprias memórias e insights.\n")
	builder.WriteString("═══════════════════════════════════════════════════════════\n")

	return builder.String()
}

// GetPromptForGemini retorna o prompt completo para ser usado com Gemini
func (u *UnifiedRetrieval) GetPromptForGemini(ctx context.Context, idosoID int64, currentText, previousText string) (string, error) {
	unified, err := u.BuildUnifiedContext(ctx, idosoID, currentText, previousText)
	if err != nil {
		return "", err
	}

	return unified.SystemPrompt, nil
}

// SaveConversationContext salva contexto da conversa para análise futura
func (u *UnifiedRetrieval) SaveConversationContext(ctx context.Context, idosoID int64, unified *UnifiedContext, userText, assistantText string) error {
	// Salvar no Postgres (análise)
	contextData := map[string]interface{}{
		"lacanian_analysis": unified.LacanianAnalysis,
		"ethical_stance":    unified.EthicalStance,
		"gurdjieff_type":    unified.GurdjieffType,
		"user_text":         userText,
		"assistant_text":    assistantText,
	}

	query := `
		INSERT INTO analise_gemini (idoso_id, tipo, conteudo, created_at)
		VALUES ($1, 'CONTEXT', $2, CURRENT_TIMESTAMP)
	`

	contextJSON, _ := json.Marshal(contextData)
	_, err := u.db.ExecContext(ctx, query, idosoID, contextJSON)

	return err
}
