package main

import (
	"fmt"
	"os"
	"rokhelper/db"
	_ "rokhelper/docs"
	"rokhelper/routers"
	"rokhelper/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Fiber Swagger Example API
// @version 1.0
// @description This is a ServerAPi for Rokhelper.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @host localhost:3000
// @BasePath /

func main() {
	utils.Load_env()
	logger := utils.Logger{}
	url := os.Getenv("DATABASE_URL")

	// Kiểm tra nếu URL rỗng
	if url == "" {
		logger.Fatal(fmt.Errorf("DATABASE_URL is not set in environment"))
		return
	}

	// Kết nối MongoDB
	Mongo, err := db.ConnectMongo(url)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to connect to MongoDB: %w", err))
		return
	}

	// Kiểm tra kết nối
	if Mongo == nil || Mongo.Client == nil {
		logger.Fatal(fmt.Errorf("MongoDB connection is nil"))
		return
	}

	// Đảm bảo MongoDB được ngắt kết nối khi kết thúc
	defer func() {
		if err := Mongo.DisconnectMongo(); err != nil {
			logger.Error(fmt.Sprintf("Failed to disconnect MongoDB: %v", err))
		}
	}()

	logger.Info("Successfully connected to MongoDB")

	// Khởi tạo và cấu hình ứng dụng Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error(fmt.Sprintf("Unhandled error: %v", err))
			return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
		},
	})

	// Cấu hình routes
	app.Get("/docs/*", swagger.HandlerDefault)
	routers.SetupRoutes(app, Mongo)

	// Khởi chạy server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}

	logger.Info(fmt.Sprintf("Server is running on port %s", port))
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal(fmt.Errorf("server failed to start: %w", err))
	}
}
