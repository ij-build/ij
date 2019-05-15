package environment

import (
	"fmt"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"
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

	branch := gitOutput("rev-parse", "--abbrev-ref", "HEAD")
	remote := gitOutput("config", "--local", "remote.origin.url")
	commit := gitOutput("rev-parse", "HEAD")
	shortCommit := gitOutput("rev-parse", "--short", "HEAD")

	if commit != "" {
		lines = append(lines, fmt.Sprintf("GIT_BRANCH_NORMALIZED=%s", normalize(branch)))
		lines = append(lines, fmt.Sprintf("GIT_BRANCH=%s", branch))
		lines = append(lines, fmt.Sprintf("GIT_COMMIT_SHORT=%s", shortCommit))
		lines = append(lines, fmt.Sprintf("GIT_COMMIT=%s", commit))
		lines = append(lines, fmt.Sprintf("GIT_REMOTE=%s", remote))
	}

	return New(lines), nil
}

func formatNow() string {
	return time.Now().UTC().Format(BuildTimeFormat)
}

func normalize(value string) string {
	return regexp.MustCompile("[^._A-Za-z0-9-]").ReplaceAllLiteralString(value, "_")
}

func gitOutput(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
