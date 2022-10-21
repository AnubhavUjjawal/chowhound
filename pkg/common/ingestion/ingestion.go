package ingestion

import (
	"context"
)

type IngestionData struct {
	Uuid string `json:"uuid" db:"uuid"`
}

// IngestionWriter is an interface which defines the Write method for writing ingestion data to 1 or multiple destinations.
// All ingestion sources must implement this interface using composition with a concrete IngestionWriter.
type IngestionWriter interface {
	Write(ctx context.Context, data IngestionData) (result interface{}, err error)
}

// Ingestor is an interface which defines the SetupAndListen method to read ingestion data from multiple
// sources and writing it to one or multiple destinations. All ingestion data receivers should implement this interface
// using composition with a concrete Ingestor.
type Ingestor interface {
	// This method could be blocking
	SetupAndListen(handler IngestionWriter) error
}
