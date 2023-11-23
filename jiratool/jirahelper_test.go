package jiratool

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

var tp jira.BasicAuthTransport = jira.BasicAuthTransport{}

const JIRA_URL = "https://jira.paradise-soft.com.tw/"

func init() {
	tp = jira.BasicAuthTransport{
		Username: "",
		Password: "",
	}
}

func ExampleUpdateSprint() {

	client, err := jira.NewClient(tp.Client(), JIRA_URL)
	if err != nil {
		fmt.Println("jira new client error", err)
	}

	err = UpdateSprint(client, "BE1-3716", 343)
	fmt.Println(err)

	// output:
	// <nil>
}

func ExampleGetEpicList() {

	client, err := jira.NewClient(tp.Client(), JIRA_URL)
	if err != nil {
		fmt.Println("jira new client error", err)
	}

	list, err := GetEpicList(client)
	fmt.Println(list, err)

	// output:
	//
}

func ExampleGetSprintList() {

	client, err := jira.NewClient(tp.Client(), JIRA_URL)
	if err != nil {
		fmt.Println("jira new client error", err)
	}

	list, err := GetSprintList(client)
	for _, sprint := range list {
		fmt.Printf("%d %s %s \n", sprint.Id, sprint.BoardName, sprint.Name)
	}

	// output:
	//
}
