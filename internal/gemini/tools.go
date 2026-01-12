package gemini

import (
	"context"
	"database/sql"
	"eva-mind/internal/push"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/v4/messaging"
)

func GetDefaultTools() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"function_declarations": []interface{}{
				map[string]interface{}{
					"name":        "alert_family",
					"description": "Alerta a família em caso de emergência detectada na conversa com o idoso",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"reason": map[string]interface{}{
								"type":        "string",
								"description": "Motivo do alerta (ex: 'Paciente relatou dor no peito', 'Idoso parece confuso')",
							},
							"severity": map[string]interface{}{
								"type":        "string",
								"description": "Severidade do alerta: critica, alta, media, baixa",
								"enum":        []string{"critica", "alta", "media", "baixa"},
							},
						},
						"required": []string{"reason"},
					},
				},
				map[string]interface{}{
					"name":        "confirm_medication",
					"description": "Confirma que o idoso tomou o remédio",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"medication_name": map[string]interface{}{
								"type":        "string",
								"description": "Nome do medicamento tomado",
							},
						},
						"required": []string{"medication_name"},
					},
				},
				map[string]interface{}{
					"name":        "schedule_appointment",
					"description": "Agenda um compromisso, consulta, medicamento ou chamada para o idoso",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"timestamp": map[string]interface{}{
								"type":        "string",
								"description": "Data e hora do agendamento no formato ISO 8601 (ex: 2024-02-25T14:30:00Z)",
							},
							"type": map[string]interface{}{
								"type":        "string",
								"description": "Tipo do agendamento",
								"enum":        []string{"consulta", "medicamento", "ligacao", "atividade", "outro"},
							},
							"description": map[string]interface{}{
								"type":        "string",
								"description": "Descrição detalhada do compromisso ou tarefa",
							},
						},
						"required": []string{"timestamp", "type", "description"},
					},
				},
				map[string]interface{}{
					"name":        "call_family_webrtc",
					"description": "Inicia uma chamada de vídeo para a família do idoso",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				map[string]interface{}{
					"name":        "call_central_webrtc",
					"description": "Inicia uma chamada de vídeo de emergência para a Central EVA-Mind",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				map[string]interface{}{
					"name":        "call_doctor_webrtc",
					"description": "Inicia uma chamada de vídeo para o médico responsável",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				map[string]interface{}{
					"name":        "call_caregiver_webrtc",
					"description": "Inicia uma chamada de vídeo para o cuidador",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				map[string]interface{}{
					"name":        "open_camera_analysis",
					"description": "Ativa a câmera do dispositivo do idoso para analisar visualmente um objeto, remédio ou ambiente",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
			},
		},
	}
}

// AlertFamily envia notificação push para cuidadores com sistema de fallback
func AlertFamily(db *sql.DB, pushService *push.FirebaseService, idosoID int64, reason string) error {
	return AlertFamilyWithSeverity(db, pushService, idosoID, reason, "alta")
}

// AlertFamilyWithSeverity envia alertas com níveis de severidade
func AlertFamilyWithSeverity(db *sql.DB, pushService *push.FirebaseService, idosoID int64, reason, severity string) error {
	// 1. Buscar todos os cuidadores ativos (primários e secundários)
	query := `
		SELECT 
			c.device_token, 
			c.telefone,
			c.email,
			c.prioridade,
			i.nome 
		FROM cuidadores c
		JOIN idosos i ON i.id = c.idoso_id
		WHERE c.idoso_id = $1 AND c.ativo = true
		ORDER BY c.prioridade ASC
	`

	rows, err := db.Query(query, idosoID)
	if err != nil {
		return fmt.Errorf("failed to query caregivers: %w", err)
	}
	defer rows.Close()

	type Caregiver struct {
		Token     sql.NullString
		Phone     sql.NullString
		Email     sql.NullString
		Priority  int
		ElderName string
	}

	var caregivers []Caregiver

	for rows.Next() {
		var cg Caregiver
		err := rows.Scan(&cg.Token, &cg.Phone, &cg.Email, &cg.Priority, &cg.ElderName)
		if err != nil {
			log.Printf("Error scanning caregiver: %v", err)
			continue
		}
		caregivers = append(caregivers, cg)
	}

	if len(caregivers) == 0 {
		log.Printf("⚠️ No active caregivers found for idoso %d", idosoID)
		return fmt.Errorf("no caregivers registered")
	}

	elderName := caregivers[0].ElderName

	// 2. Registrar alerta no banco ANTES de enviar
	var alertID int64
	insertQuery := `
		INSERT INTO alertas (
			idoso_id, 
			tipo, 
			severidade,
			mensagem, 
			visualizado,
			criado_em
		) 
		VALUES ($1, 'familia', $2, $3, false, NOW())
		RETURNING id
	`

	err = db.QueryRow(insertQuery, idosoID, severity, reason).Scan(&alertID)
	if err != nil {
		log.Printf("⚠️ Failed to log alert in database: %v", err)
	} else {
		log.Printf("📝 Alert registered in DB with ID: %d", alertID)
	}

	// 3. Tentar enviar push notifications para todos os cuidadores
	var successCount int
	var tokens []string

	for _, cg := range caregivers {
		if cg.Token.Valid && cg.Token.String != "" {
			tokens = append(tokens, cg.Token.String)
		}
	}

	if len(tokens) > 0 {
		log.Printf("📱 Enviando push para %d cuidador(es)", len(tokens))

		for _, token := range tokens {
			result, err := pushService.SendAlertNotification(token, elderName, reason)

			if err == nil && result.Success {
				successCount++

				// Registrar envio no banco
				_, _ = db.Exec(`
					UPDATE alertas 
					SET enviado = true, data_envio = NOW()
					WHERE id = $1
				`, alertID)

				log.Printf("✅ Alert sent successfully to caregiver for %s", elderName)
			} else {
				log.Printf("❌ Failed to send alert to caregiver: %v", err)
			}
		}
	}

	// 4. Se NENHUM push funcionou, tentar fallbacks
	if successCount == 0 {
		log.Printf("⚠️ Nenhum push notification enviado com sucesso. Tentando fallbacks...")

		// Registrar que o alerta precisa de escalamento
		_, _ = db.Exec(`
			UPDATE alertas 
			SET 
				necessita_escalamento = true,
				tentativas_envio = tentativas_envio + 1,
				ultima_tentativa = NOW()
			WHERE id = $1
		`, alertID)

		// TODO: Implementar SMS via Twilio
		// TODO: Implementar Email
		// TODO: Implementar ligação telefônica para casos críticos

		return fmt.Errorf("all push notifications failed, alert needs escalation")
	}

	log.Printf("✅ Alert sent to %d of %d caregivers", successCount, len(tokens))

	// 5. Para alertas críticos, marcar para escalonamento automático
	if severity == "critica" {
		_, _ = db.Exec(`
			UPDATE alertas 
			SET 
				necessita_escalamento = true,
				tempo_escalamento = NOW() + INTERVAL '5 minutes'
			WHERE id = $1
		`, alertID)

		log.Printf("🚨 Alert crítico - configurado para escalonamento em 5 minutos se não visualizado")
	}

	return nil
}

// ConfirmMedication registra que o idoso tomou o remédio
func ConfirmMedication(db *sql.DB, pushService *push.FirebaseService, idosoID int64, medicationName string) error {
	// 1. Registrar no histórico
	_, err := db.Exec(`
		INSERT INTO historico_medicamentos (idoso_id, medicamento, tomado_em) 
		VALUES ($1, $2, NOW())
	`, idosoID, medicationName)

	if err != nil {
		return fmt.Errorf("failed to log medication: %w", err)
	}

	log.Printf("💊 Medication logged: %d took %s", idosoID, medicationName)

	// 2. Atualizar status do agendamento de hoje
	_, err = db.Exec(`
		UPDATE agendamentos 
		SET medicamento_confirmado = true, 
		    status = 'concluido'
		WHERE idoso_id = $1 
		  AND DATE(data_hora_agendada) = CURRENT_DATE
		  AND status = 'em_andamento'
	`, idosoID)

	if err != nil {
		log.Printf("⚠️ Failed to update schedule: %v", err)
	}

	// 3. Notificar TODOS os cuidadores ativos
	query := `
		SELECT c.device_token, i.nome 
		FROM cuidadores c
		JOIN idosos i ON i.id = c.idoso_id
		WHERE c.idoso_id = $1 AND c.ativo = true
	`

	rows, err := db.Query(query, idosoID)
	if err != nil {
		log.Printf("⚠️ Failed to query caregivers: %v", err)
		return nil
	}
	defer rows.Close()

	var elderName string
	notificationsSent := 0

	for rows.Next() {
		var token sql.NullString
		err := rows.Scan(&token, &elderName)

		if err != nil || !token.Valid || token.String == "" {
			continue
		}

		message := &messaging.Message{
			Token: token.String,
			Notification: &messaging.Notification{
				Title: "✅ Medicamento Confirmado",
				Body:  fmt.Sprintf("%s tomou %s", elderName, medicationName),
			},
			Data: map[string]string{
				"type":       "medication_confirmed",
				"medication": medicationName,
				"idosoId":    fmt.Sprintf("%d", idosoID),
				"timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
			},
			Android: &messaging.AndroidConfig{
				Priority: "normal",
				Notification: &messaging.AndroidNotification{
					Sound:        "default",
					ChannelID:    "eva_medications",
					DefaultSound: true,
					Color:        "#00FF00",
				},
			},
		}

		// ✅ Criar contexto local com timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = pushService.GetClient().Send(ctx, message)
		if err != nil {
			log.Printf("⚠️ Failed to notify caregiver: %v", err)
		} else {
			notificationsSent++
		}
	}

	if notificationsSent > 0 {
		log.Printf("✅ %d caregiver(s) notified about medication", notificationsSent)
	}

	return nil
}

// CheckUnacknowledgedAlerts verifica alertas não visualizados e escalona se necessário
func CheckUnacknowledgedAlerts(db *sql.DB, pushService *push.FirebaseService) error {
	query := `
		SELECT 
			a.id,
			a.idoso_id,
			a.mensagem,
			a.severidade,
			i.nome,
			c.telefone
		FROM alertas a
		JOIN idosos i ON i.id = a.idoso_id
		LEFT JOIN cuidadores c ON c.idoso_id = i.id AND c.prioridade = 1
		WHERE a.visualizado = false
		  AND a.necessita_escalamento = true
		  AND a.tempo_escalamento <= NOW()
		  AND a.severidade IN ('critica', 'alta')
	`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query unacknowledged alerts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alertID, idosoID int64
		var message, severity, elderName string
		var phone sql.NullString

		if err := rows.Scan(&alertID, &idosoID, &message, &severity, &elderName, &phone); err != nil {
			log.Printf("❌ Error scanning alert: %v", err)
			continue
		}

		log.Printf("🚨 ESCALANDO alerta não visualizado - ID: %d, Idoso: %s", alertID, elderName)

		// TODO: Implementar ligação telefônica via Twilio para alertas críticos não visualizados
		// if phone.Valid && phone.String != "" {
		//     twilioService.MakeCall(phone.String, elderName, message)
		// }

		// Marcar que o escalonamento foi tentado
		_, _ = db.Exec(`
			UPDATE alertas 
			SET 
				tentativas_envio = tentativas_envio + 1,
				ultima_tentativa = NOW(),
				tempo_escalamento = NOW() + INTERVAL '10 minutes'
			WHERE id = $1
		`, alertID)
	}

	return nil
}

// ScheduleAppointment insere um novo agendamento no banco de dados
func ScheduleAppointment(db *sql.DB, idosoID int64, timestampStr, tipo, descricao string) error {
	// 1. Parse convertendo string ISO para time.Time
	// Suporta formatos ISO parciais ou completos
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	var dataHora time.Time
	var err error

	for _, layout := range layouts {
		dataHora, err = time.Parse(layout, timestampStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("formato de data inválido (%s): %w", timestampStr, err)
	}

	// 2. Inserir no banco
	query := `
		INSERT INTO agendamentos (
			idoso_id, 
			tipo, 
			data_hora_agendada, 
			status, 
			prioridade, 
			dados_tarefa, 
			criado_em, 
			atualizado_em,
			max_retries,
			tentativas_realizadas
		) 
		VALUES ($1, $2, $3, 'agendado', 'media', $4, NOW(), NOW(), 3, 0)
		RETURNING id
	`

	var id int64
	err = db.QueryRow(query, idosoID, tipo, dataHora, descricao).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert appointment: %w", err)
	}

	log.Printf("📅 Appointment scheduled: ID %d for Idoso %d at %s", id, idosoID, dataHora)
	return nil
}
