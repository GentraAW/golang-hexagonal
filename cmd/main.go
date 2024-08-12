package main

import (
	"context"
	"flag"
	"go-hexagon/internal/adapter/database"
	"go-hexagon/internal/adapter/handler/rest"
	"go-hexagon/internal/adapter/repository"
	"go-hexagon/internal/adapter/routes"
	"go-hexagon/internal/core/service"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	sqlDB   *gorm.DB
	mongoDB *mongo.Client
)

func main() {
	dbType := flag.String("db", "postgres", "Database type: postgres or mongodb")
	flag.Parse()

	app := fiber.New()

	switch *dbType {
	case "postgres":
		setupPostgres(app)
	case "mongodb":
		setupMongo(app)
	default:
		log.Fatalf("Unknown database type: %s", *dbType)
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if sqlDB != nil {
			sqlDBConn, _ := sqlDB.DB()
			sqlDBConn.Close()
		}
		if mongoDB != nil {
			mongoDB.Disconnect(context.Background())
		}
		os.Exit(0)
	}()

	log.Fatal(app.Listen(":3000"))
}

func setupPostgres(app *fiber.App) {
	sqlDB = database.SetupDatabase()
	// No defer close here

	productRepo := repository.NewProductRepositoryPostgres(sqlDB)
	productService := service.NewProductService(productRepo)
	productHandler := rest.NewProductHandlerPostgres(productService)

	routes.ProductRoutesPostgres(app, productHandler)

	app.Get("/check-postgres", checkPostgres)
}

func setupMongo(app *fiber.App) {
	var err error
	mongoDB, err = database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// No defer close here

	db := mongoDB.Database("mydb")
	productRepo := repository.NewProductRepositoryMongo(db)
	productService := service.NewProductService(productRepo)
	productHandler := rest.NewProductHandlerMongo(productService)

	routes.ProductRoutesMongo(app, productHandler)

	app.Get("/check-mongo", checkMongo)
}

func checkPostgres(c *fiber.Ctx) error {
	sqlDBConn, err := sqlDB.DB()
	if err != nil {
		return c.Status(500).SendString("Failed to get PostgreSQL database connection")
	}

	if err := sqlDBConn.Ping(); err != nil {
		return c.Status(500).SendString("Failed to connect to PostgreSQL")
	}
	return c.SendString("Successfully connected to PostgreSQL")
}

func checkMongo(c *fiber.Ctx) error {
	if err := mongoDB.Ping(context.Background(), nil); err != nil {
		return c.Status(500).SendString("Failed to connect to MongoDB")
	}
	return c.SendString("Successfully connected to MongoDB")
}
