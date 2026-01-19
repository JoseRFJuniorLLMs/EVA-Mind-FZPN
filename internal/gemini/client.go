package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"eva-mind/internal/config"
	"eva-mind/internal/tools"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	// ✅ Importar tools
)

// AudioCallback é chamado quando áudio PCM é recebido do Gemini
type AudioCallback func(audioBytes []byte)

// ToolCallCallback é chamado quando uma ferramenta precisa ser executada
type ToolCallCallback func(name string, args map[string]interface{}) map[string]interface{}

// TranscriptCallback é chamado quando há transcrição de áudio (Input ou Output)
type TranscriptCallback func(role, text string)

// Client gerencia a conexão WebSocket com Gemini Live API
type Client struct {
	conn         *websocket.Conn
	mu           sync.Mutex
	cfg          *config.Config
	onAudio      AudioCallback
	onToolCall   ToolCallCallback
	onTranscript TranscriptCallback
}

// NewClient cria um novo cliente Gemini usando WebSocket direto
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	url := fmt.Sprintf("wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent?key=%s", cfg.GoogleAPIKey)

	// 🔍 DEBUG: Log API Key (Masked)
	maskedKey := "N/A"
	if len(cfg.GoogleAPIKey) > 8 {
		maskedKey = cfg.GoogleAPIKey[:4] + "..." + cfg.GoogleAPIKey[len(cfg.GoogleAPIKey)-4:]
	}
	log.Printf("🔐 Gemini Config: Key=%s Model=%s", maskedKey, cfg.ModelID)

	// Using query param for now as primary method, but adding header is good practice if supported by library
	// The gorilla/websocket Dialer.DialContext takes headers as the third argument.
	// But EVA-Mind (working) uses nil, so we revert to nil to match it exactly.

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no websocket: %w", err)
	}

	return &Client{conn: conn, cfg: cfg}, nil
}

// SetCallbacks configura os retornos de áudio, ferramentas e transcrição
func (c *Client) SetCallbacks(onAudio AudioCallback, onToolCall ToolCallCallback, onTranscript TranscriptCallback) {
	c.onAudio = onAudio
	c.onToolCall = onToolCall
	c.onTranscript = onTranscript
}

// SendSetup envia a configuração inicial da sessão
func (c *Client) SendSetup(instructions string, voiceSettings map[string]interface{}, memories []string, initialAudio string, toolsDef []tools.FunctionDeclaration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	parsedMemories := []interface{}{}
	for _, m := range memories {
		parsedMemories = append(parsedMemories, m)
	}

	setup := map[string]interface{}{
		"setup": map[string]interface{}{
			"model": fmt.Sprintf("models/%s", c.cfg.ModelID),
			"generationConfig": map[string]interface{}{
				"responseModalities": []string{"AUDIO"},
				"speechConfig": map[string]interface{}{
					"voiceConfig": map[string]interface{}{
						"prebuiltVoiceConfig": map[string]interface{}{
							"voiceName": "Puck", // Voz padrão definida
						},
					},
				},
				"temperature": 0.6,
			},
			"systemInstruction": map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": instructions,
					},
				},
			},
		},
	}

	// ✅ Injetar Tools se houver
	if len(toolsDef) > 0 {
		toolsList := []interface{}{}
		for _, t := range toolsDef {
			toolsList = append(toolsList, t)
		}

		setup["setup"].(map[string]interface{})["tools"] = []interface{}{
			map[string]interface{}{
				"functionDeclarations": toolsList,
			},
		}
		log.Printf("🛠️ [SETUP] %d tools enviadas para o Gemini", len(toolsDef))
	}

	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🔧 CONFIGURANDO GEMINI")
	log.Printf("🎙️ Input: 16kHz PCM16 Mono")
	log.Printf("🔊 Output: 24kHz PCM16 Mono (padrão Gemini)")
	if len(memories) > 0 {
		log.Printf("🧠 Memórias carregadas: %d", len(memories))
	}
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return c.conn.WriteJSON(setup)
}

// StartSession é um alias para SendSetup (wrapper depreciado)
func (c *Client) StartSession(instructions string, tools []interface{}, memories []string, voiceName string) error {
	// Adaptador simples para manter compatibilidade se necessário, mas ideal é atualizar chamadas
	// Passando nil para toolsDef e map vazio para voiceSettings
	return c.SendSetup(instructions, nil, memories, "", nil)
}

