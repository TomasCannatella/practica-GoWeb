package repository

import (
	"ejercicio1/internal"
	"ejercicio1/internal/storage"
	"fmt"
)

type ProductMap struct {
	db     map[int]internal.Product
	lastId int
}

var s storage.StorageJson = storage.StorageJson{
	Path: "./products.json",
}

func NewProductMap(p map[int]internal.Product) *ProductMap {
	if p == nil {
		products, err := s.ReadAll()
		if err != nil {
			fmt.Println(err)
			return &ProductMap{
				db:     make(map[int]internal.Product, 0),
				lastId: 0,
			}
		}
		return &ProductMap{
			db:     products,
			lastId: len(products),
		}
	}

	return &ProductMap{
		db:     p,
		lastId: len(p),
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
	err = s.WriteAll((*p).db)
	return
}

func (p *ProductMap) Update(product *internal.Product) (err error) {
	_, ok := (*p).db[product.Id]
	if !ok {
		return internal.ErrProductNotFound
	}

	(*p).db[(*product).Id] = *product
	err = s.WriteAll((*p).db)
	return
}

func (p *ProductMap) Delete(id int) (err error) {
	_, ok := (*p).db[id]
	if !ok {
		return internal.ErrProductNotFound
	}

	delete(p.db, id)
	err = s.WriteAll((*p).db)
	return
}
