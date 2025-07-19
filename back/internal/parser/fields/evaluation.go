package fields

import (
	// "fmt"
	"docx-converter-demo/internal/types"
	"strings"
)

// ParseEvaluation: ดึง section การประเมินผล (Course Evaluation)
func ParseEvaluation(lines []string, i int, output *types.Output, tableState *types.ParseTableState) int {
	if len(output.Evaluation) > 0 {
		return i
	}

	// 1. หา header
	start := -1
	for j := i; j < len(lines); j++ {
		line := strings.TrimSpace(lines[j])
		if line == "" {
			continue
		}
		if (strings.HasPrefix(line, "2.4") || strings.HasPrefix(line, "4.")) &&
			(strings.Contains(line, "การประเมินผลตลอดหลักสูตร") || strings.Contains(line, "Course Evaluation") ||
				strings.Contains(line, "การวัดและประเมินผล")) {
			start = j
			break
		}
	}
	if start == -1 {
		// fmt.Println(">> ไม่เจอ header 'การประเมินผลตลอดหลักสูตร'")
		return i
	}

	var (
		blocks            []string
		paragraph         []string
		k                 int
		inTable           bool
		inDashTable       bool
		dashTableSepCount int
		inSpaceTable      bool
	)

	flushTable := func() {
		// สำหรับ flush ตารางค้างไว้ก่อน break หรือหลังจบ loop
		if inTable && tableState.TableIndex < len(tableState.Tables) {
			blocks = append(blocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inTable = false
		}
		if inDashTable && tableState.TableIndex < len(tableState.Tables) {
			blocks = append(blocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inDashTable = false
		}
		if inSpaceTable && tableState.TableIndex < len(tableState.Tables) {
			blocks = append(blocks, tableState.Tables[tableState.TableIndex])
			tableState.TableIndex++
			inSpaceTable = false
		}
		if len(paragraph) > 0 {
			blocks = append(blocks, strings.Join(paragraph, " "))
			paragraph = []string{}
		}
	}

	for k = start + 1; k < len(lines); k++ {
		line := strings.TrimSpace(lines[k])

		// ===== stop section เมื่อเจอ header 3.x คำสำคัญสำหรับการสืบค้น =====
		if strings.HasPrefix(line, "3.") && strings.Contains(line, "คำสำคัญสำหรับการสืบค้น") {
			flushTable()
			break
		}

		// ===== ascii-table: +...+ =====
		if !inTable && asciiTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
			if len(paragraph) > 0 {
				blocks = append(blocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			// fmt.Printf(">> [DEBUG] เจอ ascii table (+) ที่ line[%d], ดึง html table #%d (ParseEvaluation)\n", k, tableState.TableIndex)
			blocks = append(blocks, tableState.Tables[tableState.TableIndex])
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
				blocks = append(blocks, strings.Join(paragraph, " "))
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
					// fmt.Printf(">> [DEBUG] เจอ dash table (-) ที่ line[%d], ดึง html table #%d (ParseEvaluation)\n", k, tableState.TableIndex)
					blocks = append(blocks, tableState.Tables[tableState.TableIndex])
					tableState.TableIndex++
					inDashTable = false
					continue
				}
				continue
			}
			continue
		}

		// ===== space-table: ----------------- ----------------- ... =====
		if !inSpaceTable && spaceTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
			if len(paragraph) > 0 {
				blocks = append(blocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inSpaceTable = true
			continue
		}
		if inSpaceTable {
			// ถ้าเจอขีด space-table อีกที = ปิด
			if spaceTableLine.MatchString(line) && tableState.TableIndex < len(tableState.Tables) {
				// fmt.Printf(">> [DEBUG] เจอ space-table (----) ที่ line[%d], ดึง html table #%d (ParseEvaluation)\n", k, tableState.TableIndex)
				blocks = append(blocks, tableState.Tables[tableState.TableIndex])
				tableState.TableIndex++
				inSpaceTable = false
				continue
			}
			continue
		}

		// ===== paragraph: เก็บข้อความปกติ =====
		if line == "" {
			if len(paragraph) > 0 {
				blocks = append(blocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}

	// === หลังจบ loop ให้ flush ตารางหรือ paragraph ที่ค้างไว้ ===
	flushTable()

	output.Evaluation = blocks

	// Debug
	// fmt.Println("------ DEBUG Evaluation Result ------")
	// for idx, b := range blocks {
	// 	fmt.Printf("Eval Block %d: %.100s\n", idx+1, strings.ReplaceAll(b, "\n", "\\n"))
	// }
	// fmt.Println("-------------------------------------")

	return start
}
