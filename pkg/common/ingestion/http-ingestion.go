package ingestion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type SimpleHttpIngestionWriter struct {
	url string
}

func (h *SimpleHttpIngestionWriter) Write(ctx context.Context, data IngestionData) (interface{}, error) {
	postData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	// We can check response status codes here for errors, but I am leaving that for the sake of simplicity
	res, err := http.DefaultClient.Do(req)
	return res.StatusCode, err
}

type SimpleHttpIngestor struct {
	path string
	port int
}

func (h *SimpleHttpIngestor) getServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}

func (h *SimpleHttpIngestor) setupHandler(writer IngestionWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var data IngestionData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = writer.Write(r.Context(), data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (h *SimpleHttpIngestor) setupMux(writer IngestionWriter) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(h.path, h.setupHandler(writer))
	return mux
}

func (h *SimpleHttpIngestor) SetupAndListen(handler IngestionWriter) error {
	if handler == nil {
		return errors.New("handler cannot be nil")
	}
	mux := h.setupMux(handler)
	server := h.getServer(fmt.Sprintf("%d", h.port), mux)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error running http server: %s\n", err)
		}
	}
	return nil
}

func NewSimpleHttpIngestor(path string, port int) Ingestor {
	return &SimpleHttpIngestor{
		path: path,
		port: port,
	}
}

func NewSimpleHttpIngestionWriter(url string) IngestionWriter {
	return &SimpleHttpIngestionWriter{
		url: url,
	}
}
