package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(hd http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body := r.Body
		defer body.Close()

		method := r.Method
		timeRequest := time.Now().Format(time.RFC3339)
		url := r.URL.String()

		file, err := os.OpenFile("./logger.log", os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			return
		}
		defer file.Close()
		sizeBody, err := getSizeBody(body)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write([]byte("method: " + method + " | time: " + timeRequest + " | url: " + url + " | size body: " + strconv.Itoa(sizeBody) + "\n"))

		hd.ServeHTTP(w, r)
	})
}

// funcion para obtener el tamaño de la peticion
func getSizeBody(body io.ReadCloser) (int, error) {
	size, err := io.ReadAll(body)
	if err != nil {
		return 0, errors.New("error al obtener el tamaño del body")
	}
	return len(size), nil
}
