package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	requestBody := RequestBodyConfigPrey{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodGet, "/hunter/configure-prey", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	recorder := httptest.NewRecorder()

	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     3.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

	h := NewHunter(ht, pr)

	h.ConfigurePrey(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
}

// Configuração incorreta do caçador, problemas com o corpo da solicitação. Retorna um 400
func TestHunter_IncorrectConfigure(t *testing.T) {

	t.Run("Configuração incorreta do caçador", func(t *testing.T) {
		ws := hunter.WhiteShark{}
		pt := prey.Tuna{}

		h := NewHunter(&ws, &pt)
		configureHunter := h.ConfigureHunter()
		statusCode := http.StatusBadRequest
		body := "invalid request body"

		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		require.NoError(t, err)

		res := httptest.NewRecorder()

		configureHunter(res, req)
		expectedBody := map[string]interface{}{"status": "Bad Request", "message": "invalid request body"}

		var response map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, statusCode, res.Code)
		assert.Equal(t, expectedBody, response)
	})
}

// O tubarão consegue capturar a presa, exibe uma mensagem de "caça concluída" e dados com informações sobre se ela foi capturada, além do intervalo de tempo. Retorna um 200
func TestHunter_HunterSuccess(t *testing.T) {

	t.Run("O tubarão consegue capturar a presa", func(t *testing.T) {
		ws := RequestBodyConfigHunter{
			Speed: 10.5,
			Position: &positioner.Position{
				X: 100,
				Y: 200,
			},
		}

		pt := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

		tutuba := hunter.NewHunterMock()
		tutuba.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 5.0, nil
		}

		h := NewHunter(tutuba, pt)
		body, err := json.Marshal(map[string]interface{}{"hunter": ws, "prey": pt})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		require.NoError(t, err)

		res := httptest.NewRecorder()

		h.Hunt()(res, req)

		expectedBody := map[string]interface{}{
			"status":       "Success",
			"message":      "prey hunted successfully",
			"preyCaptured": true,
			"duration":     5.0,
		}

		var response map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, response)
	})
}

// O tubarão não consegue caçar a presa, exibe uma mensagem de "caçada concluída", exibindo também os mesmos dados acima. Retorna um 200
func TestHunter_HunterNotSuccess(t *testing.T) {
	t.Run("O tubarão consegue capturar a presa", func(t *testing.T) {
		ws := RequestBodyConfigHunter{
			Speed: 1.5,
			Position: &positioner.Position{
				X: 100,
				Y: 200,
			},
		}

		pt := prey.NewTuna(11.5, &positioner.Position{X: 100.0, Y: 200.0, Z: 0.0})

		tutuba := hunter.NewHunterMock()
		tutuba.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 5.0, errors.New("prey not captured")
		}

		h := NewHunter(tutuba, pt)
		body, err := json.Marshal(map[string]interface{}{"hunter": ws, "prey": pt})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		require.NoError(t, err)

		res := httptest.NewRecorder()

		h.Hunt()(res, req)

		expectedBody := map[string]interface{}{
			"status":       "Success",
			"message":      "prey hunted successfully",
			"preyCaptured": false,
			"duration":     5.0,
		}

		var response map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, response)
	})
}

// Caso haja um erro diferente daquele que o caçador não conseguiu caçar, bem como um estagiário, retorne um 500
func TestHunter_HunterError(t *testing.T) {

	t.Run("o caçador não conseguiu caçar", func(t *testing.T) {
		ws := RequestBodyConfigHunter{
			Speed: 11.5,
			Position: &positioner.Position{
				X: 100,
				Y: 200,
			},
		}

		pt := prey.NewTuna(11.5, &positioner.Position{X: 100.0, Y: 200.0, Z: 0.0})

		tutuba := hunter.NewHunterMock()
		tutuba.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0, errors.New("error")
		}

		h := NewHunter(tutuba, pt)
		body, err := json.Marshal(map[string]interface{}{"hunter": ws, "prey": pt})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		require.NoError(t, err)

		res := httptest.NewRecorder()

		h.Hunt()(res, req)

		expectedBody := map[string]interface{}{
			"status":  "Error",
			"message": "internal server error",
		}

		var response map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, expectedBody, response)
	})

}
