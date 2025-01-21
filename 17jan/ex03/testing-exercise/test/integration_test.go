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
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	requestBody := handler.RequestBodyConfigPrey{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Não foi possível criar o corpo da requisição: %v", err)
	}

	recorder := httptest.NewRecorder()

	position := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: position,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     3.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

	h := handler.NewHunter(ht, pr)

	server := httptest.NewServer(http.HandlerFunc(h.ConfigurePrey))
	defer server.Close()
	req, err := http.NewRequest(http.MethodPost, server.URL+"/hunter/configure-prey", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	h.ConfigurePrey(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "A presa está configurada corretamente", recorder.Body.String())
}

func TestHunter_ConfigureHunter(t *testing.T) {
	requestBody := handler.RequestBodyConfigHunter{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Não foi possível criar o corpo da requisição: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	recorder := httptest.NewRecorder()

	ps := positioner.NewPositionerDefault()

	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: ps,
	})

	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     3.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})

	pr := prey.NewTuna(0.4, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})

	h := handler.NewHunter(ht, pr)

	server := httptest.NewServer(http.HandlerFunc(h.ConfigureHunter()))
	defer server.Close()
	req, err = http.NewRequest(http.MethodPost, server.URL+"/hunter/configure-prey", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	funcHandler := h.ConfigureHunter()

	funcHandler(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "O caçador está configurado corretamente", recorder.Body.String())
}
