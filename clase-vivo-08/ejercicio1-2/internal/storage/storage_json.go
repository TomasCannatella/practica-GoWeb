package storage

import (
	"ejercicio1/internal"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type StorageJson struct {
	Path string
}

func (s *StorageJson) ReadAll() (products map[int]internal.Product, err error) {
	// slice auxiliar
	var sliceProduct []internal.Product
	// abro el json
	file, err := os.Open(s.Path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	// paso el json al slice
	readerProducts := json.NewDecoder(file)
	err = readerProducts.Decode(&sliceProduct)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	//instancio el mapa de mi repository con el tamaño del slice auxiliar
	products = make(map[int]internal.Product, len(sliceProduct))

	//recorro el slice auxiliar y almaceno los datos en el mapa
	for _, product := range sliceProduct {
		products[product.Id] = product
	}
	return products, nil
}

func (s *StorageJson) WriteAll(products map[int]internal.Product) (err error) {

	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer file.Close()

	var sliceProducts = make([]internal.Product, 0)
	for _, product := range products {
		sliceProducts = append(sliceProducts, product)
	}

	err = json.NewEncoder(file).Encode(&sliceProducts)
	if err != nil {
		return
	}
	log.Println("write product json")
	return
}
