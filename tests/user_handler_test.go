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

var userService = mocks.UserService{}
var userHandler = handlers.NewUserHandler(&userService)

var userModel = models.User{
	ID:       1,
	Name:     "Aqsa",
	Email:    "aqsa@gmail.com",
	Password: "$2a$10$4zVcIOVpRSuzLYmtuSOX/uKixrIqfIaXiMkpPUXFDBtvUTfF3o8hK",
}

var userRequest = models.UserRequest{
	Name:     "Aqsa",
	Email:    "aqsa@gmail.com",
	Password: "12345678",
}

func TestGetAllUser(t *testing.T) {
	t.Run("GetAll | Success", func(t *testing.T) {
		userService.On("GetAll").Return([]models.User{userModel}, nil).Once()

		app.Get("/users", userHandler.GetAll)

		req := httptest.NewRequest("GET", "/users", nil)

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "successfully get all users", bodyResponse.Message)
	})

	t.Run("GetAll | Success but empty", func(t *testing.T) {
		userService.On("GetAll").Return([]models.User{}, nil).Once()

		app.Get("/users", userHandler.GetAll)

		req := httptest.NewRequest("GET", "/users", nil)

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 204, resp.StatusCode)
		assert.Equal(t, "", bodyResponse.Message)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Update | Success", func(t *testing.T) {
		userService.On("Update", "1", userModel).Return(userModel, nil).Once()
		userService.On("GetByCondition", "id", "1").Return(userModel, nil).Once()

		app.Put("/users/:id", userHandler.Update)

		userReq, _ := json.Marshal(userRequest)

		req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(userReq))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "successfully update user", bodyResponse.Message)
	})

	t.Run("Update | Error, bad request body", func(t *testing.T) {
		userService.On("Update", "1", userModel).Return(userModel, nil).Once()
		userService.On("GetByCondition", "id", "1").Return(userModel, nil).Once()

		app.Put("/users/:id", userHandler.Update)

		req := httptest.NewRequest("PUT", "/users/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "invalid request", bodyResponse.Message)
	})

	t.Run("Update | Error, bad request id param", func(t *testing.T) {
		userService.On("Update", "1", userModel).Return(models.User{}, errors.New("error")).Once()
		userService.On("GetByCondition", "id", "2").Return(models.User{}, errors.New("error")).Once()

		app.Put("/users/:id", userHandler.Update)

		userReq, _ := json.Marshal(userRequest)

		req := httptest.NewRequest("PUT", "/users/2", bytes.NewBuffer(userReq))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 300000)

		body, _ := io.ReadAll(resp.Body)

		bodyResponse := ResponseFormat{}

		json.Unmarshal(body, &bodyResponse)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "user is not found", bodyResponse.Message)
	})
}
