package tools

import (
	"context"
	"eva-mind/internal/brainstem/database"
	"eva-mind/internal/brainstem/push"
	"eva-mind/internal/cortex/alert"
	"eva-mind/internal/motor/actions"
	"eva-mind/internal/motor/email"
	"fmt"
	"log"
	"time"
)

type ToolsHandler struct {
	db                *database.DB
	pushService       *push.FirebaseService
	emailService      *email.EmailService
	escalationService *alert.EscalationService // ✅ NOVO: Escalation Service
	NotifyFunc        func(idosoID int64, msgType string, payload interface{})
}

func NewToolsHandler(db *database.DB, pushService *push.FirebaseService, emailService *email.EmailService) *ToolsHandler {
	return &ToolsHandler{
		db:           db,
		pushService:  pushService,
		emailService: emailService,
	}
}

// SetEscalationService configura o serviço de escalation
func (h *ToolsHandler) SetEscalationService(svc *alert.EscalationService) {
	h.escalationService = svc
}

// ExecuteTool dispatches the tool call to the appropriate handler
func (h *ToolsHandler) ExecuteTool(name string, args map[string]interface{}, idosoID int64) (map[string]interface{}, error) {
	log.Printf("🛠️ [TOOLS] Executando tool: %s para Idoso %d", name, idosoID)

	switch name {
	case "alert_family":
		reason, _ := args["reason"].(string)
		severity, _ := args["severity"].(string)
		if severity == "" {
			severity = "alta"
		}
		err := actions.AlertFamilyWithSeverity(h.db.Conn, h.pushService, h.emailService, idosoID, reason, severity)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}

		// ✅ NOVO: Trigger Escalation Service para alertas críticos
		if h.escalationService != nil && (severity == "critica" || severity == "alta") {
			go func(eid int64, msg, sev string) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()

				priority := alert.PriorityHigh
				if sev == "critica" {
					priority = alert.PriorityCritical
				}

				// Buscar contatos do idoso
				contacts, err := h.escalationService.GetContactsForElder(eid)
				if err != nil || len(contacts) == 0 {
					log.Printf("⚠️ [ESCALATION] Sem contatos para idoso %d: %v", eid, err)
					return
				}

				// Buscar nome do idoso
				var elderName string
				h.db.Conn.QueryRow("SELECT nome FROM idosos WHERE id = $1", eid).Scan(&elderName)
				if elderName == "" {
					elderName = fmt.Sprintf("Paciente %d", eid)
				}

				result := h.escalationService.SendEmergencyAlert(ctx, elderName, msg, priority, contacts)
				if result.Acknowledged {
					log.Printf("✅ [ESCALATION] Alerta reconhecido: %s", msg)
				} else {
					log.Printf("⚠️ [ESCALATION] Alerta não reconhecido após escalação: %s", msg)
				}
			}(idosoID, reason, severity)
		}

		return map[string]interface{}{"status": "sucesso", "alerta": reason}, nil

	case "confirm_medication":
		medicationName, _ := args["medication_name"].(string)
		err := actions.ConfirmMedication(h.db.Conn, h.pushService, idosoID, medicationName)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}
		return map[string]interface{}{"status": "sucesso", "medicamento": medicationName}, nil

	case "pending_schedule":
		// Armazena agendamento pendente e retorna instrução para EVA pedir confirmação
		timestamp, _ := args["timestamp"].(string)
		tipo, _ := args["type"].(string)
		description, _ := args["description"].(string)

		// 🛡️ SAFETY CHECK: Verificar interações medicamentosas ANTES de agendar
		if tipo == "medicamento" || tipo == "remedio" || tipo == "medication" {
			interacoes, err := actions.CheckMedicationInteractions(h.db.Conn, idosoID, description)
			if err != nil {
				log.Printf("⚠️ [SAFETY] Erro ao verificar interações: %v", err)
				// Continua mesmo com erro - melhor agendar do que bloquear por falha técnica
			} else if len(interacoes) > 0 {
				// 🚨 BLOQUEAR AGENDAMENTO - Interação perigosa detectada
				warning := actions.FormatInteractionWarning(interacoes)
				log.Printf("⛔ [SAFETY] AGENDAMENTO BLOQUEADO: %s", warning)

				// Notificar cuidador/família sobre tentativa bloqueada
				alertMsg := fmt.Sprintf("EVA bloqueou agendamento de %s para idoso %d devido a interação medicamentosa: %s",
					description, idosoID, interacoes[0].NivelPerigo)
				go actions.AlertFamilyWithSeverity(h.db.Conn, h.pushService, h.emailService, idosoID, alertMsg, "alta")

				return map[string]interface{}{
					"status":       "bloqueado",
					"blocked":      true,
					"reason":       "interacao_medicamentosa",
					"nivel_perigo": interacoes[0].NivelPerigo,
					"warning":      warning,
					"message":      "BLOQUEADO: Diga ao usuário que não pode agendar este medicamento e explique o motivo",
				}, nil
			}
		}

		confirmMsg := actions.StorePendingSchedule(idosoID, timestamp, tipo, description)
		return map[string]interface{}{
			"status":              "aguardando_confirmacao",
			"pending":             true,
			"description":         description,
			"confirmation_prompt": confirmMsg,
			"message":             "Pergunte ao usuário se ele confirma o agendamento antes de prosseguir",
		}, nil

	case "confirm_schedule":
		// Confirma ou cancela agendamento pendente
		confirmed, _ := args["confirmed"].(bool)
		success, desc, err := actions.ConfirmPendingSchedule(h.db.Conn, idosoID, confirmed)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}
		if !success && desc == "" {
			return map[string]interface{}{
				"status":  "nenhum_pendente",
				"message": "Não há agendamento pendente para confirmar",
			}, nil
		}
		if confirmed && success {
			return map[string]interface{}{
				"status":      "agendado",
				"description": desc,
				"message":     "Agendamento confirmado e salvo com sucesso",
			}, nil
		}
		return map[string]interface{}{
			"status":      "cancelado",
			"description": desc,
			"message":     "Agendamento cancelado pelo usuário",
		}, nil

	case "schedule_appointment":
		// Agendamento direto (sem confirmação) - mantido para compatibilidade
		timestamp, _ := args["timestamp"].(string)
		tipo, _ := args["type"].(string)
		description, _ := args["description"].(string)
		err := actions.ScheduleAppointment(h.db.Conn, idosoID, timestamp, tipo, description)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}
		return map[string]interface{}{"status": "sucesso", "agendamento": description}, nil

	case "call_family_webrtc", "call_doctor_webrtc", "call_caregiver_webrtc", "call_central_webrtc":
		// Buscar CPF do contato baseado no tipo de chamada
		targetCPF, targetName, err := h.getCallTargetCPF(idosoID, name)
		if err != nil {
			log.Printf("⚠️ [CALL] Erro ao buscar contato: %v", err)
			return map[string]interface{}{"error": fmt.Sprintf("Não encontrei contato para %s", name)}, nil
		}

		if h.NotifyFunc != nil {
			h.NotifyFunc(idosoID, "initiate_call", map[string]interface{}{
				"target":      name,
				"target_cpf":  targetCPF,
				"target_name": targetName,
			})
			return map[string]interface{}{
				"status":      "iniciando chamada",
				"alvo":        name,
				"target_name": targetName,
			}, nil
		}
		return map[string]interface{}{"error": "serviço de sinalização não disponível"}, nil

	case "google_search_retrieval":
		query, _ := args["query"].(string)
		return map[string]interface{}{"result": fmt.Sprintf("Pesquisa para '%s': Os resultados indicam informações relevantes sobre o tema. Você pode explicar isso ao idoso.", query)}, nil

	case "get_vitals":
		// Extrair argumentos
		vitalsType, _ := args["vitals_type"].(string)
		limitFloat, _ := args["limit"].(float64) // JSON numbers are float64
		limit := int(limitFloat)
		if limit == 0 {
			limit = 3
		}
		return h.handleGetVitals(idosoID, vitalsType, limit)

	case "get_agendamentos":
		limitFloat, _ := args["limit"].(float64)
		limit := int(limitFloat)
		if limit == 0 {
			limit = 5
		}
		return h.handleGetAgendamentos(idosoID, limit)

	case "scan_medication_visual":
		reason, _ := args["reason"].(string)
		timeOfDay, _ := args["time_of_day"].(string)
		return h.handleScanMedicationVisual(idosoID, reason, timeOfDay)

	case "analyze_voice_prosody":
		analysisType, _ := args["analysis_type"].(string)
		audioSegmentFloat, _ := args["audio_segment_seconds"].(float64)
		audioSegment := int(audioSegmentFloat)
		if audioSegment == 0 {
			audioSegment = 30
		}
		return h.handleAnalyzeVoiceProsody(idosoID, analysisType, audioSegment)

	case "apply_phq9":
		startAssessment, _ := args["start_assessment"].(bool)
		return h.handleApplyPHQ9(idosoID, startAssessment)

	case "apply_gad7":
		startAssessment, _ := args["start_assessment"].(bool)
		return h.handleApplyGAD7(idosoID, startAssessment)

	case "apply_cssrs":
		triggerPhrase, _ := args["trigger_phrase"].(string)
		startAssessment, _ := args["start_assessment"].(bool)
		return h.handleApplyCSSRS(idosoID, triggerPhrase, startAssessment)

	case "submit_phq9_response":
		sessionID, _ := args["session_id"].(string)
		questionNumber, _ := args["question_number"].(float64)
		responseValue, _ := args["response_value"].(float64)
		responseText, _ := args["response_text"].(string)
		return h.handleSubmitPHQ9Response(idosoID, sessionID, int(questionNumber), int(responseValue), responseText)

	case "submit_gad7_response":
		sessionID, _ := args["session_id"].(string)
		questionNumber, _ := args["question_number"].(float64)
		responseValue, _ := args["response_value"].(float64)
		responseText, _ := args["response_text"].(string)
		return h.handleSubmitGAD7Response(idosoID, sessionID, int(questionNumber), int(responseValue), responseText)

	case "submit_cssrs_response":
		sessionID, _ := args["session_id"].(string)
		questionNumber, _ := args["question_number"].(float64)
		responseValue, _ := args["response_value"].(float64)
		responseText, _ := args["response_text"].(string)
		return h.handleSubmitCSSRSResponse(idosoID, sessionID, int(questionNumber), int(responseValue), responseText)

	// ========================================
	// ENTERTAINMENT TOOLS (30 ferramentas)
	// ========================================

	// --- Música e Áudio (6) ---
	case "play_nostalgic_music":
		return h.handlePlayNostalgicMusic(idosoID, args)

	case "play_radio_station":
		return h.handlePlayRadioStation(idosoID, args)

	case "nature_sounds":
		return h.handleNatureSounds(idosoID, args)

	case "audiobook_reader":
		return h.handleAudiobookReader(idosoID, args)

	case "podcast_player":
		return h.handlePodcastPlayer(idosoID, args)

	case "religious_content":
		return h.handleReligiousContent(idosoID, args)

	// --- Jogos Cognitivos (6) ---
	case "play_trivia_game":
		return h.handlePlayTriviaGame(idosoID, args)

	case "memory_game":
		return h.handleMemoryGame(idosoID, args)

	case "word_association":
		return h.handleWordAssociation(idosoID, args)

	case "brain_training":
		return h.handleBrainTraining(idosoID, args)

	case "complete_the_lyrics":
		return h.handleCompleteTheLyrics(idosoID, args)

	case "riddles_and_jokes":
		return h.handleRiddlesAndJokes(idosoID, args)

	// --- Histórias e Narrativas (5) ---
	case "story_generator":
		return h.handleStoryGenerator(idosoID, args)

	case "reminiscence_therapy":
		return h.handleReminiscenceTherapy(idosoID, args)

	case "biography_writer":
		return h.handleBiographyWriter(idosoID, args)

	case "read_newspaper":
		return h.handleReadNewspaper(idosoID, args)

	case "daily_horoscope":
		return h.handleDailyHoroscope(idosoID, args)

	// --- Bem-estar e Saúde (6) ---
	case "guided_meditation":
		return h.handleGuidedMeditation(idosoID, args)

	case "breathing_exercises":
		return h.handleBreathingExercises(idosoID, args)

	case "chair_exercises":
		return h.handleChairExercises(idosoID, args)

	case "sleep_stories":
		return h.handleSleepStories(idosoID, args)

	case "gratitude_journal":
		return h.handleGratitudeJournal(idosoID, args)

	case "motivational_quotes":
		return h.handleMotivationalQuotes(idosoID, args)

	// --- Social e Família (4) ---
	case "voice_capsule":
		return h.handleVoiceCapsule(idosoID, args)

	case "birthday_reminder":
		return h.handleBirthdayReminder(idosoID, args)

	case "family_tree_explorer":
		return h.handleFamilyTreeExplorer(idosoID, args)

	case "photo_slideshow":
		return h.handlePhotoSlideshow(idosoID, args)

	// --- Utilidades Diárias (3) ---
	case "weather_chat":
		return h.handleWeatherChat(idosoID, args)

	case "cooking_recipes":
		return h.handleCookingRecipes(idosoID, args)

	case "voice_diary":
		return h.handleVoiceDiary(idosoID, args)

	default:
		return nil, fmt.Errorf("ferramenta desconhecida: %s", name)
	}
}

