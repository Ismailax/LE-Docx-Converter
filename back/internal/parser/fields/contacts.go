package fields

import (
	// "fmt"
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

var contactSectionPattern = regexp.MustCompile(`^\s*(\d+\.)?\s*(ผู้ประสานงาน|ติดต่อสอบถาม|ข้อมูลในการติดต่อสอบถาม|ผู้ประสานหลักสูตร)`)
var nextSectionPattern = regexp.MustCompile(
	`^\s*\d+\.\s*(เงื่อนไข|คุณสมบัติ|ความรู้พื้นฐาน|หมวดหมู่|เอกสาร|ส่วนลด|หมายเหตุ|ข้อกำหนด|ระยะเวลา|กำหนดการ|วิธีการสมัคร|อื่นๆ|รายละเอียด|รูปแบบ|การรับสมัคร|การจัดอบรม)`,
)

// ParseContacts ดึงข้อมูลผู้ประสานงานหลักสูตรทั้งหมด
func ParseContacts(lines []string, i int, output *types.Output) int {
	// loop หา section header ที่ตรงกับ contactSectionPattern
	start := -1
	for j := i; j < len(lines); j++ {
		clean := strings.TrimSpace(lines[j])
		if contactSectionPattern.MatchString(clean) {
			start = j
			// fmt.Printf(">> [DEBUG][ParseContacts] พบหัวข้อ contact ที่ [%d]: %q\n", j, clean)
			break
		}
	}
	if start == -1 {
		// fmt.Println(">> [DEBUG][ParseContacts] ไม่พบหัวข้อ contact")
		return i
	}

	// หา block ข้อมูล contact ถัดจากหัวข้อ
	var contactBlock []string
	for j := start + 1; j < len(lines); j++ {
		c := strings.TrimSpace(lines[j])
		if nextSectionPattern.MatchString(c) {
			// fmt.Printf(">> [DEBUG][ParseContacts] พบหัวข้อถัดไปที่ [%d]: %q\n", j, c)
			break
		}
		contactBlock = append(contactBlock, lines[j])
	}

	if len(contactBlock) > 0 {
		// fmt.Printf(">> [DEBUG][ParseContacts] เก็บ contact block %d lines\n", len(contactBlock))
		output.Contacts = utils.ParseContactBlock(contactBlock)
	} else {
		// fmt.Println(">> [DEBUG][ParseContacts] ไม่พบ contact block")
	}
	return start
}
