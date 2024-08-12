package rest

import (
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandlerPostgres struct {
	Service *service.ProductService
}

func NewProductHandlerPostgres(service *service.ProductService) *ProductHandlerPostgres {
	return &ProductHandlerPostgres{Service: service}
}

func (h *ProductHandlerPostgres) CreateProduct(c *fiber.Ctx) error {
	product := new(entity.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Service.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandlerPostgres) UpdateProduct(c *fiber.Ctx) error {
	var product entity.Product
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Service.UpdateProduct(&product); err != nil {
		 if err.Error() == "ID: tidak tersedia" {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ID: tidak tersedia"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandlerPostgres) GetProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	product, err := h.Service.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandlerPostgres) ListProducts(c *fiber.Ctx) error {
	products, err := h.Service.ListProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandlerPostgres) DeleteProduct(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
    }

    // Panggil service untuk melakukan delete
    if err := h.Service.DeleteProduct(uint(id)); err != nil {
        if err.Error() == "ID: tidak tersedia" {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ID: tidak tersedia"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}