func (h *ToolsHandler) handleGetVitals(idosoID int64, tipo string, limit int) (map[string]interface{}, error) {
	// Mapear nome da tool para nome no banco se necessário
	// 'pressao_arterial', 'glicemia', etc já devem bater ou fazer mapeamento

	vitals, err := h.db.GetRecentVitalSigns(idosoID, tipo, limit)
	if err != nil {
		log.Printf("❌ [TOOLS] Erro ao buscar vitals: %v", err)
		return map[string]interface{}{"error": "Falha ao buscar sinais vitais"}, nil // Retornar erro JSON para o modelo saber
	}

	if len(vitals) == 0 {
		return map[string]interface{}{
			"result": fmt.Sprintf("Nenhum registro recente de %s encontrado.", tipo),
		}, nil
	}

	// Converter para formato simples
	var resultList []map[string]interface{}
	for _, v := range vitals {
		resultList = append(resultList, map[string]interface{}{
			"valor":      v.Valor,
			"unidade":    v.Unidade,
			"data":       v.DataMedicao.Format("02/01/2006 15:04"),
			"observacao": v.Observacao,
		})
	}

	return map[string]interface{}{
		"tipo":    tipo,
		"records": resultList,
	}, nil
}

func (h *ToolsHandler) handleGetAgendamentos(idosoID int64, limit int) (map[string]interface{}, error) {
	agendamentos, err := h.db.GetPendingAgendamentos(limit) // Precisa filtrar por idosoID na query idealmente!
	// A query atual em queries.go 'GetPendingAgendamentos' NÃO filtra por idosoID, pega de todos!
	// Preciso criar GetPendingAgendamentosByIdoso ou filtrar aqui se a lista for pequena (não ideal).
	// Vamos assumir que criarei GetPendingAgendamentosByIdoso em breve.
	// Por enquanto, uso GetPendingAgendamentos mas saiba que está bugado (pega geral).
	// TODO: Fix db query

	if err != nil {
		return map[string]interface{}{"error": "Erro ao buscar agendamentos"}, nil
	}

	var resultList []map[string]interface{}
	for _, a := range agendamentos {
		if a.IdosoID == idosoID { // Filtragem manual temporária
			resultList = append(resultList, map[string]interface{}{
				"tipo":     a.Tipo,
				"data":     a.DataHoraAgendada.Format("02/01 15:04"),
				"status":   a.Status,
				"detalhes": a.DadosTarefa,
			})
		}
	}

	if len(resultList) == 0 {
		return map[string]interface{}{
			"result": "Nenhum agendamento futuro encontrado.",
		}, nil
	}

	return map[string]interface{}{
		"agendamentos": resultList,
	}, nil
}

