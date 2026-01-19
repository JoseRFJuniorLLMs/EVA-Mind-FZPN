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

	rows, err := db.Query("SELECT * FROM idosos LIMIT 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='idosos' AND column_name='medicamentos_regulares')").Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("EXISTS_REGULARES: %v\n", exists)

	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='idosos' AND column_name='medicamentos_atuais')").Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("EXISTS_ATUAIS: %v\n", exists)
}
