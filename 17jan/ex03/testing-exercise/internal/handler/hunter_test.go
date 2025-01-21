package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupHunter() (*handler.Hunter, *httptest.ResponseRecorder) {
	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     5.0,
		Position:  &positioner.Position{X: 10.0, Y: 20.0, Z: 30.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.6, &positioner.Position{X: 10.0, Y: 20.0, Z: 30.0})

	hnd := handler.NewHunter(ht, pr)

	recorder := httptest.NewRecorder()

	return hnd, recorder
}

func createRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	return http.NewRequest(method, url, &buf)
}

func TestHunter_ConfigurePrey(t *testing.T) {
	t.Run("Configure prey successfully", func(t *testing.T) {
		hnd, recorder := setupHunter()

		requestBody := handler.RequestBodyConfigPrey{
			Speed: 10.5,
			Position: &positioner.Position{
				X: 100,
				Y: 200,
			},
		}

		req, err := createRequest(http.MethodGet, "/hunter/configure-prey", requestBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		hnd.ConfigurePrey(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
	})

	t.Run("Configure prey with bad request", func(t *testing.T) {
		hnd, recorder := setupHunter()

		req, err := createRequest(http.MethodGet, "/hunter/configure-prey", "invalid")
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		hnd.ConfigurePrey(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Erro ao decodificar JSON")
	})
}

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("Configure hunter successfully", func(t *testing.T) {
		hnd, recorder := setupHunter()

		requestBody := handler.RequestBodyConfigHunter{
			Speed: 15.0,
			Position: &positioner.Position{
				X: 50,
				Y: 75,
				Z: 100,
			},
		}

		req, err := createRequest(http.MethodPost, "/hunter/configure-hunter", requestBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		hnd.ConfigureHunter()(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "O hunter está configurado corretamente", recorder.Body.String())
	})

	t.Run("Configure hunter with bad request", func(t *testing.T) {
		hnd, recorder := setupHunter()

		req, err := createRequest(http.MethodPost, "/hunter/configure-hunter", "invalid")
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		hnd.ConfigureHunter()(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Erro ao decodificar JSON")
	})
}

func TestHunter_Hunt(t *testing.T) {
	t.Run("Hunt successfully", func(t *testing.T) {
		hnd, recorder := setupHunter()

		req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		hnd.Hunt()(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "A presa foi caçada em")
	})
}
