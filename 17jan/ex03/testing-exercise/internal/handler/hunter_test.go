package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_ConfigurePrey(t *testing.T) {
	shark := hunter.WhiteShark{}
	tuna := prey.Tuna{}
	hd := handler.NewHunter(&shark, &tuna)
	f := hd.ConfigurePrey
	body := handler.RequestBodyConfigPrey{
		Speed: 5.0,
		Position: &positioner.Position{
			X: 10,
			Y: 20,
			Z: 30,
		},
	}
	expectedCode := http.StatusOK

	b, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	res := httptest.NewRecorder()

	f(res, req)

	require.Equal(t, expectedCode, res.Result().StatusCode)
}

func TestHandler_ConfigureHunter(t *testing.T) {
	shark := hunter.WhiteShark{}
	tuna := prey.Tuna{}
	hd := handler.NewHunter(&shark, &tuna)
	f := hd.ConfigureHunter()
	body := "some invalid request body"
	expectedCode := http.StatusBadRequest

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	res := httptest.NewRecorder()

	f(res, req)

	require.Equal(t, expectedCode, res.Result().StatusCode)
}

func TestHandler_Hunt(t *testing.T) {
	shark := hunter.NewHunterMock()
	t.Run("the hunt succeeds", func(t *testing.T) {
		shark.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 10.0, nil
		}
		hd := handler.NewHunter(shark, nil)
		f := hd.Hunt()
		expectedCode := http.StatusOK
		expectedBody := `
		{
			"message":   "hunt executed successfully",
			"success":   true,
			"time_took": 10.0
		}
	`

		req := httptest.NewRequest("POST", "/", nil)
		res := httptest.NewRecorder()

		f(res, req)

		require.Equal(t, expectedCode, res.Result().StatusCode)
		require.JSONEq(t, expectedBody, res.Body.String())
	})
	t.Run("the hunt fails", func(t *testing.T) {
		shark.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0.0, hunter.ErrCanNotHunt
		}
		hd := handler.NewHunter(shark, nil)
		f := hd.Hunt()
		expectedCode := http.StatusOK
		expectedBody := `
		{
			"message":   "hunt executed successfully",
			"success":   false,
			"time_took": 0.0
		}
	`

		req := httptest.NewRequest("POST", "/", nil)
		res := httptest.NewRecorder()

		f(res, req)

		require.Equal(t, expectedCode, res.Result().StatusCode)
		require.JSONEq(t, expectedBody, res.Body.String())
	})
}
