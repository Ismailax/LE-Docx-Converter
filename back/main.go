package main

import (
	"docx-converter-demo/api"
	"fmt"

	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
    app := fiber.New()

    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    frontendUrl := os.Getenv("FRONTEND_URL")

    app.Use(cors.New(cors.Config{
        AllowOrigins: frontendUrl,
        AllowMethods: "GET,POST,OPTIONS",
    }))

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("docx-converter-demo ready!")
    })

    app.Post("/convert", api.UploadAndConvertHandler)

    port := "8080"
    fmt.Println("Server listening on port " + port)
    app.Listen(":" + port)
}