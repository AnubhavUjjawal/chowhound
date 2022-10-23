package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
)

type StdoutIngestionWriter struct{}

func (w *StdoutIngestionWriter) Write(ctx context.Context, data IngestionData) (interface{}, error) {
	dataJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data, will print json", err)
		fmt.Println(data)
	} else {
		fmt.Println(string(dataJson))
	}
	return nil, nil
}

func NewStdoutIngestionWriter() *StdoutIngestionWriter {
	return &StdoutIngestionWriter{}
}
