package rest

import (
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/service"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandlerMongo struct {
    Service *service.ProductService
}

func NewProductHandlerMongo(service *service.ProductService) *ProductHandlerMongo {
    return &ProductHandlerMongo{Service: service}
}

func (h *ProductHandlerMongo) CreateProduct(c *fiber.Ctx) error {
    product := new(entity.Product)
    if err := c.BodyParser(product); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := h.Service.CreateProduct(product); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandlerMongo) UpdateProduct(c *fiber.Ctx) error {
    product := new(entity.Product)
    if err := c.BodyParser(product); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := h.Service.UpdateProduct(product); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandlerMongo) GetProductByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
    }

    product, err := h.Service.GetProductByID(uint(id))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandlerMongo) ListProducts(c *fiber.Ctx) error {
    products, err := h.Service.ListProducts()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandlerMongo) DeleteProduct(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
    }

    if err := h.Service.DeleteProduct(uint(id)); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusNoContent).Send(nil)
}
