package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractTablesFromHTML(htmlPath string) []string {
	file, err := os.Open(htmlPath)
	if err != nil {
		log.Fatalf("❌ Failed to open HTML file: %v", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatalf("❌ Failed to parse HTML: %v", err)
	}

	var tables []string
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		normalizeTable(s)
		html, err := goquery.OuterHtml(s)
		if err != nil {
			log.Printf("⚠️ failed to get table outer html: %v", err)
			return
		}
		// ทำให้กระชับ
		clean := strings.ReplaceAll(html, "\n", "")
		clean = strings.TrimSpace(clean)
		tables = append(tables, clean)
	})

	return tables
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
