package utils

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

// ThaiMonthMap ใช้แปลงเดือนภาษาไทยเป็นตัวเลข
var ThaiMonthMap = map[string]string{
	"มกราคม":     "01",
	"กุมภาพันธ์": "02",
	"มีนาคม":     "03",
	"เมษายน":     "04",
	"พฤษภาคม":    "05",
	"มิถุนายน":   "06",
	"กรกฎาคม":    "07",
	"สิงหาคม":    "08",
	"กันยายน":    "09",
	"ตุลาคม":     "10",
	"พฤศจิกายน":  "11",
	"ธันวาคม":    "12",
}

// ThaiYearToAD แปลงปี พ.ศ. → ค.ศ.
func ThaiYearToAD(year string) string {
	y, _ := strconv.Atoi(year)
	return strconv.Itoa(y - 543)
}

// ParseThaiDateTime แปลงข้อความภาษาไทยเป็น time string
func ParseThaiDateTime(input string) string {
	// ดึง timezone
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("❌ Failed to load location: %v", err)
	}

	// ตรวจจับวันที่และเวลา
	dateRe := regexp.MustCompile(`วันที่\s*(\d{1,2})\s*([ก-๙]+)\s*(\d{4})`)
	timeRe := regexp.MustCompile(`เวลา\s*(\d{1,2})[.:](\d{2})`)

	dateMatch := dateRe.FindStringSubmatch(input)
	timeMatch := timeRe.FindStringSubmatch(input)

	if len(dateMatch) != 4 {
		return ""
	}

	day := dateMatch[1]
	if len(day) == 1 {
		day = "0" + day
	}
	monthName := dateMatch[2]
	month, ok := ThaiMonthMap[monthName]
	if !ok {
		return ""
	}
	year := ThaiYearToAD(dateMatch[3])

	hour, minute := "00", "00"
	if len(timeMatch) == 3 {
		hour, minute = timeMatch[1], timeMatch[2]
	}

	datetimeStr := fmt.Sprintf("%s-%s-%s %s:%s:00", year, month, day, hour, minute)
	t, err := time.ParseInLocation("2006-01-02 15:04:05", datetimeStr, loc)
	if err != nil {
		log.Printf("⚠️ Cannot parse date: %s", datetimeStr)
		return ""
	}

	return t.Format("2006-01-02 15:04:05")
}
