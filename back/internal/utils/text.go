package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func CleanText(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\t", "")
}

func Atoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func IsLikelyEnglish(s string) bool {
	hasEn := false
	for _, r := range s {
		if unicode.In(r, unicode.Thai) {
			return false
		}
		if unicode.IsLetter(r) && unicode.In(r, unicode.Latin) {
			hasEn = true
		}
	}
	return hasEn
}

// ลบวงเล็บรอบ EN
func TrimBracket(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return strings.TrimSpace(s[1 : len(s)-1])
	}
	return s
}

// แยกไทย-อังกฤษ ใน 1 คำ (ถ้ามีทั้งสอง เช่น "A B ไทย English C D")
func SplitThaiEnglish(word string) []string {
	word = strings.TrimSpace(word)
	if word == "" {
		return nil
	}
	parts := []string{}
	buf := ""
	lastIsThai := false
	for i, r := range word {
		isThai := unicode.In(r, unicode.Thai)
		if i == 0 {
			lastIsThai = isThai
			buf += string(r)
			continue
		}
		if unicode.IsSpace(r) {
			buf += string(r)
			continue
		}
		if isThai != lastIsThai && buf != "" {
			parts = append(parts, strings.TrimSpace(buf))
			buf = ""
		}
		buf += string(r)
		lastIsThai = isThai
	}
	if buf != "" {
		parts = append(parts, strings.TrimSpace(buf))
	}

	// ลบ empty/space
	final := []string{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			final = append(final, p)
		}
	}
	return final
}

// === Combine multi-line parenthesis ===
func CombineParenthesisLines(lines []string) []string {
	var result []string
	var buffer string
	open := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !open && strings.Count(line, "(") > strings.Count(line, ")") {
			buffer = line
			open = true
			continue
		}
		if open {
			buffer += " " + line
			if strings.Count(buffer, "(") == strings.Count(buffer, ")") {
				result = append(result, buffer)
				buffer = ""
				open = false
			}
			continue
		}
		result = append(result, line)
	}
	if buffer != "" {
		result = append(result, buffer)
	}
	return result
}

func BufferHasParenthesis(buf []string) bool {
	for _, line := range buf {
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			return true
		}
	}
	return false
}

// ตรวจว่าทุกบรรทัดเป็น bullet ธรรมดา (ไม่มีคอมม่า/วงเล็บ)
func IsAllBullet(lines []string, bulletRegex *regexp.Regexp) bool {
	for _, line := range lines {
		if !bulletRegex.MatchString(line) {
			return false
		}
		if strings.ContainsAny(line, ",()") {
			return false
		}
	}
	return true
}
