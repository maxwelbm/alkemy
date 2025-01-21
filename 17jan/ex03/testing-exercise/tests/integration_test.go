package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testdoubles/internal/application"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	go func() {
		app := application.NewApplicationDefault(":8080")
		if err := app.SetUp(); err != nil {
			fmt.Println(err)
			return
		}

		if err := app.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}()
}

func TestConfigureHunter(t *testing.T) {
	t.Run("case 1 : bad request - hunter not created", func(t *testing.T) {
		//given
		var requestBody []byte

		//when
		res, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "{\"status\":\"Bad Request\",\"message\":\"Erro ao decodificar JSON: EOF\"}", bodyString)
	})
	t.Run("case 2 : hunter updated", func(t *testing.T) {
		//given
		sampleBody := handler.RequestBodyConfigHunter{
			Speed: 10.4,
			Position: &positioner.Position{
				X: 10,
				Y: 5,
				Z: 2,
			},
		}

		requestBody, err := json.Marshal(sampleBody)
		if err != nil {
			fmt.Println("Error coding request body:", err)
			return
		}
		//when
		res, err := http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "O caçador está configurado corretamente", bodyString)
	})
}

func TestConfigurePrey(t *testing.T) {
	t.Run("case 1 : bad request - prey not created", func(t *testing.T) {
		//given
		var requestBody []byte

		//when
		res, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "{\"status\":\"Bad Request\",\"message\":\"Erro ao decodificar JSON: EOF\"}", bodyString)
	})
	t.Run("case 2 : prey updated", func(t *testing.T) {
		//given
		sampleBody := handler.RequestBodyConfigPrey{
			Speed: 70.4,
			Position: &positioner.Position{
				X: 19,
				Y: 5,
				Z: 9,
			},
		}

		requestBody, err := json.Marshal(sampleBody)
		if err != nil {
			fmt.Println("Error coding request body:", err)
			return
		}
		//when
		res, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "A presa está configurada corretamente", bodyString)
	})
}

func TestHunt(t *testing.T) {
	t.Run("case 1 : hunt without configuring prey and hunter", func(t *testing.T) {
		//given
		var requestBody []byte
		//when
		res, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.Contains(t, bodyString, "caça concluída sem sucesso após")
	})
	t.Run("case 2 : sucessful hunt after configuring hunter and prey", func(t *testing.T) {
		//given
		sampleBodyHunter := handler.RequestBodyConfigHunter{
			Speed: 70.4,
			Position: &positioner.Position{
				X: 19,
				Y: 5,
				Z: 9,
			},
		}

		requestBodyHunter, err := json.Marshal(sampleBodyHunter)
		if err != nil {
			fmt.Println("Error coding request body:", err)
			return
		}
		_, err = http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewReader(requestBodyHunter))
		assert.NoError(t, err)
		sampleBodyPrey := handler.RequestBodyConfigPrey{
			Speed: 5.4,
			Position: &positioner.Position{
				X: 0,
				Y: 0,
				Z: 0,
			},
		}

		requestBodyPrey, err := json.Marshal(sampleBodyPrey)
		if err != nil {
			fmt.Println("Error coding request body:", err)
			return
		}
		_, err = http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewReader(requestBodyPrey))
		assert.NoError(t, err)
		//when
		var requestBody []byte
		res, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.Contains(t, bodyString, "caça concluída com sucesso")
	})
	t.Run("case 3 : failed hunt after configuring hunter and prey", func(t *testing.T) {
		//given
		sampleBodyHunter := handler.RequestBodyConfigHunter{
			Speed: 10.4,
			Position: &positioner.Position{
				X: 19,
				Y: 5,
				Z: 9,
			},
		}

		requestBodyHunter, err := json.Marshal(sampleBodyHunter)
		if err != nil {
			fmt.Println("Error coding request body:", err)
			return
		}
		_, err = http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", bytes.NewReader(requestBodyHunter))
		assert.NoError(t, err)
		sampleBodyPrey := handler.RequestBodyConfigPrey{
			Speed: 50.4,
			Position: &positioner.Position{
				X: 0,
				Y: 0,
				Z: 0,
			},
		}

		requestBodyPrey, err := json.Marshal(sampleBodyPrey)
		if err != nil {
			fmt.Println("Error coding request body:", err)
			return
		}
		_, err = http.Post("http://localhost:8080/hunter/configure-prey", "application/json", bytes.NewReader(requestBodyPrey))
		assert.NoError(t, err)
		//when
		var requestBody []byte
		res, err := http.Post("http://localhost:8080/hunter/hunt", "application/json", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		//then
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		bodyString := string(body)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.Contains(t, bodyString, "caça concluída sem sucesso")
	})
}
