/*
Ejercicio 1 - Dominios

Es momento de organizar nuestra API, seguiremos la siguiente estructura de carpetas.
El package cmd representa los puntos de entrada de nuestra app. Por otro lado,
en el package internal tendremos nuestro domain.
Por ahora solo tenemos uno: product. Luego las implementaciones se van organizando
en distintos packages, como repository, service, handler y application.
*/
package main

import (
	"ejercicio1/internal/application"
	"fmt"
)

func main() {
	app := application.NewDefaultHTTP(":8080")
	err := app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("End program", err)
}
