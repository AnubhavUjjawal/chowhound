package main

import (
	"github.com/AnubhavUjjawal/chowhound/pkg/ingestion"
)

func main() {
	path := "/ingest"
	port := 5000
	ingestor := ingestion.NewSimpleHttpIngestor(path, port)
	ingestor.SetupAndListen(&ingestion.StdoutIngestionWriter{})
}
