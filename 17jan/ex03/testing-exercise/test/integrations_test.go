package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testdoubles/internal/application"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var app *application.ApplicationDefault
var serverReady = false

func startServer(t *testing.T) {
	cfg := &application.ConfigApplicationDefault{
		Addr: "127.0.0.1:8080",
	}
	app = application.NewApplicationDefault(cfg.Addr)

	if err := app.SetUp(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	for i := 0; i < 10; i++ {
		res, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", nil)
		if err == nil {
			res.Body.Close()
			serverReady = true
			return
		}
		time.Sleep(time.Millisecond * 500)
	}

	t.Fatal("Servidor não está respondendo")
}

func TestMain(m *testing.M) {
	startServer(nil)

	exitCode := m.Run()

	if !serverReady {
		log.Println("Erro: O servidor não ficou pronto.")
	}

	os.Exit(exitCode)
}

func TestIntegrationsConfigHunt(t *testing.T) {
	if !serverReady {
		t.Skip("Servidor não está pronto para receber requisições")
	}

	requestBody := handler.RequestBodyConfigHunter{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 0,
			Y: 0,
		},
	}
	body, _ := json.Marshal(requestBody)
	res, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	responseBody := new(bytes.Buffer)

	require.NoError(t, err)
	defer res.Body.Close()

	_, err = responseBody.ReadFrom(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)

	require.Equal(t, "The hunter configure sucess", responseBody.String())
}

func TestIntegrationsConfigPrey(t *testing.T) {
	if !serverReady {
		t.Skip("Servidor não está pronto para receber requisições")
	}

	requestBody := handler.RequestBodyConfigPrey{
		Speed: 1.5,
		Position: &positioner.Position{
			X: 0,
			Y: 0,
		},
	}
	body, _ := json.Marshal(requestBody)
	res, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	responseBody := new(bytes.Buffer)

	require.NoError(t, err)
	defer res.Body.Close()

	_, err = responseBody.ReadFrom(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)

	require.Equal(t, "A presa está configurada corretamente", responseBody.String())
}

func TestIntegrations(t *testing.T) {
	if !serverReady {
		t.Skip("Servidor não está pronto para receber requisições")
	}

	res, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", nil)
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	// responseBody := new(bytes.Buffer)

	var responseBody struct {
		Message  string `json:"message"`
		Sucess   bool   `json:"sucess"`
		Duration int    `json:"duration"`
	}

	err = json.NewDecoder(res.Body).Decode(&responseBody)
	require.NoError(t, err)
	defer res.Body.Close()

	// _, err = responseBody.ReadFrom(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "prey hunted", responseBody.Message)
	// require.True(t, responseBody.Sucess)
}
