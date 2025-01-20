package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
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
	requestBody := RequestBodyConfigHunter{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
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

	h := NewHunter(ht, pr)

	h.ConfigureHunter(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "A cacador está configurada corretamente", recorder.Body.String())
}

func TestHunter_ConfigureHunterFailed(t *testing.T) {
	requestBody := map[string]string{
		"speed": "teste corpo errado",
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

	h := NewHunter(ht, pr)

	h.ConfigureHunter(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHunt(t *testing.T) {
	ht := hunter.NewHunterMock()
	ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
		return 100.0, nil
	}
	hd := NewHunter(ht, nil)
	request := httptest.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	hd.Hunt(response, request)
	expectedCode := http.StatusOK
	require.Equal(t, expectedCode, response.Code)
	require.Equal(t, "\"A caçada foi um sucesso terminou em 100.00\"", strings.TrimSpace(response.Body.String()))
}

func TestHuntFailed(t *testing.T) {
	ht := hunter.NewHunterMock()
	ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
		return 100.0, hunter.ErrCanNotHunt
	}
	hd := NewHunter(ht, nil)
	request := httptest.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	hd.Hunt(response, request)
	expectedCode := http.StatusOK
	require.Equal(t, expectedCode, response.Code)
	require.Equal(t, "\"A caçada foi um fracasso terminou em 100.00\"", strings.TrimSpace(response.Body.String()))
}

func TestHuntWrong(t *testing.T) {
	ht := hunter.NewHunterMock()
	ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
		return 100.0, errors.New("erro diferente")
	}
	hd := NewHunter(ht, nil)
	request := httptest.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	hd.Hunt(response, request)
	expectedCode := http.StatusInternalServerError
	require.Equal(t, expectedCode, response.Code)
}
