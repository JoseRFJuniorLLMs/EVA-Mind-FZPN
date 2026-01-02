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

// AudioCallback é chamado quando áudio PCM é recebido do Gemini
type AudioCallback func(audioBytes []byte)

// ToolCallCallback é chamado quando uma ferramenta precisa ser executada
type ToolCallCallback func(name string, args map[string]interface{}) map[string]interface{}

// Client gerencia a conexão WebSocket com Gemini Live API
type Client struct {
	conn       *websocket.Conn
	mu         sync.Mutex
	cfg        *config.Config
	onAudio    AudioCallback
	onToolCall ToolCallCallback
}

// NewClient cria um novo cliente Gemini usando WebSocket direto
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	url := fmt.Sprintf("wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent?key=%s", cfg.GoogleAPIKey)

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no websocket: %w", err)
	}

	return &Client{conn: conn, cfg: cfg}, nil
}

// SetCallbacks configura os retornos de áudio e ferramentas (exigido pelo main.go:295)
func (c *Client) SetCallbacks(onAudio AudioCallback, onToolCall ToolCallCallback) {
	c.onAudio = onAudio
	c.onToolCall = onToolCall
}

// SendSetup envia configuração inicial (exigido pelo signaling/websocket.go)
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
	return c.conn.WriteJSON(setupMsg)
}

// StartSession é um alias para SendSetup
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

// ReadResponse lê a próxima resposta bruta do WebSocket
func (c *Client) ReadResponse() (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.conn.ReadJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// HandleResponses processa o loop de mensagens (exigido pelo main.go:318)
func (c *Client) HandleResponses(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := c.ReadResponse()
			if err != nil {
				return err
			}

			// Processar Conteúdo (Áudio/Texto)
			if serverContent, ok := resp["serverContent"].(map[string]interface{}); ok {
				if modelTurn, ok := serverContent["modelTurn"].(map[string]interface{}); ok {
					if parts, ok := modelTurn["parts"].([]interface{}); ok {
						for _, p := range parts {
							part := p.(map[string]interface{})
							if inlineData, ok := part["inlineData"].(map[string]interface{}); ok {
								audioB64 := inlineData["data"].(string)
								audioBytes, _ := base64.StdEncoding.DecodeString(audioB64)
								if c.onAudio != nil {
									c.onAudio(audioBytes)
								}
							}
						}
					}
				}
			}

			// Processar Ferramentas
			if toolCall, ok := resp["toolCall"].(map[string]interface{}); ok {
				c.handleToolCalls(toolCall)
			}
		}
	}
}

func (c *Client) handleToolCalls(toolCall map[string]interface{}) {
	if fcList, ok := toolCall["functionCalls"].([]interface{}); ok {
		for _, f := range fcList {
			fc := f.(map[string]interface{})
			name := fc["name"].(string)
			args := fc["args"].(map[string]interface{})

			if c.onToolCall != nil {
				result := c.onToolCall(name, args)
				c.SendToolResponse(name, result)
			}
		}
	}
}

// SendToolResponse envia o resultado da função de volta ao Gemini
func (c *Client) SendToolResponse(name string, result map[string]interface{}) error {
	msg := map[string]interface{}{
		"tool_response": map[string]interface{}{
			"function_responses": []map[string]interface{}{
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