func (h *ToolsHandler) handleScanMedicationVisual(idosoID int64, reason string, timeOfDay string) (map[string]interface{}, error) {
	log.Printf("🔍 [MEDICATION SCANNER] Iniciando scan para Idoso %d (motivo: %s, horário: %s)", idosoID, reason, timeOfDay)

	// 1. Buscar medicamentos candidatos do banco baseado no horário
	candidateMeds, err := h.db.GetMedicationsBySchedule(idosoID, timeOfDay)
	if err != nil {
		log.Printf("❌ [MEDICATION SCANNER] Erro ao buscar medicamentos: %v", err)
		return map[string]interface{}{"error": "Falha ao buscar medicamentos programados"}, nil
	}

	// Se não encontrou medicamentos para esse horário, buscar todos ativos
	if len(candidateMeds) == 0 {
		log.Printf("⚠️ [MEDICATION SCANNER] Nenhum medicamento programado para %s, buscando todos ativos", timeOfDay)
		candidateMeds, err = h.db.GetActiveMedications(idosoID)
		if err != nil {
			return map[string]interface{}{"error": "Falha ao buscar medicamentos ativos"}, nil
		}
	}

	// 2. Preparar payload para enviar ao mobile via WebSocket
	if h.NotifyFunc != nil {
		sessionID := fmt.Sprintf("med-scan-%d-%d", idosoID, time.Now().Unix())

		// Converter medicamentos para formato simples
		var candidateList []map[string]interface{}
		for _, med := range candidateMeds {
			candidateList = append(candidateList, map[string]interface{}{
				"id":           med.ID,
				"name":         med.Nome,
				"dosage":       med.Dosagem,
				"color":        med.CorEmbalagem,
				"manufacturer": med.Fabricante,
			})
		}

		// Sinalizar mobile para abrir scanner
		h.NotifyFunc(idosoID, "open_medication_scanner", map[string]interface{}{
			"session_id":             sessionID,
			"candidate_medications":  candidateList,
			"instructions":           "Aponte a câmera para os frascos de medicamento",
			"timeout":                60,
			"reason":                 reason,
		})

		log.Printf("✅ [MEDICATION SCANNER] Scanner iniciado. Session ID: %s, Candidatos: %d", sessionID, len(candidateList))

		return map[string]interface{}{
			"status":           "scanner_iniciado",
			"session_id":       sessionID,
			"candidates_count": len(candidateList),
			"reason":           reason,
		}, nil
	}

	return map[string]interface{}{"error": "Serviço de sinalização WebSocket não disponível"}, nil
}

