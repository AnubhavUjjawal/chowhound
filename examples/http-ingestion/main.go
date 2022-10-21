// This package demonstrates how to use the ingestion package to write data to an http endpoint.
package main

import (
	"context"
	"time"

	"github.com/AnubhavUjjawal/chowhound/pkg/common/ingestion"
	"github.com/google/uuid"
)

type MockSrcWithHttpIngestionWriter struct {
	ingestion.IngestionWriter
}

func (r *MockSrcWithHttpIngestionWriter) Run() error {
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		_, err := r.Write(ctx, ingestion.IngestionData{
			Uuid: uuid.New().String(),
		})
		if err != nil {
			return err
		}
		// mocking data generation
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	path := "/ingest"
	port := 5000
	url := "http://localhost:5000/ingest"
	ingestor := ingestion.NewSimpleHttpIngestor(path, port)
	// Since the ingestor is blocking, we need to run it in a goroutine
	// some of the ingested data will be lost if the ingestor is not listening yet
	// That's okay. We are just demonstrating how to use the ingestion package here.

	// In a real world scenario, the ingestor should be running remotely on a different server
	// before the ingestion source starts
	go ingestor.SetupAndListen(&ingestion.StdoutIngestionWriter{})

	src := MockSrcWithHttpIngestionWriter{
		IngestionWriter: ingestion.NewSimpleHttpIngestionWriter(url),
	}
	src.Run()
}
