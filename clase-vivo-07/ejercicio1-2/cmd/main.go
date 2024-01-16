/*
Ejercicio 1: Test de éxito
Es el momento de aplicar testing, para esto vamos a utilizar el paquete httptest.
En primera instancia probaremos los casos de éxito, es decir, casos en que las respuestas son correctas.
Como mínimo se deben tratar los siguientes casos:

	ENDOINT		  |            STATUS CODE              |            DESCRIPCIÓN

-	GET /products         | 200 Ok Lista de todos los productos | Se espera obtener todos los productos guardados.
-	GET /products/{id}    | 200 Ok Producto esperado con ese id | Obtener el producto con el id solicitado.
-	POST /products        | 201 Created Producto añadido        | Se añade un producto en la API y se devuelve el mismo en el cuerpo de la respuesta.
-	DELETE /products/{id} | 204 No Content Sin cuerpo           | Se elimina el producto con dicho id, y no es necesario retornar nada.

Ejercicio 2: Test de fallo
Ahora vamos a probar que la API responda de manera correcta cuando ocurre un error en la consulta.
Como mínimo se deben tratar los siguientes casos:

			ENDOINT		  |            STATUS CODE              |            DESCRIPCIÓN
	  - GET|PUT|PATCH|DELETE | 400 Bad Request Mensaje informando el error | El cliente envía un id erróneo, por ejemplo, un código alfanumérico.
	    /products/{id}

	  - GET|PATCH|DELETE | 404 Not Found Mensaje informando el error | Se solicita un id que no existe en los productos listados.
	    /products/{id}

	  - POST|PUT|PATCH|DELETE | 401 Unauthorized Mensaje informando | Se intenta realizar una acción sin enviar el token de acceso.
	    /products/{id}
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
