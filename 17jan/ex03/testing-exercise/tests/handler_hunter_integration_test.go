//go:build integration

package tests_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerIntegration_ConfigurePrey(t *testing.T) {
	cfgPrey := handler.RequestBodyConfigPrey{
		Speed: 5.0,
		Position: &positioner.Position{
			X: 5,
			Y: 5,
			Z: 5,
		},
	}
	b, err := json.Marshal(cfgPrey)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/hunter/configure-prey", bytes.NewReader(b))
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	b, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "A presa está configurada corretamente", string(b))
}

func TestHandlerIntegration_ConfigureHunter(t *testing.T) {
	cfgHunter := handler.RequestBodyConfigHunter{
		Speed: 5.0,
		Position: &positioner.Position{
			X: 5,
			Y: 5,
			Z: 5,
		},
	}
	expectedBody := `
	{
		"message": "hunter configured successfully"
	}
	`

	b, err := json.Marshal(cfgHunter)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/hunter/configure-hunter", bytes.NewReader(b))
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	b, err = io.ReadAll(res.Body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.JSONEq(t, expectedBody, string(b))
}

func TestHandlerIntegration_Hunt(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedBody := `{"message":"hunt executed successfully","success":true,"time_took":1.7320508075688774}`

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/hunter/hunt", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, expectedStatus, res.StatusCode)
	require.Equal(t, expectedBody, string(b))
}
