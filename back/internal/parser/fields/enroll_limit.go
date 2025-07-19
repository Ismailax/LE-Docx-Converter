package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

func ParseEnrollLimit(lines []string, i int, output *types.Output) int {
	if output.EnrollLimit != 0 {
		return i
	}

	// หา header ก่อน
	headerA := "1.4 จำนวนรับสมัคร"
	headerB := "1.4 จำนวนผู้เข้าร่วมอบรม"

	j := i
	for ; j < len(lines); j++ {
		clean := utils.CleanText(lines[j])
		if strings.Contains(clean, headerA) || strings.Contains(clean, headerB) {
			break
		}
	}
	if j >= len(lines) {
		return i
	}
	clean := utils.CleanText(lines[j])
	header := headerA
	if strings.Contains(clean, headerB) {
		header = headerB
	}
	content := strings.TrimSpace(strings.TrimPrefix(clean, header))

	// ถ้าบรรทัดเดียวกันไม่มีข้อมูล ดูบรรทัดถัดไป
	if content == "" && j+1 < len(lines) {
		j++
		content = utils.CleanText(lines[j])
	}

	// วิเคราะห์เนื้อหา
	if strings.Contains(content, "ไม่จำกัด") {
		output.EnrollLimit = 999999999
	} else {
		nums := regexp.MustCompile(`\d+`).FindAllString(content, -1)
		if len(nums) > 0 {
			output.EnrollLimit = utils.Atoi(nums[0])
		}
	}
	return j
}
