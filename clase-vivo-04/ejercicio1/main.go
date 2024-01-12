/*
Ejercicio 1: Dominios
Es momento de organizar nuestra API, seguiremos la siguiente estructura de carpetas.
El package cmd representa los puntos de entrada de nuestra app. Por otro lado,
en el package internal tendremos nuestro domain. Por ahora solo tenemos uno: product.
Luego las implementaciones se van organizando en distintos packages, como repository, service, handler y application.
*/
package main

func main() {
	// var p repository.ProductRepository = *repository.NewProductRepository()

	// router := chi.NewRouter()

	// router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(200)
	// 	w.Write([]byte("Pong"))
	// })

	// router.Route("/product", func(r chi.Router) {
	// 	r.Get("/", p.GetAll)

	// 	r.Get("/{id}", p.GetById)

	// 	r.Get("/search", p.GetPriceGreaterThan)

	// 	r.Post("/", p.Create())
	// })

	// http.ListenAndServe(":8080", router)
}
