package handlers

import (
	"product/models"
	"product/services"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) ProductHandler {
	return ProductHandler{
		productService,
	}
}

func (ph *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := ph.productService.GetAll()

	if err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	if len(products) == 0 {
		return response(c, fiber.StatusNoContent, "", nil)
	}

	var productsResponse []models.ProductResponse
	for _, product := range products {
		productsResponse = append(productsResponse, product.ConvertToResponse())
	}

	return response(c, fiber.StatusOK, "successfully get all products", productsResponse)
}

func (ph *ProductHandler) Create(c *fiber.Ctx) error {
	productRequest := models.ProductRequest{}
	c.BodyParser(&productRequest)

	if err := productRequest.Validate(); err != nil {
		return response(c, fiber.StatusBadRequest, "invalid request", nil)
	}

	product, err := ph.productService.Create(productRequest.ConvertToProduct())

	if product.ID == 0 || err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	return response(c, fiber.StatusOK, "successfully create product", product.ConvertToResponse())
}

func (ph *ProductHandler) Update(c *fiber.Ctx) error {
	productRequest := models.ProductRequest{}
	id := c.Params("id")

	c.BodyParser(&productRequest)

	if err := productRequest.Validate(); err != nil {
		return response(c, fiber.StatusBadRequest, "invalid request", nil)
	}

	product, err := ph.productService.GetByCondition("id", id)

	if err != nil {
		return response(c, fiber.StatusBadRequest, "product is not found", nil)
	}

	product.Name = productRequest.Name
	product.Description = productRequest.Description
	product.Price = productRequest.Price
	product.Stock = productRequest.Stock

	product, err = ph.productService.Update(id, product)

	if err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	return response(c, fiber.StatusOK, "successfully update product", product.ConvertToResponse())
}

func (ph *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	product, _ := ph.productService.GetByCondition("id", id)

	if product.ID == 0 {
		return response(c, fiber.StatusBadRequest, "product is not found", nil)
	}

	if err := ph.productService.Delete(product); err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	return response(c, fiber.StatusOK, "successfully delete product", nil)
}
