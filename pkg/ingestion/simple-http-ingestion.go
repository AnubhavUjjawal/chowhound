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

func PostData(url string, ctx context.Context, data interface{}) (int, error) {
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	return resp.StatusCode, nil
}

func (h *SimpleHttpIngestionWriter) Write(ctx context.Context, data IngestionData) (interface{}, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return PostData(h.url, ctxWithTimeout, data)
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

func (h *SimpleHttpIngestor) getHandlerFunc(processor IngestionDataProcessor, writer IngestionWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			log.Println("only POST is supported")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var data IngestionData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Println("error decoding json", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if processor != nil {
			data, err = processor(data)
			if err != nil {
				log.Println("error processing data", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// fmt.Println("data", data.SourceId)
		_, err = writer.Write(r.Context(), data)
		if err != nil {
			log.Println("error writing data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (h *SimpleHttpIngestor) getMux(handler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	// We should add healthchecks, but is out of scope for "Simple"HttpIngestor for now.
	mux.Handle(h.path, handler)
	return mux
}

func (h *SimpleHttpIngestor) SetupAndListen(processor IngestionDataProcessor, writer IngestionWriter) error {
	if writer == nil {
		return errors.New("handler cannot be nil")
	}
	mux := h.getMux(h.getHandlerFunc(processor, writer))
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
