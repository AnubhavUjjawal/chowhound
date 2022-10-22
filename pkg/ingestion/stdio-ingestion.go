package ingestion

import (
	"context"
	"fmt"
)

type StdoutIngestionWriter struct{}

func (w *StdoutIngestionWriter) Write(ctx context.Context, data IngestionData) (interface{}, error) {
	fmt.Println("received data: ", data)
	return nil, nil
}

func NewStdoutIngestionWriter() *StdoutIngestionWriter {
	return &StdoutIngestionWriter{}
}
