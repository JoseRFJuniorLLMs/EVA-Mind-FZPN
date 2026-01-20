package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
)

func main() {
	godotenv.Load(".env")

	host := "104.248.219.200" // Mesmo IP do Postgres
	port := 6334

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("--- LISTAGEM DE COLEÇÕES QDRANT ---")
	collections, err := client.ListCollections(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range collections {
		info, err := client.GetCollectionInfo(ctx, c)
		if err == nil {
			fmt.Printf("Coleção: %s | Pontos: %d\n", c, info.PointsCount)
		}
	}
}
