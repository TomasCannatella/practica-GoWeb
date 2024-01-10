package repository

import (
	"ejercicio2/internal/product"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductRepository struct {
	products []product.Product
}

func (p *ProductRepository) LoadJson() {
	/*
		Opcion 1 - Utilizar os.Open y luego utilizar el newDecoder
	*/
	file, err := os.Open("./products.json")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	readerProducts := json.NewDecoder(file)
	err = readerProducts.Decode(&p.products)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	/*
		Opcion 2 - Utilizar readFile donde te devuelve el file en slice byte
	*/
	// file, err := os.ReadFile("./products.json")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = json.Unmarshal(file, &p.products)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

func (p *ProductRepository) GetAll(w http.ResponseWriter, r *http.Request) {

	productsJson := p.products
	jsonData, err := json.Marshal(productsJson)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

	// writer := json.NewEncoder(w)
	// err := writer.Encode(p.products)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Error del servidor"))
	// }
}

func (p *ProductRepository) GetById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	var returnProduct product.Product
	for _, product := range p.products {
		if product.Id == idInt {
			returnProduct = product
			break
		}
	}

	if returnProduct.Id != 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		jsonData, err := json.Marshal(returnProduct)
		if err != nil {
			return
		}
		w.Write(jsonData)
	}
}

func (p *ProductRepository) GetPriceGreaterThan(w http.ResponseWriter, r *http.Request) {
	var products []product.Product = make([]product.Product, 0)

	priceGt := r.URL.Query().Get("priceGt")

	priceFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("err"))
		return
	}
	for _, product := range p.products {
		if priceFloat <= product.Price {
			products = append(products, product)
		}
	}

	jsonData, err := json.Marshal(products)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("err"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
