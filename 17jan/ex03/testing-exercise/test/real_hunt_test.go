//go:build special
// +build special

package test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRealHunt(t *testing.T) {
	// Configure Prey
	preyBody := `{
		"speed": 4.0,
		"position": {
		  "X": 0.1,
		  "Y": 0.4,
		  "Z": 3.1
		}
	  }`
	res, err := http.Post("http://localhost:8080/hunter/configure-prey", "application/json", strings.NewReader(preyBody))
	require.NoError(t, err)
	res.Body.Close()

	// Configure Hunter
	hunterBody := `{
		"speed": 1000.0,
		"position": {
		  "X": 0.5,
		  "Y": 0.5,
		  "Z": 0.5
		}
	  }`
	res, err = http.Post("http://localhost:8080/hunter/configure-hunter", "application/json", strings.NewReader(hunterBody))
	require.NoError(t, err)
	res.Body.Close()

	// Hunt
	res, err = http.Post("http://localhost:8080/hunter/hunt", "application/json", nil)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}
