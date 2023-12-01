package jiratool

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func GetIssueByKey(client *jira.Client, key string) (*jira.Issue, error) {
	issue, _, err := client.Issue.Get(key, nil)
	return issue, err
}

func GetPerentSummrayIfStory(client *jira.Client, issue *jira.Issue) (summary string, desc string, err error) {
	// 檢查 Issue 是否有父單
	if issue.Fields.Parent != nil {
		parentKey := issue.Fields.Parent.Key
		// 獲取父單的詳細信息
		parentIssue, err := GetIssueByKey(client, parentKey)
		if err != nil {
			return "", "", err
		}
		// 檢查父單是否為 Story 類型
		if parentIssue.Fields.Type.Name == ISSUE_TYPE_STORY && issue.Fields.Type.Name == ISSUE_TYPE_SUB_TASK {
			// 返回父單的 Summary
			return parentIssue.Fields.Summary, parentIssue.Fields.Description, nil
		}
	}
	return issue.Fields.Summary, issue.Fields.Description, nil
}

// 添加關聯單
func AddRelated(client *jira.Client, createdKey, ReleatdKey string) error {

	if ReleatdKey != "" {
		// 添加關聯
		_, err := client.Issue.AddLink(&jira.IssueLink{
			Type: jira.IssueLinkType{
				Name: "Relates",
			},
			InwardIssue:  &jira.Issue{Key: createdKey},
			OutwardIssue: &jira.Issue{Key: ReleatdKey},
		})
		if err != nil {
			fmt.Println("jira add link err:", err)
			return err
		}
	}
	return nil
}

// 更新Sprint
func UpdateSprint(client *jira.Client, issueKey string, sprintId int) error {
	if sprintId == 0 {
		return nil
	}
	_, err := client.Issue.UpdateIssue(issueKey, map[string]any{
		"fields": map[string]any{
			SprintFieldLabel: sprintId, // 修改sprint
		},
	})
	return err
}

// 更新Epic
func UpdateEpic(client *jira.Client, issueKey, epicIssueKey string) error {
	_, err := client.Issue.UpdateIssue(issueKey, map[string]any{
		"fields": map[string]any{
			EpicFieldLabel: epicIssueKey, // 修改epic
		},
	})
	return err
}

// 更新資深人員估時時間
func UpdateEstTime(client *jira.Client, issueKey string, estTime int) error {
	_, err := client.Issue.UpdateIssue(issueKey, map[string]any{
		"fields": map[string]any{
			StaffEstTimeLabel: estTime, // 修改時間
		},
	})
	return err
}
