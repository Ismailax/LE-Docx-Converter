package fields

import (
	"docx-converter-demo/internal/types"
	"regexp"
	"strings"
)

var checkedPattern = regexp.MustCompile(`🗹\s*([^\n🗹◻□]+)`)

func ParseCategories(lines []string, i int, output *types.Output) int {
	output.Categories = []string{} // reset
	var fallbackCategories []string
	foundChecked := false

	for j := i; j < len(lines); j++ {
		if strings.Contains(lines[j], "หมวดหมู่การเรียนรู้") {
			k := j + 1
			// ข้ามบรรทัดว่าง
			for k < len(lines) && strings.TrimSpace(lines[k]) == "" {
				k++
			}
			for ; k < len(lines); k++ {
				s := strings.TrimSpace(lines[k])
				if s == "" {
					continue
				}
				// ใช้ regexp หาเฉพาะช่องที่ "🗹"
				matches := checkedPattern.FindAllStringSubmatch(s, -1)
				if len(matches) > 0 {
					foundChecked = true
					for _, match := range matches {
						cat := strings.TrimSpace(match[1])
						if cat != "" {
							output.Categories = append(output.Categories, cat)
						}
					}
				} else {
					// ถ้าไม่เจอ "🗹" เก็บไว้ fallback
					fallbackCategories = append(fallbackCategories, s)
				}
			}
			// ถ้าไม่เจอ "🗹" เลย ให้ใช้ fallback
			if !foundChecked && len(fallbackCategories) > 0 {
				for _, cat := range fallbackCategories {
					output.Categories = append(output.Categories, cat)
				}
			}
			return len(lines)
		}
	}
	return i
}
