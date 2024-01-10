package repository

import (
	"ejercicio1/internal/product"
	"encoding/json"
	"fmt"
	"os"
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
