package handler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testing"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	shark := hunter.WhiteShark{}
	tuna := prey.Tuna{}
	hd := NewHunter(&shark, &tuna)

	body := RequestBodyConfigPrey{
		Speed: 35.0,
		Position: &positioner.Position{
			X: 15,
			Y: 30,
			Z: 45,
		},
	}
	expectedCode := http.StatusOK

	req := newRequest("POST", "/", body)
	res := executeRequest(hd.ConfigurePrey, req)

	require.Equal(t, expectedCode, res.Result().StatusCode)
}

func TestHandler_ConfigureHunter(t *testing.T) {
	shark := hunter.WhiteShark{}
	tuna := prey.Tuna{}
	hd := NewHunter(&shark, &tuna)

	body := "some invalid request body"
	expectedCode := http.StatusBadRequest

	req := newRequest("POST", "/", body)
	res := executeRequest(hd.ConfigureHunter(), req)

	require.Equal(t, expectedCode, res.Result().StatusCode)
}

func TestHandler_Hunt(t *testing.T) {
	shark := hunter.NewHunterMock()

	tests := []struct {
		name         string
		huntFunc     func(pr prey.Prey) (float64, error)
		expectedCode int
		expectedBody string
	}{
		{
			name: "the hunt succeeds",
			huntFunc: func(pr prey.Prey) (float64, error) {
				return 10.0, nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `
		{
			"duration": 10,
			"message": "hunt done",
			"success": true
		}
	`,
		},
		{
			name: "the hunt fails",
			huntFunc: func(pr prey.Prey) (float64, error) {
				return 0.0, hunter.ErrCanNotHunt
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `
		{
			"message": "can not hunt the prey", 
			"status": "Internal Server Error"
		}
	`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shark.HuntFunc = tt.huntFunc
			hd := NewHunter(shark, nil)

			req := httptest.NewRequest("POST", "/", nil)
			res := executeRequest(hd.Hunt(), req)

			require.Equal(t, tt.expectedCode, res.Result().StatusCode)
			require.JSONEq(t, tt.expectedBody, res.Body.String())
		})
	}
}

func TestHandler_ConfigurePrey_InvalidJSON(t *testing.T) {
	shark := hunter.WhiteShark{}
	tuna := prey.Tuna{}
	hd := NewHunter(&shark, &tuna)

	body := "this is not a valid json"
	expectedCode := http.StatusBadRequest
	expectedBody := `{"message":"Erro ao decodificar JSON: invalid character 'h' in literal true (expecting 'r')", "status":"Bad Request"}` // Adaptar conforme a resposta real

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	res := httptest.NewRecorder()

	hd.ConfigurePrey(res, req)

	require.Equal(t, expectedCode, res.Result().StatusCode)

	// Verifica se a resposta contém a mensagem de erro esperada
	require.JSONEq(t, expectedBody, res.Body.String())
}

func newRequest(method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			panic(err)
		}
	}
	return httptest.NewRequest(method, url, &buf)
}

func executeRequest(handlerFunc func(http.ResponseWriter, *http.Request), req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	handlerFunc(res, req)
	return res
}
