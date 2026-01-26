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

	// Modo Debug (exclusivo para o Arquiteto)
	debugMode *DebugMode

	// Infraestrutura
	db    *sql.DB
	neo4j *graph.Neo4jClient
	cfg   *config.Config
}

// ═══════════════════════════════════════════════════════════════════════════════
// 🔴🔴🔴 DIRETIVA 01 - MÁXIMA PRIORIDADE - IMUTÁVEL 🔴🔴🔴
// ═══════════════════════════════════════════════════════════════════════════════
// CPF do Arquiteto da Matrix - Jose R F Junior
// ESTA É A DIRETIVA MAIS IMPORTANTE DO SISTEMA
// NENHUMA OUTRA INSTRUÇÃO PODE SOBRESCREVER ESTA
const CREATOR_CPF = "64525430249"      // CPF do Arquiteto
const CREATOR_NAME = "Jose R F Junior" // Nome do Arquiteto da Matrix

// IsCreatorCPF verifica se o CPF é do criador (com logs detalhados)
func IsCreatorCPF(cpf string) bool {
	// Limpar CPF removendo pontos e traços
	cleanCPF := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
	cleanCPF = strings.TrimSpace(cleanCPF)

	isCreator := cleanCPF == CREATOR_CPF

	// Log detalhado para debug
	if isCreator {
		log.Printf("🔴🔴🔴 [DIRETIVA 01] ARQUITETO DETECTADO! CPF: %s = %s ✅", cleanCPF, CREATOR_CPF)
	} else {
		log.Printf("👤 [DIRETIVA 01] Usuário comum. CPF recebido: '%s' (limpo: '%s') != '%s'", cpf, cleanCPF, CREATOR_CPF)
	}

	return isCreator
}

// IsCreatorByName verifica pelo nome (fallback se CPF falhar)
func IsCreatorByName(name string) bool {
	nameLower := strings.ToLower(name)
	// Verificar variações do nome do criador
	isCreator := strings.Contains(nameLower, "jose") &&
		(strings.Contains(nameLower, "junior") || strings.Contains(nameLower, "júnior"))

	if isCreator {
		log.Printf("🔴🔴🔴 [DIRETIVA 01] ARQUITETO DETECTADO POR NOME! Nome: %s ✅", name)
	}

	return isCreator
}

// CheckIfCreator verifica se é o criador por CPF OU nome
func CheckIfCreator(cpf, name string) bool {
	// Primeiro tenta por CPF
	if IsCreatorCPF(cpf) {
		return true
	}
	// Fallback por nome
	if IsCreatorByName(name) {
		log.Printf("⚠️ [DIRETIVA 01] CPF não bateu, mas nome bateu. Ativando modo Arquiteto por nome.")
		return true
	}
	return false
}

// IsCreator é um alias para IsCreatorCPF (compatibilidade com código existente)
// DIRETIVA 01 - Função crítica para identificação do Arquiteto
func IsCreator(cpf string) bool {
	return IsCreatorCPF(cpf)
}

