package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"
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

	h.ConfigurePrey()(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
}

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("case 1: configure hunter successfully", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/hunter/configure-hunter", strings.NewReader(`{"speed": 0.0, "position": {"x": 0.0, "y": 0.0, "z": 0.0}}`))
		hd.ConfigureHunter().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, `{"message":"Hunter configured"}`, recorder.Body.String())
	})

	t.Run("case 2: configure hunter, bad request", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", strings.NewReader(`{peed 0.0, "position": {"x": 0.0, "y": 0.0, "z": 0.0}}`))
		hd.ConfigureHunter().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, `{"status":"Bad Request","message":"Incorrect hunter configuration, issues with request body"}`, recorder.Body.String())
	})
}

func TestHunter_Hunt(t *testing.T) {

	t.Run("case 1: hunter can hunt prey", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(prey prey.Prey) (duration float64, err error) { return 2000, nil }
		hd := handler.NewHunter(ht, prey.NewPreyStub())

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		hd.Hunt().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, `{"duration":2000,"message":"hunting completed"}`, recorder.Body.String())
	})

	t.Run("case 2: hunter can not hunt prey", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(prey prey.Prey) (duration float64, err error) { return 0, hunter.ErrCanNotHunt }
		hd := handler.NewHunter(ht, prey.NewPreyStub())

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		hd.Hunt().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, `{"duration":0,"error":"can not hunt the prey","message":"hunting completed"}`, recorder.Body.String())
	})

	t.Run("case 3: internal server error", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(prey prey.Prey) (duration float64, err error) { return 0, errors.New("internal server error") }
		hd := handler.NewHunter(ht, prey.NewPreyStub())

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		hd.Hunt().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, `{"status":"Internal Server Error","message":"internal server error"}`, recorder.Body.String())
	})
}
