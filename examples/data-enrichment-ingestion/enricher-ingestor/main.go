package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/AnubhavUjjawal/chowhound/pkg/ingestion"
)

func DataEnricher(data ingestion.IngestionData) (ingestion.IngestionData, error) {
	log.Println("Data received for enrichment", data)

	// Lang can be set by a NLP module
	data.Lang = "en"

	// SourceId can be set by a database lookup, with combination inspecting the data using NLP
	// We are just mocking the implementation here
	product := []string{"product1", "product2"}
	for _, product := range product {
		if strings.Contains(data.Data, product) {
			data.SourceId = product
			break
		}
	}
	return data, nil
}

func main() {
	// Accepts data on port 3000, enriches it, sends it to port 5000
	path := "/ingest"
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalln("Error converting port to int", err)
	}

	destUrl, ok := os.LookupEnv("DEST_URL")
	if !ok {
		log.Println("DEST_URL not set")
	}
	ingestor := ingestion.NewSimpleHttpIngestor(path, port)
	// writer := ingestion.NewStdoutIngestionWriter()
	writer := ingestion.NewSimpleHttpIngestionWriter(destUrl)
	ingestor.SetupAndListen(DataEnricher, writer)
}
