package internal_test

import (
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
	// arrange
	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})
	// - hunter
	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     0.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})
	// - prey
	pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
	// - handler
	hd := handler.NewHunter(ht, pr)

	hdFunc := hd.ConfigurePrey()

	cases := []struct {
		testName       string
		request        *http.Request
		expectedCode   int
		expectedString string
	}{
		{"case 1: configure prey successfully", httptest.NewRequest(http.MethodPost, "/hunter", strings.NewReader(
			`{"speed": 90.0,"position": {"X": 0.1,"Y": 0.4,"Z": 3.1}}`)), http.StatusOK, `{"data":null,"message":"A presa está configurada corretamente"}`},
		{"case 2: configure prey fails", httptest.NewRequest(http.MethodPost, "/hunter", strings.NewReader("")), http.StatusBadRequest, `{"message":"Erro ao decodificar JSON: EOF", "status":"Bad Request"}`},
	}
	// act

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			response := httptest.NewRecorder()
			hdFunc(response, c.request)

			// assert
			require.Equal(t, c.expectedCode, response.Code)
			require.JSONEq(t, c.expectedString, response.Body.String())
		})
	}

}

func TestConfigureHunter(t *testing.T) {
	// arrange
	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})
	// - hunter
	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     0.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})
	// - prey
	pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
	// - handler
	hd := handler.NewHunter(ht, pr)

	hdFunc := hd.ConfigureHunter()

	cases := []struct {
		testName       string
		request        *http.Request
		expectedCode   int
		expectedString string
	}{
		{"case 1: configure hunter successfully", httptest.NewRequest(http.MethodPost, "/hunter", strings.NewReader(
			`{"speed": 90.0,"position": {"X": 0.1,"Y": 0.4,"Z": 3.1}}`)), http.StatusOK, `{"data":null,"message":"hunter configured"}`},
		{"case 2: configure hunter fails", httptest.NewRequest(http.MethodPost, "/hunter", strings.NewReader("")), http.StatusBadRequest, `{"message":"Erro ao decodificar JSON: EOF", "status":"Bad Request"}`},
	}
	// act

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			response := httptest.NewRecorder()
			hdFunc(response, c.request)

			// assert
			require.Equal(t, c.expectedCode, response.Code)
			require.JSONEq(t, c.expectedString, response.Body.String())
		})
	}

}

func TestHuntSucess(t *testing.T) {
	// arrange
	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})
	// - hunter
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

	// act
	t.Run("case 1: hunt done", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/hunter", strings.NewReader(
			`{"speed": 90.0,"position": {"X": 0.1,"Y": 0.4,"Z": 3.1}}`))
		expectedString := `{"data":{"duration":0, "success":false}, "message":"hunt done"}`
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		require.Equal(t, http.StatusOK, response.Code)
		require.JSONEq(t, expectedString, response.Body.String())
	})
}

func TestHuntFails(t *testing.T) {
	// arrange
	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})
	// - hunter
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

	// act
	t.Run("case 2: hunt fails", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/hunter", strings.NewReader(""))
		expectedString := `{"data":{"duration":0, "success":false}, "shark can not catch the prey: shark can not catch the prey"}`
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		require.Equal(t, http.StatusOK, response.Code)
		require.JSONEq(t, expectedString, response.Body.String())
	})
}
