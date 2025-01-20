package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("cenario 1: sucesso de config. hunter", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		hd := NewHunter(ht, nil)
		hdFunc := hd.ConfigureHunter()

		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`{"speed": 10.0, "position": {"X": 0.0, "Y": 0.0, "Z": 0.0}}`,
		))

		response := httptest.NewRecorder()
		hdFunc(response, request)

		expectedCode := http.StatusOK
		expectedMessage := `O caçador está configurado corretamente`
		require.Equal(t, expectedCode, response.Code)
		require.Equal(t, expectedMessage, response.Body.String())
	})

	t.Run("cenario 2: falha ao config. hunter", func(t *testing.T) {
		hd := NewHunter(nil, nil)
		hdFunc := hd.ConfigureHunter()

		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`{"speed": AAA}`,
		))
		response := httptest.NewRecorder()
		hdFunc(response, request)

		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"Erro ao decodificar JSON: invalid character 'A' looking for beginning of value",
		)
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}

func TestHunter_Hunt(t *testing.T) {
	t.Run("cenario 1: caça bem sucedida", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 100.0, nil
		}
		hd := NewHunter(ht, nil)
		hdFunc := hd.Hunt()

		request := httptest.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		expectedCode := http.StatusOK
		expectedBody := `{"duration":100, "message":"caçada concluída com sucesso.", "result:":"presa capturada."}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("cenario 2: caça mal sucedida", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0.0, hunter.ErrCanNotHunt
		}
		hd := NewHunter(ht, nil)
		hdFunc := hd.Hunt()

		request := httptest.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		expectedCode := http.StatusOK
		expectedBody := `{"duration":0, "message":"caçada concluída.", "result:":"presa fugiu."}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("cenario 3: caça mal sucedida - internal", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0.0, errors.New("internal server error")
		}
		hd := NewHunter(ht, nil)

		request := httptest.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()
		hd.Hunt()(response, request)

		expectedCode := http.StatusInternalServerError
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"internal server error",
		)
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}
