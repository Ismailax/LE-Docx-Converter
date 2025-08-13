package main

import (
	"log"

	"docx-converter-demo/internal/api"
	"docx-converter-demo/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.MustLoad()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Origins(cfg),
		AllowMethods: "GET,POST,OPTIONS",
	}))

	api.RegisterRoutes(app)

	log.Printf("Server listening on %s (FRONTEND_URL=%s)", config.Addr(cfg), cfg.FrontendURL)
	log.Fatal(app.Listen(config.Addr(cfg)))
}
