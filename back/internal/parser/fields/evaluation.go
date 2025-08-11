package fields

import (
	"docx-converter-demo/internal/types"
	"strings"
)

// helper: ดึง html-table ถัดไปของ "evaluation"
func takeNextEvalTable(state *types.ParseTableState) (string, bool) {
	if state.EvaluationIdx < len(state.Tables.Evaluation) {
		t := state.Tables.Evaluation[state.EvaluationIdx]
		state.EvaluationIdx++
		return t, true
	}
	return "", false
}

func ParseEvaluation(lines []string, i int, output *types.Output, tableState *types.ParseTableState) int {
	if len(output.Evaluation) > 0 {
		return i
	}

	// หา header "การประเมินผล/Course Evaluation/การวัดและประเมินผล"
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

	flush := func() {
		if inTable {
			if t, ok := takeNextEvalTable(tableState); ok {
				blocks = append(blocks, t)
			}
			inTable = false
		}
		if inDashTable {
			if t, ok := takeNextEvalTable(tableState); ok {
				blocks = append(blocks, t)
			}
			inDashTable = false
		}
		if inSpaceTable {
			if t, ok := takeNextEvalTable(tableState); ok {
				blocks = append(blocks, t)
			}
			inSpaceTable = false
		}
		if len(paragraph) > 0 {
			blocks = append(blocks, strings.Join(paragraph, " "))
			paragraph = []string{}
		}
	}

	for k = start + 1; k < len(lines); k++ {
		line := strings.TrimSpace(lines[k])

		// stop เมื่อเข้าสู่ 3. คำสำคัญสำหรับการสืบค้น
		if strings.HasPrefix(line, "3.") && strings.Contains(line, "คำสำคัญสำหรับการสืบค้น") {
			flush()
			break
		}

		// ascii-table
		if !inTable && asciiTableLine.MatchString(line) {
			if len(paragraph) > 0 {
				blocks = append(blocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inTable = true
			continue
		}
		if inTable {
			if line == "" || (!strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "|")) {
				if t, ok := takeNextEvalTable(tableState); ok {
					blocks = append(blocks, t)
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
					if t, ok := takeNextEvalTable(tableState); ok {
						blocks = append(blocks, t)
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
				blocks = append(blocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			inSpaceTable = true
			continue
		}
		if inSpaceTable {
			if spaceTableLine.MatchString(line) {
				if t, ok := takeNextEvalTable(tableState); ok {
					blocks = append(blocks, t)
				}
				inSpaceTable = false
				continue
			}
			continue
		}

		// paragraph
		if line == "" {
			if len(paragraph) > 0 {
				blocks = append(blocks, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
			continue
		}
		paragraph = append(paragraph, line)
	}

	flush()
	output.Evaluation = blocks
	return start
}
