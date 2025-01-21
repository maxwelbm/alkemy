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
	type wantRes struct {
		statusCode   int
		responseBody string
	}
	tests := []struct {
		name     string
		reqBody  string
		mockFunc func(pr prey.Prey) (float64, error)
		want     wantRes
	}{
		{
			name: "Given a correct prey then return status code 200 with no errror",
			reqBody: `{
						"speed": 4.0, 
						"position": {
							"X": 0.1,
							"Y": 0.4,
							"Z": 3.1
						}
					}`,
			mockFunc: func(pr prey.Prey) (float64, error) {
				return 1.0, nil
			},
			want: wantRes{
				statusCode:   http.StatusOK,
				responseBody: "A presa está configurada corretamente",
			},
		},
		{
			name: "Given an invalid prey then return status code 400 with error",
			reqBody: `{
						"speed": "1", 
						"position": {
							"X": 0.1,
							"Y": 0.4,
							"Z": 3.1
						}
					}`,
			want: wantRes{
				statusCode:   http.StatusBadRequest,
				responseBody: "Erro ao decodificar JSON",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/hunter/configure-prey", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			ht := hunter.NewHunterMock()
			ht.HuntFunc = tt.mockFunc
			pr := prey.NewPreyStub()
			h := handler.NewHunter(ht, pr)

			h.ConfigurePrey(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.want.responseBody)
		})
	}
}

func TestHunter_ConfigureHunter(t *testing.T) {
	type wantRes struct {
		statusCode   int
		responseBody string
	}
	tests := []struct {
		name    string
		reqBody string
		want    wantRes
	}{
		{
			name: "Given a correct hunter then return status code 200 with no errror",
			reqBody: `{
						"speed": 4.0, 
						"position": {
							"X": 0.1,
							"Y": 0.4,
							"Z": 3.1
						}
					}`,
			want: wantRes{
				statusCode:   http.StatusOK,
				responseBody: "o caçador está configurado corretamente",
			},
		},
		{
			name: "Given an invalid hunter then return status code 400 with error",
			reqBody: `{
						"speed": 1.0, 
						"position": {
							"X": "0.1",
							"Z": 3.1
						}
					}`,
			want: wantRes{
				statusCode:   http.StatusBadRequest,
				responseBody: "requisição inválida",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			pos := positioner.NewPositionerDefault()
			sim := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				Positioner: pos,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     3.0,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sim,
			})

			pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
			h := handler.NewHunter(ht, pr)

			h.ConfigureHunter()(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.want.responseBody)
		})
	}
}

func TestHunterPrey_Hunt(t *testing.T) {
	type wantRes struct {
		statusCode   int
		responseBody string
	}
	tests := []struct {
		name     string
		reqBody  string
		mockFunc func(pr prey.Prey) (float64, error)
		want     wantRes
	}{
		{
			name: "Given a success hunt then return status code 200 with no errror",
			mockFunc: func(pr prey.Prey) (float64, error) {
				return 1.0, nil
			},
			want: wantRes{
				statusCode:   http.StatusOK,
				responseBody: "caça concluida",
			},
		},
		{
			name: "Given an invalid hunt then return status code 500 with error",
			mockFunc: func(pr prey.Prey) (float64, error) {
				return 0, errors.New("internal error")
			},
			want: wantRes{
				statusCode:   http.StatusInternalServerError,
				responseBody: "internal error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/hunt/hunt", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			ht := hunter.NewHunterMock()
			ht.HuntFunc = tt.mockFunc
			pr := prey.NewPreyStub()
			h := handler.NewHunter(ht, pr)

			h.Hunt()(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.want.responseBody)
		})
	}
}
