package handler

import (
	"errors"
	"fmt"
	"net/http"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/platform/web/request"
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

// ConfigurePrey configures the prey for the hunter.
func (h *Hunter) ConfigurePrey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var configPrey *RequestBodyConfigHunter
		// process
		err := request.JSON(r, &configPrey)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		// configure
		h.pr.Configure(configPrey.Speed, configPrey.Position)

		// response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "configured prey",
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
		var configHunter *RequestBodyConfigHunter

		// process
		err := request.JSON(r, &configHunter)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		h.ht.Configure(configHunter.Speed, configHunter.Position)

		// response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "configured hunter",
		})
	}
}

// Hunt hunts the prey.
func (h *Hunter) Hunt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// process
		dur, err := h.ht.Hunt(h.pr)

		// response
		if err != nil {
			switch {
			case errors.Is(err, hunter.ErrCanNotHunt):
				response.JSON(w, http.StatusOK, map[string]any{
					"message": "hunt completed",
					"status":  "hunter can not hunt the prey after " + fmt.Sprintf("%f", dur) + " seconds",
				})
			default:
				response.Error(w, http.StatusInternalServerError, "error hunting")
			}
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "hunt completed",
			"status":  "hunter hunted the prey in " + fmt.Sprintf("%f", dur) + " seconds",
		})
	}
}
