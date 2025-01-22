package handler

import (
	"encoding/json"
	"errors"
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

// Hunter is the HTTP handler that manages "hunting" endpoints.
type Hunter struct {
	ht hunter.Hunter // a "real" or "mock" that implements Hunter
	pr prey.Prey     // a "real" or "stub" that implements Prey
}

// RequestBodyConfigPrey is an struct to configure the prey for the hunter in JSON format.
type RequestBodyConfigPrey struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// ConfigurePrey configures the prey for the hunter.
func (h *Hunter) ConfigurePrey(w http.ResponseWriter, r *http.Request) {
	log.Println("call ConfigurePrey")

	var hunterConfig RequestBodyConfigPrey
	if err := json.NewDecoder(r.Body).Decode(&hunterConfig); err != nil {
		response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

	response.Text(w, http.StatusOK, "A presa está configurada corretamente")
}

// RequestBodyConfigHunter is an struct to configure the hunter in JSON format.
type RequestBodyConfigHunter struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// ConfigureHunter configures the hunter.
// Se falhar o parse do JSON, retornamos 400.
// Caso contrário, chama h.ht.Configure(...) e retorna 200.
func (h *Hunter) ConfigureHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("call ConfigureHunter")

		var cfg RequestBodyConfigHunter
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			response.Error(w, http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
			return
		}

		h.ht.Configure(cfg.Speed, cfg.Position)
		response.Text(w, http.StatusOK, "O tubarão está configurado corretamente")
	}
}

// Hunt hunts the prey.
// - Se nenhum erro => capturou => 200 com "caça concluída, capturada:true, tempo: X"
// - Se err == hunter.ErrCanNotHunt => 200 com "caçada concluída, capturada:false, tempo: X"
// - Se outro erro => 500
func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("call Hunt")

		duration, err := h.ht.Hunt(h.pr)
		switch {
		case err == nil:
			msg := fmt.Sprintf("caça concluída, capturada: true, tempo: %.2f", duration)
			response.Text(w, http.StatusOK, msg)
		case errors.Is(err, hunter.ErrCanNotHunt):
			msg := fmt.Sprintf("caçada concluída, capturada: false, tempo: %.2f", duration)
			response.Text(w, http.StatusOK, msg)
		default:
			// Se houver um erro "diferente de can not hunt", retorna 500
			response.Error(w, http.StatusInternalServerError, "Erro interno na caçada: "+err.Error())
		}
	}
}
