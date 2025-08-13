package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"docx-converter-demo/internal/parser"
	"docx-converter-demo/internal/utils"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

// ใช้ cleanup ที่รับ paths หลายตัว
func cleanup(paths ...string) {
	for _, p := range paths {
		os.Remove(p)
	}
	log.Println("[INFO] Cleaned up temporary files:", paths)
}

func UploadAndConvertHandler(c *fiber.Ctx) error {
	log.Println("[INFO] Received /convert request")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("[ERROR] No file received: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "File is required"})
	}

	// Check .docx
	if filepath.Ext(fileHeader.Filename) != ".docx" {
		log.Printf("[ERROR] Invalid file extension: %s", fileHeader.Filename)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Only .docx files are allowed"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("[ERROR] Cannot open uploaded file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Cannot open file"})
	}
	defer file.Close()

	tmpDir := "tmp"
	os.MkdirAll(tmpDir, 0755)

	// ใช้ uuid สำหรับไฟล์
	id := uuid.New().String()
	inputPath := filepath.Join(tmpDir, fmt.Sprintf("input-%s.docx", id))
	plainOutput := filepath.Join(tmpDir, fmt.Sprintf("output-%s.txt", id))
	htmlOutput := filepath.Join(tmpDir, fmt.Sprintf("output-%s.html", id))

	// Save uploaded file
	log.Printf("[INFO] Saving uploaded file to %s", inputPath)
	out, err := os.Create(inputPath)
	if err != nil {
		log.Printf("[ERROR] Cannot save file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Cannot save file"})
	}
	_, err = io.Copy(out, file)
	out.Close()
	if err != nil {
		log.Printf("[ERROR] Cannot save file (copy): %v", err)
		cleanup(inputPath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Cannot save file"})
	}

	// Run pandoc plain
	log.Printf("[INFO] Running pandoc to generate plain text: %s", plainOutput)
	if err := utils.RunPandocDocker(inputPath, plainOutput, "plain"); err != nil {
		log.Printf("[ERROR] Pandoc plain error: %v", err)
		cleanup(inputPath, plainOutput)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "pandoc error (plain): " + err.Error()})
	}

	// Run pandoc html
	log.Printf("[INFO] Running pandoc to generate HTML: %s", htmlOutput)
	if err := utils.RunPandocDocker(inputPath, htmlOutput, "html"); err != nil {
		log.Printf("[ERROR] Pandoc HTML error: %v", err)
		cleanup(inputPath, plainOutput, htmlOutput)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "pandoc error (html): " + err.Error()})
	}

	// Parse
	log.Println("[INFO] Running parser...")
	jsonBytes, err := parser.ParseDocToJSON(plainOutput, htmlOutput)
	if err != nil {
		log.Printf("[ERROR] Parse error: %v", err)
		cleanup(inputPath, plainOutput, htmlOutput)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Parse error: " + err.Error()})
	}

	// Success
	log.Println("[INFO] Returning JSON result")
	c.Set("Content-Type", "application/json")
	c.Send(jsonBytes)

	// Cleanup async
	go cleanup(inputPath, plainOutput, htmlOutput)
	return nil
}
