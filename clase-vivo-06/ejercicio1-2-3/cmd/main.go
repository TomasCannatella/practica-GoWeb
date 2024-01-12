/*
Ejercicio 1: Token de acceso
Vamos a definir un token de acceso en nuestra API, para eso crearemos una variable de entorno que
contenga un token secreto con el nombre de Token y el valor del mismo.

Luego, vamos a implementar un control de acceso a las acciones que modifiquen nuestros datos,
utilizaremos los métodos POST, PUT, PATCH y DELETE; para esto debemos leer el header de las
request que recibamos en estos métodos y validar que se encuentre el token con el valor que definimos.

En Postman podemos agregar contenido al header, veamos cómo: dentro de nuestra consulta debemos hacer
clic en la pestaña de Headers. Luego, veremos una lista que nos permite agregar campos a la cabecera
de nuestra consulta, el campo KEY contiene el nombre de la variable y VALUE contiene el valor de dicha variable.

Ejercicio 2: Paquete store
En lugar de trabajar sobre una lista cargada en memoria, ahora tenemos que definir un paquete storage que nos servirá
como interfaz para modificar el archivo .json de productos. Este paquete debe estar dentro de una carpeta storage.

Debe implementar funciones de lectura (permitirá traer todos los productos) y escritura (escribir todos los productos).

Este paquete se debe iniciar en el main.go de nuestra API y ser utilizado por Repository, implementando los respectivos
métodos en las interfaces que definimos en el ejercicio anterior.

Ejercicio 3: Manejo de responses
Finalmente dentro del paquete web crearemos una función para setear una response de error estándar de nuestra API.
Dentro de él, tendremos funciones para las condiciones de éxito y otras para las de fallo. Utilizar la siguiente estructura.
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
