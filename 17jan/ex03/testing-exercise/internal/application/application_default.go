package application

import (
	"log"
	"net/http"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/positioner"
	"testdoubles/internal/prey"
	"testdoubles/internal/simulator"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApplicationDefault struct {
	rt   *chi.Mux
	addr string
}

func NewApplicationDefault(addr string) *ApplicationDefault {
	defaultRouter := chi.NewRouter()
	defaultAddr := ":8080"
	if addr != "" {
		defaultAddr = addr
	}

	return &ApplicationDefault{
		rt:   defaultRouter,
		addr: defaultAddr,
	}
}

func (a *ApplicationDefault) TearDown() (err error) {
	return
}

func (a *ApplicationDefault) SetUp() (err error) {
	log.Println("call SetUp")

	ps := positioner.NewPositionerDefault()
	sm := simulator.NewCatchSimulatorDefault(&simulator.ConfigCatchSimulatorDefault{
		Positioner:     ps,
		MaxTimeToCatch: 10000.00,
	})
	ht := hunter.NewWhiteShark(hunter.ConfigWhiteShark{
		Speed:     0.0,
		Position:  &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0},
		Simulator: sm,
	})
	pr := prey.NewTuna(0.0, &positioner.Position{X: 0.0, Y: 0.0, Z: 0.0})
	hd := handler.NewHunter(ht, pr)

	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)
	a.rt.Route("/hunter", func(r chi.Router) {
		r.Post("/configure-prey", hd.ConfigurePrey)
		r.Post("/configure-hunter", hd.ConfigureHunter())
		r.Post("/hunt", hd.Hunt())
	})

	return
}
func (a *ApplicationDefault) Run() (err error) {
	err = http.ListenAndServe(a.addr, a.rt)
	return
}
