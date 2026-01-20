package gemini

import (
	"database/sql"
	"eva-mind/internal/brainstem/push"
	"fmt"
	"log"
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
		map[string]interface{}{
			"function_declarations": []interface{}{
				map[string]interface{}{
					"name":        "manage_calendar_event",
					"description": "Gerencia eventos no Google Calendar (cria ou lista eventos)",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"action": map[string]interface{}{
								"type":        "string",
								"description": "Ação a realizar: 'create' ou 'list'",
								"enum":        []string{"create", "list"},
							},
							"summary": map[string]interface{}{
								"type":        "string",
								"description": "Título do evento (Obrigatório para 'create'). Ex: 'Consulta Cardiologista'",
							},
							"description": map[string]interface{}{
								"type":        "string",
								"description": "Descrição detalhada do evento",
							},
							"start_time": map[string]interface{}{
								"type":        "string",
								"description": "Horário de início (ISO 8601). Ex: '2024-12-25T14:00:00-03:00'",
							},
							"end_time": map[string]interface{}{
								"type":        "string",
								"description": "Horário de término (ISO 8601). Ex: '2024-12-25T15:00:00-03:00'",
							},
						},
						"required": []string{"action"},
					},
				},
			},
		},
		// 1. Gmail - send_email
		map[string]interface{}{
			"name":        "send_email",
			"description": "Envia um email usando Gmail do usuário",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"to": map[string]interface{}{
						"type":        "string",
						"description": "Email do destinatário",
					},
					"subject": map[string]interface{}{
						"type":        "string",
						"description": "Assunto do email",
					},
					"body": map[string]interface{}{
						"type":        "string",
						"description": "Corpo do email",
					},
				},
				"required": []string{"to", "subject", "body"},
			},
		},
		// 2. Drive - save_to_drive
		map[string]interface{}{
			"name":        "save_to_drive",
			"description": "Salva um documento no Google Drive do usuário",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Nome do arquivo",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Conteúdo do documento",
					},
					"folder": map[string]interface{}{
						"type":        "string",
						"description": "Nome da pasta (opcional)",
					},
				},
				"required": []string{"filename", "content"},
			},
		},
		// 3. Sheets - manage_health_sheet
		map[string]interface{}{
			"name":        "manage_health_sheet",
			"description": "Gerencia planilha de saúde no Google Sheets",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Ação: 'create' ou 'append'",
						"enum":        []string{"create", "append"},
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Título da planilha (para create)",
					},
					"data": map[string]interface{}{
						"type":        "object",
						"description": "Dados de saúde (date, time, blood_pressure, glucose, medication, notes)",
					},
				},
				"required": []string{"action"},
			},
		},
		// 4. Docs - create_health_doc
		map[string]interface{}{
			"name":        "create_health_doc",
			"description": "Cria um documento de saúde no Google Docs",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Título do documento",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Conteúdo do documento",
					},
				},
				"required": []string{"title", "content"},
			},
		},
		// 5. Maps - find_nearby_places
		map[string]interface{}{
			"name":        "find_nearby_places",
			"description": "Busca lugares próximos (farmácias, hospitais, restaurantes)",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"place_type": map[string]interface{}{
						"type":        "string",
						"description": "Tipo: pharmacy, hospital, restaurant, etc",
					},
					"location": map[string]interface{}{
						"type":        "string",
						"description": "Localização (lat,lng)",
					},
					"radius": map[string]interface{}{
						"type":        "integer",
						"description": "Raio em metros (padrão: 5000)",
					},
				},
				"required": []string{"place_type", "location"},
			},
		},
		// 6. YouTube - search_videos
		map[string]interface{}{
			"name":        "search_videos",
			"description": "Busca vídeos no YouTube",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Termo de busca",
					},
					"max_results": map[string]interface{}{
						"type":        "integer",
						"description": "Número máximo de resultados (padrão: 5)",
					},
				},
				"required": []string{"query"},
			},
		},
		// 7. Spotify - play_music
		map[string]interface{}{
			"name":        "play_music",
			"description": "Toca musica no spotify",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "musica ou artista",
					},
				},
				"required": []string{"query"},
			},
		},
		// 8. Uber - request_ride
		map[string]interface{}{
			"name":        "request_ride",
			"description": "Solicita corrida uber",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"startLat": map[string]interface{}{
						"type":        "number",
						"description": "startLat",
					},
					"startLng": map[string]interface{}{
						"type":        "number",
						"description": "startLng",
					},
					"endLat": map[string]interface{}{
						"type":        "number",
						"description": "endLat",
					},
					"endLng": map[string]interface{}{
						"type":        "number",
						"description": "endLng",
					},
				},
				"required": []string{"startLat", "startLng", "endLat", "endLng"},
			},
		},
		// 9. Google Fit - get_health_data
		map[string]interface{}{
			"name":        "get_health_data",
			"description": "Recupera dados de saúde do Google Fit",
			"parameters": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		// 10. WhatsApp - send_whatsapp
		map[string]interface{}{
			"name":        "send_whatsapp",
			"description": "Envia mensagem whatsapp",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"to": map[string]interface{}{
						"type":        "string",
						"description": "Numero destino",
					},
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Mensagem",
					},
				},
				"required": []string{"to", "message"},
			},
		},
		// ✅ Google Search Tool (Integrada ao modelo)
		map[string]interface{}{
			"google_search_retrieval": map[string]interface{}{
				"dynamic_retrieval_config": map[string]interface{}{
					"mode":              "MODE_DYNAMIC",
					"dynamic_threshold": 0.5, // Ajuste para equilibrar pesquisa e resposta direta
				},
			},
		},
		// ✅ SQL Select Tool (Database Access)
		map[string]interface{}{
			"function_declarations": []interface{}{
				map[string]interface{}{
					"name":        "run_sql_select",
					"description": "Executa uma consulta SQL SELECT (apenas leitura) no banco de dados para responder perguntas sobre o sistema. Use valid PostgreSQL syntax.",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"query": map[string]interface{}{
								"type":        "string",
								"description": "A consulta SQL SELECT a ser executada. Ex: 'SELECT count(*) FROM idosos'",
							},
						},
						"required": []string{"query"},
					},
				},
				// ✅ Change Voice Tool (Runtime)
				map[string]interface{}{
					"name":        "change_voice",
					"description": "Altera a voz da assistente (EVA) em tempo real. Vozes disponiveis: Puck, Charon, Kore, Fenrir, Aoede",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"voice_name": map[string]interface{}{
								"type":        "string",
								"description": "Nome da voz desejada (Puck, Charon, Kore, Fenrir, Aoede)",
								"enum":        []string{"Puck", "Charon", "Kore", "Fenrir", "Aoede"},
							},
						},
						"required": []string{"voice_name"},
					},
				},
			},
		},
	}
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
