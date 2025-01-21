package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
)

type RequestBodyConfigPrey struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

type RequestBodyConfigHunter struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

type Hunter struct {
	hunter hunter.Hunter
	prey   prey.Prey
}

func NewHunter(h hunter.Hunter, p prey.Prey) *Hunter {
	return &Hunter{
		hunter: h,
		prey:   p,
	}
}

func (h *Hunter) ConfigurePrey(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBodyConfigPrey
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.prey.Configure(reqBody.Speed, reqBody.Position)

	fmt.Fprint(w, "A presa está configurada corretamente")
}

func (h *Hunter) ConfigureHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody RequestBodyConfigHunter
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.prey.Configure(reqBody.Speed, reqBody.Position)

		fmt.Fprint(w, "O caçador está configurado corretamente")
	}
}

func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		duration, err := h.hunter.Hunt(h.prey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message":  "caça concluída",
			"duration": duration,
		}
		json.NewEncoder(w).Encode(response)
	}
}
