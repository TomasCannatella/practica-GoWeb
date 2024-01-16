package storage

import "ejercicio1/internal"

type Storage interface {
	ReadAll() (products map[int]internal.Product, err error)
	WriteAll(products map[int]internal.Product) (err error)
}
