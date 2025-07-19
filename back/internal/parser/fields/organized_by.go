package fields

import (
	"docx-converter-demo/internal/types"
	"strings"
)

func ParseOrganizedBy(lines []string, i int, output *types.Output) int {
	if output.OrganizedBy != "" {
		return i
	}
	// หา "ดำเนินการโดย"
	j := i
	for ; j < len(lines); j++ {
		if strings.Contains(lines[j], "ดำเนินการโดย") {
			break
		}
	}
	if j >= len(lines) {
		return i
	}
	// ตัด "ดำเนินการโดย" หน้า prefix
	line := strings.TrimSpace(lines[j])
	if idx := strings.Index(line, "ดำเนินการโดย"); idx != -1 {
		line = strings.TrimSpace(line[idx+len("ดำเนินการโดย"):])
	}
	j++

	// รวบรวมบรรทัดต่อเนื่องจนกว่าจะเจอหัวข้อใหม่
	var orgLines []string
	if line != "" {
		orgLines = append(orgLines, line)
	}
	for ; j < len(lines); j++ {
		s := strings.TrimSpace(lines[j])
		if s == "" {
			continue
		}
		// หยุดถ้าเจอ "ผู้รับผิดชอบหลักสูตร" หรือหัวข้อใหม่ (ตาม pattern)
		if strings.Contains(s, "ผู้รับผิดชอบหลักสูตร") {
			break
		}
		orgLines = append(orgLines, s)
	}
	output.OrganizedBy = strings.Join(orgLines, " ")
	return j
}
