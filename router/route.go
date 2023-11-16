package router

import (
	"product/config"
	"product/handlers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type HandlerList struct {
	UserHandler    handlers.UserHandler
	ProductHandler handlers.ProductHandler
}

func (hl *HandlerList) InitRoute(app *fiber.App) {
	app.Post("/login", hl.UserHandler.Login)
	app.Post("/register", hl.UserHandler.Register)

	userJWTMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Cfg.JWT_SECRET_KEY)},
	})

	user := app.Group("/users")
	user.Get("", hl.UserHandler.GetAll)
	user.Put("/:id", userJWTMiddleware, hl.UserHandler.Update)

	product := app.Group("/products")
	product.Get("", hl.ProductHandler.GetAll)
	product.Post("", userJWTMiddleware, hl.ProductHandler.Create)
	product.Put("/:id", userJWTMiddleware, hl.ProductHandler.Update)
	product.Delete("/:id", userJWTMiddleware, hl.ProductHandler.Delete)
}
