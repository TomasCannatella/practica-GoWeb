package handler_test

import (
	"context"
	"ejercicio1/internal"
	"ejercicio1/internal/handler"
	"ejercicio1/internal/repository"
	"ejercicio1/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetAllProduct(t *testing.T) {

	t.Run("success 01 - get all products", func(t *testing.T) {
		// arrange

		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "Pail For Lid 1537",
				Quantity:    497,
				CodeValue:   "",
				IsPublished: false,
				Expiration:  "11/11/2021",
				Price:       505.33,
			},
		}
		rp := repository.NewProductMap(db)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)

		hdFunc := hd.GetAll()

		// act
		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"data":[{"Id":1,"Name":"Pail For Lid 1537","Quantity":497,"CodeValue":"","IsPublished":false,"Expiration":"11/11/2021","Price":505.33}],
		"message":"all products"}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})
}

func TestGetByIdProduct(t *testing.T) {
	t.Run("success 01 - found product", func(t *testing.T) {
		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "Pail For Lid 1537",
				Quantity:    497,
				CodeValue:   "",
				IsPublished: false,
				Expiration:  "11/11/2021",
				Price:       505.33,
			},
		}

		rp := repository.NewProductMap(db)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)
		hdFunc := hd.GetById

		req := httptest.NewRequest("GET", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"data":{"id":1,"name":"Pail For Lid 1537","quantity":497,"code_value":"","is_published":false,"expiration":"11/11/2021","price":505.33},
				"message":"product found"}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())

	})

	t.Run("fail 02 - invalid id", func(t *testing.T) {

		//arrange
		hd := handler.NewDefaultProduct(nil)
		hdFunc := hd.GetById

		//act
		req := httptest.NewRequest("GET", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{
			"status": "Bad Request",
			"message": "invalid id"
		}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})

	t.Run("fail 03 - product not found", func(t *testing.T) {
		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "Pail For Lid 1537",
				Quantity:    497,
				CodeValue:   "",
				IsPublished: false,
				Expiration:  "11/11/2021",
				Price:       505.33,
			},
		}
		//arrange
		rp := repository.NewProductMap(db)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)
		hdFunc := hd.GetById

		//act
		req := httptest.NewRequest("GET", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
				"status": "Not Found",
				"message": "product not found"
			}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())

	})
}

func TestCreateProduct(t *testing.T) {

	t.Run("success 01 - save product ", func(t *testing.T) {
		rp := repository.NewProductMap(make(map[int]internal.Product))
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)
		hdFunc := hd.Create()

		req := httptest.NewRequest("POST", "/products", strings.NewReader(
			`{
				"name": "Chicken",
				"quantity": 1,
				"code_value": "AF42AA9XDa",
				"is_published": true,
				"expiration": "10/11/2009",
				"price": 10.50
			}`,
		))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		expectedStatusCode := http.StatusCreated
		expectedBody := `{"data": {
			"Id": 1,
			"Name": "Chicken",
			"Quantity": 1,
			"CodeValue": "AF42AA9XDa",
			"IsPublished": true,
			"Expiration": "10/11/2009",
			"Price": 10.5
		},
		"message": "product successfully saved"}`

		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, res.Code, expectedStatusCode)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, res.Header(), expectedHandler)
	})

	t.Run("fail 02 - invalid token", func(t *testing.T) {
		hd := handler.NewDefaultProduct(nil)

		req := httptest.NewRequest("POST", "/products", nil)
		req.Header.Set("token", "invalid-token")
		res := httptest.NewRecorder()

		hd.Create()(res, req)

		expectedStatudCode := http.StatusUnauthorized
		expectedBody := `{
			"status":"Unauthorized",
			"message": "invalid token"
		}`

		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedStatudCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})

}

