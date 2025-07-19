package fields

import (
	// "fmt"
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

var bulletRegex = regexp.MustCompile(`^\s*([-–•*●▪])\s*`)

// split โดยจับคู่ "xxx (yyy)" ไม่แยก space ข้างใน
func splitKeywords(line string) []string {
	// fmt.Println(">> [DEBUG] splitKeywords : line =", line)
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	var keywords []string

	// [1] ดึง "ชื่อไทย (en1, en2, ...)"
	re := regexp.MustCompile(`([^\(,]+)\(([^\)]+)\)`)
	matches := re.FindAllStringSubmatch(line, -1)
	rest := line
	for _, m := range matches {
		th := strings.TrimSpace(m[1])
		enAll := m[2]
		if th != "" {
			keywords = append(keywords, th)
		}
		// split comma ในวงเล็บ
		for en := range strings.SplitSeq(enAll, ",") {
			en = strings.TrimSpace(en)
			// ตัด ( และ ) ทั้งหัว-ท้าย
			en = strings.Trim(en, "() ")
			if en != "" {
				keywords = append(keywords, en)
			}
		}
		rest = strings.Replace(rest, m[0], ",", 1)
	}
	// [2] ที่เหลือ split ด้วยคอมม่า
	for k := range strings.SplitSeq(rest, ",") {
		k = strings.TrimSpace(k)
		k = strings.Trim(k, "() ")
		if k != "" {
			for _, sub := range utils.SplitThaiEnglish(k) {
				sub = strings.TrimSpace(utils.TrimBracket(sub))
				sub = strings.Trim(sub, "() ")
				if sub != "" {
					keywords = append(keywords, sub)
				}
			}
		}
	}
	// fmt.Println(">> [DEBUG]   final =", keywords)
	return keywords
}

func ParseKeywords(lines []string, i int, output *types.Output) int {
	if len(output.Keywords) > 0 {
		return i
	}
	start := -1
	for j := i; j < len(lines); j++ {
		line := utils.CleanText(lines[j])
		if line == "" {
			continue
		}
		if (strings.Contains(line, "keyword") && (strings.HasPrefix(line, "3.1") || strings.HasPrefix(line, "1."))) ||
			(strings.Contains(line, "keyword") && len(line) < 50) {
			start = j
			break
		}
	}
	if start == -1 {
		return i
	}
	var buffer []string
	for k := start + 1; k < len(lines); k++ {
		line := utils.CleanText(lines[k])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "3.2") || strings.HasPrefix(line, "2.") ||
			strings.HasPrefix(line, "4.") || strings.Contains(line, "คำอธิบายหลักสูตรอย่างย่อ") {
			break
		}
		buffer = append(buffer, line)
	}

	// fmt.Println(">> [DEBUG] buffer (all lines under keywords):", buffer)
	buffer = utils.CombineParenthesisLines(buffer)

	seen := map[string]bool{}

	// --- กรณี: bullet แบบแยกบรรทัด และไม่มีวงเล็บและไม่มีคอมม่าในแต่ละบรรทัด
	if !utils.BufferHasParenthesis(buffer) && utils.IsAllBullet(buffer, bulletRegex) {
		for _, line := range buffer {
			// raw := line
			// ลบ bullet
			line = strings.TrimSpace(bulletRegex.ReplaceAllString(line, ""))
			// fmt.Println(">> [DEBUG] Processing bullet line:", raw)
			if line == "" {
				continue
			}
			word := strings.TrimSpace(line)
			if word != "" && !seen[word] {
				// fmt.Printf(">> [DEBUG]    add: '%s'\n", word)
				output.Keywords = append(output.Keywords, word)
				seen[word] = true
			} else if word != "" {
				// fmt.Printf(">> [DEBUG]    skip duplicate: '%s'\n", word)
			}
		}
		return start
	}

	// --- กรณีเดิม ---
	if !utils.BufferHasParenthesis(buffer) {
		joined := strings.Join(buffer, " ")
		// fmt.Println(">> [DEBUG] join all lines (no parenthesis):", joined)
		for _, word := range splitKeywords(joined) {
			if word != "" && !seen[word] {
				// fmt.Printf(">> [DEBUG]    add: '%s'\n", word)
				output.Keywords = append(output.Keywords, word)
				seen[word] = true
			} else if word != "" {
				// fmt.Printf(">> [DEBUG]    skip duplicate: '%s'\n", word)
			}
		}
	} else {
		for _, line := range buffer {
			// raw := line
			line = strings.TrimSpace(bulletRegex.ReplaceAllString(line, ""))
			// fmt.Println(">> [DEBUG] Processing line:", raw)
			for _, word := range splitKeywords(line) {
				if word != "" && !seen[word] {
					// fmt.Printf(">> [DEBUG]    add: '%s'\n", word)
					output.Keywords = append(output.Keywords, word)
					seen[word] = true
				} else if word != "" {
					// fmt.Printf(">> [DEBUG]    skip duplicate: '%s'\n", word)
				}
			}
		}
	}
	// for i := range output.Keywords {
	// 	fmt.Printf(">> [DEBUG] output.Keywords[%d] = '%s'\n", i, output.Keywords[i])
	// }
	return start
}
