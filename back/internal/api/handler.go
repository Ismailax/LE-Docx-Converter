package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"docx-converter-demo/internal/config"
	"docx-converter-demo/internal/parser"
	"docx-converter-demo/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var cfg = config.MustLoad()
var MaxUploadMB = cfg.MaxUploadMB
var maxUploadBytes = config.MaxUploadBytes(cfg)

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
	log.Println("[INFO] Received /convert/:id request")

	// ดึง course_id จาก path parameter
	courseID := c.Params("id")
	if courseID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "course_id is required in URL (e.g., /convert/123456)")
	}

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
		return fail(fiber.StatusUnsupportedMediaType, fmt.Sprintf("Only .docx files are allowed (got %s)", fileHeader.Filename), nil)
	}

	// เช็กขนาดไฟล์จาก header ก่อน (ถ้ามี)
	// หมายเหตุ: ไฟล์ multipart ส่วนใหญ่จะมีค่า Size ให้ใช้งาน
	if fileHeader.Size > 0 && fileHeader.Size > maxUploadBytes {
		return fail(fiber.StatusRequestEntityTooLarge,
			fmt.Sprintf("File too large (max %d MB)", MaxUploadMB), nil)
	}

	// เปิดไฟล์ที่อัปโหลด
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
	defer out.Close()

	// คัดลอกข้อมูลจากไฟล์ที่อัปโหลดไปยังไฟล์ปลายทาง โดยจำกัดขนาดด้วย io.LimitedReader
	// เพื่อป้องกันการอัปโหลดไฟล์ขนาดใหญ่เกินกำหนดจริงๆ
	// (กรณีที่ client ไม่ส่ง Content-Length หรือส่งค่าไม่ถูกต้องมา)
	limited := &io.LimitedReader{R: file, N: maxUploadBytes + 1}
	written, err := io.Copy(out, limited)
	if err != nil {
		return fail(fiber.StatusInternalServerError, "Cannot save uploaded file (copy failed)", err)
	}
	if written > maxUploadBytes {
		_ = os.Remove(inputPath)
		return fail(fiber.StatusRequestEntityTooLarge,
			fmt.Sprintf("File too large (max %d MB)", MaxUploadMB), nil)
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

	// แทรก course_id ลง JSON
	var result map[string]any
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return fail(fiber.StatusInternalServerError, "JSON unmarshal failed", err)
	}
	result["course_id"] = courseID

	// marshal ใหม่
	finalBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fail(fiber.StatusInternalServerError, "JSON marshal failed", err)
	}

	// ส่งผลลัพธ์สำเร็จ
	log.Println("[INFO] Returning JSON result")
	c.Type("json")
	if err := c.Send(finalBytes); err != nil {
		return fail(fiber.StatusInternalServerError, "Send response failed", err)
	}

	// ล้างไฟล์แบบ async หลังส่งแล้ว
	go cleanup(toClean...)
	return nil
}
