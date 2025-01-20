package handler_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup() (*handler.Hunter, *httptest.Server) {
	ps := positioner.NewPositionerDefault()
	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{Positioner: ps})

	// Configure o tubarão e a presa em posições adequadas
	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     20.0,                                         // Velocidade do tubarão
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0}, // Posição do tubarão
		Simulator: sm,
	})

	// Coloque a presa a uma distância que o tubarão pode alcançar
	pr := prey.NewTuna(0.0, &positioner.Position{X: 5.0, Y: 0.0, Z: 0.0}) // Posição da presa

	hd := handler.NewHunter(ht, pr)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/configure-prey":
			hd.ConfigurePrey(w, r)
		case "/hunt":
			hd.Hunt()(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	return hd, server
}

func TestIntegration_Hunter(t *testing.T) {
	hd, server := setup()
	defer server.Close()

	t.Run("Integration Test - ConfigurePrey", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"speed": 20.00, "position": {"X": 10.0, "Y":10.0, "Z": 10.00}}`))
		response, err := http.Post(server.URL+"/configure-prey", "application/json", body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		expectedBody := `{"message": "Prey set up"}`
		actualBody := readResponseBody(response)
		require.JSONEq(t, expectedBody, actualBody)
	})

	t.Run("Integration Test - Hunt", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, server.URL+"/hunt", nil)
		response := httptest.NewRecorder()

		log.Println("Running Hunt")
		hd.Hunt()(response, request)

		require.Equal(t, http.StatusOK, response.Code)

		expectedBody := `{"message":"hunt complete","data":{"success":true,"duration":0}}` // Ajustar conforme a lógica real
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}

// Função auxiliar para ler o corpo de resposta
func readResponseBody(res *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf.String()
}
