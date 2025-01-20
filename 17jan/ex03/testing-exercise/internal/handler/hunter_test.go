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

// A presa está configurada corretamente, indicando uma mensagem "presa configurada". Retorna um 200.
func TestHunter_ConfigurePrey(t *testing.T) {
	requestBody := handler.RequestBodyConfigPrey{
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

	h := handler.NewHunter(ht, pr)

	h.ConfigurePrey(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
}

// Configuração incorreta do caçador, problemas com o corpo da solicitação. Retorna um 400
func TestHunter_ConfigureHunter_InvalidRequest(t *testing.T) {
	requestBody := handler.RequestBodyConfigHunter{
		Speed: 0,
		Position: nil,
	}
	body, _ := json.Marshal(requestBody)

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

	h := handler.NewHunter(ht, pr)

	h.ConfigureHunter(recorder, req)

	expectedBody := `{"status":"Bad Request","message":"caçador está configurado incorretamente"}`
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, expectedBody, recorder.Body.String())
}

// O tubarão consegue capturar a presa, exibe uma mensagem de "caça concluída" e dados com informações sobre se ela foi capturada, além do intervalo de tempo. Retorna um 200
func TestHunter_Hunt_Success(t *testing.T) {
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
		Position:  &positioner.Position{X: 1.0, Y: 1.0, Z: 1.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.6, &positioner.Position{X: 1.0, Y: 1.0, Z: 1.0})

	h := handler.NewHunter(ht, pr)

	h.Hunt(recorder, req)

	expectedBody := `{"data":{"duration":0,"success":true},"message":"caça concluída"}`
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expectedBody, recorder.Body.String())
}

// O tubarão não consegue caçar a presa, exibe uma mensagem de "caçada concluída", exibindo também os mesmos dados acima. Retorna um 200
func TestHunter_Hunt_Failure(t *testing.T) {
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
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.6, &positioner.Position{X: 10.0, Y: 10.0, Z: 10.0})

	h := handler.NewHunter(ht, pr)

	h.Hunt(recorder, req)

	expectedBody := `{"data":{"duration":0,"success":false},"message":"caça concluída"}`
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expectedBody, recorder.Body.String())
}
