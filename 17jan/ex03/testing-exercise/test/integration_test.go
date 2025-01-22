package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func startTestServer() *http.Server {
// 	router := chi.NewRouter()

// 	router.Post("/hunter/configure-prey", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("The prey is configured correctl"))
// 	})

// 	server := &http.Server{
// 		Addr:    ":8080",
// 		Handler: router,
// 	}

// 	go func() {
// 		log.Println("Starting server...")
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Server failed to start: %v", err)
// 		}
// 	}()

// 	time.Sleep(100 * time.Millisecond)

// 	return server
// }

func TestPreyEndpoint(t *testing.T) {
	// server := startTestServer()
	// defer func() {
	// 	log.Println("Shutting down server...")
	// 	if err := server.Close(); err != nil {
	// 		log.Fatalf("Failed to close server: %v", err)
	// 	}
	// }()

	payload := `{"speed": 10, "position":{"X": 100, "Y": 200, "Z": 300}}`

	resp, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "The prey is configured correctly")
}

func TestHunterEndpoint(t *testing.T) {
	payload := `{"speed": 10, "position":{"X": 100, "Y": 200, "Z": 300}}`

	resp, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "The shark is configured correctly")
}

func TestHuntEndpoint(t *testing.T) {
	payload := `{0, nil}`

	resp, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, string(body), "Internal error: ")
}
