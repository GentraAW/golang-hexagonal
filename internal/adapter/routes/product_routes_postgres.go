package routes

import (
	"go-hexagon/internal/adapter/handler/rest"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutesPostgres(app *fiber.App, productHandler *rest.ProductHandlerPostgres) {
	app.Get("/products", productHandler.ListProducts)
	app.Get("/products/:id", productHandler.GetProductByID)
	app.Post("/products", productHandler.CreateProduct)
	app.Put("/products/:id", productHandler.UpdateProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)

	
}