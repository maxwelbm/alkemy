package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHunter_ConfigurePrey(t *testing.T) {
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

	t.Run("Foi possivel configurar presa", func(t *testing.T) {
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

		h.ConfigurePrey(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
	})

	t.Run("Não foi possivel configurar presa por request body invalido", func(t *testing.T) {
		body := []byte(`{"key: "value"}`)
		req, err := http.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewReader(body))

		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.ConfigurePrey(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestHunter_ConfigureHunter(t *testing.T) {
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

	t.Run("Foi possivel configurar caçador", func(t *testing.T) {
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

		h.ConfigureHunter(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "O caçador está configurado corretamente", recorder.Body.String())
	})

	t.Run("Não foi possivel configurar caçador por request body invalido", func(t *testing.T) {
		body := []byte(`{"key: "value"}`)
		req, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(body))

		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.ConfigurePrey(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

}

func TestHunter_Hunt(t *testing.T) {
	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner:     ps,
		MaxTimeToCatch: 1000,
	})

	t.Run("Caça concluída e capturada", func(t *testing.T) {
		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     10.0,
			Position:  &positioner.Position{X: 1, Y: 2, Z: 3},
			Simulator: sm,
		})

		pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

		h := NewHunter(ht, pr)

		req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		resBody := ResponseBodyHunt{
			Message:  "Caça concluída",
			Duration: 0.3897559777889522,
		}

		expectedBody, _ := json.Marshal(resBody)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, string(expectedBody), recorder.Body.String())
	})

	t.Run("Caça concluída e nao capturada", func(t *testing.T) {
		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     1.0,
			Position:  &positioner.Position{X: 1, Y: 2, Z: 3},
			Simulator: sm,
		})

		pr := prey.NewTuna(10.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

		h := NewHunter(ht, pr)

		req, err := http.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		if err != nil {
			t.Fatalf("Não foi possível criar a requisição: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		resBody := ResponseBodyHunt{
			Message: "Caça concluída",
		}

		expectedBody, _ := json.Marshal(resBody)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, string(expectedBody), recorder.Body.String())
	})

}
