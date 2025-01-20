package handler

import (
	"encoding/json"
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
	ht hunter.Hunter
	pr prey.Prey
}

// RequestBodyConfigPrey is a struct to configure the prey for the hunter in JSON format.
type RequestBodyConfigPrey struct {
	Speed    float64              `json:"speed"`
	Position *positioner.Position `json:"position"`
}

// ConfigurePrey configures the prey for the hunter.
func (h *Hunter) ConfigurePrey(w http.ResponseWriter, r *http.Request) {
	log.Println("call ConfigurePrey")

	// request
	var hunterConfig RequestBodyConfigPrey
	err := json.NewDecoder(r.Body).Decode(&hunterConfig)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "problems with the request body")
		return
	}

	// process
	h.ht.Configure(hunterConfig.Speed, hunterConfig.Position)

	// response
	response.JSON(w, http.StatusOK, map[string]string{"message": "Prey set up"})
}

// Hunt hunts the prey and returns the result.
func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("call Hunt")

		// process the hunt
		duration, err := h.ht.Hunt(h.pr)
		if err != nil {
			log.Println("error during hunt:", err)
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// build response
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"message": "hunt complete",
			"data": map[string]interface{}{
				"success":  true, // Assume true for success if no error
				"duration": duration,
			},
		})
	}
}

// ConfigureHunter configures the hunter.
func (h *Hunter) ConfigureHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process

		// response
	}
}
