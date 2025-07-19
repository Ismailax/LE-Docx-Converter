package utils

import (
	// "fmt"
	"docx-converter-demo/internal/types"
	"regexp"
	"strings"
	"unicode"
)

// ParseContactBlock : แปลง block ข้อความเป็น []Contact (debug ครบทุกจุด, ตรงตามแนวทาง null/"")
func ParseContactBlock(lines []string) []types.Contact {
	var contacts []types.Contact
	var current types.Contact
	var phones, emails, websites []string
	address, department := "", ""
	foundName := false
	foundDept := false

	contactHeaderPattern := regexp.MustCompile(`^(\d+\))?\s*ชื่อ-สกุล`)
	departmentKeywords := []string{"คณะ", "ภาควิชา", "สถาบัน", "หน่วย", "มหาวิทยาลัย"}

	// fmt.Printf(">> [DEBUG][ParseContactBlock] เริ่ม parse contact block (%d lines)\n", len(lines))

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// fmt.Printf("   [line %d] %q\n", i, line)
		if line == "" {
			continue
		}

		// 1. ถ้าเจอ contact header (ชื่อ-สกุล)
		if contactHeaderPattern.MatchString(line) {
			// fmt.Printf("   >> [DEBUG] พบชื่อ-สกุล: %q\n", line)
			// Flush department contact ก่อน (ถ้ามี department)
			if foundDept && department != "" {
				// fmt.Printf("   >> [DEBUG] FLUSH contact (department): dept=%q, address=%q, phones=%v, emails=%v, websites=%v\n", department, address, phones, emails, websites)
				contacts = append(contacts, types.Contact{
					Prefix:     "",
					Name:       "",
					Surname:    "",
					Position:   "",
					Department: department,
					Address:    address,
					Phones:     copyOrNil(phones),
					Email:      SelectContactEmail(emails),
					Websites:   copyOrNil(websites),
				})
				// reset
				address, department = "", ""
				phones, emails, websites = nil, nil, nil
				foundDept = false
			}
			// Flush person contact ก่อน (ถ้ามีข้อมูลคนเดิม)
			if foundName && (current.Name != "" || current.Surname != "") {
				current.Phones = copyOrNil(phones)
				current.Email = SelectContactEmail(emails)
				current.Websites = copyOrNil(websites)
				current.Address = address
				// fmt.Printf("   >> [DEBUG] FLUSH contact (person): %+v\n", current)
				contacts = append(contacts, current)
				current = types.Contact{}
				phones, emails, websites, address = nil, nil, nil, ""
			}
			foundName = true
			nameLine := contactHeaderPattern.ReplaceAllString(line, "")
			nameLine = strings.TrimSpace(nameLine)
			for nameLine == "" && i+1 < len(lines) {
				i++
				next := strings.TrimSpace(lines[i])
				if next != "" {
					nameLine = next
					break
				}
			}
			current.Prefix, current.Name, current.Surname, current.Position = SplitContactFullName(nameLine)
			// fmt.Printf("   >> [DEBUG] prefix: %q, name: %q, surname: %q\n", current.Prefix, current.Name, current.Surname)
			continue
		}

		// 2. ถ้าเจอหัวข้อ "ผู้ประสานงาน" แต่ไม่ใช่ "ชื่อ-สกุล"
		if strings.Contains(line, "ผู้ประสานงาน") && !contactHeaderPattern.MatchString(line) {
			// fmt.Printf("   >> [DEBUG] พบหัวข้อย่อย 'ผู้ประสานงาน'\n")
			if foundDept && department != "" {
				// fmt.Printf("   >> [DEBUG] FLUSH contact (department): dept=%q, address=%q, phones=%v, emails=%v, websites=%v\n", department, address, phones, emails, websites)
				contacts = append(contacts, types.Contact{
					Prefix:     "",
					Name:       "",
					Surname:    "",
					Position:   "",
					Department: department,
					Address:    address,
					Phones:     copyOrNil(phones),
					Email:      SelectContactEmail(emails),
					Websites:   copyOrNil(websites),
				})
				address, department = "", ""
				phones, emails, websites = nil, nil, nil
				foundDept = false
			}
			foundName = false
			continue
		}

		// 3. อื่นๆ (parse field)
		if strings.HasPrefix(line, "เบอร์โทร") || strings.HasPrefix(line, "เบอร์โทรศัพท์") {
			phone := strings.TrimPrefix(line, "เบอร์โทร")
			phone = strings.TrimPrefix(phone, "ศัพท์")
			phone = strings.TrimSpace(phone)
			// fmt.Printf("   >> [DEBUG] พบเบอร์โทร: %q\n", phone)
			for _, part := range regexp.MustCompile(`[ ,หรือ]+`).Split(phone, -1) {
				part = strings.TrimSpace(part)
				if part != "" {
					phones = append(phones, part)
					// fmt.Printf("   >> [DEBUG] เพิ่ม phone: %q\n", part)
				}
			}
			// fmt.Printf("   >> [DEBUG] phones = %v\n", phones)
		} else if after, ok := strings.CutPrefix(line, "อีเมล"); ok {
			emailLine := after
			emailLine = strings.TrimSpace(emailLine)
			ems := ExtractContactEmails(emailLine)
			// fmt.Printf("   >> [DEBUG] พบอีเมล: %q\n", emailLine)
			for _, e := range ems {
				if strings.HasSuffix(strings.ToLower(e), "@cmu.ac.th") {
					emails = append(emails, e)
					// fmt.Printf("   >> [DEBUG] เพิ่ม email: %q\n", e)
				}
			}
			// fmt.Printf("   >> [DEBUG] emails = %v\n", emails)
		} else if after, ok := strings.CutPrefix(line, "เว็บไซต์"); ok {
			web := after
			web = strings.TrimSpace(web)
			if web != "" {
				websites = append(websites, web)
				// fmt.Printf("   >> [DEBUG] เพิ่ม website: %v\n", web)
			}
			// รองรับเว็บไซต์มากกว่า 1 บรรทัด
			for i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), "http") {
				i++
				w := strings.TrimSpace(lines[i])
				websites = append(websites, w)
				// fmt.Printf("   >> [DEBUG] ต่อเว็บไซต์: %v\n", w)
			}
		} else if after, ok := strings.CutPrefix(line, "ที่อยู่"); ok {
			address = strings.TrimSpace(after)
			// fmt.Printf("   >> [DEBUG] พบที่อยู่: %q\n", address)
		} else if after0, ok0 := strings.CutPrefix(line, "ตำแหน่ง"); ok0 {
			current.Position = strings.TrimSpace(after0)
			// fmt.Printf("   >> [DEBUG] พบตำแหน่ง: %q\n", current.Position)
		} else if hasAnyKeyword(line, departmentKeywords) && !foundName {
			department = line
			foundDept = true
			// fmt.Printf("   >> [DEBUG] พบ department: %q\n", department)
		} else if foundName {
			// อื่นๆสำหรับคน (ถ้ามี field เพิ่มเติม)
		}
		if strings.HasPrefix(line, "http") {
			websites = append(websites, line)
			// fmt.Printf("   >> [DEBUG] เพิ่ม website (default): %v\n", line)
		}
	}

	// FLUSH ค้างอยู่
	if foundName && (current.Name != "" || current.Surname != "") {
		current.Phones = copyOrNil(phones)
		current.Email = SelectContactEmail(emails)
		current.Websites = copyOrNil(websites)
		current.Address = address
		// fmt.Printf("   >> [DEBUG] FLUSH contact (person สุดท้าย): %+v\n", current)
		contacts = append(contacts, current)
	}
	if foundDept && department != "" {
		// fmt.Printf("   >> [DEBUG] FLUSH contact (department สุดท้าย): dept=%q, address=%q, phones=%v, emails=%v, websites=%v\n", department, address, phones, emails, websites)
		contacts = append(contacts, types.Contact{
			Prefix:     "",
			Name:       "",
			Surname:    "",
			Position:   "",
			Department: department,
			Address:    address,
			Phones:     copyOrNil(phones),
			Email:      SelectContactEmail(emails),
			Websites:   copyOrNil(websites),
		})
	}
	if len(contacts) == 0 {
		// fmt.Println(">> [DEBUG][ParseContactBlock] ไม่พบ contact เลย")
		return nil
	}
	// fmt.Printf(">> [DEBUG][ParseContactBlock] สรุป contacts = %+v\n", contacts)
	return contacts
}

