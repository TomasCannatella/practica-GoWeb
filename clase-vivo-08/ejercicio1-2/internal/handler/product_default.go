package handler

import (
	"ejercicio1/internal"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
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

type ProductJSON struct {
	Id          int     `json:"int"`
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
			response.Error(w, http.StatusNoContent, "products not found")
			return
		}

		var productsJson []ProductJSON
		for _, product := range products {
			productJson := ProductJSON{
				Id:          product.Id,
				Name:        product.Name,
				Quantity:    product.Quantity,
				CodeValue:   product.CodeValue,
				IsPublished: product.IsPublished,
				Expiration:  product.Expiration,
				Price:       product.Price,
			}
			productsJson = append(productsJson, productJson)
		}

		//response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "all products",
			"data":    productsJson,
		})
	}
}

func (d *DefaultProduct) GetById(w http.ResponseWriter, r *http.Request) {
	//request
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	//process
	product, err := d.sv.GetById(idInt)
	if err != nil {
		if errors.Is(err, internal.ErrProductNotFound) {
			response.Error(w, http.StatusNotFound, "product not found")
			return
		}
	}
	productJson := ProductJSON{
		Id:          product.Id,
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  product.Expiration,
		Price:       product.Price,
	}
	//Response
	response.JSON(w, http.StatusOK,
		map[string]any{
			"message": "product found",
			"data":    productJson,
		})
}

func (d *DefaultProduct) GetPriceGreaterThan(w http.ResponseWriter, r *http.Request) {
	//request
	priceGt := r.URL.Query().Get("priceGt")
	priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "invalid format price")
		return
	}

	//process
	products, err := d.sv.GetPriceGreaterThan(priceGtFloat)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			response.Error(w, http.StatusInternalServerError, "products not found")
			return
		case errors.Is(err, internal.ErrInvalidFormat):
			response.Error(w, http.StatusBadRequest, "the price must be greater than 0")
		default:
			response.Error(w, http.StatusInternalServerError, "Unkown error")
		}

		return
	}

	var productsJson []ProductJSON
	for _, product := range products {
		productJson := ProductJSON{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}
		productsJson = append(productsJson, productJson)
	}

	//Response
	response.JSON(w, http.StatusOK,
		map[string]any{
			"message": "products found",
			"data":    productsJson,
		})

}

func (d *DefaultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request
		var body BodyRequestProduct
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fmt.Println(err)
			response.Error(w, http.StatusBadRequest, "invalid body")
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
				response.Error(w, http.StatusBadRequest, "invalid body")
			case errors.Is(err, internal.ErrProductNameAlredyExist):
				response.Error(w, http.StatusConflict, "product alredy exist with this name")
			case errors.Is(err, internal.ErrProductCodeValueAlredyExist):
				response.Error(w, http.StatusConflict, "product alredy exist with this code value")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		productJson := ProductJSON{
			Id:          newProduct.Id,
			Name:        newProduct.Name,
			Quantity:    newProduct.Quantity,
			CodeValue:   newProduct.CodeValue,
			IsPublished: newProduct.IsPublished,
			Expiration:  newProduct.Expiration,
			Price:       newProduct.Price,
		}

		//Response
		response.JSON(w, http.StatusCreated,
			map[string]any{
				"message": "product successfully saved",
				"data":    productJson,
			})
	}
}

func (d *DefaultProduct) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid format id")
			return
		}

		//request
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		//
		var bodyMap map[string]any
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		// - validate body
		if err := ValidateKeyExistante(bodyMap, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body: "+err.Error())
			return
		}

		var body BodyRequestProduct
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
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

			if errors.Is(err, internal.ErrProductNotFound) {
				response.Error(w, http.StatusNotFound, "product not found")
				return
			}

			response.Error(w, http.StatusInternalServerError, "unkown error")
			return
		}

		productJson := ProductJSON{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		//Response
		response.JSON(w, http.StatusOK,
			map[string]any{
				"message": "product successfully updated",
				"data":    productJson,
			})

	}
}

func (d *DefaultProduct) UpdatePartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid format id")
			return
		}

		//Get original product
		product, err := d.sv.GetById(idInt)
		if err != nil {
			if errors.Is(err, internal.ErrProductNotFound) {
				response.Error(w, http.StatusNotFound, "product not found")
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
			response.Error(w, http.StatusBadRequest, "invalid format id")
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
			if errors.Is(err, internal.ErrProductNotFound) {
				response.Error(w, http.StatusNotFound, "product not found")
				return
			}

			response.Error(w, http.StatusInternalServerError, "unkown error")
			return
		}

		productJson := ProductJSON{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		//Response
		response.JSON(w, http.StatusOK,
			map[string]any{
				"message": "product successfully updated",
				"data":    productJson,
			})

	}
}

func (d *DefaultProduct) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid format id")
			return
		}

		err = d.sv.Delete(idInt)
		if err != nil {
			if errors.Is(err, internal.ErrProductNotFound) {
				response.Error(w, http.StatusNotFound, "product not found")
				return
			}

			response.Error(w, http.StatusInternalServerError, "unkown error")
		}

		//Response
		response.Text(w, http.StatusNoContent, "")
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