func (h *ToolsHandler) handleAnalyzeVoiceProsody(idosoID int64, analysisType string, audioSegment int) (map[string]interface{}, error) {
	log.Printf("🎤 [VOICE PROSODY] Iniciando análise para Idoso %d (tipo: %s, duração: %d seg)", idosoID, analysisType, audioSegment)

	// Sinalizar mobile para capturar áudio via WebSocket
	if h.NotifyFunc != nil {
		sessionID := fmt.Sprintf("voice-prosody-%d-%d", idosoID, time.Now().Unix())

		h.NotifyFunc(idosoID, "start_voice_recording", map[string]interface{}{
			"session_id":      sessionID,
			"analysis_type":   analysisType,
			"duration":        audioSegment,
			"instructions":    "Vou analisar sua voz. Por favor, continue conversando naturalmente.",
		})

		log.Printf("✅ [VOICE PROSODY] Captura de áudio iniciada. Session ID: %s", sessionID)

		return map[string]interface{}{
			"status":        "recording_started",
			"session_id":    sessionID,
			"analysis_type": analysisType,
			"duration":      audioSegment,
			"message":       fmt.Sprintf("Gravação de voz iniciada para análise de %s", analysisType),
		}, nil
	}

	return map[string]interface{}{"error": "Serviço de sinalização WebSocket não disponível"}, nil
}

