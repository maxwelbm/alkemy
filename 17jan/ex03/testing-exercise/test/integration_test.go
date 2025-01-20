package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testdoubles/internal/application"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	go func() {
		server := application.NewApplicationDefault(":8080")

		err := server.SetUp()
		if err != nil {
			panic(err)
		}

		err = server.Run()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func TestConfigurePreyRequest(t *testing.T) {
	pr := handler.RequestBodyConfigPrey{
		Speed: 2.1,
		Position: &positioner.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}
	body, err := json.Marshal(pr)
	require.NoError(t, err)

	r, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer r.Body.Close()

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(r.Body)
	require.Equal(t, http.StatusOK, r.StatusCode)
	expectedMessage := "A presa está configurada corretamente"
	require.Equal(t, expectedMessage, responseBody.String())
}

func TestConfigureHunterRequest(t *testing.T) {
	ht := handler.RequestBodyConfigHunter{
		Speed: 100.5,
		Position: &positioner.Position{
			X: 0.0,
			Y: 0.0,
			Z: 0.0,
		},
	}
	hunterBody, err := json.Marshal(ht)
	require.NoError(t, err)

	r, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewBuffer(hunterBody))
	require.NoError(t, err)
	defer r.Body.Close()

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(r.Body)
	require.Equal(t, http.StatusOK, r.StatusCode)
	expectedMessage := "O cacador está configurado corretamente"
	require.Equal(t, expectedMessage, responseBody.String())
}

func TestHuntRequest(t *testing.T) {
	r, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", nil)
	require.NoError(t, err, "Erro ao enviar requisição para a API")
	defer r.Body.Close()

	var responseBody map[string]any
	err = json.NewDecoder(r.Body).Decode(&responseBody)
	require.NoError(t, err, "Erro ao decodificar a resposta da API")

	require.Contains(t, responseBody, "message", "A resposta não contém a chave 'message'")
	require.Contains(t, responseBody, "success", "A resposta não contém a chave 'success'")
	require.Contains(t, responseBody, "duration", "A resposta não contém a chave 'duration'")

	require.Equal(t, "caca concluida", responseBody["message"], "Mensagem incorreta na resposta")
	require.IsType(t, true, responseBody["success"], "O campo 'success' não é do tipo booleano")
	require.True(t, responseBody["success"].(bool), "O campo 'success' não é verdadeiro")
	require.NotNil(t, responseBody["duration"], "O campo 'duration' é nulo")
}
