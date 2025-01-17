package handler_test

import (
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
)

func TestHunter_ConfigurePrey(t *testing.T) {
	type wanted struct {
		statusCode int
		resBody    string
	}
	tests := []struct {
		name    string
		reqBody string
		mockFn  func(pr prey.Prey) (float64, error)
		wanted  wanted
	}{
		{
			name: "200 - success",
			reqBody: `{
				  "speed": 4.0,
				  "position": {
				    "X": 0.1,
				    "Y": 0.4,
				    "Z": 3.1
				  }
				}`,
			mockFn: func(pr prey.Prey) (float64, error) {
				return 1.0, nil
			},
			wanted: wanted{
				statusCode: http.StatusOK,
				resBody:    "A presa está configurada corretamente",
			},
		},
		{
			name:    "400 - When request body is invalid",
			reqBody: `{"speed", "4.0", "position", {"X", "0.1", "Y", "0.4", "Z", "3.1"}`,
			wanted: wanted{
				statusCode: http.StatusBadRequest,
				resBody:    "Erro ao decodificar JSON",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			r := httptest.NewRequest(http.MethodGet, "/hunter/configure-prey", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			ht := hunter.NewHunterMock()
			ht.HuntFunc = tt.mockFn
			pr := prey.NewPreyStub()
			h := handler.NewHunter(ht, pr)

			// act
			h.ConfigurePrey(w, r)

			// assert
			assert.Equal(t, tt.wanted.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wanted.resBody)
		})
	}
}

func TestConfigureHunter(t *testing.T) {
	type wanted struct {
		statusCode int
		resBody    string
	}

	tests := []struct {
		name    string
		reqBody string
		wanted  wanted
	}{
		{
			name: "200 - success",
			reqBody: `{
				"speed": 4.0,
				"position": {
					"x": 0.1,
					"y": 0.4
				}
			}`,
			wanted: wanted{
				statusCode: http.StatusOK,
				resBody:    "O caçador esta configurado corretamente",
			},
		},
		{
			name: "400 - When request body is invalid",
			reqBody: `{
				"speed": "4.0",
				"position": {
					"x": 0.1,
					"y": 0.4
				} 			
			}`,
			wanted: wanted{
				statusCode: http.StatusBadRequest,
				resBody:    "cannot unmarshal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			r := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

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

			// act
			h.ConfigureHunter()(w, r)

			// assert
			assert.Equal(t, tt.wanted.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wanted.resBody)
		})
	}
}

func TestHunt(t *testing.T) {
	type wanted struct {
		statusCode int
		resBody    string
	}
	tests := []struct {
		name    string
		reqBody string
		mockFn  func(pr prey.Prey) (float64, error)
		wanted  wanted
	}{
		{
			name: "200 - success",
			mockFn: func(pr prey.Prey) (float64, error) {
				return 1.0, nil
			},
			wanted: wanted{
				statusCode: http.StatusOK,
				resBody:    "A caça",
			},
		},
		{
			name: "500 - When simulator fails",
			mockFn: func(pr prey.Prey) (float64, error) {
				return 0, errors.New("simulator error")
			},
			wanted: wanted{
				statusCode: http.StatusInternalServerError,
				resBody:    "simulator error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			r := httptest.NewRequest(http.MethodPost, "/hunt/hunt", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			ht := hunter.NewHunterMock()
			ht.HuntFunc = tt.mockFn
			pr := prey.NewPreyStub()
			h := handler.NewHunter(ht, pr)

			// act
			h.Hunt()(w, r)

			// assert
			assert.Equal(t, tt.wanted.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wanted.resBody)
		})
	}
}
