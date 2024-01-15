package fynetool

import (
	"jiratool/conf"
	"jiratool/jiratool"
	"jiratool/lib"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/andygrunwald/go-jira"
)

const title = "關聯單產生器@ricky in 2023"

// 初始化視窗
func InitFyneApp() fyne.Window {
	a := app.New()
	// 設置字體
	a.Settings().SetTheme(&lib.MyTheme{})

	// 初始化視窗顯示內容
	w := a.NewWindow(title)

	// 設定寬高
	w.Resize(fyne.NewSize(600, 400))

	return w
}

// 產生各個元件
func SetttingWidget(w fyne.Window) fyne.Window {
	// 初始化要設定的Issue參數
	issueInfo := jiratool.IssueInfo{
		EpicKey: "BE1-191",
	}

	// 輸入框設定
	entry := widget.NewEntry()
	entry.SetPlaceHolder("請輸入欲關聯的單號:")
	entry.Resize(fyne.Size{Width: 1000, Height: 5000})

	// 錯誤標籤
	errorLabel := widget.NewLabel("")

	// 顯示完成的Issue
	showIssueUrl := widget.NewEntry()
	showIssueUrl.Disable()

	// 連結
	hyperlink := widget.NewHyperlink("", nil)
	hyperlink.Hidden = true
	hyperlink.Resize(fyne.NewSize(3, 1))

	epicList, _ := jiratool.GetEpicList(client)
	sprintList, _ := jiratool.GetSprintList(client)
	sprintList = append([]jiratool.Sprint{{Name: "空", Id: 0}}, sprintList...)
	fixversionList, _ := jiratool.GetUnreleasedVersions(client)

	// 建立下拉式選單 (EPIC)
	selectorEpic := widget.NewSelect(
		getEpicIssueNames(epicList),
		selectEpicIssue(epicList, &issueInfo.EpicKey),
	)
	selectorEpic.SetSelected("平台-彩票")

	// 建立下拉式選單 (Sprint)
	selectorSprint := widget.NewSelect(
		getSprintNames(sprintList),
		selectSprintIssue(sprintList, &issueInfo.SprintId),
	)
	defaultSprit := ""
	if len(sprintList) > 0 {
		defaultSprit = sprintList[0].Name
	}
	selectorSprint.SetSelected(defaultSprit)

	// 建立下拉選單 (fixversion)
	selectorVersion := widget.NewSelect(
		getVerionNames(fixversionList),
		selectVersionIssue(fixversionList, &issueInfo.VersionId),
	)
	selectorVersion.SetSelected(setDefaultVersion(fixversionList))

	// 創建產生按鈕
	btnCreateRelated := widget.NewButton("創建關聯單", CreateRelatedIssue(errorLabel, showIssueUrl, entry, hyperlink, &issueInfo))
	btnCreate := widget.NewButton("創單", CreateIssue(errorLabel, showIssueUrl, entry, hyperlink, &issueInfo))
	btnSubCreate := widget.NewButton("創子單", CreateSubIssue(errorLabel, showIssueUrl, entry, hyperlink, &issueInfo))
	// btnConfigReloead := widget.NewButton("Reload設定", func() {})

	btnLayout := container.New(
		layout.NewGridLayout(5),
		btnCreateRelated,
		btnCreate,
		btnSubCreate,
	)

	selectorLayout := container.New(
		layout.NewGridLayout(2),
		selectorEpic,
		selectorSprint,
	)

	w.SetContent(container.NewVBox(
		widget.NewLabel("請輸入欲關聯的單號"),
		entry,
		btnLayout,
		errorLabel,
		selectorLayout,
		selectorVersion,
		hyperlink,
	))

	return w
}

func getEpicIssueNames(EpicIssue []jira.Issue) []string {
	s := make([]string, len(EpicIssue))
	for i, v := range EpicIssue {
		s[i] = v.Fields.Summary
	}
	return s
}

func selectEpicIssue(EpicIssue []jira.Issue, EpicIssueKey *string) func(string) {
	return func(s string) {
		for _, v := range EpicIssue {
			if s == v.Fields.Summary {
				*EpicIssueKey = v.Key
				return
			}
		}
		*EpicIssueKey = ""
	}
}

