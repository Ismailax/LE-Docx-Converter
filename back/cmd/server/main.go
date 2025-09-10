package main

import (
	"log"

	"docx-converter-demo/internal/api"
	"docx-converter-demo/internal/config"
	"docx-converter-demo/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.MustLoad()

	app := fiber.New(fiber.Config{
		BodyLimit: int(config.MaxUploadBytes(cfg)), // Set max body size to max upload size
	})

	// Middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(config.Origins(cfg)))

	// Routes
	api.RegisterRoutes(app)

	log.Fatal(app.Listen(config.Addr(cfg)))
}
