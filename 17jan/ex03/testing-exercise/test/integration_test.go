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

func TestHunter_ConfigurePrey_Integration(t *testing.T) {
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
			name: "Given a valid prey then configure prey with no error",
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
				responseBody: "A presa está configurada corretamente",
			},
		},
		{
			name: "Given an invalid prey then return error",
			reqBody: `{
						"speed": "4.0",
						"position": {
							"X": 0.1,
							"Y": 0.4,
							"Z": 3.1
						}
					}`,
			want: wantRes{
				statusCode:   http.StatusBadRequest,
				responseBody: `{"status":"Bad Request","message":"Erro ao decodificar JSON"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/hunter/configure-prey", bytes.NewReader([]byte(tt.reqBody)))
			rec := httptest.NewRecorder()
			pos := positioner.NewPositionerDefault()
			sim := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				Positioner: pos,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     4.0,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sim,
			})
			pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
			h := handler.NewHunter(ht, pr)

			router := chi.NewRouter()
			router.HandleFunc("/hunter/configure-prey", h.ConfigurePrey)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.want.statusCode, rec.Code)
			assert.Equal(t, tt.want.responseBody, rec.Body.String())
		})
	}
}

func TestHunter_ConfigureHunter_Integration(t *testing.T) {
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
			name: "Given a valid hunter then configure hunter with no error",
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
				responseBody: `{"message":"o caçador está configurado corretamente"}`,
			},
		},
		{
			name: "Give an invalid hunter then return error",
			reqBody: `{
						"speed": "4.0",
						"position": {
							"X": 0.1,
							"Y": 0.4,
							"Z": 3.1
						}
					}`,
			want: wantRes{
				statusCode:   http.StatusBadRequest,
				responseBody: `{"status":"Bad Request","message":"requisição inválida"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/hunter/configure-hunter", bytes.NewReader([]byte(tt.reqBody)))
			rec := httptest.NewRecorder()
			pos := positioner.NewPositionerDefault()
			sim := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				Positioner: pos,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     4.0,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sim,
			})
			pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
			h := handler.NewHunter(ht, pr)

			router := chi.NewRouter()
			router.HandleFunc("/hunter/configure-hunter", h.ConfigureHunter())
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.want.statusCode, rec.Code)
			assert.Equal(t, tt.want.responseBody, rec.Body.String())
		})
	}
}

func TestHunter_Hunt_Integration(t *testing.T) {
	type wantRes struct {
		statusCode   int
		responseBody string
	}
	tests := []struct {
		name         string
		want         wantRes
		tunaPosition *positioner.Position
	}{
		{
			name: "Given a hunt that hunter can hunt the prey then return success with no error",
			want: wantRes{
				statusCode:   http.StatusOK,
				responseBody: `{"hunt_duration":1.5,"message":"caça concluida","success":true}`,
			},
			tunaPosition: &positioner.Position{X: 3.0, Y: 0.0, Z: 0.0},
		},
		{
			name: "Given a hunt that hunter cannot hunt the prey then return error",
			want: wantRes{
				statusCode:   http.StatusInternalServerError,
				responseBody: `{"status":"Internal Server Error","message":"internal error"}`,
			},
			tunaPosition: &positioner.Position{X: 90.0, Y: 300.0, Z: 4000.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/hunter/hunt", nil)
			rec := httptest.NewRecorder()
			pos := positioner.NewPositionerDefault()
			sim := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				MaxTimeToCatch: 3,
				Positioner:     pos,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     3,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sim,
			})
			pr := prey.NewTuna(1, tt.tunaPosition)
			h := handler.NewHunter(ht, pr)

			router := chi.NewRouter()
			router.HandleFunc("/hunter/hunt", h.Hunt())
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.want.statusCode, rec.Code)
			assert.Equal(t, tt.want.responseBody, rec.Body.String())
		})
	}
}
