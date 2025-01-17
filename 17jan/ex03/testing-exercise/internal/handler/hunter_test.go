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

func TestHunter_ConfigurePrey(t *testing.T) {
	requestBody := handler.RequestBodyConfigPrey{
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

	h := handler.NewHunter(ht, pr)

	h.ConfigurePrey(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
}

func TestHunter_ConfigurePrey_StatusBadRequest(t *testing.T) {
	requestBody := "s"
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

	h := handler.NewHunter(ht, pr)

	h.ConfigurePrey(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Erro ao decodificar JSON")
}

func TestHunter_ConfigureHunter(t *testing.T) {
	requestBody := handler.RequestBodyConfigHunter{
		Speed: 15.0,
		Position: &positioner.Position{
			X: 50,
			Y: 75,
			Z: 100,
		},
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

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

	h := handler.NewHunter(ht, pr)

	h.ConfigureHunter()(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "O hunter está configurado corretamente", recorder.Body.String())
}

func TestHunter_ConfigureHunter_StatusBadRequest(t *testing.T) {
	requestBody := "a"
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

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

	h := handler.NewHunter(ht, pr)

	h.ConfigureHunter()(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Erro ao decodificar JSON")
}

func TestHunter_Hunt(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

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

	h := handler.NewHunter(ht, pr)

	h.Hunt()(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "A presa foi caçada em")
}
