package jiratool

// 根據不同系統可能有差異性
const (
	SprintFieldLabel = "customfield_10101"
	EpicFieldLabel   = "customfield_10102"
)

const (
	ISSUE_TYPE_TASK     = "Task"
	ISSUE_TYPE_SUB_TASK = "Sub-Task"
)

var project = ""
var label = "BEPFT_P_XUNYA"
var component = "xunya"

func SetProject(s string) {
	project = s
}

func SetLabel(s string) {
	label = s
}

func SetComponent(s string) {
	component = s
}
