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

func Test_Hunter_ConfigureHunter(t *testing.T) {
	t.Run("case 1: configure hunter, bad request", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/configure-hunter", strings.NewReader(`{"speed 0.0, "position": {"x": 0.0, "y": 0.0, "z": 0.0}}`))
		hd.ConfigureHunter().ServeHTTP(res, req)
		require.Equal(t, 400, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
	})
}

func Test_Hunter_Hunt(t *testing.T) {
	t.Run("case 1: hunter can hunt pray", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(prey prey.Prey) (duration float64, err error) {
			return 2000, nil
		}
		pr := prey.NewPreyStub()

		hd := handler.NewHunter(ht, pr)
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/hunt", nil)
		hd.Hunt().ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
		require.Equal(t, `{"duration":2000,"message":"caça concluída"}`, strings.TrimSpace(res.Body.String()))

	})
	t.Run("case 2: hunter can't hunt pray", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(prey prey.Prey) (duration float64, err error) {
			return 0, hunter.ErrCanNotHunt
		}
		pr := prey.NewPreyStub()

		hd := handler.NewHunter(ht, pr)
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/hunt", nil)
		hd.Hunt().ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
		require.Equal(t, `{"duration":0,"error":"can not hunt the prey","message":"caça concluída"}`, strings.TrimSpace(res.Body.String()))
	})
	t.Run("case 3: internal server error", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(prey prey.Prey) (duration float64, err error) {
			return 0, errors.New("internal server error")
		}
		pr := prey.NewPreyStub()

		hd := handler.NewHunter(ht, pr)
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/hunt", nil)
		hd.Hunt().ServeHTTP(res, req)
		require.Equal(t, 500, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
	})
}
