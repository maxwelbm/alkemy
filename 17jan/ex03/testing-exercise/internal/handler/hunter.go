package handler

import (
	"encoding/json"
	"errors"
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
func (h *Hunter) ConfigurePrey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("call ConfigurePrey")

		// request
		var hunterConfig RequestBodyConfigPrey
		err := json.NewDecoder(r.Body).Decode(&hunterConfig)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
			return
		}

		// process
		h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

		// response
		response.Text(w, http.StatusOK, "A presa está configurada corretamente")
	}
}

// RequestBodyConfigHunter is an struct to configure the hunter in JSON format.
type RequestBodyConfigHunter struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// ConfigureHunter configures the hunter.
func (h *Hunter) ConfigureHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var hunterConfig RequestBodyConfigHunter
		err := json.NewDecoder(r.Body).Decode(&hunterConfig)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Incorrect hunter configuration, issues with request body")
			return
		}

		// process
		h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

		// response
		response.JSON(w, http.StatusOK, map[string]any{"message": "Hunter configured"})
	}
}

// Hunt hunts the prey.
func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		duration, err := h.ht.Hunt(h.pr)

		// response
		if err != nil {
			if errors.Is(err, hunter.ErrCanNotHunt) {
				response.JSON(w, http.StatusOK, map[string]any{"message": "hunting completed", "duration": duration, "error": "can not hunt the prey"})
				return
			} else {
				response.Error(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}

		response.JSON(w, http.StatusOK, map[string]any{"message": "hunting completed", "duration": duration})
	}
}
