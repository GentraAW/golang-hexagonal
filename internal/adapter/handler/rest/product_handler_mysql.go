package rest

import (
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandlerMySQL struct {
	Service *service.ProductService
}

func NewProductHandlerMySQL(service *service.ProductService) *ProductHandlerMySQL {
	return &ProductHandlerMySQL{Service: service}
}

func (h *ProductHandlerMySQL) CreateProduct(c *fiber.Ctx) error {
	product := new(entity.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Service.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    product.MySQLID,
		"name":  product.Name,
		"stock": product.Stock,
	})
}

func (h *ProductHandlerMySQL) UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	productID := uint(id)

	existingProduct, err := h.Service.GetProductByID(productID)
	if err != nil {
		if err.Error() == "ID Not Found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ID Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Parse the updated product details
	var product entity.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the MySQLID of the product to be updated
	product.MySQLID = existingProduct.MySQLID

	// Update product
	if err := h.Service.UpdateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    product.MySQLID,
		"name":  product.Name,
		"stock": product.Stock,
	})
}

func (h *ProductHandlerMySQL) GetProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Convert ID to uint
	productID := uint(id)

	// Call service to get product by ID
	product, err := h.Service.GetProductByID(productID)
	if err != nil {
		if err.Error() == "ID not found" || err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ID Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    product.MySQLID,
		"name":  product.Name,
		"stock": product.Stock,
	})
}

func (h *ProductHandlerMySQL) ListProducts(c *fiber.Ctx) error {
	products, err := h.Service.ListProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(products) == 0 {
		return c.Status(fiber.StatusOK).JSON([]fiber.Map{})
	}

	var productResponses []fiber.Map
	for _, product := range products {
		productResponses = append(productResponses, fiber.Map{
			"id":    product.MySQLID,
			"name":  product.Name,
			"stock": product.Stock,
		})
	}

	return c.Status(fiber.StatusOK).JSON(productResponses)
}

func (h *ProductHandlerMySQL) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	if err := h.Service.DeleteProduct(uint(id)); err != nil {
		if err.Error() == "ID not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ID Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}
