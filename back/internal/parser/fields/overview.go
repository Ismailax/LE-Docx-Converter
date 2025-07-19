package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"strings"
)

// ParseOverview: รองรับ paragraph array
func ParseOverview(lines []string, i int, output *types.Output) int {
	if len(output.Overview) > 0 {
		return i
	}

	start := -1
	// 1. หา header ที่เป็นหัวข้อ overview จริง
	for j := i; j < len(lines); j++ {
		line := utils.CleanText(lines[j])
		if line == "" {
			continue
		}
		// ข้าม header รวม
		if strings.HasPrefix(line, "3.") && strings.Contains(line, "keyword") && strings.Contains(line, "คำอธิบายหลักสูตรอย่างย่อ") {
			continue
		}
		if (strings.Contains(line, "คำอธิบายหลักสูตรอย่างย่อ") && (strings.HasPrefix(line, "3.2") || strings.HasPrefix(line, "2."))) ||
			(strings.Contains(line, "คำอธิบายหลักสูตรอย่างย่อ") && len(line) < 50) {
			start = j
			break
		}
	}
	if start == -1 {
		return i
	}

	// เก็บ paragraph ต่อเนื่อง (คั่นย่อหน้าด้วยบรรทัดว่าง)
	var paragraph []string
	var overviewLines []string
	for k := start + 1; k < len(lines); k++ {
		line := utils.CleanText(lines[k])
		if strings.HasPrefix(line, "4.") || strings.HasPrefix(line, "5.") ||
			strings.Contains(line, "ช่วงวัน-เวลา") || strings.Contains(line, "ประเภทของหลักสูตร") {
			break
		}
		// ถ้าเจอบรรทัดว่าง = จบ paragraph เดิม
		if line == "" {
			if len(paragraph) > 0 {
				overviewLines = append(overviewLines, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}
	// กรณีย่อหน้าสุดท้าย (ไม่ได้จบบรรทัดว่าง)
	if len(paragraph) > 0 {
		overviewLines = append(overviewLines, strings.Join(paragraph, " "))
	}

	output.Overview = overviewLines
	return start
}
