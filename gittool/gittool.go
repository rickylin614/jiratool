package gittool

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func GitClone(projectUrl string) (string, error) {
	// Clone the Git repository
	repoURL := projectUrl
	repoNameWithGit := path.Base(projectUrl)
	repoName := strings.TrimSuffix(repoNameWithGit, ".git")

	cloneCmd := exec.Command("git", "clone", repoURL)
	cloneCmd.Dir = "."
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		return repoName, err
	}

	// Change to the repository directory
	os.Chdir(repoName)
	defer os.Chdir("..")

	// Pull the latest changes
	pullCmd := exec.Command("git", "pull")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	if err := pullCmd.Run(); err != nil {
		return repoName, err
	}

	fmt.Printf("Successfully pulled the latest changes for %s\n", projectUrl)
	return repoName, nil
}

func GitAdd(repoName string, filename string) error {
	if err := os.Chdir(repoName); err != nil {
		return err
	}
	defer os.Chdir("..")

	target := fmt.Sprintf("./%s", filename)

	gitAddCmd := exec.Command("git", "add", target)
	gitAddCmd.Stdout = os.Stdout
	gitAddCmd.Stderr = os.Stderr

	if err := gitAddCmd.Run(); err != nil {
		return err
	}
	return nil
}

func GitCommitAndPush(repoName string) error {
	if err := os.Chdir(repoName); err != nil {
		return err
	}
	defer os.Chdir("..")

	// Commit changes with current date and time
	commitMessage := time.Now().Format("2006-01-02 15:04:05")
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	if err := commitCmd.Run(); err != nil {
		return err
	}

	// Push changes
	pushCmd := exec.Command("git", "push")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	return pushCmd.Run()
}

func AddTag(repoName string) error {
	nextVerion, err := getNextVersion(repoName)
	if err != nil {
		return err
	}

	tagMessage := time.Now().Format("2006-01-02 15:04:05")
	cmd := exec.Command("git", "tag", "-a", nextVerion, "-m", tagMessage)
	cmd.Dir = repoName
	if err := cmd.Run(); err != nil {
		return err
	}

	// 推送 tag 到远程仓库
	pushCmd := exec.Command("git", "push", "origin", nextVerion)
	pushCmd.Dir = repoName

	return pushCmd.Run()
}

func getNextVersion(repoName string) (string, error) {
	// 获取当前版本号
	currentVersion, err := getCurrentVersion(repoName)
	if err != nil {
		return "", err
	}

	// 解析当前版本号
	major, minor, patch, err := parseVersion(currentVersion)
	if err != nil {
		return "", err
	}

	// 递增 patch 版本号
	nextPatch := patch + 1

	// 构建下一个版本号
	nextVersion := fmt.Sprintf("v%d.%d.%d", major, minor, nextPatch)

	return nextVersion, nil
}

func getCurrentVersion(repoName string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = repoName

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func parseVersion(version string) (int, int, int, error) {
	var major, minor, patch int
	_, err := fmt.Sscanf(version, "v%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return 0, 0, 0, err
	}
	return major, minor, patch, nil
}
