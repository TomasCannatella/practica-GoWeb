/*
Ejercicio 2 : Creando un servidor
Vamos a levantar un servidor en el puerto 8080. Para probar nuestros endpoints haremos uso de postman.

Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
Crear una ruta /products/:id que nos devuelva un producto por su id.
Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
*/
package main

import (
	"ejercicio2/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	var p repository.ProductRepository = repository.ProductRepository{}
	p.LoadJson()

	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Pong"))
	})

	router.Route("/product", func(r chi.Router) {
		r.Get("/", p.GetAll)

		r.Get("/{id}", p.GetById)

		r.Get("/search", p.GetPriceGreaterThan)
	})

	http.ListenAndServe(":8080", router)
}
