package application

import (
	hand "ejercicio1/internal/handler"
	"ejercicio1/internal/middleware"
	repo "ejercicio1/internal/repository"
	"ejercicio1/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DefaultHTTP struct {
	addr  string
	token string
}

func NewDefaultHTTP(addr, token string) *DefaultHTTP {
	return &DefaultHTTP{
		addr:  addr,
		token: token,
	}
}

func (h *DefaultHTTP) Run() (err error) {
	//inicio de dependencias
	auth := middleware.NewAuthenticator(h.token)
	logger := middleware.NewLogger()

	// - repository
	rp := repo.NewProductMap(nil)
	// - service
	sv := service.NewProductDefault(rp)

	// - handler
	hd := hand.NewDefaultProduct(sv)

	router := chi.NewRouter()
	router.Use(logger.Log)
	router.Route("/product", func(r chi.Router) {
		//Get
		r.Get("/", hd.GetAll())
		r.Get("/{id}", hd.GetById)
		r.Get("/search", hd.GetPriceGreaterThan)

		r.Route("/", func(r chi.Router) {
			r.Use(auth.Auth)

			r.Post("/", hd.Create())

			r.Put("/{id}", hd.Update())

			r.Patch("/{id}", hd.UpdatePartial())

			r.Delete("/{id}", hd.Delete())
		})
	})

	http.ListenAndServe(h.addr, router)
	return
}
