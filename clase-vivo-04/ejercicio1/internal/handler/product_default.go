package repository

import (
	"ejercicio1/internal"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DefaultProduct struct {
	sv internal.ProductService
}

func NewDefaultProduct(sv internal.ProductService) *DefaultProduct {
	return &DefaultProduct{
		sv: sv,
	}
}

type BodyRequestProduct struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (d *DefaultProduct) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request
		// ....

		//process
		products, err := d.sv.GetAll()
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("products not found"))
		}

		//response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "all products",
			"data":    products,
		})
	}
}

func (d *DefaultProduct) GetById(w http.ResponseWriter, r *http.Request) {

	//request
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	//process
	product, err := d.sv.GetById(idInt)
	if err != nil {
		if errors.Is(err, internal.ErrProductNotFound) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("product not found"))
			return
		}
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "product found",
		"data":    product,
	})
}

func (d *DefaultProduct) GetPriceGreaterThan(w http.ResponseWriter, r *http.Request) {
	//request
	priceGt := r.URL.Query().Get("priceGt")
	priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid format price"))
		return
	}

	//process
	products, err := d.sv.GetPriceGreaterThan(priceGtFloat)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("products not found"))
			return
		case errors.Is(err, internal.ErrInvalidFormat):
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("the price must be greater than 0"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unkown error"))
		}

		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "products found",
		"data":    products,
	})

}

func (d *DefaultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//response
		var body BodyRequestProduct
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
		}

		//fmt.Println(newProductJSON)
		// Serializar los datos y armar la estructura
		newProduct := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		//process
		err := d.sv.Save(&newProduct)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrFieldRequire), errors.Is(err, internal.ErrInvalidFormat):
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid body"))
			case errors.Is(err, internal.ErrProductAlredyExist):
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("product alredy exist "))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			return
		}

		//request
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product successfully saved",
			"data":    newProduct,
		})
	}
}
