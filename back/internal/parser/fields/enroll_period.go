package fields

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"regexp"
	"strings"
)

var enrollPattern = regexp.MustCompile(`(เปิดรับสมัคร|ปิดรับสมัคร)[^\d]*(วันที่)?\s*([0-9๐-๙]+ [^\d]+ \d{4})(?: เวลา ([0-9:.]+( น\.?)?))?`)

func ParseEnrollPeriod(lines []string, i int, output *types.Output) int {
	header := "ช่วงวัน-เวลาของการรับสมัคร"

	// 1. loop หา header ให้เจอ
	j := i
	for ; j < len(lines); j++ {
		if strings.Contains(lines[j], header) {
			break
		}
	}
	if j >= len(lines) {
		// ไม่เจอหัวข้อ
		return i
	}

	// เตรียม slice สำหรับหลายรอบ
	var startList, endList []string
	var curStart, curEnd string

	// 2. loop เก็บรายละเอียดจนเจอหัวข้อถัดไป
	for k := j + 1; k < len(lines); k++ {
		line := utils.CleanText(lines[k])
		if line == "" {
			continue
		}
		// เจอหัวข้อถัดไป
		if strings.Contains(line, "ชำระค่าธรรมเนียม") || strings.HasPrefix(line, "5.") {
			break
		}
		// match "เปิดรับสมัคร ..." หรือ "ปิดรับสมัคร ..."
		matches := enrollPattern.FindStringSubmatch(line)
		if len(matches) > 0 {
			what := matches[1] // เปิดรับสมัคร / ปิดรับสมัคร
			date := matches[3]
			time := ""
			if len(matches) >= 5 {
				time = matches[4]
			}
			fullDate := "วันที่ " + date
			if time != "" {
				fullDate += " เวลา " + time
			}
			if what == "เปิดรับสมัคร" {
				curStart = utils.ParseThaiDateTime(fullDate)
			} else if what == "ปิดรับสมัคร" {
				curEnd = utils.ParseThaiDateTime(fullDate)
			}
		}

		// ถ้าเจอทั้ง curStart และ curEnd ให้เก็บรอบใหม่
		if curStart != "" && curEnd != "" {
			startList = append(startList, curStart)
			endList = append(endList, curEnd)
			curStart, curEnd = "", ""
		}
	}

	// ถ้ายังเหลือรอบค้างอยู่ (เช่นมีรอบสุดท้ายที่ไม่มีปิด)
	if curStart != "" {
		startList = append(startList, curStart)
	}
	if curEnd != "" && len(endList) < len(startList) {
		endList = append(endList, curEnd)
	}

	output.StartEnroll = startList
	output.EndEnroll = endList

	return j
}
