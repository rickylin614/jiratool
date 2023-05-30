package fynetool

import (
	"errors"
	"fmt"
	"jiratool/conf"
	"jiratool/jiratool"

	"github.com/andygrunwald/go-jira"
)

var initError error
var client *jira.Client

var EpicIssue []jira.Issue = make([]jira.Issue, 0)
var Springs []jiratool.Sprint = make([]jiratool.Sprint, 0)

func InitDataList() {
	if conf.GetConfig().UserName == "" || conf.GetConfig().Pwd == "" || conf.GetConfig().Project == "" || conf.GetConfig().Jiraurl == "" {
		initError = errors.New("設定檔案讀取失敗, 請設定後重新啟動程序")
	} else {
		jiratool.SetProject(conf.GetConfig().Project)
		jiratool.SetLabel(conf.GetConfig().Label)
		jiratool.SetComponent(conf.GetConfig().Component)
	}

	tp := jira.BasicAuthTransport{
		Username: conf.GetConfig().UserName,
		Password: conf.GetConfig().Pwd,
	}

	var err error
	client, err = jira.NewClient(tp.Client(), conf.GetConfig().Jiraurl)
	if err != nil {
		fmt.Println("jira new client error", err)
		initError = err
	}
	if initError == nil {
		// 更新epic列表
		epicList, _ := jiratool.GetEpicList(client)
		EpicIssue = epicList
		// 更新sprint列表
		sprintList, _ := jiratool.GetSprintList(client)
		Springs = sprintList
	}

}
