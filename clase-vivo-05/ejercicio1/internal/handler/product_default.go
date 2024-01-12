package repository

import (
	"ejercicio1/internal"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (d *DefaultProduct) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid format id"))
			return
		}

		//response
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		//
		var bodyMap map[string]any
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		// - validate body
		if err := ValidateKeyExistante(bodyMap, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body: " + err.Error()))
			return
		}

		var body BodyRequestProduct
		if err := json.Unmarshal(bytes, &body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		product := internal.Product{
			Id:          idInt,
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		err = d.sv.Update(&product)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			if errors.Is(err, internal.ErrProductNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("product not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unkown error"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product successfully updated",
			"data":    product,
		})

	}
}

func (d *DefaultProduct) UpdatePartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid format id"))
			return
		}

		//Get original product
		product, err := d.sv.GetById(idInt)
		if err != nil {
			if errors.Is(err, internal.ErrProductNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("product not found"))
				return
			}
			return
		}

		//serialize product to format request
		reqBody := BodyRequestProduct{
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		//updated  requestProduct data with requested information
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid format id"))
			return
		}

		//serialize produc to internal product
		product = internal.Product{
			Id:          idInt,
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  reqBody.Expiration,
			Price:       reqBody.Price,
		}

		//update product
		err = d.sv.Update(&product)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")

			if errors.Is(err, internal.ErrProductNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("product not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unkown error"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product successfully updated",
			"data":    product,
		})

	}
}

func (d *DefaultProduct) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid format id"))
			return
		}

		err = d.sv.Delete(idInt)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")

			if errors.Is(err, internal.ErrProductNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("product not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unkown error"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product successfully deleted",
		})
	}
}

func ValidateKeyExistante(mp map[string]any, keys ...string) (err error) {
	for _, k := range keys {
		if _, ok := mp[k]; !ok {
			return fmt.Errorf("key %s not found", k)
		}
	}
	return
}
