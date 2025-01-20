package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/prey"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitHunter_Hunter_ConfigurePrey(t *testing.T) {
	t.Run("case 1: configure prey successfully", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/configure-prey", strings.NewReader(`{"speed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`))
		hd.ConfigurePrey().ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Equal(t, `A presa está configurada corretamente`, strings.TrimSpace(res.Body.String()))
	})
	t.Run("case 2: configure prey with bad JSON", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/configure-prey", strings.NewReader(`peed": 4.0,"position":{"X": 0.1, "Y": 0.4,"Z": 3.1}}`))
		hd.ConfigurePrey().ServeHTTP(res, req)
		require.Equal(t, 400, res.Code)
		require.Equal(t, `{"status":"Bad Request","message":"Erro ao decodificar JSON: invalid character 'p' looking for beginning of value"}`, strings.TrimSpace(res.Body.String()))
	})

}
func TestUnitHunter_ConfigureHunter(t *testing.T) {
	t.Run("case 1: configure hunter successfully", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/configure-hunter", strings.NewReader(`{"speed": 0.0, "position": {"x": 0.0, "y": 0.0, "z": 0.0}}`))
		hd.ConfigureHunter().ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
	})
	t.Run("case 2: configure hunter, bad request", func(t *testing.T) {
		hd := handler.NewHunter(hunter.NewHunterMock(), prey.NewPreyStub())
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/configure-hunter", strings.NewReader(`{peed 0.0, "position": {"x": 0.0, "y": 0.0, "z": 0.0}}`))
		hd.ConfigureHunter().ServeHTTP(res, req)
		require.Equal(t, 400, res.Code)
		require.Equal(t, "application/json", res.Header().Get("content-type"))
	})
}

func TestUnitHunter_Hunt(t *testing.T) {
	cases := []struct {
		name         string
		mock         func(prey prey.Prey) (duration float64, err error)
		expectedCode int
		expectedBody string
	}{
		{
			name:         "case 1: hunter can hunt pray",
			mock:         func(prey prey.Prey) (duration float64, err error) { return 2000, nil },
			expectedCode: 200,
			expectedBody: `{"duration":2000,"message":"caça concluída"}`,
		},
		{
			name:         "case 2: hunter can't hunt pray",
			mock:         func(prey prey.Prey) (duration float64, err error) { return 0, hunter.ErrCanNotHunt },
			expectedCode: 200,
			expectedBody: `{"duration":0,"error":"can not hunt the prey","message":"caça concluída"}`,
		},
		{
			name:         "case 3: internal server error",
			mock:         func(prey prey.Prey) (duration float64, err error) { return 0, errors.New("internal server error") },
			expectedCode: 500,
			expectedBody: `{"status":"Internal Server Error","message":"internal server error"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ht := hunter.NewHunterMock()
			ht.HuntFunc = c.mock
			pr := prey.NewPreyStub()
			hd := handler.NewHunter(ht, pr)
			res := httptest.NewRecorder()
			hd.Hunt().ServeHTTP(res, &http.Request{Method: http.MethodGet})
			require.Equal(t, c.expectedCode, res.Code)
			require.Equal(t, "application/json", res.Header().Get("content-type"))
			require.Equal(t, c.expectedBody, strings.TrimSpace(res.Body.String()))
		})
	}

}
