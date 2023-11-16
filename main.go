package main

import (
	"product/config"
	"product/db"
	"product/handlers"
	"product/router"
	"product/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.InitConfig()

	db := db.InitDB()

	userService := services.NewUserService(db)
	productService := services.NewProductService(db)

	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	route := router.HandlerList{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	route.InitRoute(app)

	app.Listen(":3000")
}
