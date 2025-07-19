package utils

import (
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// ExtractTablesFromHTML: ดึง <table> ทั้งหมดในไฟล์ htmlPath (return []string)
// ได้ทั้ง <table>...</table> (outer html) แต่ละตัว
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
		html, err := goquery.OuterHtml(s)
		if err != nil {
			log.Printf("⚠️ failed to get table outer html: %v", err)
			return
		}
		tables = append(tables, html)
	})

	return tables
}
