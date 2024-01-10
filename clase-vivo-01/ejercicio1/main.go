package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}

func main() {
	router := chi.NewRouter()

	router.Get("/ping", pong)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Println(err)
	}
}
