package ingestion

import (
	"context"
	"time"
)

const (
	METADATA_TYPE_NULL = "NULL"
)

type IngestionDataProcessor func(IngestionData) (IngestionData, error)

type IngestionData struct {
	// universal identifier for the ingestion data record to be created.
	// If not provided, a new UUID will be generated by the consuming ingestor.
	Uuid string `json:"uuid" db:"uuid"`

	// source type, for example, discourse, github etc.
	SourceType string `json:"source_type" db:"source_type"`

	// unique identifier for the source, for example, discourse forum name, github repo cannonical url
	SourceId string `json:"source_id" db:"source_id"`

	// the main content of the ingestion data record
	Data string `json:"data" db:"data"`

	// source specific identifier for the data, for example, discourse post id, github issue id
	// data id together with source type and source id would help us acheive idempotency
	DataId string `json:"data_id" db:"data_id"`

	// metadata for the data, for example, discourse post id, github issue id
	// it is an open ended field, and the format is up to the source.
	// in database, it should be stored as a json field (Ex: jsonb in postgres), or a json encoded string.
	Metadata interface{} `json:"metadata" db:"metadata"`

	// type of Metadata object, so that ingestor service know how to process or ignore it.
	MetadataType string `json:"metadata_type" db:"metadata_type"`

	// language of the data ingested
	Lang string `json:"lang" db:"lang"`

	// id of tenant whose data is being ingested
	TenantId string `json:"tenant_id" db:"tenant_id"`

	// timestamps of creation, updation
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IngestionWriter is an interface which defines the Write method for writing ingestion data to 1 or multiple destinations.
// All ingestion sources must implement this interface using composition with a concrete IngestionWriter.
// A destination can be defined as a file, a database, a ingestor, etc.
type IngestionWriter interface {
	Write(ctx context.Context, data IngestionData) (result interface{}, err error)
}

// Ingestor is an interface which defines the SetupAndListen method to read ingestion data from multiple
// sources and writing it to one or multiple destinations. All ingestion data receivers should implement this interface
// using composition with a concrete Ingestor.
type Ingestor interface {
	// This method could be blocking
	SetupAndListen(IngestionDataProcessor, IngestionWriter) error
}
