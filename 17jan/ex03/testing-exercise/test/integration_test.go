//go:build integration

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go test -tags integration ./test/integration_test.go

func TestHunter_ConfigurePrey(t *testing.T) {
	requestBody := handler.RequestBodyConfigPrey{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}
	body, _ := json.Marshal(requestBody)

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

	server := httptest.NewServer(http.HandlerFunc(h.ConfigurePrey))
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/hunter/configure-prey", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(t, "A presa está configurada corretamente", bodyString)
}

func TestHunter_ConfigureHunter(t *testing.T) {
	requestBody := handler.RequestBodyConfigHunter{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
		},
	}
	body, _ := json.Marshal(requestBody)

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

	server := httptest.NewServer(http.HandlerFunc(h.ConfigurePrey))
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/hunter/configure-hunter", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(t, "A presa está configurada corretamente", bodyString)
}

func TestHunter_Hunt(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(h.Hunt()))
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/hunter/hunt", nil)
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "A presa foi caçada em")
}
