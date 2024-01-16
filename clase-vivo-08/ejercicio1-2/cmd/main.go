/*
Ejercicio 1: Middleware - Request
En este momento vamos a refactorizar nuestro código para trasladar la validación del token de acceso a un middleware,
así facilitaremos su aplicación y mantenimiento, en caso de ser necesario.

Podemos definir un paquete para este trabajo o bien dentro de nuestra función main.go.

Ejercicio 2: Middleware - Response
De igual manera que implementamos un middleware para las request, vamos implementar uno para las response.
La función de esta herramienta es llevar un registro de las consultas realizadas, es decir, un logger.
El paquete Chi proporciona uno por defecto, nosotros vamos a crear uno propio. Este deberá llevar registro de:

  - Verbo utilizado. GET, POST, PUT, etc.
  - Fecha y hora. Pueden utilizar el paquete time.
  - Url de la consulta. localhost:8080/products
  - Tamaño en bytes. Tamaño de la consulta.
*/
package main

import (
	"ejercicio1/internal/application"
	"fmt"
	"os"
)

func main() {
	app := application.NewDefaultHTTP(":8080", os.Getenv("token"))
	err := app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("End program", err)
}
