/*
Ejercicio 1: Método PUT
Añadir el método PUT a nuestra API: recordemos que crea o reemplaza un recurso en su totalidad con el contenido en la request.
Tené en cuenta validar los campos que se envían, como hiciste con el método POST.
Seguimos aplicando los cambios sobre la lista cargada en memoria.

Ejercicio 2: Método PATCH
Añadir el método PATCH a nuestra API: recordemos que aplica modificaciones parciales a un recurso.
Cuando realizamos un PATCH, solo debemos indicar los campos que vamos a modificar.

Ejercicio 3: Método DELETE
Añadir el método DELETE a nuestra API: recordemos que elimina un recurso.
Una vez eliminado el recurso, no debe ser posible acceder a él.
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
