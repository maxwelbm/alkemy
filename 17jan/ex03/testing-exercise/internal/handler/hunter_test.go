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

func TestConfigureHunter(t *testing.T) {
	var body []byte
	req, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(body))
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

	h.ConfigureHunter(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHunt(t *testing.T) {
	t.Run("caçador consegue capturar a presa", func(t *testing.T) {
		var body []byte
		req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}
		recorder := httptest.NewRecorder()

		ps := positioner.NewPositionerDefault()

		sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
			MaxTimeToCatch: 10,
			Positioner:     ps,
		})

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     5.0,
			Position:  &positioner.Position{X: 10, Y: 20, Z: 30},
			Simulator: sm,
		})

		pr := prey.NewTuna(0.3, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

		h := NewHunter(ht, pr)

		h.Hunt(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "caça concluída")
		assert.Contains(t, recorder.Body.String(), "com sucesso")
	})
	t.Run("caçador não consegue capturar a presa", func(t *testing.T) {
		var body []byte
		req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}
		recorder := httptest.NewRecorder()

		ps := positioner.NewPositionerDefault()

		sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
			MaxTimeToCatch: 10,
			Positioner:     ps,
		})

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     5.0,
			Position:  &positioner.Position{X: 10, Y: 20, Z: 30},
			Simulator: sm,
		})

		pr := prey.NewTuna(6.3, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

		h := NewHunter(ht, pr)

		h.Hunt(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "caça concluída")
		assert.Contains(t, recorder.Body.String(), "sem sucesso")
	})
}
