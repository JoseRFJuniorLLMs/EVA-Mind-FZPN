package main

import (
	"context"
	"database/sql"
	"eva-mind/internal/config"
	"eva-mind/internal/infrastructure/vector"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔄 PostgreSQL → Qdrant Migration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	ctx := context.Background()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	// Connect to PostgreSQL
	dbURL := buildDatabaseURL(cfg)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ PostgreSQL connection error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ PostgreSQL ping error: %v", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")

	// Connect to Qdrant
	qdrantClient, err := vector.NewQdrantClient(cfg.QdrantHost, cfg.QdrantPort)
	if err != nil {
		log.Fatalf("❌ Qdrant connection error: %v", err)
	}
	defer qdrantClient.Close()
	fmt.Println("✅ Connected to Qdrant")
	fmt.Println()

	// Migrate memories
	fmt.Println("📦 Migrating memories...")
	memoriesCount, err := migrateMemories(ctx, db, qdrantClient)
	if err != nil {
		log.Printf("⚠️ Memories migration error: %v", err)
	} else {
		fmt.Printf("✅ Migrated %d memories\n", memoriesCount)
	}
	fmt.Println()

	// Summary
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 MIGRATION SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Memories:    %d\n", memoriesCount)
	fmt.Printf("Total:       %d\n", memoriesCount)
	fmt.Println()
	fmt.Println("✅ Migration complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func migrateMemories(ctx context.Context, db *sql.DB, qdrant *vector.QdrantClient) (int, error) {
	// Query memories with embeddings from PostgreSQL
	// Note: Adjust table/column names based on your schema
	query := `
		SELECT id, user_id, content, created_at
		FROM memories
		WHERE id > 0
		ORDER BY id
		LIMIT 100
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	count := 0
	var points []*qdrant.PointStruct
	batchSize := 100

	for rows.Next() {
		var (
			id        int64
			userID    int64
			content   string
			createdAt time.Time
		)

		if err := rows.Scan(&id, &userID, &content, &createdAt); err != nil {
			log.Printf("⚠️ Scan error (row %d): %v", id, err)
			continue
		}

		// Generate dummy embedding (768 dimensions)
		// TODO: Replace with actual embedding generation
		embedding := make([]float32, 768)
		for i := range embedding {
			embedding[i] = 0.1
		}

		// Create Qdrant point using helper
		point := vector.CreatePoint(
			uint64(id),
			embedding,
			map[string]interface{}{
				"user_id":    userID,
				"content":    content,
				"timestamp":  createdAt.Format(time.RFC3339),
				"event_type": "memory",
			},
		)

		points = append(points, point)

		// Insert batch
		if len(points) >= batchSize {
			if err := qdrant.Upsert(ctx, "memories", points); err != nil {
				log.Printf("⚠️ Batch insert error: %v", err)
			} else {
				count += len(points)
				fmt.Printf("  Migrated %d memories...\n", count)
			}
			points = nil
		}
	}

	// Insert remaining
	if len(points) > 0 {
		if err := qdrant.Upsert(ctx, "memories", points); err != nil {
			return count, fmt.Errorf("final batch error: %w", err)
		}
		count += len(points)
	}

	return count, nil
}

func buildDatabaseURL(cfg *config.Config) string {
	if cfg.DatabaseURL != "" {
		return cfg.DatabaseURL
	}

	// Build from individual components
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")
	dbName := getEnvOrDefault("DB_NAME", "eva_db")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
