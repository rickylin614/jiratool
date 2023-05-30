package conf

type BaseConfig struct {
	UserName  string `json:"user_name"`
	Pwd       string `json:"pwd"`
	Jiraurl   string `json:"jiraurl"`
	Project   string `json:"project"`
	Label     string `json:"label"`
	Component string `json:"component"`
}
