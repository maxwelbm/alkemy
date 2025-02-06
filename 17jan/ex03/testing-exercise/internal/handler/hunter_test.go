package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	t.Run("Configure Prey successfully", func(t *testing.T) {
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
	})

	t.Run("Bad request", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hunter/configure-prey", strings.NewReader(`{"speed": "0.0", "position": {"X": 0.0, "Y": 0.0, "Z": 0.0}}`))

		hd.ConfigurePrey(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
	})
}

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("Bad Request", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/configure-hunter", strings.NewReader(`{"speed": "0.0", "position": {"X": 0.0, "Y": 0.0, "Z": 0.0}}`))

		hd.ConfigureHunter(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
	})

	t.Run("Configure Hunter successfully", func(t *testing.T) {
		requestBody := handler.RequestBodyConfigHunter{
			Speed: 100.0,
			Position: &positioner.Position{
				X: 10,
				Y: 10,
				Z: 10,
			},
		}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("GET", "/configure-hunter", bytes.NewReader(body))

		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()

		hd.ConfigureHunter(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, `{"message":"O caçador está configurado corretamente"}`, res.Body.String())
	})
}

func TestHunter_Hunt(t *testing.T) {
	t.Run("Hunt Completed", func(t *testing.T) {
		pr := prey.NewPreyStub()
		ht := hunter.NewHunterMock()
		res := httptest.NewRecorder()
		hd := handler.NewHunter(ht, pr)
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) { return 100, nil }

		hd.Hunt().ServeHTTP(res, &http.Request{Method: http.MethodGet})

		require.Equal(t, http.StatusOK, res.Code)
		require.Equal(t, `{"duration":100,"message":"caça concluída"}`, strings.TrimSpace(res.Body.String()))
	})

	t.Run("Hunt Failed", func(t *testing.T) {
		pr := prey.NewPreyStub()
		ht := hunter.NewHunterMock()
		res := httptest.NewRecorder()
		hd := handler.NewHunter(ht, pr)
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) { return 0, hunter.ErrCanNotHunt }

		hd.Hunt().ServeHTTP(res, &http.Request{Method: http.MethodGet})

		require.Equal(t, http.StatusOK, res.Code)
		require.Equal(t, `{"duration":0,"error":"can not hunt the prey","message":"caça concluída"}`, strings.TrimSpace(res.Body.String()))
	})

	t.Run("Bad Request", func(t *testing.T) {
		pr := prey.NewPreyStub()
		ht := hunter.NewHunterMock()
		res := httptest.NewRecorder()
		hd := handler.NewHunter(ht, pr)
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) { return 0, errors.New("internal server error") }

		hd.Hunt().ServeHTTP(res, &http.Request{Method: http.MethodGet})

		require.Equal(t, http.StatusInternalServerError, res.Code)
		require.Equal(t, `{"status":"Internal Server Error","message":"internal server error"}`, strings.TrimSpace(res.Body.String()))
	})
}
