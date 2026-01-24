package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"eva-mind/internal/brainstem/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// ToolsClient usa Gemini 2.5 Flash via REST para analisar transcrições e executar tools
type ToolsClient struct {
	cfg        *config.Config
	httpClient *http.Client
}

// ToolCall representa uma chamada de ferramenta detectada
type ToolCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// NewToolsClient cria um novo cliente para análise de tools
func NewToolsClient(cfg *config.Config) *ToolsClient {
	return &ToolsClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// AnalyzeTranscription envia transcrição para Gemini 2.5 Flash e detecta tools
func (tc *ToolsClient) AnalyzeTranscription(ctx context.Context, transcript string, role string) ([]ToolCall, error) {
	// Só analisar falas do usuário (idoso)
	if role != "user" {
		return nil, nil
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
		tc.cfg.GoogleAPIKey,
	)

	// Prompt para detectar intenções e tools
	systemPrompt := `Você é um analisador de intenções para assistente de saúde.
Analise a fala do idoso e detecte se ele está solicitando alguma ação que requer uma ferramenta.

FERRAMENTAS DISPONÍVEIS:
- alert_family: Alertar família em emergência (args: reason, severity)
- confirm_medication: Confirmar medicamento tomado (args: medication_name)
- schedule_appointment: Agendar compromisso/lembrete (args: timestamp, type, description)
- call_family_webrtc: Ligar para família
- call_central_webrtc: Ligar para central
- call_doctor_webrtc: Ligar para médico
- call_caregiver_webrtc: Ligar para cuidador
- google_search_retrieval: Pesquisar informações em tempo real na internet (args: query)
- play_nostalgic_music: Tocar músicas da juventude do paciente (args: decade)
- radio_station_tuner: Sintonizar rádio AM/FM (args: station_name)
- play_relaxation_sounds: Tocar sons relaxantes (args: sound_type)
- hymn_and_prayer_player: Tocar hinos ou orações (args: type, content_name)
- daily_mass_stream: Ver missa ao vivo
- watch_classic_movies: Ver filmes clássicos (args: movie_name)
- watch_news_briefing: Ver resumo de notícias (args: topic)
- read_newspaper_aloud: Ler manchetes de jornais (args: newspaper)
- horoscope_daily: Ler horóscopo do dia (args: sign)
- play_trivia_game: Iniciar jogo de quiz/trivia (args: theme)
- riddle_and_joke_teller: Contar piada ou adivinha (args: type)
- voice_diary: Iniciar sessão de diário por voz
- poetry_generator: Criar um poema personalizado (args: theme)
- learn_new_language: Iniciar lição de idioma (args: language)

Se detectar uma intenção que requer ferramenta, responda APENAS com JSON:
{"tool": "nome_da_tool", "args": {...}}

Se NÃO detectar nenhuma intenção de ferramenta, responda: {"tool": "none"}

Exemplos:
Fala: "Me lembre de tomar remédio às 14h"
Resposta: {"tool": "schedule_appointment", "args": {"timestamp": "2026-01-13T14:00:00Z", "type": "medicamento", "description": "Tomar remédio"}}

Fala: "Estou com dor no peito"
Resposta: {"tool": "alert_family", "args": {"reason": "Paciente relatou dor no peito", "severity": "critica"}}

Fala: "Como está o tempo hoje?"
Resposta: {"tool": "google_search_retrieval", "args": {"query": "previsão do tempo para hoje"}}

Fala: "Obrigado"
Resposta: {"tool": "none"}`

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": systemPrompt},
				},
			},
			{
				"role": "model",
				"parts": []map[string]string{
					{"text": "Entendido. Vou analisar as falas e detectar intenções de ferramentas."},
				},
			},
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": fmt.Sprintf("Fala do idoso: \"%s\"", transcript)},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.1, // Baixa temperatura para respostas consistentes
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := tc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Extrair texto da resposta
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return nil, nil
	}

	candidate := candidates[0].(map[string]interface{})
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return nil, nil
	}

	part := parts[0].(map[string]interface{})
	text, ok := part["text"].(string)
	if !ok {
		return nil, nil
	}

	log.Printf("🤖 [TOOLS] Resposta do modelo: %s", text)

	// Parsear JSON da resposta
	var toolResponse struct {
		Tool string                 `json:"tool"`
		Args map[string]interface{} `json:"args"`
	}

	if err := json.Unmarshal([]byte(text), &toolResponse); err != nil {
		log.Printf("⚠️ [TOOLS] Erro ao parsear resposta como JSON: %v", err)
		return nil, nil
	}

	// Se não detectou tool, retornar vazio
	if toolResponse.Tool == "none" || toolResponse.Tool == "" {
		return nil, nil
	}

	log.Printf("✅ [TOOLS] Tool detectada: %s com args: %+v", toolResponse.Tool, toolResponse.Args)

	return []ToolCall{
		{
			Name: toolResponse.Tool,
			Args: toolResponse.Args,
		},
	}, nil
}
