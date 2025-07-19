package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"strings"
)

// ParseRationale: แปลง rationale เป็น array (paragraph) ตามการเว้นบรรทัด
func ParseRationale(lines []string, i int, output *types.Output) int {
	if len(output.Rationale) > 0 {
		return i
	}

	start := -1
	// หา header ที่แท้จริง
	for j := i; j < len(lines); j++ {
		line := utils.CleanText(lines[j])
		if line == "" {
			continue
		}
		if (strings.HasPrefix(line, "2.1") || strings.HasPrefix(line, "1.")) && strings.Contains(line, "หลักการและเหตุผล") {
			start = j
			break
		}
	}
	if start == -1 {
		return i
	}

	// เก็บ paragraph ต่อเนื่อง (คั่นด้วยบรรทัดว่าง)
	var paragraph []string
	var rationaleLines []string
	for k := start + 1; k < len(lines); k++ {
		line := utils.CleanText(lines[k])
		// เจอหัวข้อใหม่ ให้หยุด
		if (strings.HasPrefix(line, "2.2") || strings.HasPrefix(line, "2.")) && strings.Contains(line, "วัตถุประสงค์") {
			break
		}
		if line == "" {
			if len(paragraph) > 0 {
				rationaleLines = append(rationaleLines, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}
	// ย่อหน้าสุดท้าย (กรณีไม่ได้จบบรรทัดว่าง)
	if len(paragraph) > 0 {
		rationaleLines = append(rationaleLines, strings.Join(paragraph, " "))
	}

	output.Rationale = rationaleLines
	return start
}
