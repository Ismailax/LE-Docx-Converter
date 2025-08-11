package types

type TablesBySection struct {
	Content    []string
	Evaluation []string
	Other      []string
}

type ParseTableState struct {
	ContentIdx    int
	EvaluationIdx int
	Tables        TablesBySection
}
