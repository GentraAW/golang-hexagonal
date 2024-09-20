package handler_test

import (
	"fmt"
	"go-hexagon/internal/adapter/handler/rest"
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/service"
	"io"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Inisialisasi mock
type ProductRepositoryMock struct {
	mock.Mock
}

func getResponseBody(t *testing.T, resp *http.Response) string {
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(bodyBytes)
}

func (m *ProductRepositoryMock) List() ([]entity.Product, error) {
	args := m.Called()
	return args.Get(0).([]entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Create(product *entity.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Update(product *entity.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetByID(id interface{}) (*entity.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Delete(id interface{}) error {
	args := m.Called(id)
	return args.Error(0)
}

// -------- GET --------------
func TestListProducts_Success(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Data produk palsu
	mockProducts := []entity.Product{
		{MySQLID: 1, Name: "Product A", Stock: 100},
		{MySQLID: 2, Name: "Product B", Stock: 50},
	}

	// Atur mock untuk mengembalikan daftar produk
	productRepoMock.On("List").Return(mockProducts, nil)

	// Membuat request dan response menggunakan Fiber
	app := fiber.New()
	app.Get("/products", productHandler.ListProducts)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert bahwa respons berisi produk yang diharapkan
	expectedBody := `[{"id":1,"name":"Product A","stock":100},{"id":2,"name":"Product B","stock":50}]`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode "List" dipanggil
	productRepoMock.AssertExpectations(t)
}

func TestListProducts_Empty(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Atur mock untuk mengembalikan daftar kosong
	productRepoMock.On("List").Return([]entity.Product{}, nil)

	// Membuat request dan response menggunakan Fiber
	app := fiber.New()
	app.Get("/products", productHandler.ListProducts)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert bahwa respons adalah array kosong
	expectedBody := `[]`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode "List" dipanggil
	productRepoMock.AssertExpectations(t)
}

// -------- POST ------------
func TestCreateProduct_Success(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Set expectation: Panggil metode Create dengan produk baru
	productRepoMock.On("Create", mock.AnythingOfType("*entity.Product")).Return(nil)

	// Membuat request untuk produk baru
	app := fiber.New()
	app.Post("/products", productHandler.CreateProduct)

	reqBody := `{"name": "Product A", "stock": 100}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Assert bahwa produk berhasil dibuat dengan nilai yang benar
	expectedBody := `{"id":0,"name":"Product A","stock":100}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode Create dipanggil
	productRepoMock.AssertExpectations(t)
}

func TestCreateProduct_InvalidInput(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Membuat request dengan input tidak valid (tanpa field `name`)
	app := fiber.New()
	app.Post("/products", productHandler.CreateProduct)

	reqBody := `{"name": "Product A", "stock": "AAAA"}` // Input tidak valid
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Assert pesan error dalam respons
	expectedBody := `{"error":"json: cannot unmarshal string into Go struct field Product.stock of type int"}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode Create tidak dipanggil (karena input tidak valid)
	productRepoMock.AssertNotCalled(t, "Create")
}

// ------------ Update --------------
func TestUpdateProduct_Success(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Produk yang ada di database
	existingProduct := &entity.Product{MySQLID: 1, Name: "Old Product", Stock: 50}

	// Setup mock untuk GetByID dan Update
	productRepoMock.On("GetByID", uint(1)).Return(existingProduct, nil)
	productRepoMock.On("Update", mock.AnythingOfType("*entity.Product")).Return(nil)

	// Membuat request untuk update produk
	app := fiber.New()
	app.Put("/products/:id", productHandler.UpdateProduct)

	reqBody := `{"name": "Updated Product ABC", "stock": 100}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert bahwa produk berhasil di-update
	expectedBody := `{"id":1,"name":"Updated Product ABC","stock":100}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode GetByID dan Update dipanggil
	productRepoMock.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Setup mock untuk GetByID (produk tidak ditemukan)
	productRepoMock.On("GetByID", uint(1)).Return(nil, fmt.Errorf("ID Not Found"))

	// Membuat request untuk update produk yang tidak ada
	app := fiber.New()
	app.Put("/products/:id", productHandler.UpdateProduct)

	reqBody := `{"name": "Product A", "stock": 100}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Assert pesan error
	expectedBody := `{"error":"ID Not Found"}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode GetByID dipanggil tapi Update tidak
	productRepoMock.AssertCalled(t, "GetByID", uint(1))
	productRepoMock.AssertNotCalled(t, "Update")
}

// ------------- GET BY ID ---------------
func TestGetProductByID_Success(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Produk yang ada di database
	existingProduct := &entity.Product{MySQLID: 1, Name: "Product A", Stock: 100}

	// Setup mock untuk GetByID
	productRepoMock.On("GetByID", uint(1)).Return(existingProduct, nil)

	// Membuat request untuk mengambil produk berdasarkan ID
	app := fiber.New()
	app.Get("/products/:id", productHandler.GetProductByID)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert bahwa produk berhasil dikembalikan
	expectedBody := `{"id":1,"name":"Product A","stock":100}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode GetByID dipanggil
	productRepoMock.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Setup mock untuk GetByID (produk tidak ditemukan)
	productRepoMock.On("GetByID", uint(1)).Return(nil, fmt.Errorf("ID not found"))

	// Membuat request untuk mengambil produk yang tidak ada
	app := fiber.New()
	app.Get("/products/:id", productHandler.GetProductByID)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Assert pesan error
	expectedBody := `{"error":"ID Not Found"}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode GetByID dipanggil
	productRepoMock.AssertExpectations(t)
}

// ------------- DELETE ---------------
func TestDeleteProductByID_Success(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Setup mock untuk Delete
	productRepoMock.On("Delete", uint(1)).Return(nil)

	// Membuat request untuk menghapus produk berdasarkan ID
	app := fiber.New()
	app.Delete("/products/:id", productHandler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert pesan sukses
	expectedBody := `{"message":"Product deleted successfully"}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode Delete dipanggil
	productRepoMock.AssertExpectations(t)
}

func TestDeleteProductByID_NotFound(t *testing.T) {
	// Inisialisasi mock repository dan service
	productRepoMock := new(ProductRepositoryMock)
	productService := service.NewProductService(productRepoMock)
	productHandler := rest.NewProductHandlerMySQL(productService)

	// Setup mock untuk Delete (produk tidak ditemukan)
	productRepoMock.On("Delete", uint(1)).Return(fmt.Errorf("ID not found"))

	// Membuat request untuk menghapus produk yang tidak ada
	app := fiber.New()
	app.Delete("/products/:id", productHandler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Assert pesan error
	expectedBody := `{"error":"ID Not Found"}`
	assert.JSONEq(t, expectedBody, getResponseBody(t, resp))

	// Assert bahwa metode Delete dipanggil
	productRepoMock.AssertExpectations(t)
}
