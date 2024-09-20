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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	sqlDB   *gorm.DB
	mongoDB *mongo.Client
)

func main() {
	dbType := flag.String("db", "mysql", "Database type: mysql or mongodb")
	flag.Parse()

	app := fiber.New()

	switch *dbType {
	case "mysql":
		setupMySQL(app)
	case "mongodb":
		setupMongo(app)
	default:
		log.Fatalf("Unknown database type: %s", *dbType)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down")
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

func setupMySQL(app *fiber.App) {
	dsn := "root:@tcp(127.0.0.1:3306)/db_store_go?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	sqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	productRepo := repository.NewProductRepositoryMySQL(sqlDB)
	productService := service.NewProductService(productRepo)
	productHandler := rest.NewProductHandlerMySQL(productService)

	routes.ProductRoutesMySQL(app, productHandler)

	app.Get("/check-mysql", checkMySQL)
}

func setupMongo(app *fiber.App) {
	var err error
	mongoDB, err = database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := mongoDB.Database("mydb")
	productRepo := repository.NewProductRepositoryMongo(db)
	productService := service.NewProductService(productRepo)
	productHandler := rest.NewProductHandlerMongo(productService)

	routes.ProductRoutesMongodb(app, productHandler)

	app.Get("/check-mongo", checkMongo)
}

func checkMySQL(c *fiber.Ctx) error {
	sqlDBConn, err := sqlDB.DB()
	if err != nil {
		return c.Status(500).SendString("Failed to get MySQL database connection")
	}

	if err := sqlDBConn.Ping(); err != nil {
		return c.Status(500).SendString("Failed to connect to MySQL")
	}
	return c.SendString("Successfully connected to MySQL")
}

func checkMongo(c *fiber.Ctx) error {
	if err := mongoDB.Ping(context.Background(), nil); err != nil {
		return c.Status(500).SendString("Failed to connect to MongoDB")
	}
	return c.SendString("Successfully connected to MongoDB")
}
