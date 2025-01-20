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
func (h *Hunter) ConfigurePrey(w http.ResponseWriter, r *http.Request) {
	log.Println("call ConfigurePrey")

	// request
	var preyConfig RequestBodyConfigPrey
	err := json.NewDecoder(r.Body).Decode(&preyConfig)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// process
	h.ht.Configure(preyConfig.Speed, preyConfig.Position)

	// response
	response.Text(w, http.StatusOK, "A presa está configurada corretamente")
}

// RequestBodyConfigHunter is an struct to configure the hunter in JSON format.
type RequestBodyConfigHunter struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// ConfigureHunter configures the hunter.
func (h *Hunter) ConfigureHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("call ConfigureHunter")

		// request
		var hunterConfig RequestBodyConfigHunter
		err := json.NewDecoder(r.Body).Decode(&hunterConfig)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
			return
		}

		// process
		h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

		// response
		response.Text(w, http.StatusOK, "O caçador está configurado corretamente")
	}
}

// Hunt hunts the prey.
func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		duration, err := h.ht.Hunt(h.pr)

		if err == nil {
			response.JSON(w, http.StatusOK, map[string]any{"message": "caça concluída", "duration": duration})
		}

		if err != nil {
			if errors.Is(err, hunter.ErrCanNotHunt) {
				response.JSON(w, http.StatusOK, map[string]any{"message": "caça concluída", "duration": duration, "error": "can not hunt the prey"})
			} else {
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}
	}
}
