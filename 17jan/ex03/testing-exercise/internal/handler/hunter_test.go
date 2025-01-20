package handler

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/prey"
	"testing"
)

func TestHandlerConfigurePrey(t *testing.T) {

	t.Run("case 1: prey configured successfully", func(t *testing.T) {
		// arrange
		// - prey: stub
		pr := prey.NewPreyStub()
		// - handler
		hd := NewHunter(nil, pr)
		hdFunc := hd.ConfigurePrey()

		// act
		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`{"speed": 5.9, "position": {"X": 0.0, "Y": 0.0, "Z": 0.0}}`,
		))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"A presa está configurada corretamente","data":null}`
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
	})

	t.Run("case 2: prey configured failed - invalid request body", func(t *testing.T) {
		// arrange
		// - handler
		hd := NewHunter(nil, nil)
		hdFunc := hd.ConfigurePrey()

		// act
		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`invalid request body`,
		))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"Erro ao decodificar JSON: invalid character 'i' looking for beginning of value",
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
	})
}