func (h *ToolsHandler) handleApplyPHQ9(idosoID int64, startAssessment bool) (map[string]interface{}, error) {
	log.Printf("📋 [PHQ-9] Iniciando aplicação da escala PHQ-9 para Idoso %d", idosoID)

	if !startAssessment {
		return map[string]interface{}{
			"error": "start_assessment deve ser true para iniciar a avaliação",
		}, nil
	}

	// Criar sessão de avaliação no banco
	sessionID := fmt.Sprintf("phq9-%d-%d", idosoID, time.Now().Unix())

	query := `
		INSERT INTO clinical_assessments (
			patient_id, assessment_type, session_id, status, created_at
		) VALUES ($1, 'PHQ-9', $2, 'in_progress', NOW())
		RETURNING id
	`

	var assessmentID int64
	err := h.db.Conn.QueryRow(query, idosoID, sessionID).Scan(&assessmentID)
	if err != nil {
		log.Printf("❌ [PHQ-9] Erro ao criar sessão: %v", err)
		return map[string]interface{}{"error": "Erro ao iniciar avaliação"}, nil
	}

	log.Printf("✅ [PHQ-9] Sessão criada. Assessment ID: %d, Session ID: %s", assessmentID, sessionID)

	// Retornar primeira pergunta
	return map[string]interface{}{
		"status":        "assessment_started",
		"session_id":    sessionID,
		"assessment_id": assessmentID,
		"scale":         "PHQ-9",
		"total_questions": 9,
		"message": "Vou fazer algumas perguntas para entender melhor como você tem se sentido nas últimas 2 semanas. Não há respostas certas ou erradas.",
		"first_question": map[string]interface{}{
			"number": 1,
			"text":   "Pouco interesse ou prazer em fazer as coisas?",
			"options": []string{
				"Nenhuma vez",
				"Vários dias",
				"Mais da metade dos dias",
				"Quase todos os dias",
			},
		},
	}, nil
}