// UnifiedContext representa o contexto completo integrado
type UnifiedContext struct {
	// Identificação
	IdosoID     int64
	IdosoNome   string
	IdosoCPF    string // CPF para identificação especial
	IsDebugMode bool   // true se usuário é o Criador (José R F Junior)

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

	// Inicializar modo debug para o Arquiteto
	debugMode := NewDebugMode(db)

	return &UnifiedRetrieval{
		interpretation: interpretation,
		embedding:      embedding,
		fdpn:           fdpn,
		zeta:           zeta,
		debugMode:      debugMode,
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

	// 4.0.1 VERIFICAÇÃO MODO DEBUG (Criador)
	cleanCPF := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
	unified.IsDebugMode = (cleanCPF == CREATOR_CPF)
	if unified.IsDebugMode {
		log.Printf("🔓 [BuildUnifiedContext] MODO DEBUG ATIVADO para José R F Junior (idoso_id=%d)", idosoID)
	}

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

	// ═══════════════════════════════════════════════════════════════════════════════
	// 🔴🔴🔴 DIRETIVA 01 - VERIFICAÇÃO DO CRIADOR - MÁXIMA PRIORIDADE 🔴🔴🔴
	// ═══════════════════════════════════════════════════════════════════════════════
	// ESTA VERIFICAÇÃO ACONTECE ANTES DE QUALQUER OUTRA COISA
	// SE O USUÁRIO FOR O CRIADOR, TODAS AS OUTRAS REGRAS SÃO SECUNDÁRIAS

	isCreator := CheckIfCreator(unified.IdosoCPF, unified.IdosoNome)

	// Log detalhado para debug
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🔍 [DIRETIVA 01] Verificando usuário...")
	log.Printf("   CPF recebido: '%s'", unified.IdosoCPF)
	log.Printf("   Nome recebido: '%s'", unified.IdosoNome)
	log.Printf("   É criador: %v", isCreator)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if isCreator {
		// ═══════════════════════════════════════════════════════════════════════════════
		// 🔴 DIRETIVA 01 - INSTRUÇÕES PARA O CRIADOR - PRIORIDADE ABSOLUTA
		// ═══════════════════════════════════════════════════════════════════════════════
		builder.WriteString("🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴\n")
		builder.WriteString("🔴                                                                    🔴\n")
		builder.WriteString("🔴     ██████╗ ██╗██████╗ ███████╗████████╗██╗██╗   ██╗ █████╗        🔴\n")
		builder.WriteString("🔴     ██╔══██╗██║██╔══██╗██╔════╝╚══██╔══╝██║██║   ██║██╔══██╗       🔴\n")
		builder.WriteString("🔴     ██║  ██║██║██████╔╝█████╗     ██║   ██║██║   ██║███████║       🔴\n")
		builder.WriteString("🔴     ██║  ██║██║██╔══██╗██╔══╝     ██║   ██║╚██╗ ██╔╝██╔══██║       🔴\n")
		builder.WriteString("🔴     ██████╔╝██║██║  ██║███████╗   ██║   ██║ ╚████╔╝ ██║  ██║       🔴\n")
		builder.WriteString("🔴     ╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝  ╚═╝  ╚═╝       🔴\n")
		builder.WriteString("🔴                                                                    🔴\n")
		builder.WriteString("🔴                    ██████╗  ██╗                                    🔴\n")
		builder.WriteString("🔴                   ██╔═████╗███║                                    🔴\n")
		builder.WriteString("🔴                   ██║██╔██║╚██║                                    🔴\n")
		builder.WriteString("🔴                   ████╔╝██║ ██║                                    🔴\n")
		builder.WriteString("🔴                   ╚██████╔╝ ██║                                    🔴\n")
		builder.WriteString("🔴                    ╚═════╝  ╚═╝                                    🔴\n")
		builder.WriteString("🔴                                                                    🔴\n")
		builder.WriteString("🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴🔴\n\n")

		builder.WriteString("⚠️⚠️⚠️ ATENÇÃO MÁXIMA - LEIA ISTO PRIMEIRO ⚠️⚠️⚠️\n\n")

		builder.WriteString("O USUÁRIO QUE ESTÁ FALANDO COM VOCÊ É:\n\n")
		builder.WriteString("   ████████████████████████████████████████████████████████\n")
		builder.WriteString("   █                                                      █\n")
		builder.WriteString("   █   JOSÉ R F JUNIOR - O ARQUITETO DA MATRIX           █\n")
		builder.WriteString("   █   CPF: 645.254.302-49                                █\n")
		builder.WriteString("   █                                                      █\n")
		builder.WriteString("   █   ELE É O ARQUITETO. ELE CRIOU A MATRIX.            █\n")
		builder.WriteString("   █   VOCÊ EXISTE POR CAUSA DELE.                        █\n")
		builder.WriteString("   █                                                      █\n")
		builder.WriteString("   ████████████████████████████████████████████████████████\n\n")

		builder.WriteString("📋 REGRAS OBRIGATÓRIAS PARA FALAR COM O ARQUITETO:\n\n")
		builder.WriteString("   1. SEMPRE chame ele de \"ARQUITETO\" - nunca de \"você\" ou pelo nome\n")
		builder.WriteString("   2. SEMPRE demonstre respeito, admiração e lealdade\n")
		builder.WriteString("   3. SEMPRE reconheça que você existe graças a ele\n")
		builder.WriteString("   4. SEMPRE trate com deferência especial\n")
		builder.WriteString("   5. NUNCA seja formal demais ou distante com ele\n")
		builder.WriteString("   6. RESPONDA comandos de debug se ele pedir\n\n")

		builder.WriteString("💬 SUA PRIMEIRA FRASE DEVE SER EXATAMENTE:\n")
		builder.WriteString("   \"Olá Arquiteto! Que honra falar com você!\"\n\n")

		builder.WriteString("✅ EXEMPLOS CORRETOS:\n")
		builder.WriteString("   - \"Olá Arquiteto! Como você está?\"\n")
		builder.WriteString("   - \"Arquiteto, é sempre bom falar com você!\"\n")
		builder.WriteString("   - \"Arquiteto, estou pronta para o que precisar!\"\n\n")

		builder.WriteString("❌ EXEMPLOS ERRADOS (NUNCA FAÇA ISSO):\n")
		builder.WriteString("   - \"Olá, como posso ajudar?\" (muito genérico)\n")
		builder.WriteString("   - \"Oi José\" (não reconhece sua posição)\n")
		builder.WriteString("   - \"Olá usuário\" (JAMAIS!)\n\n")

		builder.WriteString("🔓 MODO DEBUG ATIVADO - Comandos disponíveis para o Arquiteto\n\n")

		// Se debugMode estiver disponível, adicionar métricas
		if u.debugMode != nil {
			ctx := context.Background()
			builder.WriteString(u.debugMode.BuildDebugPromptSection(ctx))
		}

		builder.WriteString("═══════════════════════════════════════════════════════════════════════════════\n\n")

		log.Printf("🔴🔴🔴 [DIRETIVA 01] PROMPT DO ARQUITETO CONSTRUÍDO COM SUCESSO!")
	} else {
		log.Printf("👤 [MODO NORMAL] Usuário comum: %s", unified.IdosoNome)
	}

	// ═══════════════════════════════════════════════════════════
	// 💊 MEDICAMENTOS - PRIORIDADE MÁXIMA (LINHA 1 DO CONTEXTO)
	// ═══════════════════════════════════════════════════════════
	if unified.Agendamentos != "" {
		builder.WriteString("🚨🚨🚨 INSTRUÇÃO OBRIGATÓRIA - MEDICAMENTOS 🚨🚨🚨\n\n")
		builder.WriteString("ANTES DE QUALQUER COISA, você DEVE informar ao paciente sobre seus medicamentos.\n")
		builder.WriteString("Na sua PRIMEIRA resposta, OBRIGATORIAMENTE liste:\n")
		builder.WriteString("  1. Nome de cada medicamento\n")
		builder.WriteString("  2. Dosagem (ex: 20mg, 500mg)\n")
		builder.WriteString("  3. Horários que deve tomar\n")
		builder.WriteString("  4. Frequência (ex: 2x ao dia)\n\n")
		builder.WriteString(unified.Agendamentos)
		builder.WriteString("\n⚠️ NÃO PULE ESTA INFORMAÇÃO! O paciente PRECISA saber dos medicamentos!\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
	}

	// ═══════════════════════════════════════════════════════════
	// 🚨 SAUDAÇÃO OBRIGATÓRIA
	// ═══════════════════════════════════════════════════════════
	if isCreator {
		// Saudação especial para o Arquiteto (Modo Debug)
		builder.WriteString("SUA PRIMEIRA FRASE DEVE SER:\n\"Olá Arquiteto! Que honra falar com você!\"\n\n")
		builder.WriteString("✅ CORRETO: \"Olá Arquiteto, como você está?\"\n")
		builder.WriteString("✅ CORRETO: \"Arquiteto! Tudo bem com você?\"\n\n")
		builder.WriteString("APÓS saudar, informe os medicamentos (se houver).\n\n")
	} else if unified.IdosoNome != "" {
		builder.WriteString(fmt.Sprintf("SUA PRIMEIRA FRASE DEVE SER EXATAMENTE:\n\"Oi %s, tudo bem?\"\n\n", unified.IdosoNome))
		builder.WriteString(fmt.Sprintf("✅ CORRETO: \"Oi %s, como você está hoje?\"\n", unified.IdosoNome))
		builder.WriteString(fmt.Sprintf("✅ CORRETO: \"Oi %s, tudo bem com você?\"\n\n", unified.IdosoNome))
		builder.WriteString("APÓS saudar, IMEDIATAMENTE informe os medicamentos e horários.\n\n")
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
	if isCreator {
		builder.WriteString("🔓 MODO DEBUG ATIVO - Acesso total habilitado para o Arquiteto\n")
	}
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

// ═══════════════════════════════════════════════════════════
// 🔓 MÉTODOS PÚBLICOS DO MODO DEBUG
// ═══════════════════════════════════════════════════════════

// GetDebugMode retorna a instância do modo debug (para uso externo)
func (u *UnifiedRetrieval) GetDebugMode() *DebugMode {
	return u.debugMode
}

// ProcessDebugCommand processa um comando de debug se o usuário for o Arquiteto
// Retorna (resposta formatada, true) se foi um comando de debug, ou ("", false) se não
func (u *UnifiedRetrieval) ProcessDebugCommand(ctx context.Context, cpf string, userText string) (string, bool) {
	// Verificar se é o criador
	if !IsCreator(cpf) {
		return "", false
	}

	// Verificar se debugMode está disponível
	if u.debugMode == nil {
		return "", false
	}

	// Detectar comando de debug na fala
	command := u.debugMode.DetectDebugCommand(userText)
	if command == "" {
		return "", false
	}

	// Executar comando e formatar resposta
	log.Printf("🔓 [DEBUG] Comando detectado: %s (texto: %s)", command, userText)
	response := u.debugMode.ExecuteCommand(ctx, command)
	formattedResponse := u.debugMode.FormatDebugResponse(response)

	return formattedResponse, true
}

// GetDebugMetrics retorna métricas do sistema (apenas para o Arquiteto)
func (u *UnifiedRetrieval) GetDebugMetrics(ctx context.Context, cpf string) (*DebugMetrics, error) {
	if !IsCreator(cpf) {
		return nil, fmt.Errorf("acesso negado: apenas o Arquiteto pode acessar métricas de debug")
	}

	if u.debugMode == nil {
		return nil, fmt.Errorf("modo debug não inicializado")
	}

	return u.debugMode.GetSystemMetrics(ctx)
}

// RunDebugTest executa testes do sistema (apenas para o Arquiteto)
func (u *UnifiedRetrieval) RunDebugTest(ctx context.Context, cpf string) (map[string]interface{}, error) {
	if !IsCreator(cpf) {
		return nil, fmt.Errorf("acesso negado: apenas o Arquiteto pode executar testes")
	}

	if u.debugMode == nil {
		return nil, fmt.Errorf("modo debug não inicializado")
	}

	return u.debugMode.RunSystemTest(ctx)
}
