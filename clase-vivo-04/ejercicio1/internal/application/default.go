package application

import (
	hand "ejercicio1/internal/handler"
	repo "ejercicio1/internal/repository"
	"ejercicio1/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewDefaultHTTP(addr string) *DefaultHTTP {
	return &DefaultHTTP{
		addr: addr,
	}
}

type DefaultHTTP struct {
	addr string
}

func (h *DefaultHTTP) Run() (err error) {
	//inicio de dependencias
	// - repository
	rp := repo.NewProductMap()
	// - service
	sv := service.NewProductDefault(rp)

	// - handler
	hd := hand.NewDefaultProduct(sv)

	router := chi.NewRouter()

	router.Route("/product", func(r chi.Router) {
		//Get
		r.Get("/", hd.GetAll())
		r.Get("/{id}", hd.GetById)
		r.Get("/search", hd.GetPriceGreaterThan)

		//Post
		r.Post("/", hd.Create())
	})

	http.ListenAndServe(h.addr, router)
	return
}
