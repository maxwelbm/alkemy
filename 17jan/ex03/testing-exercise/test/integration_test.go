//go:build integration
// +build integration

package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	cases := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "case 1: configure prey successfully",
			body:         `{"speed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`,
			expectedCode: http.StatusOK,
			expectedBody: "A presa está configurada corretamente",
		},
		{
			name:         "case 2: configure prey with bad JSON",
			body:         `peed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`,
			expectedCode: 400,
			expectedBody: `{"status":"Bad Request","message":"Erro ao decodificar JSON: invalid character 'p' looking for beginning of value"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/hunter/configure-prey", bytes.NewReader([]byte(c.body)))
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

			router := chi.NewRouter()
			router.HandleFunc("/hunter/configure-prey", h.ConfigurePrey())
			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.expectedCode, recorder.Code)
			assert.Equal(t, c.expectedBody, recorder.Body.String())
		})
	}
}

func TestHunter_ConfigureHunter(t *testing.T) {
	cases := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "case 1: configure hunter successfully",
			body:         `{"speed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"caçador configurado"}`,
		},
		{
			name:         "case 2: configure hunter with bad JSON",
			body:         `peed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`,
			expectedCode: 400,
			expectedBody: `{"status":"Bad Request","message":"Erro ao decodificar JSON: invalid character 'p' looking for beginning of value"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/hunter/configure-hunter", bytes.NewReader([]byte(c.body)))
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

			router := chi.NewRouter()
			router.HandleFunc("/hunter/configure-hunter", h.ConfigureHunter())
			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.expectedCode, recorder.Code)
			assert.Equal(t, c.expectedBody, recorder.Body.String())
		})
	}
}

func TestHunter_Hunt(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		expectedBody string
		tunaPosition *positioner.Position
	}{
		{
			name:         "case 1: hunter can hunt pray",
			expectedCode: 200,
			expectedBody: `{"duration":1,"message":"caça concluída"}`,
			tunaPosition: &positioner.Position{X: 3.0, Y: 0.0, Z: 0.0},
		},
		{
			name:         "case 2: hunter can't hunt pray",
			expectedCode: 200,
			expectedBody: `{"duration":0,"error":"can not hunt the prey","message":"caça concluída"}`,
			tunaPosition: &positioner.Position{X: 900.0, Y: 2000.0, Z: 30000.0},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/hunter/hunt", nil)
			recorder := httptest.NewRecorder()
			ps := positioner.NewPositionerDefault()
			sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				MaxTimeToCatch: 3,
				Positioner:     ps,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     4,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sm,
			})
			pr := prey.NewTuna(1, c.tunaPosition)
			h := handler.NewHunter(ht, pr)

			router := chi.NewRouter()
			router.HandleFunc("/hunter/hunt", h.Hunt())
			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.expectedCode, recorder.Code)
			assert.Equal(t, c.expectedBody, recorder.Body.String())
		})
	}

}
