package internal

import (
	"errors"
)

var (
	ErrProductNameAlredyExist      = errors.New("product name alredy exist")
	ErrProductCodeValueAlredyExist = errors.New("product code value alredy exist")
	ErrProductAlredyExist          = errors.New("movie alredy exist")
	ErrProductNotFound             = errors.New("product not found")
)

type ProductRepository interface {
	Save(product *Product) (err error)
	GetAll() (products []Product, err error)
	GetById(id int) (product Product, err error)
	GetPriceGreaterThan(priceGt float64) (products []Product, err error)
}
