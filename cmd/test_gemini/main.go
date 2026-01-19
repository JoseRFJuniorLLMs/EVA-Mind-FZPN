package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado, tentando variáveis de ambiente...")
	}

	apiKey := strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	if apiKey == "" {
		log.Fatal("❌ GOOGLE_API_KEY não encontrada")
	}

	maskedKey := "N/A"
	if len(apiKey) > 8 {
		maskedKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
	}
	fmt.Printf("🔐 Usando API Key: %s\n", maskedKey)

	url := fmt.Sprintf("wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent?key=%s", apiKey)

	fmt.Println("🌐 Conectando a:", "wss://generativelanguage.googleapis.com...")

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	headers := make(http.Header)
	headers.Add("x-goog-api-key", apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, resp, err := dialer.DialContext(ctx, url, headers)
	if err != nil {
		if resp != nil {
			log.Printf("❌ Falha na conexão. Status Code: %d", resp.StatusCode)
		}
		log.Fatalf("❌ Erro ao conectar: %v", err)
	}
	defer conn.Close()

	fmt.Println("✅ Conexão WebSocket estabelecida com SUCESSO!")
	fmt.Println("✅ Autenticação aceita pela API do Gemini.")
}
