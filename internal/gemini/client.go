package gemini

import (
	"context"
	"encoding/base64"
	"eva-mind/internal/config"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// AudioCallback é chamado quando áudio é recebido do Gemini
type AudioCallback func(audioBytes []byte)

// ToolCallCallback é chamado quando uma ferramenta precisa ser executada
type ToolCallCallback func(name string, args map[string]interface{}) map[string]interface{}

// Client gerencia a conexão WebSocket com Gemini Live API
type Client struct {
	conn          *websocket.Conn
	mu            sync.Mutex
	cfg           *config.Config
	onAudio       AudioCallback
	onToolCall    ToolCallCallback
	setupComplete bool
}

// NewClient cria um novo cliente Gemini usando WebSocket direto
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	url := fmt.Sprintf("wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent?key=%s", cfg.GoogleAPIKey)
	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, cfg: cfg}, nil
}

// SendSetup envia configuração inicial para o Gemini
func (c *Client) SendSetup(instructions string, tools []interface{}) error {
	setupMsg := map[string]interface{}{
		"setup": map[string]interface{}{
			"model": fmt.Sprintf("models/%s", c.cfg.ModelID),
			"generation_config": map[string]interface{}{
				"response_modalities": []string{"AUDIO"},
				"speech_config": map[string]interface{}{
					"voice_config": map[string]interface{}{
						"prebuilt_voice_config": map[string]string{
							"voice_name": "Aoede",
						},
					},
				},
			},
			"system_instruction": map[string]interface{}{
				"parts": []map[string]string{
					{"text": instructions},
				},
			},
			"tools": tools,
		},
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.conn.WriteJSON(setupMsg); err != nil {
		return fmt.Errorf("failed to send setup: %w", err)
	}

	return nil
}

// StartSession é um alias para SendSetup (compatibilidade)
func (c *Client) StartSession(instructions string, tools []interface{}) error {
	return c.SendSetup(instructions, tools)
}

// SendAudio envia dados de áudio PCM para o Gemini
func (c *Client) SendAudio(audioData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(audioData)

	msg := map[string]interface{}{
		"realtime_input": map[string]interface{}{
			"media_chunks": []map[string]string{
				{
					"mime_type": "audio/pcm;rate=16000",
					"data":      encoded,
				},
			},
		},
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

// ReadResponse lê a próxima resposta do Gemini (SÍNCRONO - SEM TIMEOUT)
func (c *Client) ReadResponse() (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.conn.ReadJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SetCallbacks configura os callbacks para áudio e tool calls
func (c *Client) SetCallbacks(onAudio AudioCallback, onToolCall ToolCallCallback) {
	c.onAudio = onAudio
	c.onToolCall = onToolCall
}

// HandleResponses processa respostas do Gemini em loop (versão melhorada)
func (c *Client) HandleResponses(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := c.ReadResponse()
			if err != nil {
				return fmt.Errorf("erro na leitura do websocket: %w", err)
			}

			// 1. Confirmação de Setup (Acontece uma vez no início)
			if _, ok := resp["setupComplete"]; ok {
				c.setupComplete = true
				fmt.Println("✅ Gemini Live: Conexão estabelecida e configurada.")
				continue
			}

			// 2. Conteúdo do Servidor (Áudio e Transcrição)
			if serverContent, ok := resp["serverContent"].(map[string]interface{}); ok {

				// Verifica interrupção (O usuário começou a falar)
				if _, interrupted := serverContent["interrupted"]; interrupted {
					fmt.Println("⚠️ IA Interrompida pelo usuário.")
					// TODO: Parar de tocar o áudio atual no player
					continue
				}

				if modelTurn, ok := serverContent["modelTurn"].(map[string]interface{}); ok {
					parts, _ := modelTurn["parts"].([]interface{})
					for _, part := range parts {
						p, ok := part.(map[string]interface{})
						if !ok {
							continue
						}

						// Áudio Nativo (Base64 -> Bytes)
						if inlineData, ok := p["inlineData"].(map[string]interface{}); ok {
							audioBase64, _ := inlineData["data"].(string)
							audioBytes, err := base64.StdEncoding.DecodeString(audioBase64)
							if err == nil && c.onAudio != nil {
								// ENVIAR PARA O PLAYER OU FRONTEND
								// O áudio de saída do Gemini Live é PCM 24kHz Mono
								c.onAudio(audioBytes)
							}
						}

						// Transcrição de texto
						if text, ok := p["text"].(string); ok {
							fmt.Printf("💬 Gemini: %s\n", text)
						}
					}
				}
			}

			// 3. Chamada de Ferramentas (Tool Calls)
			if toolCall, ok := resp["toolCall"].(map[string]interface{}); ok {
				c.handleToolCalls(toolCall)
			}

			// 4. Erros do Servidor
			if serverError, ok := resp["error"].(map[string]interface{}); ok {
				return fmt.Errorf("erro do servidor Gemini: %v", serverError["message"])
			}
		}
	}
}

// handleToolCalls processa chamadas de ferramentas
func (c *Client) handleToolCalls(toolCall map[string]interface{}) {
	functionCalls, ok := toolCall["functionCalls"].([]interface{})
	if !ok {
		return
	}

	for _, fc := range functionCalls {
		call, ok := fc.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := call["name"].(string)
		args, _ := call["args"].(map[string]interface{})

		fmt.Printf("🛠️ Executando ferramenta: %s com args: %v\n", name, args)

		if c.onToolCall != nil {
			result := c.onToolCall(name, args)
			// TODO: Enviar toolResponse de volta para o Gemini
			c.sendToolResponse(name, result)
		}
	}
}

// sendToolResponse envia resposta de ferramenta de volta ao Gemini
func (c *Client) sendToolResponse(name string, result map[string]interface{}) error {
	msg := map[string]interface{}{
		"toolResponse": map[string]interface{}{
			"functionResponses": []map[string]interface{}{
				{
					"name":     name,
					"response": result,
				},
			},
		},
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

// Close fecha a conexão
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
