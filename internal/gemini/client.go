package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"eva-mind/internal/config"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// AudioCallback é chamado quando áudio PCM é recebido do Gemini
type AudioCallback func(audioBytes []byte)

// TranscriptCallback é chamado quando há transcrição de áudio (Input ou Output)
type TranscriptCallback func(role, text string)

// Client gerencia a conexão WebSocket com Gemini Live API
type Client struct {
	conn         *websocket.Conn
	mu           sync.Mutex
	cfg          *config.Config
	onAudio      AudioCallback
	onTranscript TranscriptCallback
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

// SetCallbacks configura os retornos de áudio e transcrição
func (c *Client) SetCallbacks(onAudio AudioCallback, onTranscript TranscriptCallback) {
	c.onAudio = onAudio
	c.onTranscript = onTranscript
}

// SendSetup envia configuração inicial com memórias episódicas
func (c *Client) SendSetup(instructions string, tools []interface{}, memories []string, voiceName string) error {
	// Enriquecer instruções com memórias relevantes
	enrichedInstructions := instructions

	if len(memories) > 0 {
		enrichedInstructions += "\n\n=== MEMÓRIAS RELEVANTES DO PACIENTE ===\n"
		for i, mem := range memories {
			enrichedInstructions += fmt.Sprintf("%d. %s\n", i+1, mem)
		}
		enrichedInstructions += "=== FIM DAS MEMÓRIAS ===\n\n"
		enrichedInstructions += "IMPORTANTE: Use essas memórias para contextualizar suas respostas e demonstrar que você se lembra do paciente.\n"
	}

	// ✅ CORRETO: Gemini SEMPRE retorna 24kHz quando usa response_modalities: ["AUDIO"]
	// NÃO existe campo sample_rate_hertz na API!
	// 🚨 PROTECTION: User requested to DISABLE TOOLS temporarily to fix Error 1008.
	// A delegação será feita via Texto/Prompt.

	// Default voice fallback
	if voiceName == "" {
		voiceName = "Aoede"
	}

	setupMsg := map[string]interface{}{
		"setup": map[string]interface{}{
			"model": fmt.Sprintf("models/%s", c.cfg.ModelID),
			"generation_config": map[string]interface{}{
				"response_modalities": []string{"AUDIO"},
				"speech_config": map[string]interface{}{
					"voice_config": map[string]interface{}{
						"prebuilt_voice_config": map[string]string{
							"voice_name": voiceName,
						},
					},
				},
			},
			"system_instruction": map[string]interface{}{
				"parts": []map[string]string{
					{"text": enrichedInstructions},
				},
			},
		},
	}

	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🔧 CONFIGURANDO GEMINI")
	log.Printf("🎙️ Input: 16kHz PCM16 Mono")
	log.Printf("🔊 Output: 24kHz PCM16 Mono (padrão Gemini)")
	log.Printf("🗣️ Voz: %s", voiceName)
	if len(memories) > 0 {
		log.Printf("🧠 Memórias carregadas: %d", len(memories))
	}
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(setupMsg)
}

// StartSession é um alias para SendSetup
func (c *Client) StartSession(instructions string, tools []interface{}, memories []string, voiceName string) error {
	return c.SendSetup(instructions, tools, memories, voiceName)
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
		}
	}
}

// Close fecha a conexão
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
