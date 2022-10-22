package ingestion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type SimpleHttpIngestionWriter struct {
	url string
}

func (h *SimpleHttpIngestionWriter) Write(ctx context.Context, data IngestionData) (interface{}, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	postData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctxWithTimeout)
	// We can check response status codes here for errors, but it is a "Simple"HttpIngestionWriter after all
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
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
	// We should add healthchecks, but is out of scope for "Simple"HttpIngestor for now.
	mux.Handle(h.path, h.setupHandler(writer))
	return mux
}

func (h *SimpleHttpIngestor) SetupAndListen(writer IngestionWriter) error {
	if writer == nil {
		return errors.New("handler cannot be nil")
	}
	mux := h.setupMux(writer)
	server := h.getServer(fmt.Sprintf("%d", h.port), mux)
	log.Println("starting server on port", h.port)
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
