package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"strings"
)

// ParseObjective: แปลง objective เป็น array (แต่ละบรรทัด/paragraph)
func ParseObjective(lines []string, i int, output *types.Output) int {
	if len(output.Objective) > 0 {
		return i
	}

	start := -1
	// หา header
	for j := i; j < len(lines); j++ {
		line := utils.CleanText(lines[j])
		if line == "" {
			continue
		}
		// เงื่อนไขที่ยืดหยุ่น: ต้องมี "วัตถุประสงค์" และขึ้นต้น 2.2 หรือ 2. หรืออื่น ๆ (กัน format)
		if (strings.HasPrefix(line, "2.2") || strings.HasPrefix(line, "2.")) && strings.Contains(line, "วัตถุประสงค์") {
			start = j
			break
		}
	}
	if start == -1 {
		return i
	}

	var paragraph []string
	var objectiveLines []string
	for k := start + 1; k < len(lines); k++ {
		line := utils.CleanText(lines[k])
		// หยุดถ้าเจอหัวข้อใหม่ (2.3 หรือ 3. หรือ 3.1 ฯลฯ)
		if (strings.HasPrefix(line, "2.3") || strings.HasPrefix(line, "3.")) && strings.Contains(line, "เนื้อหาของหลักสูตร") {
			break
		}
		// แยก paragraph ด้วยบรรทัดว่าง
		if line == "" {
			if len(paragraph) > 0 {
				objectiveLines = append(objectiveLines, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}
	if len(paragraph) > 0 {
		objectiveLines = append(objectiveLines, strings.Join(paragraph, " "))
	}
	output.Objective = objectiveLines
	return start
}
