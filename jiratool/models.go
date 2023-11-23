package jiratool

type Sprint struct {
	Name      string `json:"name"`
	Id        int    `json:"id"`
	StateKey  string `json:"state_key"`
	BoardName string `json:"board_name"`
	Date      string `json:"date"`
}

type IssueInfo struct {
	SprintId  int    `json:"sprint_id"`
	IssueId   int    `json:"issue_id"`
	EpicKey   string `json:"epic_key"`
	VersionId string `json:"version_id"`
}

type Version struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}
