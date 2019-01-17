package environment

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/rw_grim/govcs"
	"bitbucket.org/rw_grim/govcs/vcs"
)

const BuildTimeFormat = "2006-01-02T15:04:05-0700"

func Default() (Environment, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	lines := []string{
		fmt.Sprintf("UID=%s", user.Uid),
		fmt.Sprintf("GID=%s", user.Gid),
		fmt.Sprintf("BUILD_TIME=%s", formatNow()),
	}

	if repo, err := getVCS(); err == nil {
		name := strings.ToUpper(repo.Name())

		lines = append(lines, fmt.Sprintf("%s_BRANCH_NORMALIZED=%s", name, normalize(repo.Branch())))
		lines = append(lines, fmt.Sprintf("%s_BRANCH=%s", name, repo.Branch()))
		lines = append(lines, fmt.Sprintf("%s_COMMIT_SHORT=%s", name, repo.ShortCommit()))
		lines = append(lines, fmt.Sprintf("%s_COMMIT=%s", name, repo.Commit()))
		lines = append(lines, fmt.Sprintf("%s_REMOTE=%s", name, repo.Remote("")))
	}

	return New(lines), nil
}

func formatNow() string {
	return time.Now().UTC().Format(BuildTimeFormat)
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
