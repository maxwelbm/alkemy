package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
)

func TestHunter_ConfigurePrey(t *testing.T) {
	preyConfig := handler.RequestBodyConfigPrey{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	if err := encoder.Encode(preyConfig); err != nil {
		t.Fatalf("Erro ao encodar JSON para configuração da presa: %v", err)
	}

	testRecorder := httptest.NewRecorder()

	pos := positioner.NewPositionerDefault()
	catchSim := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: pos,
	})

	shark := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     3.0,
		Position:  &positioner.Position{X: 0, Y: 0, Z: 0},
		Simulator: catchSim,
	})
	tuna := prey.NewTuna(0.4, &positioner.Position{X: 0, Y: 0, Z: 0})

	hunterHandler := handler.NewHunter(shark, tuna)

	testServer := httptest.NewServer(http.HandlerFunc(hunterHandler.ConfigurePrey))
	defer testServer.Close()

	req, errReq := http.NewRequest(
		http.MethodPost,
		testServer.URL+"/hunter/configure-prey",
		bytes.NewReader(buffer.Bytes()),
	)
	if errReq != nil {
		t.Fatalf("Falha ao criar request de teste para configurar a presa: %v", errReq)
	}

	hunterHandler.ConfigurePrey(testRecorder, req)

	require.Equal(t, http.StatusOK, testRecorder.Code)
	require.Equal(t, "A presa está configurada corretamente", testRecorder.Body.String())
}

func TestHunter_ConfigureHunter(t *testing.T) {
	hunterConfig := handler.RequestBodyConfigHunter{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(hunterConfig); err != nil {
		t.Fatalf("Erro ao encodar JSON para configuração do hunter: %v", err)
	}

	req, errReq := http.NewRequest(http.MethodPost, "/hunter/configure-hunter", &buffer)
	if errReq != nil {
		t.Fatalf("Não foi possível criar a requisição inicial: %v", errReq)
	}

	rec := httptest.NewRecorder()

	pos := positioner.NewPositionerDefault()
	sim := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner: pos,
	})

	whiteShark := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     3.0,
		Position:  &positioner.Position{X: 0, Y: 0, Z: 0},
		Simulator: sim,
	})
	tuna := prey.NewTuna(0.4, &positioner.Position{X: 0, Y: 0, Z: 0})

	hunterHandler := handler.NewHunter(whiteShark, tuna)

	testServer := httptest.NewServer(http.HandlerFunc(hunterHandler.ConfigureHunter()))
	defer testServer.Close()

	req, errReq = http.NewRequest(http.MethodPost, testServer.URL+"/hunter/configure-prey", bytes.NewReader(buffer.Bytes()))
	if errReq != nil {
		t.Fatalf("Não foi possível criar a requisição final: %v", errReq)
	}

	configFunction := hunterHandler.ConfigureHunter()

	configFunction(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "O tubarão está configurado corretamente", rec.Body.String())
}
