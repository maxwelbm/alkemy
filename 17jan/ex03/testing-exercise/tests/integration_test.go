//go:build special
// +build special

package integrations_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testdoubles/internal/handler"
	"testdoubles/internal/positioner"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationPost(t *testing.T) {
	// Config Prey
	requestBodyPrey := handler.RequestBodyConfigPrey{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
			Z: 100,
		},
	}

	localHost := "http://localhost:8080/hunter"
	content := "application/json"

	requestBodyPreyJSON, err := json.Marshal(requestBodyPrey)
	require.NoError(t, err)

	res, err := http.Post(localHost+"/configure-prey", content, bytes.NewReader(requestBodyPreyJSON))
	require.NoError(t, err)

	res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	// Config Hunter
	requestBodyHunter := handler.RequestBodyConfigHunter{
		Speed: 10.5,
		Position: &positioner.Position{
			X: 100,
			Y: 200,
			Z: 100,
		},
	}

	requestBodyHunterJSON, err := json.Marshal(requestBodyHunter)
	require.NoError(t, err)

	res, err = http.Post(localHost+"/configure-hunter", content, bytes.NewReader(requestBodyHunterJSON))
	require.NoError(t, err)

	res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	// Hunt
	res, err = http.Post(localHost+"/hunt", content, nil)
	require.NoError(t, err)

	res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
}
