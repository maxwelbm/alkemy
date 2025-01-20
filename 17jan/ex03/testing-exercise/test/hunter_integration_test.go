package test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testdoubles/internal/application"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	url_base   = "http://localhost:8080/hunter"
	app_json   = "application/json"
	ErrSendReq = errors.New("error while sending request to server")
	ErrReadRes = errors.New("error while reading the body of the response")
)

func init() {
	go func() {
		// Creating a new application
		server := application.NewApplicationDefault(":8080")

		// Setting up the application
		err := server.SetUp()
		if err != nil {
			panic(err)
		}

		// Starting the application
		err = server.Run()
		if err != nil {
			panic(err)
		}
	}()
}
func TestIntegrationHunterConfigPrey(t *testing.T) {
	// Sending a request to the server
	bodyReq := `{
		"speed": 100.0,
		"position": {
			"X": 10.0,
			"Y": 20.0,
			"Z": 30.0
		}
	}`
	resp, err := http.Post(url_base+"/configure-prey", app_json, strings.NewReader(bodyReq))
	if err != nil {
		panic(ErrSendReq)
	}
	defer resp.Body.Close()

	// Reading the body of the response
	bodyRes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(ErrReadRes)
	}

	// Assertions
	expectedBody := `{"message": "configured prey"}`
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.JSONEq(t, expectedBody, string(bodyRes))
}

func TestIntegrationHunterConfigureHunter(t *testing.T) {
	// Sending a request to the server
	bodyReq := `{
		"speed": 100.0,
		"position": {
			"X": 10.0,
			"Y": 20.0,
			"Z": 30.0
		}
	}`
	resp, err := http.Post(url_base+"/configure-hunter", app_json, strings.NewReader(bodyReq))
	if err != nil {
		panic(ErrSendReq)
	}
	defer resp.Body.Close()

	// Reading the body of the response
	bodyRes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(ErrReadRes)
	}

	// Assertions
	expectedBody := `{"message": "configured hunter"}`
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.JSONEq(t, expectedBody, string(bodyRes))
}

func TestIntegrationHunterHunt(t *testing.T) {
	bodyReqPrey := `{
		"speed": 5.0,
		"position": {
		  "X": 1.0,
		  "Y": 2.0,
		  "Z": 3.0
		}
	}`
	bodyReqHunter := `{
		"speed": 10.0,
		"position": {
		  "X": 1.5,
		  "Y": 3.5,
		  "Z": 5.5
		}
	}`
	// Configuring the prey
	_, err := http.Post(url_base+"/configure-prey", app_json, strings.NewReader(bodyReqPrey))
	if err != nil {
		panic(ErrSendReq)
	}

	// Configuring the hunter
	_, err = http.Post(url_base+"/configure-hunter", app_json, strings.NewReader(bodyReqHunter))
	if err != nil {
		panic(ErrSendReq)
	}

	// Sending a request to the server
	resp, err := http.Post(url_base+"/hunt", app_json, nil)
	if err != nil {
		panic(ErrSendReq)
	}
	defer resp.Body.Close()

	// Reading the body of the response
	bodyRes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(ErrReadRes)
	}

	// Assertions
	expectedBody := `{"message":"hunt completed", "status":"hunter hunted the prey in 0.591608 seconds"}`
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.JSONEq(t, expectedBody, string(bodyRes))
}
