package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"product/handlers"
	"product/models"
	"product/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

var productService = mocks.ProductService{}
var productHandler = handlers.NewProductHandler(&productService)

var productModel = models.Product{
	ID:          1,
	Name:        "Permen",
	Description: "permen terenak",
	Price:       1000,
	Stock:       10000,
}

var productRequest = models.ProductRequest{
	Name:        "Permen",
	Description: "permen terenak",
	Price:       1000,
	Stock:       10000,
}

func TestGetAllProduct(t *testing.T) {
	t.Run("GetAll | Success", func(t *testing.T) {
		productService.On("GetAll").Return([]models.Product{productModel}, nil).Once()

		app.Get("/products", productHandler.GetAll)

		req := httptest.NewRequest("GET", "/products", nil)

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "successfully get all products", bodyResponse.Message)
	})

	t.Run("GetAll | Success but empty", func(t *testing.T) {
		productService.On("GetAll").Return([]models.Product{}, nil).Once()

		app.Get("/products", productHandler.GetAll)

		req := httptest.NewRequest("GET", "/products", nil)

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 204, resp.StatusCode)
		assert.Equal(t, "", bodyResponse.Message)
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("Create | Success", func(t *testing.T) {
		productService.On("Create", productRequest.ConvertToProduct()).Return(productModel, nil).Once()

		app.Post("/create", productHandler.Create)

		productReq, _ := json.Marshal(productRequest)

		req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(productReq))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "successfully create product", bodyResponse.Message)
	})

	t.Run("Create | Error, bed request body", func(t *testing.T) {
		productService.On("Create", "1", productModel).Return(productModel, nil).Once()

		app.Post("/products/:id", productHandler.Create)

		req := httptest.NewRequest("POST", "/products/2", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "invalid request", bodyResponse.Message)
	})

	t.Run("Create | Error internal server error", func(t *testing.T) {
		productService.On("Create", productRequest.ConvertToProduct()).Return(models.Product{}, errors.New("error")).Once()

		app.Post("/create", productHandler.Create)

		productReq, _ := json.Marshal(productRequest)

		req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(productReq))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 500, resp.StatusCode)
		assert.Equal(t, "Upps Sorry, There is something wrong in server", bodyResponse.Message)
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Update | Success", func(t *testing.T) {
		productService.On("Update", "1", productModel).Return(productModel, nil).Once()
		productService.On("GetByCondition", "id", "1").Return(productModel, nil).Once()

		app.Put("/products/:id", productHandler.Update)

		productReq, _ := json.Marshal(productRequest)

		req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(productReq))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "successfully update product", bodyResponse.Message)
	})

	t.Run("Update | Error, bed request body", func(t *testing.T) {
		productService.On("Update", "1", productModel).Return(productModel, nil).Once()
		productService.On("GetByCondition", "id", "1").Return(productModel, nil).Once()

		app.Put("/products/:id", productHandler.Update)

		req := httptest.NewRequest("PUT", "/products/2", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "invalid request", bodyResponse.Message)
	})

	t.Run("Update | Error, bad request id param", func(t *testing.T) {
		productService.On("Update", "1", productModel).Return(models.Product{}, errors.New("error")).Once()
		productService.On("GetByCondition", "id", "2").Return(models.Product{}, errors.New("error")).Once()

		app.Put("/products/:id", productHandler.Update)

		productReq, _ := json.Marshal(productRequest)

		req := httptest.NewRequest("PUT", "/products/2", bytes.NewBuffer(productReq))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "product is not found", bodyResponse.Message)
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("Delete | Success", func(t *testing.T) {
		productService.On("Delete", productModel).Return(nil).Once()
		productService.On("GetByCondition", "id", "1").Return(productModel, nil).Once()

		app.Delete("/products/:id", productHandler.Delete)

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "successfully delete product", bodyResponse.Message)
	})

	t.Run("Delete | Error, bad request id param", func(t *testing.T) {
		productService.On("Delete", models.Product{}).Return(nil).Once()
		productService.On("GetByCondition", "id", "2").Return(models.Product{}, errors.New("error")).Once()

		app.Delete("/products/:id", productHandler.Delete)

		req := httptest.NewRequest("DELETE", "/products/2", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "product is not found", bodyResponse.Message)
	})
}
