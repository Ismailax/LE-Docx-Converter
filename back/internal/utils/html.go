package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"docx-converter-demo/internal/types"

	"github.com/PuerkitoBio/goquery"
)

func ExtractTablesFromHTML(htmlPath string) types.TablesBySection {
	file, err := os.Open(htmlPath)
	if err != nil {
		log.Fatalf("❌ Failed to open HTML file: %v", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatalf("❌ Failed to parse HTML: %v", err)
	}

	var out types.TablesBySection
	current := "other"

	// เดินเอกสารตามลำดับ DOM: อัปเดต section เมื่อเจอหัวข้อ, เก็บ table ตาม section
	doc.Find("body *").Each(func(_ int, s *goquery.Selection) {
		if sec, ok := detectSection(s); ok {
			current = sec // "content" | "evaluation"
			return
		}

		if goquery.NodeName(s) == "table" {
			// normalize + เก็บ HTML
			normalizeTable(s)
			html, err := goquery.OuterHtml(s)
			if err != nil {
				log.Printf("⚠️ failed to get table outer html: %v", err)
				return
			}
			clean := strings.TrimSpace(strings.ReplaceAll(html, "\n", ""))

			switch current {
			case "content":
				out.Content = append(out.Content, clean)
			case "evaluation":
				out.Evaluation = append(out.Evaluation, clean)
			default:
				out.Other = append(out.Other, clean)
			}
		}
	})

	return out
}

// ---- Section detection ----

func detectSection(s *goquery.Selection) (string, bool) {
	tag := strings.ToLower(goquery.NodeName(s))
	text := strings.TrimSpace(s.Text())
	lt := strings.ToLower(text)

	// ผู้สมัครมักทำหัวข้อเป็น h1..h6 หรือ <p><strong>...</strong></p>
	isHeadingTag := tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6"

	// เดาถ้าเป็น <p> ที่ขึ้นต้นด้วยเลขหัวข้อ + มี strong ก็ถือเป็นหัวข้อ
	if tag == "p" && s.Find("strong").Length() > 0 {
		strongText := strings.TrimSpace(s.Find("strong").First().Text())
		if strongText != "" {
			// normalize spaces
			norm := func(s string) string {
				s = strings.ToLower(strings.TrimSpace(s))
				s = strings.Join(strings.Fields(s), " ")
				return s
			}
			if strings.HasPrefix(norm(lt), norm(strongText)) {
				isHeadingTag = true
			}
		}
	}

	// คีย์เวิร์ดที่มักใช้ในไฟล์จริง
	contentKeys := []string{
		"โครงสร้างหรือเนื้อหาของหลักสูตร",
		"เนื้อหาของหลักสูตร",
		"รายละเอียดของหลักสูตร",
		"course content",
	}
	evalKeys := []string{
		"การวัดและประเมินผล",
		"การประเมินผลตลอดหลักสูตร",
		"course evaluation",
		"เกณฑ์การให้ลำดับขั้น",
		"ตารางแสดงสัดส่วนการประเมิน",
	}

	// ถ้าไม่ใช่หัวข้อเลย ก็ไม่ต้องเช็คคีย์เวิร์ด
	if !isHeadingTag {
		return "", false
	}

	for _, k := range contentKeys {
		if strings.Contains(lt, strings.ToLower(k)) {
			return "content", true
		}
	}
	for _, k := range evalKeys {
		if strings.Contains(lt, strings.ToLower(k)) {
			return "evaluation", true
		}
	}
	return "", false
}

// --- helpers ---

func normalizeTable(t *goquery.Selection) {
	// 1) บังคับ table border + border-collapse
	addAttrIfAbsent(t, "border", "1")
	appendStyle(t, "border-collapse: collapse;")

	// 2) ดึงความกว้างจาก colgroup
	colWidths := extractColWidths(t)

	// 3) ใส่ width ลงในแต่ละ td/th ตามลำดับคอลัมน์
	t.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		colIdx := 0
		tr.ChildrenFiltered("th,td").Each(func(_ int, cell *goquery.Selection) {
			// รองรับ colspan
			colspan := 1
			if val, exists := cell.Attr("colspan"); exists {
				if n, err := atoi(val); err == nil && n > 0 {
					colspan = n
				}
			}
			// ความกว้างคอลัมน์แรกของ cell
			if colIdx < len(colWidths) && colWidths[colIdx] != "" {
				appendStyle(cell, "width: "+colWidths[colIdx]+";")
			}
			colIdx += colspan

			// 4) ภายใน cell: แปลง <p> เป็น <br/> (unwrap)
			unwrapParagraphs(cell)
		})
	})
}

func extractColWidths(t *goquery.Selection) []string {
	var widths []string
	t.Find("colgroup col").Each(func(_ int, col *goquery.Selection) {
		style, _ := col.Attr("style")
		width := parseWidthFromStyle(style) // ดึงเฉพาะ "width: xx%"
		widths = append(widths, width)
	})
	return widths
}

func parseWidthFromStyle(style string) string {
	// style รูปแบบ "width: 31%;" หรือมีหลายคำสั่ง
	style = strings.ToLower(style)
	parts := strings.Split(style, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "width:") {
			return strings.TrimSpace(strings.TrimPrefix(p, "width:"))
		}
	}
	return ""
}

func appendStyle(sel *goquery.Selection, extra string) {
	if extra == "" {
		return
	}
	style, _ := sel.Attr("style")
	if !strings.Contains(style, extra) {
		if style != "" && !strings.HasSuffix(style, ";") {
			style += ";"
		}
		style += extra
		sel.SetAttr("style", style)
	}
}

func addAttrIfAbsent(sel *goquery.Selection, attr, value string) {
	if _, ok := sel.Attr(attr); !ok {
		sel.SetAttr(attr, value)
	}
}

func atoi(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}

func unwrapParagraphs(cell *goquery.Selection) {
	// แทน <p> ... </p> ด้วย innerHTML + <br/>
	cell.Find("p").Each(func(_ int, p *goquery.Selection) {
		html, _ := p.Html()
		// แทน p ด้วย span+br (หรือแค่ใส่เนื้อหาแล้วตามด้วย <br/>)
		p.ReplaceWithHtml(html + "<br/>")
	})
}
