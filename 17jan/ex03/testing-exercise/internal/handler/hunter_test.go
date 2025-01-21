package handler

import (
	"bytes"
	"encoding/json"
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

func Test_ConfigureHunterTestHunte_Ok(t *testing.T) {
	ht := hunter.NewHunterMock()
	hd := NewHunter(ht, nil)
	hdFunc := hd.ConfigureHunterTestHunter()

	req := httptest.NewRequest("POST", "/hunter/configure-hunter", strings.NewReader(
		`{"speed": 10.0, "position": {"X": 1.0, "Y": 2.0, "Z": 3.0}}`,
	))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	hdFunc(response, req)
	println(response.Body.String())

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Hunt_Ok(t *testing.T) {

	ht := hunter.NewHunterMock()

	hd := NewHunter(ht, nil)
	hdFunc := hd.Hunt()

	req := httptest.NewRequest("POST", "/hunter/configure-hunter", nil)

	response := httptest.NewRecorder()
	hdFunc(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

}
