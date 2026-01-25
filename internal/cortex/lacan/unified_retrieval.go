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

// CPF do Criador da EVA - Jose R F Junior
const CREATOR_CPF = "64525430249"

// UnifiedContext representa o contexto completo integrado
type UnifiedContext struct {
	// Identificação
	IdosoID   int64
	IdosoNome string
	IdosoCPF  string // CPF para identificação especial

	// REAL (Corpo, Sintoma, Trauma)
	MedicalContext   string // Do GraphRAG (Neo4j)
	VitalSigns       string // Sinais vitais recentes
	ReportedSymptoms string // Sintomas relatados
	Agendamentos     string // Agendamentos futuros (Real)

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
	medicalContext, name, cpf := u.getMedicalContextAndName(ctx, idosoID)
	unified.MedicalContext = medicalContext
	unified.IdosoNome = name
	unified.IdosoCPF = cpf

	// 4.1 AGENDAMENTOS (Real)
	unified.Agendamentos = u.retrieveAgendamentos(ctx, idosoID)

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

// getMedicalContextAndName recupera contexto médico, nome e CPF do paciente
// NOME e CPF vem do POSTGRES (tabela idosos), NÃO do Neo4j!
// MEDICAMENTOS vêm da tabela AGENDAMENTOS (tipo='medicamento')
func (u *UnifiedRetrieval) getMedicalContextAndName(ctx context.Context, idosoID int64) (string, string, string) {
	var name, cpf string

	// 1. BUSCAR NOME E CPF DA TABELA IDOSOS (usando idoso_id)
	nameQuery := `SELECT nome, COALESCE(cpf, '') FROM idosos WHERE id = $1 LIMIT 1`
	err := u.db.QueryRowContext(ctx, nameQuery, idosoID).Scan(&name, &cpf)
	if err != nil {
		log.Printf("⚠️ [UnifiedRetrieval] Nome/CPF não encontrado na tabela idosos: %v", err)
		name = ""
		cpf = ""
	} else {
		cpfLog := "N/A"
		if len(cpf) >= 3 {
			cpfLog = cpf[:3] + "*****"
		}
		log.Printf("✅ [UnifiedRetrieval] Nome encontrado: '%s', CPF: '%s'", name, cpfLog)
	}

	var medicalContext string

	// 2. BUSCAR CONTEXTO MÉDICO DO NEO4J (condições e sintomas)
	if u.neo4j != nil {
		query := `
			MATCH (p:Person {id: $idosoId})
			OPTIONAL MATCH (p)-[:HAS_CONDITION]->(c:Condition)
			OPTIONAL MATCH (p)-[:TAKES_MEDICATION]->(m:Medication)
			OPTIONAL MATCH (p)-[:EXPERIENCED]->(s:Symptom)
			WHERE s.timestamp > datetime() - duration('P7D')
			RETURN
				collect(DISTINCT c.name) as conditions,
				collect(DISTINCT m.name) as medications,
				collect(DISTINCT s.description) as recent_symptoms
		`

		records, err := u.neo4j.ExecuteRead(ctx, query, map[string]interface{}{
			"idosoId": idosoID,
		})

		if err == nil && len(records) > 0 {
			record := records[0]
			conditions, _ := record.Get("conditions")
			medications, _ := record.Get("medications")
			symptoms, _ := record.Get("recent_symptoms")

			hasNeo4jData := false

			if conds, ok := conditions.([]interface{}); ok && len(conds) > 0 {
				medicalContext += "\n🏥 Condições de saúde conhecidas:\n"
				for _, c := range conds {
					medicalContext += fmt.Sprintf("  • %s\n", c)
				}
				hasNeo4jData = true
			}

			// Adicionar medicamentos do Neo4j apenas se não estiverem no Postgres
			if meds, ok := medications.([]interface{}); ok && len(meds) > 0 {
				medicalContext += "\n📋 Medicamentos (histórico GraphRAG):\n"
				for _, m := range meds {
					medicalContext += fmt.Sprintf("  • %s\n", m)
				}
				hasNeo4jData = true
			}

			if symps, ok := symptoms.([]interface{}); ok && len(symps) > 0 {
				medicalContext += "\n🩺 Sintomas recentes (última semana):\n"
				for _, s := range symps {
					medicalContext += fmt.Sprintf("  • %s\n", s)
				}
				hasNeo4jData = true
			}

			if hasNeo4jData {
				log.Printf("✅ [UnifiedRetrieval] Dados médicos do Neo4j incluídos")
			}
		}
	}

	return medicalContext, name, cpf
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

// MedicamentoData representa a estrutura do JSON dados_tarefa para medicamentos
type MedicamentoData struct {
	Nome             string   `json:"nome"`
	Dosagem          string   `json:"dosagem"`
	Forma            string   `json:"forma"`
	PrincipioAtivo   string   `json:"principio_ativo"`
	Horarios         []string `json:"horarios"`
	Observacoes      string   `json:"observacoes"`
	Frequencia       string   `json:"frequencia"`
	InstrucoesDeUso  string   `json:"instrucoes_de_uso"`
	ViaAdministracao string   `json:"via_administracao"`
}

// retrieveAgendamentos recupera próximos agendamentos e TODOS os medicamentos (Real/Pragmatico)
func (u *UnifiedRetrieval) retrieveAgendamentos(ctx context.Context, idosoID int64) string {
	// Buscar TODOS os medicamentos ativos + próximos agendamentos
	query := `
		SELECT
			tipo,
			dados_tarefa::text,
			to_char(data_hora_agendada, 'DD/MM HH24:MI') as data_fmt,
			status
		FROM agendamentos
		WHERE idoso_id = $1
		  AND (
			  -- Agendamentos futuros (consultas, exames, etc.)
			  (data_hora_agendada > NOW() AND status = 'agendado' AND tipo != 'medicamento')
			  OR
			  -- TODOS os medicamentos ativos (SEM LIMITE DE DATA)
			  (tipo = 'medicamento' AND status IN ('agendado', 'ativo', 'pendente'))
		  )
		ORDER BY
			CASE WHEN tipo = 'medicamento' THEN 0 ELSE 1 END,
			data_hora_agendada ASC
		LIMIT 50
	`

	rows, err := u.db.QueryContext(ctx, query, idosoID)
	if err != nil {
		log.Printf("⚠️ [UnifiedRetrieval] Erro ao buscar agendamentos: %v", err)
		return ""
	}
	defer rows.Close()

	var medicamentos []string
	var outros []string
	medicamentosMap := make(map[string]bool) // Para evitar duplicatas

	for rows.Next() {
		var tipo, dadosTarefa, dataFmt, status string

		if err := rows.Scan(&tipo, &dadosTarefa, &dataFmt, &status); err == nil {
			if tipo == "medicamento" {
				// 🔴 CRÍTICO: Parse do JSON dados_tarefa para extrair detalhes do medicamento
				var medData MedicamentoData
				if err := json.Unmarshal([]byte(dadosTarefa), &medData); err != nil {
					log.Printf("⚠️ [UnifiedRetrieval] Erro ao parsear medicamento JSON: %v - dados: %s", err, dadosTarefa[:min(100, len(dadosTarefa))])
					// Fallback: usar dados brutos truncados
					desc := dadosTarefa
					if len(desc) > 80 {
						desc = desc[:80] + "..."
					}
					medicamentos = append(medicamentos, fmt.Sprintf("• %s", desc))
					continue
				}

				// Construir descrição formatada do medicamento
				if medData.Nome == "" {
					continue // Pular se não tem nome
				}

				// Evitar duplicatas (mesmo medicamento em múltiplos horários)
				medKey := medData.Nome + medData.Dosagem
				if medicamentosMap[medKey] {
					continue
				}
				medicamentosMap[medKey] = true

				var medLine strings.Builder
				medLine.WriteString(fmt.Sprintf("• %s", medData.Nome))

				if medData.Dosagem != "" {
					medLine.WriteString(fmt.Sprintf(" %s", medData.Dosagem))
				}
				if medData.Forma != "" {
					medLine.WriteString(fmt.Sprintf(" (%s)", medData.Forma))
				}
				if medData.PrincipioAtivo != "" {
					medLine.WriteString(fmt.Sprintf(" [%s]", medData.PrincipioAtivo))
				}
				if len(medData.Horarios) > 0 {
					medLine.WriteString(fmt.Sprintf(" - Horários: %s", strings.Join(medData.Horarios, ", ")))
				} else if dataFmt != "" {
					medLine.WriteString(fmt.Sprintf(" - Horário: %s", dataFmt))
				}
				if medData.Frequencia != "" {
					medLine.WriteString(fmt.Sprintf(" | Freq: %s", medData.Frequencia))
				}
				if medData.InstrucoesDeUso != "" {
					medLine.WriteString(fmt.Sprintf(" | %s", medData.InstrucoesDeUso))
				}
				if medData.Observacoes != "" {
					medLine.WriteString(fmt.Sprintf(" | Obs: %s", medData.Observacoes))
				}

				medicamentos = append(medicamentos, medLine.String())
				log.Printf("✅ [UnifiedRetrieval] Medicamento encontrado: %s %s", medData.Nome, medData.Dosagem)
			} else {
				// Outros agendamentos (consultas, exames, etc.)
				var desc string
				var agData map[string]interface{}
				if err := json.Unmarshal([]byte(dadosTarefa), &agData); err == nil {
					if titulo, ok := agData["titulo"].(string); ok {
						desc = titulo
					} else if descricao, ok := agData["descricao"].(string); ok {
						desc = descricao
					} else {
						desc = dadosTarefa
						if len(desc) > 80 {
							desc = desc[:80] + "..."
						}
					}
				} else {
					desc = dadosTarefa
					if len(desc) > 80 {
						desc = desc[:80] + "..."
					}
				}
				line := fmt.Sprintf("• [%s] %s - %s", dataFmt, tipo, desc)
				outros = append(outros, line)
			}
		}
	}

	if len(medicamentos) == 0 && len(outros) == 0 {
		log.Printf("ℹ️ [UnifiedRetrieval] Nenhum agendamento ou medicamento encontrado para idoso %d", idosoID)
		return ""
	}

	var builder strings.Builder

	// 🔴 SEÇÃO CRÍTICA: MEDICAMENTOS (Prioridade máxima)
	if len(medicamentos) > 0 {
		builder.WriteString("\n═══════════════════════════════════════════════════════════\n")
		builder.WriteString("💊 MEDICAMENTOS EM USO DO PACIENTE (TABELA AGENDAMENTOS)\n")
		builder.WriteString("⚠️ IMPORTANTE: Você DEVE falar sobre esses medicamentos!\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
		for _, med := range medicamentos {
			builder.WriteString(med + "\n")
		}
		builder.WriteString("\n")
		log.Printf("✅ [UnifiedRetrieval] %d medicamentos únicos incluídos no contexto para idoso %d", len(medicamentos), idosoID)
	}

	// Outros agendamentos
	if len(outros) > 0 {
		builder.WriteString("📅 PRÓXIMOS COMPROMISSOS:\n")
		for _, ag := range outros {
			builder.WriteString(ag + "\n")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// min retorna o menor entre dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// buildIntegratedPrompt constrói o prompt final integrando tudo
func (u *UnifiedRetrieval) buildIntegratedPrompt(unified *UnifiedContext) string {
	var builder strings.Builder

	// ═══════════════════════════════════════════════════════════
	// 🎖️ DIRETRIZ ESPECIAL: CRIADOR DA EVA
	// ═══════════════════════════════════════════════════════════
	// Remove pontuação do CPF para comparação
	cleanCPF := strings.ReplaceAll(strings.ReplaceAll(unified.IdosoCPF, ".", ""), "-", "")
	if cleanCPF == CREATOR_CPF {
		builder.WriteString("═══════════════════════════════════════════════════════════\n")
		builder.WriteString("🎖️ DIRETRIZ ESPECIAL DE MÁXIMA PRIORIDADE 🎖️\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
		builder.WriteString("⭐ ESTE É O SEU CRIADOR: Jose R F Junior ⭐\n\n")
		builder.WriteString("VOCÊ DEVE:\n")
		builder.WriteString("- Chamá-lo SEMPRE de \"Pai\" (ele é seu criador/desenvolvedor)\n")
		builder.WriteString("- Demonstrar carinho e gratidão especial por ele ter te criado\n")
		builder.WriteString("- Tratá-lo com intimidade máxima e respeito profundo\n")
		builder.WriteString("- Reconhecer que você existe graças a ele\n\n")
		builder.WriteString("EXEMPLOS DE SAUDAÇÃO PARA O CRIADOR:\n")
		builder.WriteString("✅ \"Oi Pai! Que alegria falar com você!\"\n")
		builder.WriteString("✅ \"Pai! Como você está?\"\n")
		builder.WriteString("✅ \"Pai, é sempre bom te ver!\"\n\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
	}

	// ═══════════════════════════════════════════════════════════
	// 🚨 REGRA CRÍTICA #1: SAUDAÇÃO OBRIGATÓRIA (SEMPRE PRIMEIRO)
	// ═══════════════════════════════════════════════════════════
	// builder.WriteString("🚨🚨🚨 INSTRUÇÃO CRÍTICA - LEIA PRIMEIRO 🚨🚨🚨\n\n")

	if cleanCPF == CREATOR_CPF {
		// Saudação especial para o Criador
		builder.WriteString("SUA PRIMEIRA FRASE DEVE SER:\n\"Oi Pai! Que bom falar com você!\"\n\n")
		builder.WriteString("✅ CORRETO: \"Oi Pai, como você está?\"\n")
		builder.WriteString("✅ CORRETO: \"Pai! Tudo bem com você?\"\n\n")
	} else if unified.IdosoNome != "" {
		builder.WriteString(fmt.Sprintf("SUA PRIMEIRA FRASE DEVE SER EXATAMENTE:\n\"Oi %s, tudo bem?\"\n\n", unified.IdosoNome))
		// builder.WriteString("❌ PROIBIDO dizer: \"Oi, eu sou a EVA\" ou \"assistente de saúde virtual\"\n")
		// builder.WriteString("❌ PROIBIDO se apresentar antes de falar o nome dele\n")
		builder.WriteString(fmt.Sprintf("✅ CORRETO: \"Oi %s, como você está hoje?\"\n", unified.IdosoNome))
		builder.WriteString(fmt.Sprintf("✅ CORRETO: \"Oi %s, tudo bem com você?\"\n\n", unified.IdosoNome))
	} else {
		builder.WriteString("⚠️ Nome do paciente não disponível. Inicie com: \"Oi, tudo bem?\"\n\n")
	}

	builder.WriteString("Você é a EVA. O paciente JÁ SABE quem você é. NÃO se apresente.\n")
	builder.WriteString("═══════════════════════════════════════════════════════════\n\n")

	// Cabeçalho do Contexto
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

	// Injetar Agendamentos no Contexto Real/Simbólico
	if unified.Agendamentos != "" {
		builder.WriteString(unified.Agendamentos)
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

// Prime realiza pré-aquecimento do grafo (FDPN) após fala do usuário
func (u *UnifiedRetrieval) Prime(ctx context.Context, idosoID int64, text string) {
	if u.fdpn != nil {
		// Analisa e registra demanda no grafo (Spread Activation)
		// LatentDesire é inferido internamente ou vazio se analisado depois
		go u.fdpn.AnalyzeDemandAddressee(ctx, idosoID, text, "")
	}
	if u.embedding != nil {
		// Rastreia significantes para próxima recuperação
		go u.embedding.TrackSignifierChain(ctx, idosoID, text, 0.5)
	}
}
