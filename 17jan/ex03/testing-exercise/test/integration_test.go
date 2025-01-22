package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/handler"
	"testing"

	"github.com/stretchr/testify/assert"

	"testdoubles/internal/hunter"
	"testdoubles/internal/prey"
)

// TestHunterHandler cobre os cenários especificados  1 2 3 4
func TestHunterHandler(t *testing.T) {

	t.Run("ConfigurePrey - success", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		h := handler.NewHunter(mockHunter, stubPrey)

		body := `{"speed": 4.0, "position": {"X": 0.1, "Y": 0.4, "Z": 3.1}}`
		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		h.ConfigurePrey(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "A presa está configurada corretamente")
		assert.Equal(t, 1, mockHunter.Calls.Configure, "esperava que Configure fosse chamado 1x")
	})

	t.Run("ConfigureHunter - bad request (400)", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()
		hh := handler.NewHunter(mockHunter, stubPrey)

		body := `{"speed": "not a float", "position": {}}`
		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handlerFunc := hh.ConfigureHunter()
		handlerFunc(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Erro ao decodificar JSON")
		assert.Equal(t, 0, mockHunter.Calls.Configure)
	})

	t.Run("Hunt - tubarão consegue capturar (200)", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		mockHunter.HuntFunc = func(pr prey.Prey) (float64, error) {
			return 12.5, nil
		}
		hh := handler.NewHunter(mockHunter, prey.NewPreyStub())

		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		rr := httptest.NewRecorder()

		handlerFunc := hh.Hunt()
		handlerFunc(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		body := rr.Body.String()
		assert.Contains(t, body, "caça concluída")
		assert.Contains(t, body, "capturada: true")
		assert.Contains(t, body, "tempo: 12.5")

		assert.Equal(t, 1, mockHunter.Calls.Hunt)
	})

	t.Run("Hunt - tubarão NÃO consegue capturar (200)", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		mockHunter.HuntFunc = func(pr prey.Prey) (float64, error) {
			return 30.0, hunter.ErrCanNotHunt
		}
		hh := handler.NewHunter(mockHunter, prey.NewPreyStub())

		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		rr := httptest.NewRecorder()

		handlerFunc := hh.Hunt()
		handlerFunc(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		body := rr.Body.String()
		assert.Contains(t, body, "caçada concluída")
		assert.Contains(t, body, "capturada: false")
		assert.Contains(t, body, "tempo: 30")

		assert.Equal(t, 1, mockHunter.Calls.Hunt)
	})

	t.Run("Hunt - erro inesperado (500)", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		mockHunter.HuntFunc = func(pr prey.Prey) (float64, error) {
			return 0, errors.New("erro aleatório de estagiário")
		}
		hh := handler.NewHunter(mockHunter, prey.NewPreyStub())

		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		rr := httptest.NewRecorder()

		handlerFunc := hh.Hunt()
		handlerFunc(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Erro interno na caçada:")
		assert.Equal(t, 1, mockHunter.Calls.Hunt)
	})
}
