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

func NewHunter(ht hunter.Hunter, pr prey.Prey) *Hunter {
	return &Hunter{ht: ht, pr: pr}
}

type Hunter struct {
	ht hunter.Hunter
	pr prey.Prey
}

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

	var hunterConfig RequestBodyConfigPrey
	err := json.NewDecoder(r.Body).Decode(&hunterConfig)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	h.pr.Configure(hunterConfig.Speed, hunterConfig.Position)

	response.Text(w, http.StatusOK, "A presa está configurada corretamente")
}

type RequestBodyConfigHunter struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

func (h *Hunter) ConfigureHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("call ConfigureHunter")

		var hunterConfig RequestBodyConfigHunter
		err := json.NewDecoder(r.Body).Decode(&hunterConfig)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
			return
		}

		h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

		response.Text(w, http.StatusOK, "O hunter está configurado corretamente")
	}
}

func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		duration, err := h.ht.Hunt(h.pr)

		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// response
		response.Text(w, http.StatusOK, "A presa foi caçada em "+fmt.Sprintf("%f", duration))
	}
}
