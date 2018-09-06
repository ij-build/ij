package environment

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/rw_grim/govcs"
	"bitbucket.org/rw_grim/govcs/vcs"
)

const BuildTimeFormat = "2006-01-02T15:04:05-0700"

func Default() Environment {
	buildTime := time.Now().UTC().Format(BuildTimeFormat)

	lines := []string{
		fmt.Sprintf("BUILD_TIME=%s", buildTime),
	}

	if repo, err := getVCS(); err == nil {
		name := strings.ToUpper(repo.Name())

		lines = append(lines, fmt.Sprintf("%s_BRANCH_NORMALIZED=%s", name, normalize(repo.Branch())))
		lines = append(lines, fmt.Sprintf("%s_BRANCH=%s", name, repo.Branch()))
		lines = append(lines, fmt.Sprintf("%s_COMMIT_SHORT=%s", name, repo.ShortCommit()))
		lines = append(lines, fmt.Sprintf("%s_COMMIT=%s", name, repo.Commit()))
		lines = append(lines, fmt.Sprintf("%s_REMOTE=%s", name, repo.Remote("")))
	}

	return New(lines)
}

func getVCS() (vcs.VCS, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return govcs.Detect(wd)
}

func normalize(value string) string {
	return regexp.MustCompile("[^._A-Za-z0-9-]").ReplaceAllLiteralString(value, "_")
}
