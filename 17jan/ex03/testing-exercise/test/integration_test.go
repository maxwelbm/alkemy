package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testdoubles/internal/application"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	go func() {
		app := application.NewApplicationDefault(":8080")
		if err := app.SetUp(); err != nil {
			panic(err)
		}
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	m.Run()
}

func TestIntegration_Full(t *testing.T) {
	t.Run("configurar presa", func(t *testing.T) {
		preyConfig := handler.RequestBodyConfigPrey{
			Speed: 5.0,
			Position: &positioner.Position{
				X: 10.0,
				Y: 20.0,
				Z: 0.0,
			},
		}
		body, err := json.Marshal(preyConfig)
		require.NoError(t, err)

		resp, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "A presa está configurada corretamente", string(responseBody))
	})

	t.Run("configurar caçador", func(t *testing.T) {
		hunterConfig := handler.RequestBodyConfigHunter{
			Speed: 10.0,
			Position: &positioner.Position{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
			},
		}
		body, err := json.Marshal(hunterConfig)
		require.NoError(t, err)

		resp, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		responseBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "O caçador está configurado corretamente", string(responseBody))
	})

	t.Run("realizar caçada", func(t *testing.T) {
		resp, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Contains(t, responseBody, "message")
		assert.Contains(t, responseBody, "duration")
		assert.Equal(t, "caça concluída", responseBody["message"])
	})
}

func TestIntegration_Invalid(t *testing.T) {
	t.Run("erro configuração caçador", func(t *testing.T) {
		bodyInvalido := []byte(`{"speed": "invalido"}`)

		resp, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBuffer(bodyInvalido))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
