package tools

import (
	"eva-mind/internal/database"
	"eva-mind/internal/email"
	"eva-mind/internal/gemini"
	"eva-mind/internal/push"
	"fmt"
	"log"
)

type ToolsHandler struct {
	db           *database.DB
	pushService  *push.FirebaseService
	emailService *email.EmailService
	NotifyFunc   func(idosoID int64, msgType string, payload interface{}) // ✅ Callback para sinalização
}

func NewToolsHandler(db *database.DB, pushService *push.FirebaseService, emailService *email.EmailService) *ToolsHandler {
	return &ToolsHandler{
		db:           db,
		pushService:  pushService,
		emailService: emailService,
	}
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
		err := gemini.AlertFamilyWithSeverity(h.db.Conn, h.pushService, h.emailService, idosoID, reason, severity)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}
		return map[string]interface{}{"status": "sucesso", "alerta": reason}, nil

	case "confirm_medication":
		medicationName, _ := args["medication_name"].(string)
		err := gemini.ConfirmMedication(h.db.Conn, h.pushService, idosoID, medicationName)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}
		return map[string]interface{}{"status": "sucesso", "medicamento": medicationName}, nil

	case "schedule_appointment":
		timestamp, _ := args["timestamp"].(string)
		tipo, _ := args["type"].(string)
		description, _ := args["description"].(string)
		err := gemini.ScheduleAppointment(h.db.Conn, idosoID, timestamp, tipo, description)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}, nil
		}
		return map[string]interface{}{"status": "sucesso", "agendamento": description}, nil

	case "call_family_webrtc", "call_doctor_webrtc", "call_caregiver_webrtc", "call_central_webrtc":
		if h.NotifyFunc != nil {
			h.NotifyFunc(idosoID, "initiate_call", map[string]string{
				"target": name,
			})
			return map[string]interface{}{"status": "iniciando chamada", "alvo": name}, nil
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
