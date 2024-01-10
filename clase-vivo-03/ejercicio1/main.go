/*
Ejercicio 1 : Iniciando el proyecto
Debemos crear un repositorio en github.com para poder subir nuestros avances.
Este repositorio es el que vamos a utilizar para llevar lo que realicemos durante las distintas prácticas de Go Web.

Primero debemos clonar el repositorio creado, luego iniciar nuestro proyecto de go con con el comando go mod init.
El siguiente paso será crear un archivo main.go donde deberán cargar en una slice, desde un archivo JSON, los datos de productos.
Esta slice se debe cargar cada vez que se inicie la API para realizar las distintas consultas.
*/
package main

import (
	"ejercicio1/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	var p repository.ProductRepository = *repository.NewProductRepository()

	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Pong"))
	})

	router.Route("/product", func(r chi.Router) {
		r.Get("/", p.GetAll)

		r.Get("/{id}", p.GetById)

		r.Get("/search", p.GetPriceGreaterThan)

		r.Post("/", p.Create())
	})

	http.ListenAndServe(":8080", router)
}
