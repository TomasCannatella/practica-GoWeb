package repository

import (
	"ejercicio1/internal/product"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductRepository struct {
	products map[int]product.Product
	lastId   int
}

func NewProductRepository() *ProductRepository {
	var product ProductRepository = ProductRepository{}
	product.LoadJson()
	product.lastId = len(product.products)
	return &product
}

func (p *ProductRepository) LoadJson() {
	var products []product.Product

	file, err := os.Open("./products.json")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	readerProducts := json.NewDecoder(file)
	err = readerProducts.Decode(&products)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	p.products = make(map[int]product.Product, len(products))

	for _, product := range products {
		p.products[product.Id] = product
	}
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
}

func (p *ProductRepository) GetById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	product, ok := p.products[idInt]
	if !ok {
		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("product not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonData, err := json.Marshal(product)
	if err != nil {
		return
	}
	w.Write(jsonData)
}

func (p *ProductRepository) GetPriceGreaterThan(w http.ResponseWriter, r *http.Request) {
	var products []product.Product = make([]product.Product, 0)

	priceGt := r.URL.Query().Get("priceGt")

	priceFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid format price"))
		return
	}

	for _, product := range p.products {
		if priceFloat <= product.Price {
			products = append(products, product)
		}
	}

	jsonData, err := json.Marshal(products)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error parse json"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

/*
    Ejercicio 1: Añadir un producto
	En esta ocasión vamos a añadir un producto al slice cargado en memoria.
	Dentro de la ruta /products añadimos el método POST, al cual vamos a enviar
	en el cuerpo de la request el nuevo producto.

	El mismo tiene ciertas restricciones, conozcámoslas:

		-	No es necesario pasar el Id, al momento de añadirlo se debe inferir del estado de la lista de productos,
			verificando que no se repitan ya que debe ser un campo único.

		-	Ningún dato puede estar vacío, exceptuando is_published (vacío indica un valor false).
		-	El campo code_value debe ser único para cada producto.
		- 	Los tipos de datos deben coincidir con los definidos en el planteo del problema.
		-	La fecha de vencimiento debe tener el formato: XX/XX/XXXX, además debemos verificar que día, mes y año sean valores válidos.
	Recordá: si una consulta está mal formulada por parte del cliente, el status code cae en los 4XX.

*/

func (p *ProductRepository) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Capturar la info del body
		var body product.BodyRequestProduct

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
		}

		//fmt.Println(newProductJSON)
		// Serializar los datos y armar la estructura
		newProduct := product.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}
		(*p).lastId++
		newProduct.Id = p.lastId

		// Validaciones de negocio
		if err := newProduct.ValidateBusinessRule(); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}

		// Validar que el producto no exista
		for _, product := range p.products {
			if product.Name == newProduct.Name {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("product alredy exist"))
				return
			}
		}

		//Guardar el producto
		p.products[newProduct.Id] = newProduct

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product successfully saved",
			"data":    newProduct,
		})
	}
}
