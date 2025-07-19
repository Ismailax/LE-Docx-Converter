package fields

import (
	// "fmt"
	"docx-converter-demo/internal/types"
	"regexp"
	"strings"
)

var asciiTableLine = regexp.MustCompile(`^\s*\+[-:=]+.*\+$`)
var dashTableLine = regexp.MustCompile(`^\s*-{3,}.*$`)
var spaceTableLine = regexp.MustCompile(`^(\s*-{2,}\s*){2,}$`)

// ParseContent: รองรับ ascii-table (+), dash-table (---), space-table (---- ---- ----)
func ParseContent(lines []string, i int, output *types.Output, tableState *types.ParseTableState) int {
	if len(output.Content) > 0 {
		return i
	}

	// 1. หา header
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
		// fmt.Println(">> ไม่เจอ header 'เนื้อหาของหลักสูตร'")
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

	flushTable := func() {
		if inTable && tableState.TableIndex < len(tableState.Tables) {
			contentBlocks = append(contentBlocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inTable = false
		}
		if inDashTable && tableState.TableIndex < len(tableState.Tables) {
			contentBlocks = append(contentBlocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inDashTable = false
		}
		if inSpaceTable && tableState.TableIndex < len(tableState.Tables) {
			contentBlocks = append(contentBlocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inSpaceTable = false
		}
		if len(paragraph) > 0 {
			contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
			paragraph = []string{}
		}
	}

	for k = start + 1; k < len(lines); k++ {
		line := strings.TrimSpace(lines[k])

		// ===== stop section เมื่อเจอ header ถัดไป =====
		if (strings.HasPrefix(line, "2.4") || strings.HasPrefix(line, "4.")) &&
			(strings.Contains(line, "Course Evaluation") || strings.Contains(line, "การประเมินผลตลอดหลักสูตร") || strings.Contains(line, "การวัดและประเมินผล")) {
			flushTable()
			break
		}

		// ===== ascii-table: +...+ =====
		if !inTable && asciiTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			// fmt.Printf(">> [DEBUG] เจอ ascii table (+) ที่ line[%d], ดึง html table #%d (ParseContent)\n", k, tableState.TableIndex)
			contentBlocks = append(contentBlocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inTable = true
			continue
		}
		if inTable {
			if line == "" || (!strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "|")) {
				inTable = false
				if line == "" {
					continue
				}
			} else {
				continue
			}
		}

		// ===== dash-table: ---... =====
		if !inDashTable && dashTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
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
					// fmt.Printf(">> [DEBUG] เจอ dash table (-) ที่ line[%d], ดึง html table #%d (ParseContent)\n", k, tableState.TableIndex)
					contentBlocks = append(contentBlocks, tableState.Tables[tableState.TableIndex])
					tableState.TableIndex++
					inDashTable = false
					continue
				}
				continue
			}
			continue
		}

		// ===== space-table: ---- ---- ---- =====
		if !inSpaceTable && spaceTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inSpaceTable = true
			continue
		}
		if inSpaceTable {
			if spaceTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
				// fmt.Printf(">> [DEBUG] เจอ space-table (----) ที่ line[%d], ดึง html table #%d (ParseContent)\n", k, tableState.TableIndex)
				contentBlocks = append(contentBlocks, tableState.Tables[tableState.TableIndex])
				tableState.TableIndex++
				inSpaceTable = false
				continue
			}
			continue
		}

		// ===== paragraph: เก็บข้อความปกติ =====
		if line == "" {
			if len(paragraph) > 0 {
				contentBlocks = append(contentBlocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}

	// paragraph/table สุดท้าย
	flushTable()

	output.Content = contentBlocks

	// Debug: print content blocks
	// fmt.Println("------ DEBUG Content Result ------")
	// for idx, b := range contentBlocks {
	// 	fmt.Printf("Block %d: %.100s\n", idx+1, strings.ReplaceAll(b, "\n", "\\n"))
	// }
	// fmt.Println("---------------------------------")

	return start
}
