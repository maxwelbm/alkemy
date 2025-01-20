package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurePrey(t *testing.T) {

	t.Run("The prey is configured correctly, indicating a 'prey configured' message. Returns a 200.", func(t *testing.T) {
		ps := positioner.NewPositionerDefault()
		// - catch simulator
		sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
			Positioner: ps,
		})

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     0.0,
			Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
			Simulator: sm,
		})
		// - prey
		pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
		// - handler
		hd := handler.NewHunter(ht, pr)
		hdfunc := hd.ConfigurePrey

		body := []byte(`{"speed": 10.00, "position": {"X": 10.0, "Y":10.0, "Z": 10.00}}`)
		request := httptest.NewRequest(http.MethodPost, "/configure-prey", bytes.NewReader(body))
		response := httptest.NewRecorder()
		hdfunc(response, request)

		expectedCode := http.StatusOK
		expectedBody := `{"message": "Prey set up"}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("Incorrect hunter configuration, problems with the request body. Returns a 400", func(t *testing.T) {
		ps := positioner.NewPositionerDefault()
		// - catch simulator
		sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
			Positioner: ps,
		})

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     0.0,
			Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
			Simulator: sm,
		})
		// - prey
		pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
		// - handler
		hd := handler.NewHunter(ht, pr)
		hdfunc := hd.ConfigurePrey

		request := httptest.NewRequest(http.MethodPost, "/configure-prey", strings.NewReader(
			`invalid request body`,
		))
		response := httptest.NewRecorder()
		hdfunc(response, request)

		expectedCode := http.StatusBadRequest
		expectedBody := `{"message": "problems with the request body", "status":"Bad Request"}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

}

func TestHunt(t *testing.T) {
	t.Run("The shark manages to capture the prey, displays a 'hunt complete' message and data with information about whether it was captured, in addition to the time interval. Returns a 200", func(t *testing.T) {
		ps := positioner.NewPositionerDefault()
		// - catch simulator
		sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
			Positioner: ps,
		})

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     0.0,
			Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
			Simulator: sm,
		})
		// - prey
		pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
		// - handler
		hd := handler.NewHunter(ht, pr)
		hdFunc := hd.Hunt()

		request := httptest.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)
		expectedCode := http.StatusOK
		expectedBody := `{"message":"hunt complete","data":{"success":true,"duration":100.0}}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("The shark is unable to hunt its prey, it displays a 'hunt complete' message, also displaying the same data as above. Returns a 200", func(t *testing.T) {
		ps := positioner.NewPositionerDefault()
		// - catch simulator
		sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
			Positioner: ps,
		})

		ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
			Speed:     0.0,
			Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
			Simulator: sm,
		})
		// - prey
		pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
		// - handler
		hd := handler.NewHunter(ht, pr)
		hdFunc := hd.Hunt()

		request := httptest.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)
		expectedCode := http.StatusOK
		expectedBody := `{"message":"hunt complete","data":{"success":false,"duration":101.0}}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("fail to hunt - internal server error", func(t *testing.T) {
		ht := hunter.NewHunterMock()
		ht.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
			return 0.0, errors.New("internal server error")
		}
		// - prey
		pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
		// - handler
		hd := handler.NewHunter(ht, pr)
		hdFunc := hd.Hunt()

		request := httptest.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)
		expectedCode := http.StatusInternalServerError
		require.Equal(t, expectedCode, response.Code)

	})
}
