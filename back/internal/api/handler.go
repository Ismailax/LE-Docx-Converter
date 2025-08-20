package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"docx-converter-demo/internal/parser"
	"docx-converter-demo/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ล้างไฟล์ (log สรุปให้ด้วย)
func cleanup(paths ...string) {
	for _, p := range paths {
		_ = os.Remove(p)
	}
	if len(paths) > 0 {
		log.Println("[INFO] Cleaned up temporary files:", paths)
	}
}

func UploadAndConvertHandler(c *fiber.Ctx) error {
	log.Println("[INFO] Received /convert request")

	// เก็บรายการไฟล์ที่ต้องลบตอนจบ (ทั้งกรณีสำเร็จ/ล้มเหลว)
	toClean := make([]string, 0, 3)

	// helper: log error + cleanup + ส่ง error ให้ logger middleware แสดงใน ${error}
	fail := func(status int, publicMsg string, err error) error {
		if err != nil {
			log.Printf("[ERROR] %s: %v", publicMsg, err)
		} else {
			log.Printf("[ERROR] %s", publicMsg)
		}
		cleanup(toClean...)
		return fiber.NewError(status, publicMsg)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fail(fiber.StatusBadRequest, "File is required", err)
	}

	// ตรวจนามสกุล .docx
	if filepath.Ext(fileHeader.Filename) != ".docx" {
		return fail(fiber.StatusBadRequest, fmt.Sprintf("Only .docx files are allowed (got %s)", fileHeader.Filename), nil)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fail(fiber.StatusInternalServerError, "Cannot open uploaded file", err)
	}
	defer file.Close()

	tmpDir := "tmp"
	_ = os.MkdirAll(tmpDir, 0o755)

	// ตั้งชื่อไฟล์ชั่วคราว
	id := uuid.New().String()
	inputPath := filepath.Join(tmpDir, fmt.Sprintf("input-%s.docx", id))
	plainOutput := filepath.Join(tmpDir, fmt.Sprintf("output-%s.txt", id))
	htmlOutput := filepath.Join(tmpDir, fmt.Sprintf("output-%s.html", id))
	toClean = append(toClean, inputPath, plainOutput, htmlOutput)

	// บันทึกไฟล์ที่อัปโหลด
	log.Printf("[INFO] Saving uploaded file to %s", inputPath)
	out, err := os.Create(inputPath)
	if err != nil {
		return fail(fiber.StatusInternalServerError, "Cannot save uploaded file", err)
	}
	if _, err = io.Copy(out, file); err != nil {
		_ = out.Close()
		return fail(fiber.StatusInternalServerError, "Cannot save uploaded file (copy failed)", err)
	}
	if err := out.Close(); err != nil {
		return fail(fiber.StatusInternalServerError, "Cannot finalize uploaded file", err)
	}

	// รัน pandoc -> plain
	log.Printf("[INFO] Running pandoc (plain): %s", plainOutput)
	if err := utils.RunPandocDocker(inputPath, plainOutput, "plain"); err != nil {
		return fail(fiber.StatusInternalServerError, "pandoc error (plain)", err)
	}

	// รัน pandoc -> html
	log.Printf("[INFO] Running pandoc (html): %s", htmlOutput)
	if err := utils.RunPandocDocker(inputPath, htmlOutput, "html"); err != nil {
		return fail(fiber.StatusInternalServerError, "pandoc error (html)", err)
	}

	// แปลงเป็น JSON
	log.Println("[INFO] Running parser...")
	jsonBytes, err := parser.ParseDocToJSON(plainOutput, htmlOutput)
	if err != nil {
		return fail(fiber.StatusInternalServerError, "Parse error", err)
	}

	// ส่งผลลัพธ์สำเร็จ
	log.Println("[INFO] Returning JSON result")
	c.Type("json")
	if err := c.Send(jsonBytes); err != nil {
		return fail(fiber.StatusInternalServerError, "Send response failed", err)
	}

	// ล้างไฟล์แบบ async หลังส่งแล้ว
	go cleanup(toClean...)
	return nil
}
