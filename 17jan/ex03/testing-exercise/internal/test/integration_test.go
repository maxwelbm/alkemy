package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testdoubles/internal/application"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	endpoint = "http://localhost:8080/hunter/configure-"
	appJson = "application/json"
	expectedHunterResp = "O caçador está configurado corretamente"
	expectedPreyResp = "A presa está configurada corretamente"
)

func init() {
	go func() {
		server := application.NewApplicationDefault(":8080")

		if err := server.SetUp(); err != nil {
			panic(err)
		}

		if err := server.Run(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func TestConfigurePrey(t *testing.T) {
	prey := handler.RequestBodyConfigPrey{Speed: 150.5, Position: &positioner.Position{X: 50.0, Y: 50.0, Z: 50.0},}
	
	body, err := json.Marshal(prey)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	response, err := http.Post(endpoint+"prey", appJson, bytes.NewBuffer(body))
	require.NoError(t, err)
	
	resp := new(bytes.Buffer)
	if _, err = resp.ReadFrom(response.Body); err != nil {
		t.Fatalf("Falha ao ler o body da requisição: %v", err)
	}
	
	require.Equal(t, expectedPreyResp, resp.String())
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestConfigureHunter(t *testing.T) {
	hunter := handler.RequestBodyConfigHunter{Speed: 180.5, Position: &positioner.Position{X: 40.0, Y: 40.0, Z: 40.0,},}
	
	body, err := json.Marshal(hunter)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	resp, err := http.Post(endpoint+"hunter", appJson, bytes.NewBuffer(body))
	require.NoError(t, err)
	
	response := new(bytes.Buffer)
	if _, err = response.ReadFrom(resp.Body); err != nil {
		t.Fatalf("Falha ao ler o body da requisição: %v", err)
	}
	
	require.Equal(t, expectedHunterResp, response.String())
	require.Equal(t, http.StatusOK, resp.StatusCode)
	
}
