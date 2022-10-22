// This package demonstrates how to use the ingestion package to write data to an http endpoint.
package main

import (
	"context"
	"strconv"
	"time"

	"github.com/AnubhavUjjawal/chowhound/pkg/ingestion"
	"github.com/google/uuid"
)

type MockSrcWithHttpIngestionWriter struct {
	ingestion.IngestionWriter
}

func (r *MockSrcWithHttpIngestionWriter) Run() error {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := r.Write(ctx, ingestion.IngestionData{
			Uuid:         uuid.New().String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			SourceType:   "MockSrcWithHttpIngestionWriter",
			SourceId:     "MockSrcWithHttpIngestionWriterExample",
			Metadata:     nil,
			MetadataType: ingestion.METADATA_TYPE_NULL,
			Data:         "Hello World",
			DataId:       strconv.Itoa(i),
			Lang:         "en",
			TenantId:     "mock-tenant-001",
		})
		// log.Println(reflect.TypeOf(data), data, err)
		if err != nil {
			return err
		}
		// mocking data generation delay
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	url := "http://localhost:5000/ingest"

	src := MockSrcWithHttpIngestionWriter{
		IngestionWriter: ingestion.NewSimpleHttpIngestionWriter(url),
	}
	src.Run()
}
