package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testing"

	"github.com/stretchr/testify/assert"
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

	req := httptest.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewReader(body))
	res := httptest.NewRecorder()

	mockHunter := hunter.NewHunterMock()
	mockPrey := prey.NewPreyStub()
	h := handler.NewHunter(mockHunter, mockPrey)

	h.ConfigurePrey(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "A presa está configurada corretamente", res.Body.String())
}

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("configuração bem sucedida", func(t *testing.T) {
		requestBody := handler.RequestBodyConfigHunter{
			Speed: 15.0,
			Position: &positioner.Position{
				X: 0,
				Y: 0,
			},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(body))
		res := httptest.NewRecorder()

		mockHunter := hunter.NewHunterMock()
		mockPrey := prey.NewPreyStub()
		h := handler.NewHunter(mockHunter, mockPrey)

		h.ConfigureHunter()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "O caçador está configurado corretamente", res.Body.String())
	})

	t.Run("erro na configuração - body inválido", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader([]byte("body inválido")))
		res := httptest.NewRecorder()

		mockHunter := hunter.NewHunterMock()
		mockPrey := prey.NewPreyStub()
		h := handler.NewHunter(mockHunter, mockPrey)

		h.ConfigureHunter()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}

func TestHunter_Hunt(t *testing.T) {
	t.Run("caçada bem sucedida", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)
		res := httptest.NewRecorder()

		mockHunter := hunter.NewHunterMock()
		mockPrey := prey.NewPreyStub()
		h := handler.NewHunter(mockHunter, mockPrey)

		mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 10.5, nil
		}

		h.Hunt()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var response map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "caça concluída", response["message"])
		assert.Equal(t, 10.5, response["duration"])
	})
}
