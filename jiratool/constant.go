package jiratool

// 根據不同系統可能有差異性
const (
	SprintFieldLabel  = "customfield_10101" // Sprint標籤
	EpicFieldLabel    = "customfield_10102" // Epic標籤
	StaffEstTimeLabel = "customfield_10600" // 資深人員估時
)

const (
	ISSUE_TYPE_TASK     = "Task"
	ISSUE_TYPE_SUB_TASK = "Sub-Task"
	ISSUE_TYPE_STORY    = "Story"
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
