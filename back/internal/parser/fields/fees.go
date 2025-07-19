package fields

import (
	// "fmt"
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

// Extract all fee amounts in a line, return slice of int
func extractAllFees(line string) []int {
	// fmt.Printf(">> [DEBUG] extractAllFees: line = %s\n", line)
	var result []int
	re := regexp.MustCompile(`([\d,]+)\s*บาท`)
	matches := re.FindAllStringSubmatch(line, -1)
	// fmt.Printf(">> [DEBUG]    matches = %+v\n", matches)
	for _, m := range matches {
		if len(m) >= 2 {
			fee := strings.ReplaceAll(m[1], ",", "")
			n := utils.Atoi(fee)
			if n > 0 {
				// fmt.Printf(">> [DEBUG]    found fee: %d\n", n)
				result = append(result, n)
			}
		}
	}
	return result
}

// ParseFees ดึงค่าธรรมเนียมและค่าบำรุงมหาวิทยาลัย (array + single)
func ParseFees(lines []string, i int, output *types.Output) int {
	fees := []int{}
	univFees := []int{}
	var feeSectionPattern = regexp.MustCompile(`^\s*(\d+\.)?\s*ค่าธรรมเนียม`)
	var nextSectionPattern = regexp.MustCompile(`^\s*(\d+\.)?\s*(แหล่งที่มาของงบประมาณ|ส่วนลด|หมายเหตุ)`)

	// หา section ที่เป็นหัวข้อค่าธรรมเนียม
	start := -1
	for j := i; j < len(lines); j++ {
		line := utils.CleanText(lines[j])
		if feeSectionPattern.MatchString(line) {
			start = j
			break
		}
	}
	if start == -1 {
		output.Fees = []int{0}
		output.UniversityFee = 0
		// fmt.Println(">> [DEBUG] ไม่พบ section ค่าธรรมเนียม")
		return i
	}

	// ======= [NEW] ดึงค่าธรรมเนียมและค่าบำรุงจากหัวข้อ =======
	startLine := utils.CleanText(lines[start])
	// fmt.Printf(">> [DEBUG] startLine = %s\n", startLine)
	if strings.Contains(startLine, "บาท") {
		// กรณีพิเศษ: มี "ไม่รวมค่าบำรุงมหาวิทยาลัย" ในหัวข้อ
		if strings.Contains(startLine, "ไม่รวมค่าบำรุงมหาวิทยาลัย") {
			parts := strings.Split(startLine, "(")
			if len(parts) > 0 {
				fs := extractAllFees(parts[0])
				// fmt.Printf(">> [DEBUG]    [หัวข้อ] ส่วนธรรมเนียม = %v\n", fs)
				fees = append(fees, fs...)
			}
			if len(parts) > 1 {
				us := extractAllFees(parts[1])
				// fmt.Printf(">> [DEBUG]    [หัวข้อ] ส่วนค่าบำรุง = %v\n", us)
				univFees = append(univFees, us...)
			}
		} else {
			// กรณีปกติ ดึงค่าธรรมเนียมทั้งหมด
			vals := extractAllFees(startLine)
			for _, v := range vals {
				fees = append(fees, v)
			}
			// fmt.Printf(">> [DEBUG]    [หัวข้อ] บรรทัดค่าธรรมเนียม = %v\n", vals)
		}
	}

	// fmt.Printf(">> [DEBUG] เริ่มอ่านบรรทัดค่าธรรมเนียม ตั้งแต่ index = %d\n", start+1)

	for j := start + 1; j < len(lines); j++ {
		content := utils.CleanText(lines[j])
		// fmt.Printf(">> [DEBUG] line = %s\n", content)

		// จบ section
		if nextSectionPattern.MatchString(content) {
			// fmt.Printf(">> [DEBUG] พบจบ section ค่าธรรมเนียม ที่ index = %d\n", j)
			break
		}
		if content == "" {
			continue
		}

		// ไม่มีการเก็บค่าธรรมเนียม
		if strings.Contains(content, "ไม่มีค่าธรรมเนียม") || strings.Contains(content, "ไม่มีการเก็บค่าธรรมเนียม") {
			// fmt.Printf(">> [DEBUG]    → พบว่าไม่มีค่าธรรมเนียม\n")
			fees = append(fees, 0)
			continue
		}

		// กรณีเจอ "ไม่รวมค่าบำรุงมหาวิทยาลัย"
		if strings.Contains(content, "ไม่รวมค่าบำรุงมหาวิทยาลัย") {
			parts := strings.Split(content, "(")
			if len(parts) > 0 {
				fs := extractAllFees(parts[0])
				// fmt.Printf(">> [DEBUG]    → ส่วนธรรมเนียม = %v\n", fs)
				fees = append(fees, fs...)
			}
			if len(parts) > 1 {
				us := extractAllFees(parts[1])
				// fmt.Printf(">> [DEBUG]    → ส่วนค่าบำรุง = %v\n", us)
				univFees = append(univFees, us...)
			}
			continue
		}

		// ค่าบำรุงมหาวิทยาลัย
		if strings.Contains(content, "ค่าบำรุงมหาวิทยาลัย") {
			vals := extractAllFees(content)
			if len(vals) == 0 && (strings.Contains(content, "ไม่มีค่าบำรุง") || strings.Contains(content, "ไม่มีการเก็บค่าบำรุง")) {
				// fmt.Printf(">> [DEBUG]    → ไม่มีค่าบำรุงมหาวิทยาลัย\n")
				univFees = append(univFees, 0)
			}
			for _, v := range vals {
				univFees = append(univFees, v)
			}
			// fmt.Printf(">> [DEBUG]    → บรรทัดค่าบำรุงมหาวิทยาลัย = %v\n", univFees)
			continue
		}

		// ดึงค่าธรรมเนียม (รวม bullet ด้วย)
		if strings.Contains(content, "ค่าธรรมเนียม") || strings.Contains(content, "ค่าธรรมเนียมหลักสูตร") || strings.Contains(content, "ค่าธรรมเนียมการอบรม") {
			vals := extractAllFees(content)
			for _, v := range vals {
				fees = append(fees, v)
			}
			// fmt.Printf(">> [DEBUG]    → บรรทัดค่าธรรมเนียม = %v\n", fees)
			continue
		}

		// bullet เช่น "- ราคา 500 บาท/หลักสูตรย่อย/คน"
		if strings.HasPrefix(content, "-") || strings.HasPrefix(content, "•") {
			vals := extractAllFees(content)
			if strings.Contains(content, "ค่าบำรุงมหาวิทยาลัย") {
				univFees = append(univFees, vals...)
				// fmt.Printf(">> [DEBUG]    → bullet บรรทัดค่าบำรุงมหาวิทยาลัย = %v\n", univFees)
			} else {
				fees = append(fees, vals...)
				// fmt.Printf(">> [DEBUG]    → bullet บรรทัดค่าธรรมเนียม = %v\n", fees)
			}
			continue
		}

		// ดักกรณีมีแต่เลขกับ "บาท" (สำรองสุดท้าย)
		if strings.Contains(content, "บาท") {
			vals := extractAllFees(content)
			// fmt.Printf(">> [DEBUG]    → เจอเลขกับบาท = %v\n", vals)
			fees = append(fees, vals...)
			continue
		}
	}

	// Fallback
	if len(fees) == 0 {
		// fmt.Printf(">> [DEBUG] Fallback: ไม่มีค่าธรรมเนียม พบ [0]\n")
		fees = []int{0}
	}
	// univFee ให้เอาตัวแรกที่เจอ ไม่ซ้ำซ้อน
	univFee := 0
	for _, v := range univFees {
		if v > 0 {
			univFee = v
			break
		}
	}
	// fmt.Printf(">> [DEBUG] ผลลัพธ์: Fees = %v | UniversityFee = %v\n", fees, univFee)
	output.Fees = fees
	output.UniversityFee = univFee
	return start
}
