package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/platform/web/response"
)

// NewHunter returns a new Hunter handler.
func NewHunter(ht hunter.Hunter, pr prey.Prey) *Hunter {
	return &Hunter{ht: ht, pr: pr}
}

// Hunter returns handlers to manage hunting.
type Hunter struct {
	// ht is the Hunter interface that this handler will use
	ht hunter.Hunter
	// pr is the Prey interface that the hunter will hunt
	pr prey.Prey
}

// RequestBodyConfigPrey is an struct to configure the prey for the hunter in JSON format.
type RequestBodyConfigPrey struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// Example
// curl -X POST http://localhost:8080/hunter/configure-prey \
// -H "Content-Type: application/json" \
// -d '{
//   "speed": 4.0,
//   "position": {
//     "X": 0.1,
//     "Y": 0.4,
//     "Z": 3.1
//   }
// }'

// ConfigurePrey configures the prey for the hunter.
func (h *Hunter) ConfigurePrey(w http.ResponseWriter, r *http.Request) {
	log.Println("call ConfigurePrey")

	// request
	var preyConfig RequestBodyConfigPrey
	err := json.NewDecoder(r.Body).Decode(&preyConfig)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Error decoding JSON: "+err.Error())
		return
	}

	// process
	h.pr.Configure(preyConfig.Speed, preyConfig.Position)

	// response
	response.Text(w, http.StatusOK, "The prey is configured correctly")
}

// RequestBodyConfigHunter is an struct to configure the hunter in JSON format.
type RequestBodyConfigHunter struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// ConfigureHunter configures the hunter.
// if an error occurs when decoding the body, returns status code 400
// if no error occurs, returns status code 200
func (h *Hunter) ConfigureHunter(w http.ResponseWriter, r *http.Request) {
	log.Println("call ConfigureHunter")

	// request
	var hunterConfig RequestBodyConfigHunter
	err := json.NewDecoder(r.Body).Decode(&hunterConfig)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Error decoding JSON: "+err.Error())
		return
	}

	// process
	h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

	// response
	response.Text(w, http.StatusOK, "The shark is configured correctly")
}

// Hunt hunts the prey.
// if no error occurs, returns status code 200
// if ErrCanNotHunt occurs, returns status code 200
// if another error occurs, returns status code 500
func (h *Hunter) Hunt(w http.ResponseWriter, r *http.Request) {
	log.Println("call Hunt")

	duration, err := h.ht.Hunt(h.pr)

	if err == nil {
		msg := fmt.Sprintf("hunt completed with success | prey caught: true | duration: %.2f", duration)
		response.Text(w, http.StatusOK, msg)
	} else if err == hunter.ErrCanNotHunt {
		msg := fmt.Sprintf("hunt completed with success | prey caught: false | duration: %.2f", duration)
		response.Text(w, http.StatusOK, msg)
	} else {
		response.Error(w, http.StatusInternalServerError, "Internal error: "+err.Error())

	}
}
