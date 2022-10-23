package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AnubhavUjjawal/chowhound/pkg/ingestion"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type MockSrcIncompleteDataWithHttpIngestionWriter struct {
	ingestion.IngestionWriter
}

func (r *MockSrcIncompleteDataWithHttpIngestionWriter) Run() error {
	ctx := context.Background()
	product := []string{"product1", "product2"}
	for i := 0; i < 10; i++ {
		// Lang and sourceId is missing when we send the data for enrichment.
		data := ingestion.IngestionData{
			Uuid:       uuid.New().String(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			SourceType: "Twitter",
			// SourceId is missing
			// SourceId:     "",
			Metadata:     map[string]interface{}{"lorem": "ipsum dolor sit amet"},
			MetadataType: ingestion.METADATA_TYPE_NULL,
			Data:         fmt.Sprintf("@Enterpret your %s seems to have an issue", product[i%2]),
			DataId:       strconv.Itoa(i),
			// Lang is missing
			// Lang:         "",
			TenantId: "mock-tenant-001",
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
	url := "http://localhost:3000/ingest"
	newUrl, hasUrlInEnv := os.LookupEnv("INGESTOR_URL")
	if hasUrlInEnv {
		url = newUrl
	}

	src := MockSrcIncompleteDataWithHttpIngestionWriter{
		IngestionWriter: ingestion.NewSimpleHttpIngestionWriter(url),
	}
	src.Run()
}
