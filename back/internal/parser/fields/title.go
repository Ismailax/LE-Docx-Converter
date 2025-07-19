package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"strings"
)

func ParseTitle(lines []string, i int, output *types.Output) int {
	if output.TitleTH != "" {
		return i
	}
	// หาจุดเริ่มต้น
	j := i
	for ; j < len(lines); j++ {
		if strings.Contains(lines[j], "ชื่อหลักสูตร") {
			break
		}
	}
	if j >= len(lines) {
		return i
	}

	// ตัด "ชื่อหลักสูตร" และ prefix หน้า
	line := strings.TrimSpace(lines[j])
	if idx := strings.Index(line, "ชื่อหลักสูตร"); idx != -1 {
		line = strings.TrimSpace(line[idx+len("ชื่อหลักสูตร"):])
	}
	j++

	// รวบรวมชื่อไทยให้ครบ (รวมบรรทัดว่าง/บรรทัดต่อเนื่อง ถ้ายังไม่เจอ EN หรือ "ดำเนินการโดย")
	var thLines []string
	if line != "" {
		thLines = append(thLines, line)
	}
	for ; j < len(lines); j++ {
		s := strings.TrimSpace(lines[j])
		if s == "" {
			continue
		}
		// จบถ้าเจอ "ดำเนินการโดย"
		if strings.Contains(s, "ดำเนินการโดย") {
			break
		}
		// ถ้าเป็น EN ล้วน (ไม่มีอักษรไทย) —> ถือเป็นบรรทัดแรกของ EN
		if utils.IsLikelyEnglish(s) {
			break
		}
		thLines = append(thLines, s)
	}

	// รวบรวม EN (ถ้ามี EN จริง)
	var enLines []string
	for ; j < len(lines); j++ {
		s := strings.TrimSpace(lines[j])
		if s == "" {
			continue
		}
		if strings.Contains(s, "ดำเนินการโดย") {
			break
		}
		if utils.IsLikelyEnglish(s) {
			enLines = append(enLines, utils.TrimBracket(s))
		} else {
			break // จบ EN ถ้าเจอไทยอีก
		}
	}

	// Special case: EN in parenthesis at the end of last TH line
	if len(enLines) == 0 && len(thLines) > 0 {
		last := thLines[len(thLines)-1]
		open := strings.LastIndex(last, "(")
		close := strings.LastIndex(last, ")")
		// ต้องวงเล็บอยู่ท้ายบรรทัด และข้างในดูเหมือน EN
		if open != -1 && close == len(last)-1 {
			candidate := strings.TrimSpace(last[open+1 : close])
			if utils.IsLikelyEnglish(candidate) {
				thLines[len(thLines)-1] = strings.TrimSpace(last[:open])
				enLines = append(enLines, candidate)
			}
		}
	}

	output.TitleTH = strings.Join(thLines, " ")
	output.TitleEN = strings.Join(enLines, " ")
	return j
}
