package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
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

	req, err := http.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewReader(body))
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

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		var body []byte

		req, err := http.NewRequest(http.MethodGet, "/hunter/configure-hunter", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     3.0,
			Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
			Simulator: nil,
		})

		h := NewHunter(ht, nil)

		h.ConfigureHunter(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, `{"status":"Bad Request","message":"Configuração incorreta do caçador, problemas com o corpo da requisição"}`, recorder.Body.String())
	})
}

func TestHunter_Hunt(t *testing.T) {

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

	t.Run("hunt success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseBody map[string]any
		err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("Erro ao desserializar o corpo da resposta: %v", err)
		}

		assert.Equal(t, "caca concluida", responseBody["message"])
		assert.True(t, responseBody["success"].(bool))
		assert.NotNil(t, responseBody["duration"])
	})

	t.Run("hunt error", func(t *testing.T) {
		mockHunter := &hunter.HunterMock{
			HuntFunc: func(pr prey.Prey) (float64, error) {
				return 0.0, hunter.ErrCanNotHunt
			},
		}
		h.ht = mockHunter

		req, err := http.NewRequest(http.MethodPost, "/", nil)
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseBody map[string]any
		err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("Erro ao desserializar o corpo da resposta: %v", err)
		}

		assert.Equal(t, "caca concluida", responseBody["message"])
		assert.False(t, responseBody["success"].(bool))
		assert.NotNil(t, responseBody["duration"])
	})
}
