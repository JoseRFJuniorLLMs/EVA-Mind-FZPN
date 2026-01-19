package gemini

import (
	"bytes"
	"encoding/json"
	"eva-mind/internal/config"
	"fmt"
	"io"
	"net/http"
)

// AnalyzeText envia um prompt para o Gemini via REST API (não-stream)
// Útil para raciocínio (Thinking) e análise de contexto
func AnalyzeText(cfg *config.Config, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", cfg.GeminiAnalysisModel, cfg.GoogleAPIKey)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.2, // Baixa temperatura para raciocínio lógico
			"maxOutputTokens": 1024,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao criar JSON: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Gemini (%d): %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Extrair texto da resposta
	if candidates, ok := response["candidates"].([]interface{}); ok && len(candidates) > 0 {
		candidate := candidates[0].(map[string]interface{})
		if content, ok := candidate["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				part := parts[0].(map[string]interface{})
				if text, ok := part["text"].(string); ok {
					return text, nil
				}
			}
		}
	}

	return "", fmt.Errorf("nenhum texto retornado na resposta")
}