// SendAudio envia dados de áudio PCM para o Gemini
func (c *Client) SendAudio(audioData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(audioData)

	// ✅ INPUT: 16kHz (correto para captura do microfone)
	msg := map[string]interface{}{
		"realtime_input": map[string]interface{}{
			"media_chunks": []map[string]string{
				{
					"mime_type": "audio/pcm;rate=16000", // ✅ Correto para INPUT
					"data":      encoded,
				},
			},
		},
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

// SendText envia uma mensagem de texto (system note ou user message)
func (c *Client) SendText(text string) error {
	msg := map[string]interface{}{
		"client_content": map[string]interface{}{
			"turn_complete": true,
			"turns": []map[string]interface{}{
				{
					"role": "user",
					"parts": []map[string]interface{}{
						{
							"text": text,
						},
					},
				},
			},
		},
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

// SendImage envia frames de imagem (JPEG) para o Gemini (Visão Computacional)
func (c *Client) SendImage(imageData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(imageData)

	msg := map[string]interface{}{
		"realtime_input": map[string]interface{}{
			"media_chunks": []map[string]string{
				{
					"mime_type": "image/jpeg",
					"data":      encoded,
				},
			},
		},
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

// ✅ SendMessage envia uma mensagem genérica JSON para o Gemini (usado para ToolResponse e SystemNotes)
func (c *Client) SendMessage(msg interface{}) error {
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

// HandleResponses processa o loop de mensagens
func (c *Client) HandleResponses(ctx context.Context) error {
	log.Printf("👂 HandleResponses: loop iniciado")

	for {
		select {
		case <-ctx.Done():
			log.Printf("🛑 HandleResponses: contexto cancelado")
			return ctx.Err()
		default:
			resp, err := c.ReadResponse()
			if err != nil {
				log.Printf("❌ Erro ao ler resposta: %v", err)
				return err
			}

			// ✅ DEBUG: Mostrar TODAS as respostas do Gemini
			if respBytes, _ := json.Marshal(resp); len(respBytes) > 0 {
				preview := string(respBytes)
				if len(preview) > 300 {
					preview = preview[:300] + "..."
				}
				// log.Printf("📦 Gemini Response: %s", preview)
			}

			// ✅ Verificar setupComplete
			if setupComplete, ok := resp["setupComplete"].(bool); ok && setupComplete {
				log.Printf("✅ Gemini Setup Complete - Pronto para receber áudio!")
				continue
			}

			// Debug de erros
			if errMsg, ok := resp["error"]; ok {
				log.Printf("❌ Gemini Error: %v", errMsg)
				continue
			}

			// ✅ Processar áudio e transcrição
			if serverContent, ok := resp["serverContent"].(map[string]interface{}); ok {

				// ▶️ 1. Capturar Transcrição do Usuário (Input)
				if inputTrans, ok := serverContent["inputAudioTranscription"].(map[string]interface{}); ok {
					if userText, ok := inputTrans["text"].(string); ok && userText != "" {
						// log.Printf("🗣️ [CLIENT] IDOSO: %s", userText)
						if c.onTranscript != nil {
							c.onTranscript("user", userText)
						}
					}
				}

				// ▶️ 2. Capturar Transcrição da IA (Output)
				if audioTrans, ok := serverContent["audioTranscription"].(map[string]interface{}); ok {
					if aiText, ok := audioTrans["text"].(string); ok && aiText != "" {
						// log.Printf("💬 [CLIENT] EVA: %s", aiText)
						if c.onTranscript != nil {
							c.onTranscript("assistant", aiText)
						}
					}
				}

				if modelTurn, ok := serverContent["modelTurn"].(map[string]interface{}); ok {
					if parts, ok := modelTurn["parts"].([]interface{}); ok {
						for _, p := range parts {
							part, ok := p.(map[string]interface{})
							if !ok {
								continue
							}

							// ✅ Procurar por inlineData (áudio)
							if inlineData, ok := part["inlineData"].(map[string]interface{}); ok {
								if audioB64, ok := inlineData["data"].(string); ok {
									audioBytes, err := base64.StdEncoding.DecodeString(audioB64)
									if err != nil {
										log.Printf("❌ Erro ao decodificar base64: %v", err)
										continue
									}
									// ✅ CHAMAR CALLBACK
									if c.onAudio != nil {
										c.onAudio(audioBytes)
									}
								}
							}
						}
					}
				}
			}

			// ✅ Processar tool calls
			if toolCall, ok := resp["toolCall"].(map[string]interface{}); ok {
				log.Printf("🔧 Tool call detectado")
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
				// ✅ FIX: Executar tools em goroutine separada para não travar a voz
				go func(n string, a map[string]interface{}) {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("🚨 PANIC na Tool %s: %v", n, r)
							c.SendToolResponse(n, map[string]interface{}{"error": "Internal error"})
						}
					}()

					result := c.onToolCall(n, a)
					c.SendToolResponse(n, result)
				}(name, args)
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