func TestUpdateProduct(t *testing.T) {
	t.Run("success 01 - update product", func(t *testing.T) {
		//arrange
		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "Chicken",
				Quantity:    1,
				CodeValue:   "AF42AA9XD",
				IsPublished: true,
				Expiration:  "10/11/2009",
				Price:       10.5,
			},
		}

		rp := repository.NewProductMap(db)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)

		//act
		req := httptest.NewRequest("PUT", "/products", strings.NewReader(`{
			"name": "Chicken rice",
			"quantity": 1,
			"code_value": "AF429XD",
			"is_published": true,
			"expiration": "10/11/2009",
			"price": 10.50
		}`))
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()

		hd.Update()(res, req)

		//assert
		expectedStatusCode := http.StatusOK
		expectedBody := `{
			"data": {
				"Id": 1,
				"Name": "Chicken rice",
				"Quantity": 1,
				"CodeValue": "AF429XD",
				"IsPublished": true,
				"Expiration": "10/11/2009",
				"Price": 10.5
			},
			"message": "product successfully updated"
		}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, res.Code, expectedStatusCode)
		require.JSONEq(t, res.Body.String(), expectedBody)
		require.Equal(t, res.Header(), expectedHandler)

	})

	t.Run("fail 02 - invalid id", func(t *testing.T) {

		hd := handler.NewDefaultProduct(nil)

		req := httptest.NewRequest("PUT", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1a")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		hd.Update()(res, req)

		expectedStatusCode := http.StatusBadRequest
		expectedBody := `
			{
				"status": "Bad Request",
				"message": "invalid format id"
			}
		`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, res.Code, expectedStatusCode)
		require.JSONEq(t, res.Body.String(), expectedBody)
		require.Equal(t, res.Header(), expectedHandler)
	})

	t.Run("fail 03 - invalid token", func(t *testing.T) {
		hd := handler.NewDefaultProduct(nil)

		req := httptest.NewRequest("POST", "/products", nil)
		req.Header.Set("token", "invalid-token")
		res := httptest.NewRecorder()

		hd.Update()(res, req)

		expectedStatudCode := http.StatusUnauthorized
		expectedBody := `{
			"status":"Unauthorized",
			"message": "invalid token"
		}`

		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedStatudCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})
}

func TestUpdatePartialProduct(t *testing.T) {
	t.Run("success 01 - update product", func(t *testing.T) {
		//arrange
		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "Chicken",
				Quantity:    1,
				CodeValue:   "AF429XD",
				IsPublished: true,
				Expiration:  "10/11/2009",
				Price:       10.5,
			},
		}

		rp := repository.NewProductMap(db)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)

		//act
		req := httptest.NewRequest("PATCH", "/products", strings.NewReader(`{
			"quantity": 500
		}`))
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()

		hd.UpdatePartial()(res, req)

		//assert
		expectedStatusCode := http.StatusOK
		expectedBody := `{
			"data": {
				"Id": 1,
				"Name": "Chicken",
				"Quantity": 500,
				"CodeValue": "AF429XD",
				"IsPublished": true,
				"Expiration": "10/11/2009",
				"Price": 10.5
			},
			"message": "product successfully updated"
		}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, res.Code, expectedStatusCode)
		require.JSONEq(t, res.Body.String(), expectedBody)
		require.Equal(t, res.Header(), expectedHandler)

	})

	t.Run("fail 02 - invalid id", func(t *testing.T) {

		hd := handler.NewDefaultProduct(nil)

		req := httptest.NewRequest("PATCH", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1a")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		hd.UpdatePartial()(res, req)

		expectedStatusCode := http.StatusBadRequest
		expectedBody := `
			{
				"status": "Bad Request",
				"message": "invalid format id"
			}
		`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, res.Code, expectedStatusCode)
		require.JSONEq(t, res.Body.String(), expectedBody)
		require.Equal(t, res.Header(), expectedHandler)
	})

	t.Run("fail 03 - product not found", func(t *testing.T) {
		rp := repository.NewProductMap(make(map[int]internal.Product))
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)

		req := httptest.NewRequest("PATCH", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		hd.UpdatePartial()(res, req)

		expectedStatusCode := http.StatusNotFound
		expectedBody := `{
			"status": "Not Found",
			"message": "product not found"
		}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedStatusCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})

	t.Run("fail 04 - invalid token", func(t *testing.T) {
		hd := handler.NewDefaultProduct(nil)

		req := httptest.NewRequest("POST", "/products", nil)
		req.Header.Set("token", "invalid-token")
		res := httptest.NewRecorder()

		hd.UpdatePartial()(res, req)

		expectedStatudCode := http.StatusUnauthorized
		expectedBody := `{
			"status":"Unauthorized",
			"message": "invalid token"
		}`

		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedStatudCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})
}
func TestDeleteProduct(t *testing.T) {
	t.Run("success 01 - Delete Product", func(t *testing.T) {
		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "Chicken",
				Quantity:    1,
				CodeValue:   "AF42AA9XDa",
				IsPublished: true,
				Expiration:  "10/11/2009",
				Price:       10.5,
			},
		}

		rp := repository.NewProductMap(db)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)

		req := httptest.NewRequest("DELETE", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		hd.Delete()(res, req)

		expectedStatusCode := http.StatusNoContent
		expectedBody := ""
		expectedHandler := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		require.Equal(t, res.Code, expectedStatusCode)
		require.Equal(t, res.Body.String(), expectedBody)
		require.Equal(t, res.Header(), expectedHandler)
	})

	t.Run("fail 02 - invalid id", func(t *testing.T) {
		//arrange
		hd := handler.NewDefaultProduct(nil)
		hdFunc := hd.Delete()

		//act
		req := httptest.NewRequest("DELETE", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{
			"status": "Bad Request",
			"message": "invalid format id"
		}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})

	t.Run("fail 03 - product not found", func(t *testing.T) {
		//arrange
		rp := repository.NewProductMap(make(map[int]internal.Product))
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProduct(sv)
		hdFunc := hd.Delete()

		//act
		req := httptest.NewRequest("GET", "/products", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
				"status": "Not Found",
				"message": "product not found"
			}`
		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())

	})

	t.Run("fail 04 - invalid token", func(t *testing.T) {
		hd := handler.NewDefaultProduct(nil)

		req := httptest.NewRequest("POST", "/products", nil)
		req.Header.Set("token", "invalid-token")
		res := httptest.NewRecorder()

		hd.Delete()(res, req)

		expectedStatudCode := http.StatusUnauthorized
		expectedBody := `{
			"status":"Unauthorized",
			"message": "invalid token"
		}`

		expectedHandler := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedStatudCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHandler, res.Header())
	})
}