// --- ส่วน helper function เดิม copy ตามที่เคยส่ง ---

func RemoveParenthesis(s string) string {
	idx := strings.Index(s, "(")
	if idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

func SelectContactEmail(emails []string) string {
	for _, e := range emails {
		if strings.HasSuffix(strings.ToLower(e), "@cmu.ac.th") {
			return e
		}
	}
	if len(emails) > 0 {
		return emails[len(emails)-1]
	}
	return ""
}

func SplitContactFullName(full string) (prefix, name, surname string, position string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", "", "", ""
	}

	// 1. ถ้ามี (....) ปลาย string ให้ดึงออกมาเป็น position
	if idx := strings.LastIndex(full, "("); idx >= 0 && strings.HasSuffix(full, ")") {
		pos := strings.TrimSpace(full[idx+1 : len(full)-1])
		full = strings.TrimSpace(full[:idx])
		position = pos
		// fmt.Printf("   >> [DEBUG] แยก position ในวงเล็บ: %q | เหลือ: %q\n", position, full)
	}

	var prefixParts []string
	remain := full

	// 2. Greedy match plainPrefixes (ไม่มีจุด)
	plainPrefixes := []string{"ผู้ช่วยศาสตราจารย์", "รองศาสตราจารย์", "ศาสตราจารย์", "นาย", "นางสาว", "นาง", "คุณ"}
	for _, p := range plainPrefixes {
		if strings.HasPrefix(remain, p) {
			// จับ prefix พร้อมช่องว่างหลัง prefix (ถ้ามี)
			part := p
			remain = remain[len(p):]
			space := ""
			for len(remain) > 0 && remain[0] == ' ' {
				space += " "
				remain = remain[1:]
			}
			part += space
			prefixParts = append(prefixParts, part)
			break
		}
	}

	// 3. Greedy match prefix ย่อ (มีจุด) ต่อกันได้ เช่น "รศ.ดร. พญ."
	for {
		segment, ok := nextDotPrefix(remain)
		if !ok {
			break
		}
		part := segment
		remain = remain[len(segment):]
		space := ""
		for len(remain) > 0 && remain[0] == ' ' {
			space += " "
			remain = remain[1:]
		}
		part += space
		prefixParts = append(prefixParts, part)
	}

	prefix = strings.Join(prefixParts, "")
	prefix = strings.TrimSpace(prefix)

	// 4. แยกชื่อ-นามสกุล (ตัดช่องว่างหน้าหลัง)
	fields := strings.Fields(remain)
	if len(fields) >= 3 && fields[len(fields)-2] == "ณ" {
		name = fields[len(fields)-3]
		surname = fields[len(fields)-2] + " " + fields[len(fields)-1]
	} else if len(fields) >= 2 {
		name = fields[0]
		surname = strings.Join(fields[1:], " ")
	} else if len(fields) == 1 {
		name = fields[0]
	}
	return
}

func nextDotPrefix(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return "", false
	}
	for i := 2; i <= len(s); i++ {
		sub := s[:i]
		if sub[len(sub)-1] == '.' && isAllLetterOrDot(sub[:len(sub)-1]) {
			return sub, true
		}
	}
	return "", false
}

func isAllLetterOrDot(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func ExtractContactEmails(s string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	seps := []string{" และ ", ",", " "}
	for _, sep := range seps {
		parts := strings.Split(s, sep)
		var emails []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if emailRegex.MatchString(part) {
				emails = append(emails, part)
			}
		}
		if len(emails) > 0 {
			return emails
		}
	}
	return emailRegex.FindAllString(s, -1)
}

func ExtractPhones(s string) []string {
	re := regexp.MustCompile(`0\d[\d\-]{7,}`)
	return re.FindAllString(s, -1)
}

func copyOrNil(arr []string) []string {
	if len(arr) == 0 {
		return nil
	}
	return arr
}

func hasAnyKeyword(line string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(line, k) {
			return true
		}
	}
	return false
}
