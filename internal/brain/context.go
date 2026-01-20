package brain

import (
	"context"
	"database/sql"
	"eva-mind/pkg/types"
	"fmt"
	"log"
	"strings"
	"time"
)

// BuildSystemPrompt constructs the system prompt with patient data
func (s *Service) BuildSystemPrompt(idosoID int64) string {
	nomeDefault := "Paciente"

	// 1. Resilient Query
	query := `
		SELECT 
			nome, 
			EXTRACT(YEAR FROM AGE(data_nascimento)) as idade,
			nivel_cognitivo, 
			tom_voz,
			preferencia_horario_ligacao,
			medicamentos_atuais,
			condicoes_medicas,
			endereco
		FROM idosos 
		WHERE id = $1
	`

	var nome, nivelCognitivo, tomVoz string
	var idade int
	var preferenciaHorario sql.NullString
	var medicamentosAtuais, condicoesMedicas, endereco sql.NullString

	// Optional fields
	var mobilidade string = "Não informada"
	var limitacoesVisuais, familiarPrincipal, contatoEmergencia, medicoResponsavel, notasGerais sql.NullString
	var limitacoesAuditivas, usaAparelhoAuditivo, ambienteRuidoso sql.NullBool

	err := s.db.QueryRow(query, idosoID).Scan(
		&nome,
		&idade,
		&nivelCognitivo,
		&tomVoz,
		&preferenciaHorario,
		&medicamentosAtuais,
		&condicoesMedicas,
		&endereco,
	)

	if err != nil {
		log.Printf("⚠️ [Brain] Using partial data for %s due to SQL error: %v", nomeDefault, err)
		nome = nomeDefault
		idade = 0
		nivelCognitivo = "Não informado"
		tomVoz = "Suave"
	}

	// Fetch Relational Meds
	medsQuery := `
		SELECT nome, dosagem, horarios, observacoes 
		FROM medicamentos 
		WHERE idoso_id = $1 AND ativo = true
	`
	rows, errMeds := s.db.Query(medsQuery, idosoID)
	var medsList []string
	if errMeds == nil {
		defer rows.Close()
		for rows.Next() {
			var mNome, mDosagem, mHorarios, mObs string
			if err := rows.Scan(&mNome, &mDosagem, &mHorarios, &mObs); err == nil {
				medInfo := fmt.Sprintf("- %s (%s)", mNome, mDosagem)
				if mHorarios != "" {
					medInfo += fmt.Sprintf(" às %s", mHorarios)
				}
				if mObs != "" {
					medInfo += fmt.Sprintf(". Obs: %s", mObs)
				}
				medsList = append(medsList, medInfo)
			}
		}
	}

	// Fetch Agenda
	agendaQuery := `
		SELECT tipo, data_hora_agendada, dados_tarefa
		FROM agendamentos
		WHERE idoso_id = $1 
		  AND status = 'agendado'
		  AND data_hora_agendada >= NOW()
		ORDER BY data_hora_agendada ASC
	`
	rowsAgenda, errAgenda := s.db.Query(agendaQuery, idosoID)
	var agendaList []string
	if errAgenda == nil {
		defer rowsAgenda.Close()
		for rowsAgenda.Next() {
			var aTipo string
			var aData time.Time
			var aDadosJSON sql.NullString

			if err := rowsAgenda.Scan(&aTipo, &aData, &aDadosJSON); err == nil {
				dataHora := aData.Format("02/01 às 15:04")
				item := fmt.Sprintf("- [%s]: %s", dataHora, strings.Title(aTipo))
				if aDadosJSON.Valid && aDadosJSON.String != "{}" && aDadosJSON.String != "" {
					item += fmt.Sprintf(" (%s)", aDadosJSON.String)
				}
				agendaList = append(agendaList, item)
			}
		}
	}

	// 2. Base Template
	templateQuery := `SELECT template FROM prompt_templates WHERE nome = 'eva_base_v2' AND ativo = true LIMIT 1`
	var template string
	if err := s.db.QueryRow(templateQuery).Scan(&template); err != nil {
		template = `Você é a EVA, assistente de saúde virtual para {{nome_idoso}}.`
	}

	// 3. Build Dossier
	dossier := fmt.Sprintf("\n\n📋 --- FICHA COMPLETA DO PACIENTE (INFORMAÇÃO CONFIDENCIAL) ---\n")
	dossier += fmt.Sprintf("NOME: %s\n", nome)
	dossier += fmt.Sprintf("IDADE: %d anos\n", idade)
	dossier += fmt.Sprintf("ENDEREÇO: %s\n", getString(endereco, "Não completado"))

	dossier += "\n🥼 --- SAÚDE E CONDIÇÕES ---\n"
	dossier += fmt.Sprintf("Nível Cognitivo: %s\n", nivelCognitivo)
	dossier += fmt.Sprintf("Mobilidade: %s\n", mobilidade)
	dossier += fmt.Sprintf("Limitações Auditivas: %v (Usa Aparelho: %v)\n", limitacoesAuditivas.Bool, usaAparelhoAuditivo.Bool)
	dossier += fmt.Sprintf("Limitações Visuais: %s\n", getString(limitacoesVisuais, "Nenhuma"))
	dossier += fmt.Sprintf("Condições Médicas: %s\n", getString(condicoesMedicas, "Nenhuma registrada"))

	dossier += "\n💊 --- MEDICAMENTOS (FONTE OFICIAL) ---\n"
	if len(medsList) > 0 {
		dossier += "O paciente possui os seguintes medicamentos prescritos e ativos no sistema:\n"
		for _, m := range medsList {
			dossier += m + "\n"
		}
		oldMeds := getString(medicamentosAtuais, "")
		if oldMeds != "" {
			dossier += fmt.Sprintf("\n(Nota de cadastro antigo: %s)\n", oldMeds)
		}
	} else {
		medsA := getString(medicamentosAtuais, "")
		if medsA == "" {
			dossier += "Nenhum medicamento registrado no sistema.\n"
		} else {
			dossier += fmt.Sprintf("Medicamentos (Legado): %s\n", medsA)
		}
	}
	dossier += "INSTRUÇÃO: Se o paciente perguntar o que deve tomar, consulte EXCLUSIVAMENTE esta lista acima.\n"

	dossier += "\n📅 --- AGENDA COMPLETA (FUTURO) ---\n"
	if len(agendaList) > 0 {
		dossier += "O paciente tem os seguintes compromissos agendados no sistema:\n"
		for _, a := range agendaList {
			dossier += a + "\n"
		}
	} else {
		dossier += "Nenhum compromisso agendado no futuro.\n"
	}

	dossier += "\n📞 --- REDE DE APOIO ---\n"
	dossier += fmt.Sprintf("Familiar: %s\n", getString(familiarPrincipal, "Não informado"))
	dossier += fmt.Sprintf("Emergência: %s\n", getString(contatoEmergencia, "Não informado"))
	dossier += fmt.Sprintf("Médico: %s\n", getString(medicoResponsavel, "Não informado"))

	dossier += "\n📝 --- OUTRAS NOTAS ---\n"
	dossier += fmt.Sprintf("Notas Gerais: %s\n", getString(notasGerais, ""))
	dossier += fmt.Sprintf("Preferência Horário: %s\n", getString(preferenciaHorario, "Indiferente"))
	dossier += fmt.Sprintf("Ambiente Ruidoso: %v\n", ambienteRuidoso.Bool)
	dossier += fmt.Sprintf("Tom de Voz Ideal: %s\n", tomVoz)
	dossier += "--------------------------------------------------------\n"

	// 4. Replacements
	replacements := map[string]string{
		"{{nome_idoso}}":        nome,
		"{{.NomeIdoso}}":        nome,
		"{{idade}}":             fmt.Sprintf("%d", idade),
		"{{.Idade}}":            fmt.Sprintf("%d", idade),
		"{{nivel_cognitivo}}":   nivelCognitivo,
		"{{.NivelCognitivo}}":   nivelCognitivo,
		"{{tom_voz}}":           tomVoz,
		"{{.TomVoz}}":           tomVoz,
		"{{condicoes_medicas}}": getString(condicoesMedicas, ""),
		"{{.CondicoesMedicas}}": getString(condicoesMedicas, ""),
	}

	instructions := template + "\n\n" + dossier
	for old, new := range replacements {
		instructions = strings.ReplaceAll(instructions, old, new)
	}

	medsString := strings.Join(medsList, ", ")
	if medsString == "" {
		medsString = getString(medicamentosAtuais, "Nenhum")
	}
	instructions = strings.ReplaceAll(instructions, "{{medicamentos}}", medsString)
	instructions = strings.ReplaceAll(instructions, "{{.MedicamentosAtuais}}", medsString)

	// Clean tags
	tags := []string{
		"{{#limitacoes_auditivas}}", "{{/limitacoes_auditivas}}",
		"{{#usa_aparelho_auditivo}}", "{{/usa_aparelho_auditivo}}",
		"{{#primeira_interacao}}", "{{/primeira_interacao}}",
		"{{^primeira_interacao}}", "{{taxa_adesao}}",
		"{{.LimitacoesAuditivas}}", "{{.UsaAparelhoAuditivo}}",
	}
	for _, tag := range tags {
		instructions = strings.ReplaceAll(instructions, tag, "")
	}

	// 5.5. Safety Protocol
	safetyProtocol := fmt.Sprintf(`
	
	🚨 PROTOCOLO DE SEGURANÇA (INTERAÇÃO MEDICAMENTOSA):
	Sempre que o paciente mencionar um novo mal-estar (ex: tontura, dor, náusea) ou perguntar sobre um novo remédio:
	1. Verifique SILENCIOSAMENTE em sua base de conhecimento se há interação perigosa com a lista de "MEDICAMENTOS (FONTE OFICIAL)" mostrada acima.
	2. Se houver qualquer risco, ALERTE IMEDIATAMENTE o paciente de forma calma mas firme.
	3. Recomende que ele NÃO tome nada sem falar com o médico responsável: %s.
	`, getString(medicoResponsavel, "médico cadastrado"))

	// 6. Zeta Story Engine
	var storySection string
	if s.personalityService != nil && s.zetaRouter != nil {
		if state, err := s.personalityService.GetState(context.Background(), idosoID); err == nil {
			profile := &types.IdosoProfile{ID: idosoID, Name: nome}
			if story, directive, err := s.zetaRouter.SelectIntervention(context.Background(), idosoID, state.DominantEmotion, profile); err == nil && story != nil {
				storySection = fmt.Sprintf(`
📚 INTERVENÇÃO NARRATIVA (ZETA ENGINE):
%s
TÍTULO: %s
CONTEÚDO: "%s"
MORAL: %s
INSTRUÇÃO: %s
`, directive, story.Title, story.Content, story.Moral, directive)
			}
		}
	}

	// 5. Agent Protocol
	agentProtocol := `
	
	IMPORTANTE - PROTOCOLO DE FERRAMENTAS:
	Você está rodando em um modelo focado em Áudio e NÃO pode executar ferramentas nativamente.
	Se você precisar realizar uma ação (Pesquisar, Agendar, Ligar) ou buscar informações externas:
	1. Avise o usuário que vai verificar: "Só um momento, vou verificar isso..." ou "Vou agendar para você, um instante...".
	2. Em seguida, GERE IMEDIATAMENTE um comando de texto oculto no formato JSON-in-TEXT:
	   [[TOOL:google_search_retrieval:{"query": "..."}]]
	   [[TOOL:schedule_appointment:{"type": "...", "description": "...", "timestamp": "..."}]]
	   [[TOOL:alert_family:{"reason": "...", "severity": "..."}]]

	NÃO invente dados. Se não souber, use o comando de busca [[TOOL:google_search_retrieval:{"query": "..."}]].
	O sistema irá processar esse comando e te devolver a resposta.
	`

	// 7. Assemble
	finalInstructions := instructions + agentProtocol + safetyProtocol + dossier + storySection

	log.Printf("✅ [Brain] Instruções finais geradas (%d chars)", len(finalInstructions))
	return finalInstructions
}

// Helper seguro para NullString
func getString(ns sql.NullString, def string) string {
	if ns.Valid && ns.String != "" {
		return ns.String
	}
	return def
}
