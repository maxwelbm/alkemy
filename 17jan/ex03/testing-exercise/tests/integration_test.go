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
			name:         "1 - configurando presa sucesso",
			body:         `{"speed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`,
			expectedCode: http.StatusOK,
			expectedBody: "A presa está configurada corretamente",
		},
		{
			name:         "2 - Configurando tubarao erro Json presa",
			body:         `peed": 4.0,osition":{"X":}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","message":"Erro ao decodificar JSON: invalid character 'p' looking for beginning of value"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewReader([]byte(c.body)))
			recorder := httptest.NewRecorder()
			ps := positioner.NewPositionerDefault()
			sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				Positioner: ps,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     1,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sm,
			})
			pr := prey.NewTuna(1, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
			h := handler.NewHunter(ht, pr)

			router := chi.NewRouter()
			router.HandleFunc("/hunter/configure-prey", h.ConfigurePrey)
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
			name:         "1 - Configurando tubarao sucesso",
			body:         `{"speed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`,
			expectedCode: http.StatusOK,
			expectedBody: `O caçador está configurado corretamente`,
		},
		{
			name:         "2 - Configurando tubarao erro Json",
			body:         `"speed": 4.0,osition":{"X": 0.1, Y: 0.4,}}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","message":"Erro ao decodificar JSON"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader([]byte(c.body)))
			recorder := httptest.NewRecorder()
			ps := positioner.NewPositionerDefault()
			sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
				Positioner: ps,
			})
			ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
				Speed:     2,
				Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
				Simulator: sm,
			})
			pr := prey.NewTuna(1, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
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
			name:         "1 - Caçador pode consegue a presa",
			expectedCode: http.StatusOK,
			expectedBody: `1`,
			tunaPosition: &positioner.Position{X: 3.0, Y: 0.0, Z: 0.0},
		},
		{
			name:         "2 - Caçador não consegue caçar a presa",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":"Internal Server Error","message":"can not hunt the prey: shark can not catch the prey"}`,
			tunaPosition: &positioner.Position{X: 900.0, Y: 2000.0, Z: 30000.0},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
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