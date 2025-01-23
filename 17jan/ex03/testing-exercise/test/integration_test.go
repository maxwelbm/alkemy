package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_ConfigurePrey(t *testing.T) {
	payload := `{"speed":10, "position":{"X":200, "Y":200, "Z":200}}`

	resp, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBufferString((payload)))
	if err != nil {
		t.Fatalf("Erro ao fazer a request %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler a response %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(body), "A presa está configurada corretamente")
}

func TestIntegration_ConfigurePrey_InvalidPayload(t *testing.T) {
	payload := `{"speed": "aaaa", "position":{"X":200, "Y":200, "Z":200}}`

	resp, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Erro ao fazer a request %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler a response %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(body), "Erro ao decodificar JSON:")
}

func TestIntegration_ConfigureHunter(t *testing.T) {
	payload := `{"speed":20.1, "position":{"X":100, "Y":150, "Z":300}}`

	resp, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Erro ao fazer a request %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler a response %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(body), "O tubarão está configurado corretamente")
}

func TestIntegration_ConfigureHunter_InvalidPayload(t *testing.T) {
	payload := `{"speed":"aaaaa", "position":{"X":100, "Y":150, "Z":300}}`

	resp, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Erro ao fazer a request %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler a response %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(body), "Erro ao decodificar JSON:")
}

func TestIntegration_Hunt(t *testing.T) {
	payload := ""

	resp, err := http.Post("http://localhost:8080/hunter/hunt", "text/plain", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Erro ao fazer a request %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler a response %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "caçada concluída, capturada:")
}
