package rest

import (
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/service"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if product.Name == "" || product.Stock == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and Stock fields are required"})
	}

	if err := h.Service.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    product.MongoID,
		"name":  product.Name,
		"stock": product.Stock,
	})
}

func (h *ProductHandlerMongo) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi apakah ID yang diberikan adalah ObjectID MongoDB
	if !primitive.IsValidObjectID(id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Konversi ID menjadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Ambil produk yang ada dari MongoDB
	existingProduct, err := h.Service.GetProductByID(objectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Parsing request body untuk produk yang akan diupdate
	var updatedProduct entity.Product
	if err := c.BodyParser(&updatedProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input format"})
	}

	// Update field yang diberikan, jika field kosong, pertahankan nilai lama
	if updatedProduct.Name != "" {
		existingProduct.Name = updatedProduct.Name
	}
	if updatedProduct.Stock != 0 {
		existingProduct.Stock = updatedProduct.Stock
	}

	// Simpan produk yang telah diupdate ke MongoDB
	if err := h.Service.UpdateProduct(existingProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Kembalikan produk yang telah diperbarui
	return c.Status(fiber.StatusOK).JSON(existingProduct)
}

func (h *ProductHandlerMongo) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var product *entity.Product
	var err error

	if primitive.IsValidObjectID(id) {
		objectID, _ := primitive.ObjectIDFromHex(id)
		product, err = h.Service.GetProductByID(objectID)
	} else {
		uintID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}
		product, err = h.Service.GetProductByID(uint(uintID)) // Service MySQL
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    product.MongoID,
		"name":  product.Name,
		"stock": product.Stock,
	})
}

func (h *ProductHandlerMongo) ListProducts(c *fiber.Ctx) error {
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
			"id":    product.MongoID,
			"name":  product.Name,
			"stock": product.Stock,
		})
	}

	return c.Status(fiber.StatusOK).JSON(productResponses)
}

func (h *ProductHandlerMongo) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var err error
	if primitive.IsValidObjectID(id) {
		objectID, _ := primitive.ObjectIDFromHex(id)
		err = h.Service.DeleteProduct(objectID)
	} else {
		uintID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}
		err = h.Service.DeleteProduct(uint(uintID))
	}

	if err != nil {
		if err.Error() == "ID: Not Found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ID: Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}
