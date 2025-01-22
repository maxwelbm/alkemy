package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/hunter"
	"testdoubles/internal/prey"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	t.Run("Status Code 200 - ConfigurePrey", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		h := NewHunter(mockHunter, stubPrey)

		body := `{"speed": 10, "position":{"X": 100, "Y": 200, "Z": 300}}`

		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewBufferString(body))

		recorder := httptest.NewRecorder()

		h.ConfigurePrey(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "The prey is configured correctly")
	})

	t.Run("Status Code 400 - ConfigurePrey", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		h := NewHunter(mockHunter, stubPrey)

		body := `{"speed": "aaa", "position":{"X": 100, "Y": 200, "Z": 300}}`

		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-prey", bytes.NewBufferString(body))

		recorder := httptest.NewRecorder()

		h.ConfigurePrey(recorder, req)

		assert.Equal(t, 400, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Error decoding JSON: ")
	})
}

func TestHunter_ConfigureHunter(t *testing.T) {
	t.Run("Status Code 200 - ConfigureHunter", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		h := NewHunter(mockHunter, stubPrey)

		body := `{"speed": 20, "position":{"X": 100, "Y": 200, "Z": 300}}`

		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewBufferString(body))

		recorder := httptest.NewRecorder()

		h.ConfigureHunter(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "The shark is configured correctly")
	})

	t.Run("Status Code 400 - ConfigureHunter", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		h := NewHunter(mockHunter, stubPrey)

		body := `{"speed": "aaa", "position":{"X": 100, "Y": 200, "Z": 300}}`

		req := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewBufferString(body))

		recorder := httptest.NewRecorder()

		h.ConfigureHunter(recorder, req)

		assert.Equal(t, 400, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Error decoding JSON: ")
	})
}

func TestHunt(t *testing.T) {
	t.Run("Status Code 200 - Prey Caught: true", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0, nil
		}

		h := NewHunter(mockHunter, stubPrey)

		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "hunt completed with success | prey caught: true | duration: ")
	})

	t.Run("Status Code 200 - Prey Caught: false", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0, hunter.ErrCanNotHunt
		}

		h := NewHunter(mockHunter, stubPrey)

		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "hunt completed with success | prey caught: false | duration: ")
	})

	t.Run("Status Code 500 - Another Error", func(t *testing.T) {
		mockHunter := hunter.NewHunterMock()
		stubPrey := prey.NewPreyStub()

		mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			err = bytes.ErrTooLarge
			return 0, err
		}

		h := NewHunter(mockHunter, stubPrey)

		req := httptest.NewRequest(http.MethodPost, "/hunter/hunt", nil)

		recorder := httptest.NewRecorder()

		h.Hunt(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Internal error: ")
	})
}
