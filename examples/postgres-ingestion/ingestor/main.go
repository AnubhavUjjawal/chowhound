package main

import (
	"log"
	"os"
	"strconv"

	"github.com/AnubhavUjjawal/chowhound/pkg/ingestion"
)

func main() {
	path := "/ingest"
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("Error converting port to int", err)
	}
	postgresUri := os.Getenv("POSTGRES_URI")
	postgresTableName := os.Getenv("POSTGRES_TABLE_NAME")

	if postgresUri == "" || postgresTableName == "" {
		log.Fatalln("POSTGRES_URI and POSTGRES_TABLE_NAME env vars are required")
	}
	// postgresUri := "postgres://example:example@localhost:5432/example?sslmode=disable"
	// postgresTableName := "ingestion_data"
	ingestor := ingestion.NewSimpleHttpIngestor(path, port)
	writer := ingestion.NewSimplePostgresIngestionWriter(postgresUri, postgresTableName)
	ingestor.SetupAndListen(writer)
}
