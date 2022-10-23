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
		log.Fatalln("Error converting port to int", err)
	}
	ingestor := ingestion.NewSimpleHttpIngestor(path, port)
	writer := ingestion.NewStdoutIngestionWriter()
	ingestor.SetupAndListen(nil, writer)
}
