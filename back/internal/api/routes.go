package api

import "github.com/gofiber/fiber/v2"

// RegisterRoutes รวมการผูกเส้นทางทั้งหมดไว้ที่เดียว
func RegisterRoutes(app *fiber.App) {
	// health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("docx-converter-demo ready!")
	})

	// main conversion endpoint
	app.Post("/convert/:id", UploadAndConvertHandler)
}
