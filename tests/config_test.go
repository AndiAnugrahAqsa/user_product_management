package tests

import (
	"product/models"

	"github.com/gofiber/fiber/v2"
)

type ResponseFormat struct {
	Data    []models.User `json:"data"`
	Message string        `json:"message"`
}

var app = fiber.New()
