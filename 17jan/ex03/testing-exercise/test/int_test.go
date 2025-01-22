//go:build special

// go:build special
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

func init() {
	go func() {
		server := application.NewApplicationDefault(":8080")

		err := server.SetUp()
		if err != nil {
			panic(err)
		}

		err = server.Run()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func TestConfigurePreyRequest(t *testing.T) {
	pr := handler.RequestBodyConfigPrey{
		Speed: 2.1,
		Position: &positioner.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}
	body, err := json.Marshal(pr)
	require.NoError(t, err)

	res, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer res.Body.Close()

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(res.Body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	expectedMessage := "A presa está configurada corretamente"
	require.Equal(t, expectedMessage, responseBody.String())
}

func TestConfigureHunterRequest(t *testing.T) {
	ht := handler.RequestBodyConfigHunter{
		Speed: 100.5,
		Position: &positioner.Position{
			X: 0.0,
			Y: 0.0,
			Z: 0.0,
		},
	}
	hunterBody, err := json.Marshal(ht)
	require.NoError(t, err)

	res, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBuffer(hunterBody))
	require.NoError(t, err)
	defer res.Body.Close()

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(res.Body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	expectedMessage := "O caçador está configurado corretamente"
	require.Equal(t, expectedMessage, responseBody.String())
}

func TestHuntRequest(t *testing.T) {
	res, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", nil)
	require.NoError(t, err)
	defer res.Body.Close()

	var resJSON struct {
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	err = json.NewDecoder(res.Body).Decode(&resJSON)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "prey hunted", resJSON.Message)
	require.True(t, resJSON.Data["success"].(bool))
	require.NotNil(t, resJSON.Data["duration"])
}
