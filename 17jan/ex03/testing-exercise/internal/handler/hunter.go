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
		var configPrey RequestBodyConfigPrey
		err := json.NewDecoder(r.Body).Decode(&configPrey)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
			return
		}

		// process
		h.pr.Configure(configPrey.Speed, configPrey.Position)

		// response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "A presa está configurada corretamente",
			"data":    nil,
		})
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
		var configHunter RequestBodyConfigHunter
		err := json.NewDecoder(r.Body).Decode(&configHunter)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
			return
		}

		// process
		h.ht.Configure(configHunter.Speed, configHunter.Position)

		// response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "hunter configured",
			"data":    nil,
		})
	}
}

// Hunt hunts the prey.
func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		duration, err := h.ht.Hunt(h.pr)
		if err != nil {
			if errors.Is(err, hunter.ErrCanNotHunt) {
				response.JSON(w, http.StatusOK, err.Error())
				return
			} else {
				response.Error(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}
		ok := err == nil

		// response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "hunt done",
			"data": map[string]any{
				"success":  ok,
				"duration": duration,
			},
		})
	}
}
