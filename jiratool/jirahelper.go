package jiratool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"jiratool/conf"
	"net/http"
	"time"

	"github.com/andygrunwald/go-jira"
)

// 設置自定義請求的URL
func NewRequest(client *jira.Client, url string) ([]byte, error) {
	req, err := client.NewRequest(http.MethodGet, url, map[string]any{
		"_": time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req, nil)
	return ioutil.ReadAll(resp.Body)
}

// 取得Epic列表
func GetEpicList(client *jira.Client) ([]jira.Issue, error) {
	opt := &jira.SearchOptions{
		MaxResults: 1000,
		StartAt:    0,
	}

	searchString := "project = BE1 AND issuetype = Epic" // JQL語法 by chatGPT
	chunk, _, err := client.Issue.Search(searchString, opt)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// 取得推薦Sprint列表
func GetSprintList(client *jira.Client) ([]Sprint, error) {
	url := client.GetBaseURL().Scheme + `://` + client.GetBaseURL().Host + `/rest/greenhopper/1.0/sprint/picker`
	resp, err := NewRequest(client, url)
	if err != nil {
		return nil, err
	}

	data := make(map[string][]Sprint, 0)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if v, ok := data["suggestions"]; ok {
		return v, nil
	}

	return nil, errors.New("not found any Sprint")
}

// 取得fix version列表
func GetUnreleasedVersions(client *jira.Client) ([]Version, error) {
	url := fmt.Sprintf("%s://%s/rest/api/2/project/%s/versions", client.GetBaseURL().Scheme, client.GetBaseURL().Host, project)

	resp, err := NewRequest(client, url)
	if err != nil {
		return nil, err
	}

	var versions []Version
	err = json.Unmarshal(resp, &versions)
	if err != nil {
		return nil, err
	}

	var unreleasedVersions []Version
	for _, version := range versions {
		if !version.Released {
			unreleasedVersions = append(unreleasedVersions, version)
		}
	}

	return unreleasedVersions, nil
}

// 創建Issue
func CreateNewIssue(client *jira.Client, summary, desc, parentKey, fixverionId string) (*jira.Issue, error) {
	var parent *jira.Parent
	var issueType string = ISSUE_TYPE_TASK
	var isSub bool
	if parentKey != "" {
		parent = &jira.Parent{
			Key: parentKey,
		}
		issueType = ISSUE_TYPE_SUB_TASK
		isSub = true
	}
	fixversion := make([]*jira.FixVersion, 0)
	if fixverionId != "" {
		fixversion = []*jira.FixVersion{
			{ID: fixverionId},
		}
	}

	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: conf.GetConfig().UserName,
			},
			Type: jira.IssueType{
				Name:    issueType,
				Subtask: isSub,
			},
			Parent: parent,
			Project: jira.Project{
				Key: project,
			},
			Labels: []string{label},
			Components: []*jira.Component{
				{Name: "xunya"},
			},
			Summary:     summary,
			Description: desc,
			FixVersions: fixversion,
		},
	}

	i, _, err := client.Issue.Create(&issue)
	if err != nil {
		fmt.Println("jira create issue err:", err)
		return nil, err
	}
	return i, nil
}

func GeneratorRelatedIssue(client *jira.Client, relatedIssueKey string, issueinfo *IssueInfo) (*string, error) {
	if relatedIssueKey == "" {
		return nil, errors.New("未設定關聯單")
	}
	// 取得關聯單的資訊
	pmIssue, _, err := client.Issue.Get(relatedIssueKey, nil)
	if err != nil {
		fmt.Printf("jira get issue %+v client error: %s\n", pmIssue, err)
		return nil, err
	}
	summary, desc, err := GetPerentSummrayIfStory(client, pmIssue)
	if err != nil {
		fmt.Printf("jira get Perent issue %+v client error: %s\n", pmIssue, err)
		return nil, err
	}
	issue, err := CreateNewIssue(client, summary, desc, "", issueinfo.VersionId)
	if err != nil {
		return nil, err
	}
	AddRelated(client, issue.Key, relatedIssueKey)      // 添加關聯單
	UpdateEpic(client, issue.Key, issueinfo.EpicKey)    // 添加Epic
	UpdateSprint(client, issue.Key, issueinfo.SprintId) // 添加Sprint
	UpdateEstTime(client, issue.Key, 4)                 // 添加資深人員估時時間, 默認4

	return &issue.Key, nil
}

func GeneratorIssue(client *jira.Client, issueinfo *IssueInfo) (*string, error) {
	issue, err := CreateNewIssue(client, "Empty Content", "", "", issueinfo.VersionId)
	if err != nil {
		return nil, err
	}
	UpdateEpic(client, issue.Key, issueinfo.EpicKey)    // 添加Epic
	UpdateSprint(client, issue.Key, issueinfo.SprintId) // 添加Sprint
	UpdateEstTime(client, issue.Key, 4)                 // 添加資深人員估時時間, 默認4
	return &issue.Key, nil
}

func GeneratorSubIssue(client *jira.Client, SubIssueKey string, issueinfo *IssueInfo) (*string, error) {
	if SubIssueKey == "" {
		return nil, errors.New("未設定Parent Issue")
	}
	issue, err := CreateNewIssue(client, "SubTask", "", SubIssueKey, issueinfo.VersionId)
	if err != nil {
		return nil, err
	}
	UpdateEpic(client, issue.Key, issueinfo.EpicKey)    // 添加Epic
	UpdateSprint(client, issue.Key, issueinfo.SprintId) // 添加Sprint
	UpdateEstTime(client, issue.Key, 2)                 // 添加資深人員估時時間, 默認2
	return &issue.Key, nil
}
