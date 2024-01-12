package internal

import "errors"

var (
	ErrFieldRequire  = errors.New("field required")
	ErrInvalidFormat = errors.New("invalid format")
)

type ProductService interface {
	Save(product *Product) (err error)
	GetAll() (products []Product, err error)
	GetById(id int) (product Product, err error)
	GetPriceGreaterThan(priceGt float64) (products []Product, err error)
	Update(product *Product) (err error)
	Delete(id int) (err error)
}
