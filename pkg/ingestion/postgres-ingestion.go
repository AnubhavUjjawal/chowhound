package ingestion

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type SimplePostgresIngestionWriter struct {
	db        *sqlx.DB
	uri       string
	tableName string
}

func (diw *SimplePostgresIngestionWriter) Write(data IngestionData) (interface{}, error) {
	if diw.db == nil {
		db, err := sqlx.Connect("postgres", diw.uri)
		if err != nil {
			log.Fatalln(err)
		}
		diw.db = db
	}
	query := fmt.Sprintf(`INSERT INTO %s (uuid) VALUES (:uuid)`, diw.tableName)
	res, err := diw.db.NamedExec(query, data)
	return res, err
}

func NewSimplePostgresIngestionWriter(uri string, tableName string) *SimplePostgresIngestionWriter {
	return &SimplePostgresIngestionWriter{
		uri:       uri,
		tableName: tableName,
	}
}
