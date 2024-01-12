package repository

import (
	"ejercicio1/internal"
	"encoding/json"
	"fmt"
	"os"
)

type ProductMap struct {
	db     map[int]internal.Product
	lastId int
}

func NewProductMap() *ProductMap {
	var p ProductMap = ProductMap{}
	//Cargo los elementos del archivos Json
	p.LoadJson()
	//Seteo el ultimo id como el tamaño del map
	p.lastId = len(p.db)

	return &p
}

func (p *ProductMap) LoadJson() {
	// slice auxiliar
	var products []internal.Product
	// abro el json
	file, err := os.Open("./products.json")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	// paso el json al slice
	readerProducts := json.NewDecoder(file)
	err = readerProducts.Decode(&products)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	//instancio el mapa de mi repository con el tamaño del slice auxiliar
	p.db = make(map[int]internal.Product, len(products))

	//recorro el slice auxiliar y almaceno los datos en el mapa
	for _, product := range products {
		p.db[product.Id] = product
	}
}

func (p *ProductMap) GetAll() (products []internal.Product, err error) {
	for _, value := range p.db {
		products = append(products, value)
	}
	return
}

func (p *ProductMap) GetById(id int) (product internal.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}

	return
}

func (p *ProductMap) GetPriceGreaterThan(priceGt float64) (products []internal.Product, err error) {
	for _, product := range p.db {
		if priceGt <= product.Price {
			products = append(products, product)
		}
	}

	if len(products) <= 0 {
		err = internal.ErrProductNotFound
		return
	}
	return
}

func (p *ProductMap) Save(product *internal.Product) (err error) {
	//validar consistencia de los datos
	for _, v := range (*p).db {
		if v.Name == product.Name {
			return internal.ErrProductNameAlredyExist
		}

		if v.CodeValue == product.CodeValue {
			return internal.ErrProductCodeValueAlredyExist
		}

	}

	//incremento id
	(*p).lastId++
	(*product).Id = (*p).lastId

	//guardar producto
	(*p).db[(*product).Id] = *product
	return
}
