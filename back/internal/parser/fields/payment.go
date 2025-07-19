package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

var paymentDeadlinePattern = regexp.MustCompile(`ถึง(?:\s*วันที่)?\s*([0-9๐-๙]+ [^\d]+ \d{4})(?: เวลา ([0-9:.]+(?: น\.?)?))?`)

func ParsePayment(lines []string, i int, output *types.Output) int {
	header := "ช่วงวัน-เวลาของการชำระค่าธรรมเนียมในการอบรม"

	j := i
	for ; j < len(lines); j++ {
		if strings.Contains(lines[j], header) {
			break
		}
	}
	if j >= len(lines) {
		output.PaymentDeadline = nil
		return i
	}

	var deadlines []string
	foundAny := false

	for k := j + 1; k < len(lines); k++ {
		line := strings.TrimSpace(lines[k])
		if line == "" {
			continue
		}
		if strings.Contains(line, "ช่วงวัน-เวลาของการอบรม") || strings.HasPrefix(line, "6.") {
			break
		}
		if strings.Contains(line, "ไม่มีการเก็บค่าธรรมเนียม") || strings.Contains(line, "ไม่มีค่าธรรมเนียม") || strings.Contains(line, "ฟรี") {
			output.PaymentDeadline = nil
			return k
		}
		matches := paymentDeadlinePattern.FindStringSubmatch(line)
		if len(matches) > 0 {
			date := matches[1]
			time := ""
			if len(matches) >= 3 {
				time = matches[2]
			}
			dateTime := "วันที่ " + date
			if time != "" {
				dateTime += " เวลา " + time
			}
			deadlines = append(deadlines, utils.ParseThaiDateTime(dateTime))
			foundAny = true
		}
	}

	if foundAny && len(deadlines) > 0 {
		output.PaymentDeadline = deadlines
	} else {
		output.PaymentDeadline = nil
	}

	return j
}
