package gemini

import (
	"context"
	"encoding/base64"
	"eva-mind/internal/config"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	genaiClient *genai.Client
	session     *genai.ChatSession
	cfg         *config.Config
	mu          sync.Mutex

	// Canal para enviar respostas decodificadas
	RespChan chan map[string]interface{}
	stopChan chan struct{}

	// ✅ NOVO: Controle de concorrência
	ctx      context.Context
	cancel   context.CancelFunc
	audioSem chan struct{} // Limita goroutines concorrentes
	wg       sync.WaitGroup

	// ✅ NOVO: Métricas
	audioSent     int64
	audioErrors   int64
	lastErrorTime time.Time
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	log.Printf("🔌 Inicializando Gemini SDK Client...")
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GoogleAPIKey))
	if err != nil {
		log.Printf("❌ Erro ao criar cliente Gemini: %v", err)
		return nil, err
	}

	clientCtx, cancel := context.WithCancel(ctx)

	return &Client{
		genaiClient: client,
		cfg:         cfg,
		RespChan:    make(chan map[string]interface{}, 100),
		stopChan:    make(chan struct{}),
		ctx:         clientCtx,
		cancel:      cancel,
		audioSem:    make(chan struct{}, 5), // Máximo 5 goroutines concorrentes
	}, nil
}

func (c *Client) SendSetup(instructions string, tools []interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Printf("⚙️ Configurando Modelo Gemini: %s", c.cfg.ModelID)
	model := c.genaiClient.GenerativeModel(c.cfg.ModelID)

	model.SystemInstruction = genai.NewUserContent(genai.Text(instructions))

	log.Printf("🚀 Iniciando Chat Session...")
	c.session = model.StartChat()

	return nil
}

func (c *Client) SendAudio(audioData []byte) error {
	c.mu.Lock()
	session := c.session
	c.mu.Unlock()

	if session == nil {
		return fmt.Errorf("sessão não iniciada")
	}

	// ✅ Controle de concorrência
	select {
	case c.audioSem <- struct{}{}:
		// Got semaphore
	case <-c.ctx.Done():
		return fmt.Errorf("client closed")
	default:
		// Backpressure: se já tem 5 goroutines, rejeita
		log.Printf("⚠️ Too many concurrent audio requests, rejecting")
		return fmt.Errorf("too many concurrent audio requests")
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer func() { <-c.audioSem }()

		// ✅ Recovery para evitar panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("🚨 PANIC in SendAudio goroutine: %v", r)
				log.Printf("Stack trace: %s", debug.Stack())
				c.audioErrors++
			}
		}()

		// ✅ Usar contexto do client com timeout
		ctx, cancel := context.WithTimeout(c.ctx, 15*time.Second)
		defer cancel()

		iter := session.SendMessageStream(ctx, genai.Blob{
			MIMEType: "audio/pcm;rate=16000",
			Data:     audioData,
		})

		for {
			resp, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					// Normal termination
					break
				}

				// ✅ Log apropriado
				c.audioErrors++
				c.lastErrorTime = time.Now()

				// Log apenas a cada 5 segundos para não inundar
				if time.Since(c.lastErrorTime) > 5*time.Second {
					log.Printf("⚠️ Gemini stream error: %v", err)
				}
				break
			}

			c.audioSent++
			c.processResponse(resp)
		}
	}()

	return nil
}

func (c *Client) processResponse(resp *genai.GenerateContentResponse) {
	if resp == nil || len(resp.Candidates) == 0 {
		return
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil {
		return
	}

	partsSlice := []interface{}{}

	for _, part := range candidate.Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			partsSlice = append(partsSlice, map[string]interface{}{
				"text": string(txt),
			})
		} else if blob, ok := part.(genai.Blob); ok {
			// Encode bytes to base64 string for handler compatibility
			encoded := base64.StdEncoding.EncodeToString(blob.Data)

			partsSlice = append(partsSlice, map[string]interface{}{
				"inlineData": map[string]interface{}{
					"mimeType": blob.MIMEType,
					"data":     encoded,
				},
			})
		}
	}

	responseMap := map[string]interface{}{
		"serverContent": map[string]interface{}{
			"modelTurn": map[string]interface{}{
				"parts": partsSlice,
			},
		},
	}

	select {
	case c.RespChan <- responseMap:
	default:
		log.Printf("⚠️ Canal de resposta cheio, dropando msg")
	}
}

func (c *Client) ReadResponse() (map[string]interface{}, error) {
	select {
	case resp := <-c.RespChan:
		return resp, nil
	case <-time.After(100 * time.Millisecond):
		return nil, nil
	}
}

func (c *Client) Close() error {
	log.Printf("🧹 Closing Gemini client...")

	close(c.stopChan)
	c.cancel() // Cancela todas as goroutines

	// ✅ Aguarda todas as goroutines terminarem
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Printf("✅ Todas as goroutines Gemini finalizadas")
	case <-time.After(5 * time.Second):
		log.Printf("⚠️ Timeout aguardando goroutines Gemini")
	}

	if c.genaiClient != nil {
		return c.genaiClient.Close()
	}
	return nil
}

// ✅ NOVO: Método para obter estatísticas
func (c *Client) GetStats() map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	return map[string]interface{}{
		"audio_sent":      c.audioSent,
		"audio_errors":    c.audioErrors,
		"last_error_time": c.lastErrorTime,
	}
}
