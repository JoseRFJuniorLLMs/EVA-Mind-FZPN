// cmd/reembed/main.go
// ===================
// Script para re-embedar memórias existentes
// ===================

package main

import (
	"context"
	"database/sql"
	"eva-mind/internal/brainstem/config"
	"eva-mind/internal/hippocampus/memory"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("🔄 Re-embedding Script Started")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Carregar config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	// Conectar DB
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ DB connection error: %v", err)
	}
	defer db.Close()

	// Criar embedding service
	embedder := memory.NewEmbeddingService(cfg.GoogleAPIKey)

	// Buscar memórias sem embedding
	query := `
        SELECT id, content 
        FROM episodic_memories 
        WHERE embedding IS NULL 
          AND content IS NOT NULL
          AND LENGTH(content) > 10
        ORDER BY id
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("❌ Query error: %v", err)
	}
	defer rows.Close()

	// Processar em batches
	ctx := context.Background()
	total := 0
	success := 0
	failed := 0

	log.Println("\n📊 Processing memories...")

	for rows.Next() {
		var id int64
		var content string

		if err := rows.Scan(&id, &content); err != nil {
			log.Printf("⚠️ Scan error: %v", err)
			continue
		}

		total++

		// Gerar embedding
		embedding, err := embedder.GenerateEmbedding(ctx, content)
		if err != nil {
			log.Printf("❌ ID=%d failed: %v", id, err)
			failed++
			continue
		}

		// Atualizar DB
		updateQuery := `
            UPDATE episodic_memories 
            SET embedding = $1 
            WHERE id = $2
        `

		_, err = db.Exec(updateQuery, vectorToPostgres(embedding), id)
		if err != nil {
			log.Printf("❌ ID=%d update failed: %v", id, err)
			failed++
			continue
		}

		success++

		if success%10 == 0 {
			log.Printf("✅ Progress: %d/%d (%.1f%%)", success, total, float64(success)/float64(total)*100)
		}

		// Rate limit (10 RPS = 100ms)
		time.Sleep(100 * time.Millisecond)
	}

	// Summary
	log.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("📊 Re-embedding Complete")
	log.Printf("   Total: %d", total)
	log.Printf("   ✅ Success: %d", success)
	log.Printf("   ❌ Failed: %d", failed)
	log.Printf("   Success Rate: %.1f%%", float64(success)/float64(total)*100)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func vectorToPostgres(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}

	result := "["
	for i, v := range vec {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%f", v)
	}
	result += "]"

	return result
}