func getSprintNames(EpicIssue []jiratool.Sprint) []string {
	s := make([]string, len(EpicIssue))
	for i, v := range EpicIssue {
		s[i] = v.Name
	}
	return s
}

func selectSprintIssue(EpicIssue []jiratool.Sprint, SprintId *int) func(string) {
	return func(s string) {
		for _, v := range EpicIssue {
			if s == v.Name {
				*SprintId = v.Id
				return
			}
		}
		*SprintId = 0
	}
}

func getVerionNames(versions []jiratool.Version) []string {
	s := make([]string, len(versions))
	for i, v := range versions {
		s[i] = v.Name
	}
	s = append([]string{""}, s...)
	return s
}

func selectVersionIssue(versions []jiratool.Version, verionId *string) func(string) {
	return func(s string) {
		for _, v := range versions {
			if s == v.Name {
				*verionId = v.Id
				return
			}
		}
		return
	}
}

func setDefaultVersion(versions []jiratool.Version) string {
	for _, v := range versions {
		if strings.HasPrefix(v.Name, "xunya") {
			return v.Name
		}
	}
	return ""
}

// 創建關聯單
func CreateRelatedIssue(
	errorLabel *widget.Label,
	showIssueUrl *widget.Entry,
	entry *widget.Entry,
	hyperlink *widget.Hyperlink,
	issueInfo *jiratool.IssueInfo) func() {

	return func() {
		errorLabel.SetText("")
		showIssueUrl.SetText(``)
		releatedIssue := ""
		if initError != nil {
			errorLabel.SetText(initError.Error())
			return
		}

		releatedIssue = entry.Text

		str, err := jiratool.GeneratorRelatedIssue(client, releatedIssue, issueInfo) // 創關聯單
		if err != nil {
			errorLabel.SetText(`創立錯誤!! err:` + err.Error())
		} else {
			errorLabel.SetText(`創單成功 單號:` + *str)
			hyperlink.URL, _ = url.Parse(conf.GetConfig().Jiraurl + `browse/` + *str)
			hyperlink.SetText(*str)
			hyperlink.Hidden = false
		}

	}
}

// 創單
func CreateIssue(
	errorLabel *widget.Label,
	showIssueUrl *widget.Entry,
	entry *widget.Entry,
	hyperlink *widget.Hyperlink,
	issueInfo *jiratool.IssueInfo) func() {

	return func() {
		errorLabel.SetText("")
		showIssueUrl.SetText(``)
		if initError != nil {
			errorLabel.SetText(initError.Error())
			return
		}

		str, err := jiratool.GeneratorIssue(client, issueInfo) // 創關聯單
		if err != nil {
			errorLabel.SetText(`創立錯誤!! err:` + err.Error())
		} else {
			errorLabel.SetText(`創單成功 單號:` + *str)
			hyperlink.URL, _ = url.Parse(conf.GetConfig().Jiraurl + `browse/` + *str)
			hyperlink.SetText(*str)
			hyperlink.Hidden = false
		}

	}
}

// 創建子單
func CreateSubIssue(
	errorLabel *widget.Label,
	showIssueUrl *widget.Entry,
	entry *widget.Entry,
	hyperlink *widget.Hyperlink,
	issueInfo *jiratool.IssueInfo) func() {

	return func() {
		errorLabel.SetText("")
		showIssueUrl.SetText(``)
		if initError != nil {
			errorLabel.SetText(initError.Error())
			return
		}

		str, err := jiratool.GeneratorSubIssue(client, entry.Text, issueInfo) // 創關聯單
		if err != nil {
			errorLabel.SetText(`創立錯誤!! err:` + err.Error())
		} else {
			errorLabel.SetText(`創單成功 單號:` + *str)
			hyperlink.URL, _ = url.Parse(conf.GetConfig().Jiraurl + `browse/` + *str)
			hyperlink.SetText(*str)
			hyperlink.Hidden = false
		}

	}
}
