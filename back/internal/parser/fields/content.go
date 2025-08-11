package fields

import (
	"docx-converter-demo/internal/types"
	"regexp"
	"strings"
)

var asciiTableLine = regexp.MustCompile(`^\s*\+[-:=]+.*\+$`)
var dashTableLine = regexp.MustCompile(`^\s*-{3,}.*$`)
var spaceTableLine = regexp.MustCompile(`^(\s*-{2,}\s*){2,}$`)

// helper: ดึง html-table ถัดไปของ "content"
func takeNextContentTable(state *types.ParseTableState) (string, bool) {
	if state.ContentIdx < len(state.Tables.Content) {
		t := state.Tables.Content[state.ContentIdx]
		state.ContentIdx++
		return t, true
	}
	return "", false
}

func ParseContent(lines []string, i int, output *types.Output, tableState *types.ParseTableState) int {
	if len(output.Content) > 0 {
		return i
	}

	// หา header "เนื้อหาของหลักสูตร"
	start := -1
	for j := i; j < len(lines); j++ {
		line := strings.TrimSpace(lines[j])
		if line == "" {
			continue
		}
		if (strings.HasPrefix(line, "2.3") || strings.HasPrefix(line, "3.")) && strings.Contains(line, "เนื้อหาของหลักสูตร") {
			start = j
			break
		}
	}
	if start == -1 {
		return i
	}

	var (
		contentBlocks     []string
		paragraph         []string
		k                 int
		inTable           bool
		inDashTable       bool
		dashTableSepCount int
		inSpaceTable      bool
	)

	flush := func() {
		if inTable {
			if t, ok := takeNextContentTable(tableState); ok {
				contentBlocks = append(contentBlocks, t)
			}
			inTable = false
		}
		if inDashTable {
			if t, ok := takeNextContentTable(tableState); ok {
				contentBlocks = append(contentBlocks, t)
			}
			inDashTable = false
		}
		if inSpaceTable {
			if t, ok := takeNextContentTable(tableState); ok {
				contentBlocks = append(contentBlocks, t)
			}
			inSpaceTable = false
		}
		if len(paragraph) > 0 {
			contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
			paragraph = []string{}
		}
	}

	for k = start + 1; k < len(lines); k++ {
		line := strings.TrimSpace(lines[k])

		// เจอ header ถัดไปของ section ประเมินผล → ปิด
		if (strings.HasPrefix(line, "2.4") || strings.HasPrefix(line, "4.")) &&
			(strings.Contains(line, "Course Evaluation") || strings.Contains(line, "การประเมินผลตลอดหลักสูตร") || strings.Contains(line, "การวัดและประเมินผล")) {
			flush()
			break
		}

		// ascii-table
		if !inTable && asciiTableLine.MatchString(line) {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inTable = true
			continue
		}
		if inTable {
			if line == "" || (!strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "|")) {
				// ปิดแล้วหยิบ html table
				if t, ok := takeNextContentTable(tableState); ok {
					contentBlocks = append(contentBlocks, t)
				}
				inTable = false
				if line == "" {
					continue
				}
			} else {
				continue
			}
		}

		// dash-table
		if !inDashTable && dashTableLine.MatchString(line) {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inDashTable = true
			dashTableSepCount = 1
			continue
		}
		if inDashTable {
			if dashTableLine.MatchString(line) {
				dashTableSepCount++
				if dashTableSepCount == 3 {
					if t, ok := takeNextContentTable(tableState); ok {
						contentBlocks = append(contentBlocks, t)
					}
					inDashTable = false
					continue
				}
				continue
			}
			continue
		}

		// space-table
		if !inSpaceTable && spaceTableLine.MatchString(line) {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inSpaceTable = true
			continue
		}
		if inSpaceTable {
			if spaceTableLine.MatchString(line) {
				if t, ok := takeNextContentTable(tableState); ok {
					contentBlocks = append(contentBlocks, t)
				}
				inSpaceTable = false
				continue
			}
			continue
		}

		// paragraph
		if line == "" {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}

	flush()
	output.Content = contentBlocks
	return start
}
