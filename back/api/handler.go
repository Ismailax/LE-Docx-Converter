package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"docx-converter-demo/internal/parser"
	"docx-converter-demo/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func UploadAndConvertHandler(c *fiber.Ctx) error {
	// log เริ่มต้นรับ request
	log.Println("[INFO] Received /convert request")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("[ERROR] No file received: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("[ERROR] Cannot open uploaded file: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot open file")
	}
	defer file.Close()

	ts := time.Now().UnixNano()
	tmpDir := "tmp"
	os.MkdirAll(tmpDir, 0755)

	inputPath := filepath.Join(tmpDir, fmt.Sprintf("input-%d-%s", ts, fileHeader.Filename))
	plainOutput := filepath.Join(tmpDir, fmt.Sprintf("output-%d.txt", ts))
	htmlOutput := filepath.Join(tmpDir, fmt.Sprintf("output-%d.html", ts))

	// log การบันทึกไฟล์
	log.Printf("[INFO] Saving uploaded file to %s", inputPath)
	out, err := os.Create(inputPath)
	if err != nil {
		log.Printf("[ERROR] Cannot save file: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot save file")
	}
	_, err = io.Copy(out, file)
	out.Close()
	if err != nil {
		log.Printf("[ERROR] Cannot save file (copy): %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot save file")
	}

	// log การรัน pandoc
	log.Printf("[INFO] Running pandoc to generate plain text: %s", plainOutput)
	if err := utils.RunPandocDocker(inputPath, plainOutput, "plain"); err != nil {
		log.Printf("[ERROR] Pandoc plain error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "pandoc error (plain): "+err.Error())
	}
	log.Printf("[INFO] Running pandoc to generate HTML: %s", htmlOutput)
	if err := utils.RunPandocDocker(inputPath, htmlOutput, "html"); err != nil {
		log.Printf("[ERROR] Pandoc HTML error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "pandoc error (html): "+err.Error())
	}

	log.Println("[INFO] Running parser...")
	jsonBytes, err := parser.ParseDocToJSON(plainOutput, htmlOutput)
	if err != nil {
		log.Printf("[ERROR] Parse error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Parse error: "+err.Error())
	}
	log.Println("[INFO] Returning JSON result")

	c.Set("Content-Type", "application/json")
	c.Send(jsonBytes)

	go func() {
		os.Remove(inputPath)
		os.Remove(plainOutput)
		os.Remove(htmlOutput)
		log.Printf("[INFO] Cleaned up temporary files for %d", ts)
	}()

	return nil
}