package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

var (
	blRegex              = regexp.MustCompile(`^\s*([-–•*●▪])\s*`)
	checkboxCheckedRegex = regexp.MustCompile(`^[🗹☑]\s*`)
	checkboxAnyRegex     = regexp.MustCompile(`^[🗹☑□☐]\s*`)
)

func ParseTarget(lines []string, i int, output *types.Output) int {
	if len(output.Target) > 0 {
		return i
	}
	start := -1
	header := ""
	// หา header
	for j := i; j < len(lines); j++ {
		line := utils.CleanText(lines[j])
		if line == "" {
			continue
		}
		if strings.Contains(line, "กลุ่มเป้าหมายหลัก") {
			start = j
			header = "กลุ่มเป้าหมายหลัก"
			break
		}
		if strings.Contains(line, "กลุ่มเป้าหมาย") {
			start = j
			header = "กลุ่มเป้าหมาย"
			break
		}
	}
	if start == -1 {
		return i
	}

	var (
		targets         []string
		checkboxTargets []string
		paragraph       []string
		prevWasParen    bool
		hasCheckbox     bool
	)
	for k := start; k < len(lines); k++ {
		line := utils.CleanText(lines[k])
		if line == "" {
			if len(paragraph) > 0 {
				targets = append(targets, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			prevWasParen = false
			continue
		}
		// stop
		if strings.HasPrefix(line, "2.") && strings.Contains(line, "ข้อมูลเฉพาะของหลักสูตร") {
			break
		}
		// Checkbox (เฉพาะช่องที่ติ๊ก)
		if checkboxCheckedRegex.MatchString(line) {
			hasCheckbox = true
			text := strings.TrimSpace(checkboxCheckedRegex.ReplaceAllString(line, ""))
			if text != "" {
				checkboxTargets = append(checkboxTargets, text)
			}
			continue
		}
		// Checkbox ที่ไม่ได้ติ๊ก = ข้ามไปเลย!
		if checkboxAnyRegex.MatchString(line) && !checkboxCheckedRegex.MatchString(line) {
			continue
		}
		// Bullet: ลบ bullet ออก!
		if blRegex.MatchString(line) {
			text := strings.TrimSpace(blRegex.ReplaceAllString(line, ""))
			if text != "" {
				targets = append(targets, text)
			}
			continue
		}
		// วงเล็บขึ้นบรรทัดใหม่ -> รวมกับบรรทัดก่อนหน้า
		if strings.HasPrefix(line, "(") && len(targets) > 0 {
			targets[len(targets)-1] += " " + line
			prevWasParen = true
			continue
		}
		// ลำดับข้อ
		if matched, _ := regexp.MatchString(`^\d+\)`, line); matched {
			if len(paragraph) > 0 {
				targets = append(targets, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			paragraph = append(paragraph, line)
			prevWasParen = false
			continue
		}
		// บรรทัดแรก (header+text)
		if k == start && header != "" && strings.Contains(line, header) {
			idx := strings.Index(line, header)
			rest := strings.TrimSpace(line[idx+len(header):])
			// ลบ bullet ถ้ามี
			rest = strings.TrimSpace(blRegex.ReplaceAllString(rest, ""))
			rest = strings.TrimSpace(checkboxAnyRegex.ReplaceAllString(rest, ""))
			if rest != "" {
				paragraph = append(paragraph, rest)
			}
			prevWasParen = false
			continue
		}
		// กรณีปกติ (text ปกติ)
		if prevWasParen && len(targets) > 0 {
			targets[len(targets)-1] += " " + line
			prevWasParen = false
			continue
		}
		paragraph = append(paragraph, line)
		prevWasParen = false
	}

	if len(paragraph) > 0 {
		targets = append(targets, strings.Join(paragraph, " "))
	}
	for i := range targets {
		targets[i] = strings.TrimSpace(targets[i])
	}
	for i := range checkboxTargets {
		checkboxTargets[i] = strings.TrimSpace(checkboxTargets[i])
	}

	// ** ตรงนี้คือใจความสำคัญ **
	// ถ้ามี checkbox ที่ถูกติ๊กอย่างน้อย 1 ช่อง ให้เก็บเฉพาะที่ checkboxTargets เท่านั้น
	if hasCheckbox && len(checkboxTargets) > 0 {
		output.Target = checkboxTargets
	} else {
		output.Target = targets
	}
	return start
}
