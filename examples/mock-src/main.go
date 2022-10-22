// This package demonstrates how to use the ingestion package to write data to an http endpoint.
package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AnubhavUjjawal/chowhound/pkg/ingestion"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type MockSrcWithHttpIngestionWriter struct {
	ingestion.IngestionWriter
}

func (r *MockSrcWithHttpIngestionWriter) Run() error {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		data := ingestion.IngestionData{
			Uuid:         uuid.New().String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			SourceType:   "MockSrcWithHttpIngestionWriter",
			SourceId:     "MockSrcWithHttpIngestionWriterExample",
			Metadata:     map[string]interface{}{"lorem": "ipsum dolor sit amet"},
			MetadataType: ingestion.METADATA_TYPE_NULL,
			Data:         "Hello World",
			DataId:       strconv.Itoa(i),
			Lang:         "en",
			TenantId:     "mock-tenant-001",
		}
		log.Println("Writing data using ingestion writer", data)
		_, err := r.Write(ctx, data)
		if err != nil {
			log.Println("Error writing data using ingestion writer", err)
			// return err
		}
		// mocking data generation delay
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	url := "http://localhost:5000/ingest"
	newUrl, hasUrlInEnv := os.LookupEnv("INGESTOR_URL")
	if hasUrlInEnv {
		url = newUrl
	}

	src := MockSrcWithHttpIngestionWriter{
		IngestionWriter: ingestion.NewSimpleHttpIngestionWriter(url),
	}
	src.Run()
}
