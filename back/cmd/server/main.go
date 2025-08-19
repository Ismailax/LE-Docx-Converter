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

	app := fiber.New()

	// Middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(config.Origins(cfg)))

	// Routes
	api.RegisterRoutes(app)

	log.Printf("Server listening on %s (FRONTEND_URL=%s)", config.Addr(cfg), cfg.FrontendURL)
	log.Fatal(app.Listen(config.Addr(cfg)))
}
