package ingestion

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DatabaseIngestionWriter struct {
	db        *sqlx.DB
	tableName string
}

func (diw *DatabaseIngestionWriter) Write(data IngestionData) (interface{}, error) {
	query := fmt.Sprintf(`INSERT INTO %s (uuid) VALUES (:uuid)`, diw.tableName)
	res, err := diw.db.NamedExec(query, data)
	return res, err
}
