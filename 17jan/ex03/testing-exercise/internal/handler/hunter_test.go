package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestHunter_ConfigurePrey_BadRequest(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewReader([]byte(`{name}`)))
	if err != nil {
		t.Fatalf("Não é possível criar a requisição: %v", err)
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

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Erro ao decodificar JSON")

}

func TestHunter_Hunt(t *testing.T) {

	requestBodyHunter := RequestBodyConfigHunter{
		Speed: 12.0,
		Position: &positioner.Position{
			X: 0.0,
			Y: 0.0,
		},
	}
	bodyHunter, _ := json.Marshal(requestBodyHunter)

	reqHunter, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(bodyHunter))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição para configurar o caçador: %v", err)
	}

	recorderHunter := httptest.NewRecorder()

	ps := positioner.NewPositionerDefault()
	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     12.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(5.0, &positioner.Position{X: 10.0, Y: 10.0, Z: 0.0})

	h := NewHunter(ht, pr)

	h.ConfigureHunter()(recorderHunter, reqHunter)

	// start hunt
	reqHunt, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição para iniciar a caça: %v", err)
	}

	recorderHunt := httptest.NewRecorder()

	h.Hunt()(recorderHunt, reqHunt)

	assert.Equal(t, http.StatusOK, recorderHunt.Code)

	assert.Contains(t, recorderHunt.Body.String(), "Caça concluída")

}

func TestHunter_Hunt_LowerSpeed(t *testing.T) {

	requestBodyHunter := RequestBodyConfigHunter{
		Speed: 1.0,
		Position: &positioner.Position{
			X: 0.0,
			Y: 0.0,
		},
	}
	bodyHunter, _ := json.Marshal(requestBodyHunter)

	// configure hunter
	reqHunter, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(bodyHunter))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição para configurar o caçador: %v", err)
	}

	recorderHunter := httptest.NewRecorder()

	ps := positioner.NewPositionerDefault()
	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     1.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})

	// speed tuna is greater than speed hunter
	pr := prey.NewTuna(5.0, &positioner.Position{X: 10.0, Y: 10.0, Z: 0.0})

	h := NewHunter(ht, pr)

	// Configure hunter
	h.ConfigureHunter()(recorderHunter, reqHunter)

	reqHunt, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição para iniciar a caça: %v", err)
	}

	recorderHunt := httptest.NewRecorder()

	h.Hunt()(recorderHunt, reqHunt)

	assert.Equal(t, http.StatusOK, recorderHunt.Code)

	assert.Contains(t, recorderHunt.Body.String(), "Caça concluída")

}

func TestHunter_ConfigurePrey_InternalServerError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/hunter/hunt", nil)
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	recorder := httptest.NewRecorder()

	// Configure hunter
	ht := &mockHunter{
		shouldFail: true, // error should be returned
	}

	// Configure prey
	pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

	h := NewHunter(ht, pr)

	h.Hunt()(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	expectedResponse := "internal server error"
	assert.Contains(t, recorder.Body.String(), expectedResponse)
}

type mockHunter struct {
	shouldFail bool
}

func (m *mockHunter) Configure(speed float64, position *positioner.Position) {

}

func (m *mockHunter) Hunt(pr prey.Prey) (float64, error) {
	if m.shouldFail {
		return 0, fmt.Errorf("erro inesperado na caça")
	}
	return 1.50, nil
}
