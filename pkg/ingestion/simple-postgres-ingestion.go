package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SimplePostgresIngestionWriter struct {
	db        *sqlx.DB
	uri       string
	tableName string
}

func (diw *SimplePostgresIngestionWriter) Write(ctx context.Context, data IngestionData) (interface{}, error) {
	if diw.db == nil {
		db, err := sqlx.Connect("postgres", diw.uri)
		if err != nil {
			log.Println(err)
		}
		diw.db = db
	}
	query := fmt.Sprintf(
		`
		INSERT INTO %s (uuid, source_type, source_id, data, data_id, metadata, metadata_type, lang, tenant_id, created_at, updated_at)
		VALUES
		(:uuid, :source_type, :source_id, :data, :data_id, :metadata, :metadata_type, :lang, :tenant_id, :created_at, :updated_at)
		`, diw.tableName)
	marshaledMetadata, err := json.Marshal(data.Metadata)
	if err != nil {
		log.Println(err)
	}
	data.Metadata = marshaledMetadata
	res, err := diw.db.NamedExec(query, data)
	log.Println(res, err)
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	return res, err
}

func NewSimplePostgresIngestionWriter(uri string, tableName string) IngestionWriter {
	return &SimplePostgresIngestionWriter{
		uri:       uri,
		tableName: tableName,
	}
}
