package gittool

import (
	"fmt"
	"testing"
)

func Test_gitClone(t *testing.T) {
	r, err := GitClone("https://gitlab.paradise-soft.com.tw/PM/i18n_admin.git")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func Test_gitAdd(t *testing.T) {
	err := GitAdd("i18n_admin", "admin-i18n.csv")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_gitCommitAndPush(t *testing.T) {
	err := GitCommitAndPush("i18n_admin")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AddTag(t *testing.T) {
	err := AddTag("i18n_admin")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_getNextVersion(t *testing.T) {
	version, err := getCurrentVersion("i18n_admin")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(version)
}
