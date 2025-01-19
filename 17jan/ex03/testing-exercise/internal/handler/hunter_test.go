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

func TestHunter_ConfigureHunter(t *testing.T) {
	requestBody := RequestBodyConfigHunter{
		Speed: 150.0,
		Position: &positioner.Position{
			X: 300,
			Y: 300,
		},
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodGet, "/hunter/configure-hunter", bytes.NewReader(body))
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

	config := h.ConfigureHunter()
	config(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "O caçador está configurado corretamente", recorder.Body.String())
}

func TestHunter_FailConfigureHunter(t *testing.T) {
	requestBody := map[string]string{
		"speed": "wrong value",
	}

	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodGet, "/hunter/configure-hunter", bytes.NewReader(body))
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

	config := h.ConfigureHunter()
	config(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestHunter_Hunt(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	recorder := httptest.NewRecorder()

	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		MaxTimeToCatch: 100,
		Positioner:     ps,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     3.0,
		Position:  &positioner.Position{X: 0.0, Y: 50, Z: 0.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.01, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

	h := NewHunter(ht, pr)

	h.Hunt()(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, recorder.Body.String(), "A caçada foi concluída com sucesso em 16.72 segundos")
}
