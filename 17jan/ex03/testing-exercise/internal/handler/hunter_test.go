package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/hunter"
	"testdoubles/internal/prey"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerConfigurePrey(t *testing.T) {
	t.Run("case 1: prey configured successfully", func(t *testing.T) {
		// arrange
		pr := prey.NewPreyStub()
		hd := NewHunter(nil, pr)
		hdFunc := hd.ConfigurePrey()

		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`{"speed": 5.9, "position": {"X": 0.0, "Y": 0.0, "Z": 0.0}}`,
		))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		// act
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"A presa está configurada corretamente","data":null}`
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
	})

	t.Run("case 2: prey configuration failed - invalid request body", func(t *testing.T) {
		// arrange
		hd := NewHunter(nil, nil)
		hdFunc := hd.ConfigurePrey()

		// JSON malformado (falta uma chave ou valor)
		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`{"speed": "invalid", "position": {"X": 0.0, "Y": 0.0}}`, // Aqui o valor de speed é uma string
		))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		// act
		hdFunc(response, request)

		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"Erro ao decodificar JSON: json: cannot unmarshal string into Go struct field RequestBodyConfigPrey.speed of type float64",
		)

		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}

		// assert
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
	})
}

func TestHandlerConfigureHunter(t *testing.T) {
	t.Run("case 1: hunter configured successfully", func(t *testing.T) {
		// arrange
		// - hunter: mock
		ht := hunter.NewHunterMock()
		// - handler
		hd := NewHunter(ht, nil)
		hdFunc := hd.ConfigureHunter()

		// act
		request := httptest.NewRequest("POST", "/", strings.NewReader(
			`{"speed": 10.0, "position": {"X": 0.0, "Y": 0.0, "Z": 0.0}}`,
		))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"hunter configured","data":null}`
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		expectedCallConfigure := 1
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
		require.Equal(t, expectedCallConfigure, ht.Calls.Configure)
	})

	t.Run("case 2: hunter configured failed  - invalid request body", func(t *testing.T) {
		// arrange
		// - handler
		hd := NewHunter(nil, nil)
		hdFunc := hd.ConfigureHunter()

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

func TestHandlerHunt(t *testing.T) {
	t.Run("case 1: success to hunt", func(t *testing.T) {
		// arrange
		// - hunter: mock
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 100.0, nil
		}
		// - handler
		hd := NewHunter(ht, nil)
		hdFunc := hd.Hunt()

		// act
		request := httptest.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"hunt done","data":{"success":true,"duration":100.0}}`
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		expectedCallHunt := 1
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
		require.Equal(t, expectedCallHunt, ht.Calls.Hunt)
	})

	t.Run("case 2: fail to hunt - fail hunting the prey", func(t *testing.T) {
		// arrange
		// - hunter: mock
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0.0, hunter.ErrCanNotHunt
		}
		// - handler
		hd := NewHunter(ht, nil)
		hdFunc := hd.Hunt()

		// act
		request := httptest.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"hunt done","data":{"success":false,"duration":0.0}}`
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		expectedCallHunt := 1
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
		require.Equal(t, expectedCallHunt, ht.Calls.Hunt)
	})

	t.Run("case 3: fail to hunt - internal server error", func(t *testing.T) {
		// arrange
		// - hunter: mock
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0.0, errors.New("internal server error")
		}
		// - handler
		hd := NewHunter(ht, nil)

		// act
		request := httptest.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()
		hd.Hunt()(response, request)

		// assert
		expectedCode := http.StatusInternalServerError
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"internal server error",
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		expectedCallHunt := 1
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeaders, response.Header())
		require.Equal(t, expectedCallHunt, ht.Calls.Hunt)
	})
}
