//go:build integration

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testdoubles/internal/application"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const HunterEndpoint = "http://localhost:8080/hunter"

func waitForServer(timeout time.Duration) {
	deadline := time.Now().Add(timeout * time.Second)

	for time.Now().Before(deadline) {
		resp, err := http.DefaultClient.Get("http://localhost:8080/healthcheck")
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	log.Fatal("Failed to connect to handler")
}

func init() {
	go func() {
		app := application.NewApplicationDefault(":8080")
		if err := app.SetUp(); err != nil {
			log.Fatal(err)
		}

		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	waitForServer(5)
}

func TestHandlerIntegrationPrey(t *testing.T) {
	cfgPrey := handler.RequestBodyConfigPrey{
		Speed: 5.0,
		Position: &positioner.Position{
			X: 2,
			Y: 2,
			Z: 5,
		},
	}
	b, err := json.Marshal(cfgPrey)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, HunterEndpoint+"/configure-prey", bytes.NewReader(b))
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	b, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "A presa está configurada corretamente", string(b))
}

func TestHandlerIntegrationHunter(t *testing.T) {
	cfgHunter := handler.RequestBodyConfigHunter{
		Speed: 5.0,
		Position: &positioner.Position{
			X: 2,
			Y: 2,
			Z: 5,
		},
	}
	expectedBody := `
	{
		"message": "Caçador configurado com sucesso"
	}
	`
	b, err := json.Marshal(cfgHunter)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, HunterEndpoint+"/configure-hunter", bytes.NewReader(b))
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	b, err = io.ReadAll(res.Body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.JSONEq(t, expectedBody, string(b))
}

func TestHandlerIntegrationHunt(t *testing.T) {
	expectedStatus := http.StatusOK

	req, err := http.NewRequest(http.MethodPost, HunterEndpoint+"/hunt", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, expectedStatus, res.StatusCode)
}
