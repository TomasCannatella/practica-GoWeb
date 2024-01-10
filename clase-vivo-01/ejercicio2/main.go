package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Person struct {
	Name     string
	Lastname string
}

func helloPerson(w http.ResponseWriter, r *http.Request) {
	readerBodyJson := r.Body
	myDecoder := json.NewDecoder(readerBodyJson)

	var p Person
	err := myDecoder.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.Write([]byte("Hello " + p.Name + " " + p.Lastname))
}

func main() {
	router := chi.NewRouter()

	router.Post("/greetings", helloPerson)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Println(err)
	}
}