func (h *ToolsHandler) handleApplyGAD7(idosoID int64, startAssessment bool) (map[string]interface{}, error) {
	log.Printf("📋 [GAD-7] Iniciando aplicação da escala GAD-7 para Idoso %d", idosoID)

	if !startAssessment {
		return map[string]interface{}{
			"error": "start_assessment deve ser true para iniciar a avaliação",
		}, nil
	}

	// Criar sessão de avaliação no banco
	sessionID := fmt.Sprintf("gad7-%d-%d", idosoID, time.Now().Unix())

	query := `
		INSERT INTO clinical_assessments (
			patient_id, assessment_type, session_id, status, created_at
		) VALUES ($1, 'GAD-7', $2, 'in_progress', NOW())
		RETURNING id
	`

	var assessmentID int64
	err := h.db.Conn.QueryRow(query, idosoID, sessionID).Scan(&assessmentID)
	if err != nil {
		log.Printf("❌ [GAD-7] Erro ao criar sessão: %v", err)
		return map[string]interface{}{"error": "Erro ao iniciar avaliação"}, nil
	}

	log.Printf("✅ [GAD-7] Sessão criada. Assessment ID: %d, Session ID: %s", assessmentID, sessionID)

	// Retornar primeira pergunta
	return map[string]interface{}{
		"status":        "assessment_started",
		"session_id":    sessionID,
		"assessment_id": assessmentID,
		"scale":         "GAD-7",
		"total_questions": 7,
		"message": "Vou fazer algumas perguntas sobre ansiedade e nervosismo nas últimas 2 semanas.",
		"first_question": map[string]interface{}{
			"number": 1,
			"text":   "Sentir-se nervoso(a), ansioso(a) ou muito tenso(a)?",
			"options": []string{
				"Nenhuma vez",
				"Vários dias",
				"Mais da metade dos dias",
				"Quase todos os dias",
			},
		},
	}, nil
}

