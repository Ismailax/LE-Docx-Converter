package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"docx-converter-demo/internal/parser/fields"
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
)

// ParseDocToJSON : รับ path ของ plain text และ html, คืนค่า []byte (JSON) แทนการเขียนไฟล์
func ParseDocToJSON(plainTextPath string, htmlPath string) ([]byte, error) {
	file, err := os.Open(plainTextPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// EXTRACT TABLES FROM HTML
	tables := utils.ExtractTablesFromHTML(htmlPath)
	tableState := &types.ParseTableState{TableIndex: 0, Tables: tables}

	var output types.Output

	for i := 0; i < len(lines); i++ {
		i = fields.ParseTitle(lines, i, &output)
		i = fields.ParseOrganizedBy(lines, i, &output)
		i = fields.ParseEnrollLimit(lines, i, &output)
		i = fields.ParseTarget(lines, i, &output)
		i = fields.ParseRationale(lines, i, &output)
		i = fields.ParseObjective(lines, i, &output)
		i = fields.ParseContent(lines, i, &output, tableState)
		i = fields.ParseEvaluation(lines, i, &output, tableState)
		i = fields.ParseKeywords(lines, i, &output)
		i = fields.ParseOverview(lines, i, &output)
		i = fields.ParseEnrollPeriod(lines, i, &output)
		i = fields.ParsePayment(lines, i, &output)
		i = fields.ParseFees(lines, i, &output)
		i = fields.ParseContacts(lines, i, &output)
		i = fields.ParseCategories(lines, i, &output)
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, err
	}

	// แก้ unicode escape
	jsonStr := string(jsonData)
	jsonStr = strings.ReplaceAll(jsonStr, `\u003c`, "<")
	jsonStr = strings.ReplaceAll(jsonStr, `\u003e`, ">")
	jsonStr = strings.ReplaceAll(jsonStr, `\u0026`, "&")

	return []byte(jsonStr), nil
}
