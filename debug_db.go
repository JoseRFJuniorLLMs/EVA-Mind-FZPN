package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbURL := "postgres://postgres:Debian23%40@104.248.219.200:5432/eva-db?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables := []string{"episodic_memories", "analise_gemini", "analise_audio_avancada", "idosos"}
	for _, t := range tables {
		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", t)).Scan(&count)
		if err == nil {
			fmt.Printf("TABELA %s: %d registros\n", t, count)
		} else {
			fmt.Printf("TABELA %s: erro (%v)\n", t, err)
		}
	}

	fmt.Println("\n--- DETALHES DE ANALISE_AUDIO_AVANCADA ---")
	rows, err := db.Query("SELECT id, idoso_id, transcricao_principal FROM analise_audio_avancada ORDER BY criada_em DESC LIMIT 3")
	if err == nil {
		for rows.Next() {
			var id, idosoID int64
			var trans string
			rows.Scan(&id, &idosoID, &trans)
			fmt.Printf("ID: %d | Idoso: %d | Trans: %s...\n", id, idosoID, truncate(trans, 50))
		}
		rows.Close()
	}

	fmt.Println("\n--- DETALHES DE ANALISE_GEMINI ---")
	rows, err = db.Query("SELECT id, idoso_id, analise FROM analise_gemini ORDER BY criada_em DESC LIMIT 3")
	if err == nil {
		for rows.Next() {
			var id, idosoID int64
			var analise string
			rows.Scan(&id, &idosoID, &analise)
			fmt.Printf("ID: %d | Idoso: %d | Analise: %s...\n", id, idosoID, truncate(analise, 50))
		}
		rows.Close()
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