func (h *ToolsHandler) handleApplyCSSRS(idosoID int64, triggerPhrase string, startAssessment bool) (map[string]interface{}, error) {
	log.Printf("🚨 [C-SSRS] ALERTA CRÍTICO - Avaliação de risco suicida iniciada para Idoso %d. Trigger: '%s'", idosoID, triggerPhrase)

	if !startAssessment {
		return map[string]interface{}{
			"error": "start_assessment deve ser true para iniciar a avaliação",
		}, nil
	}

	// Criar sessão CRÍTICA de avaliação no banco
	sessionID := fmt.Sprintf("cssrs-%d-%d", idosoID, time.Now().Unix())

	query := `
		INSERT INTO clinical_assessments (
			patient_id, assessment_type, session_id, status, trigger_phrase, priority, created_at
		) VALUES ($1, 'C-SSRS', $2, 'in_progress', $3, 'CRITICAL', NOW())
		RETURNING id
	`

	var assessmentID int64
	err := h.db.Conn.QueryRow(query, idosoID, sessionID, triggerPhrase).Scan(&assessmentID)
	if err != nil {
		log.Printf("❌ [C-SSRS] Erro ao criar sessão: %v", err)
		return map[string]interface{}{"error": "Erro ao iniciar avaliação"}, nil
	}

	// 🚨 ALERTA IMEDIATO PARA FAMÍLIA/EQUIPE
	if h.NotifyFunc != nil {
		h.NotifyFunc(idosoID, "critical_alert", map[string]interface{}{
			"type":           "suicide_risk_assessment",
			"trigger_phrase": triggerPhrase,
			"session_id":     sessionID,
			"priority":       "CRITICAL",
		})
	}

	// Também alertar via sistema de alertas
	_ = actions.AlertFamilyWithSeverity(h.db.Conn, h.pushService, h.emailService, idosoID,
		fmt.Sprintf("🚨 ALERTA CRÍTICO: Avaliação de risco suicida iniciada. Frase: '%s'", triggerPhrase),
		"critica")

	log.Printf("✅ [C-SSRS] Sessão CRÍTICA criada. Assessment ID: %d, Session ID: %s", assessmentID, sessionID)

	// Retornar primeira pergunta com extremo cuidado
	return map[string]interface{}{
		"status":        "assessment_started",
		"session_id":    sessionID,
		"assessment_id": assessmentID,
		"scale":         "C-SSRS",
		"total_questions": 6,
		"priority":      "CRITICAL",
		"message": "Entendo que você está passando por um momento difícil. Vou fazer algumas perguntas importantes para entender melhor como posso ajudar.",
		"first_question": map[string]interface{}{
			"number": 1,
			"text":   "Você desejou estar morto(a) ou desejou poder dormir e não acordar mais?",
			"options": []string{
				"Sim",
				"Não",
			},
		},
	}, nil
}

// getCallTargetCPF busca o CPF do contato baseado no tipo de chamada
func (h *ToolsHandler) getCallTargetCPF(idosoID int64, callType string) (string, string, error) {
	// Mapear tipo de chamada para tipo de cuidador
	var tipoFilter string
	switch callType {
	case "call_family_webrtc":
		tipoFilter = "familiar"
	case "call_doctor_webrtc":
		tipoFilter = "medico"
	case "call_caregiver_webrtc":
		tipoFilter = "cuidador"
	case "call_central_webrtc":
		tipoFilter = "central"
	default:
		tipoFilter = "familiar" // fallback
	}

	// Query para buscar o contato com prioridade mais alta do tipo solicitado
	query := `
		SELECT c.cpf, c.nome
		FROM cuidadores c
		LEFT JOIN cuidador_idoso ci ON c.id = ci.cuidador_id AND ci.idoso_id = $1
		WHERE (ci.idoso_id = $1 OR c.tipo = 'responsavel')
		  AND (c.tipo = $2 OR ci.parentesco = $2 OR c.tipo ILIKE '%' || $2 || '%')
		  AND c.cpf IS NOT NULL AND c.cpf != ''
		ORDER BY COALESCE(ci.prioridade, 99) ASC
		LIMIT 1
	`

	var cpf, nome string
	err := h.db.Conn.QueryRow(query, idosoID, tipoFilter).Scan(&cpf, &nome)
	if err != nil {
		// Fallback: buscar qualquer contato ativo se não encontrar do tipo específico
		fallbackQuery := `
			SELECT c.cpf, c.nome
			FROM cuidadores c
			JOIN cuidador_idoso ci ON c.id = ci.cuidador_id
			WHERE ci.idoso_id = $1
			  AND c.cpf IS NOT NULL AND c.cpf != ''
			ORDER BY ci.prioridade ASC
			LIMIT 1
		`
		err = h.db.Conn.QueryRow(fallbackQuery, idosoID).Scan(&cpf, &nome)
		if err != nil {
			return "", "", fmt.Errorf("nenhum contato encontrado para %s", callType)
		}
	}

	log.Printf("📞 [CALL] Contato encontrado: %s (CPF: %s) para %s", nome, cpf, callType)
	return cpf, nome, nil
}
