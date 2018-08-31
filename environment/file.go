package environment

import (
	"fmt"
	"strings"
)

func NormalizeEnvironmentFile(text string) ([]string, error) {
	lines := []string{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)

		if line == "" || line[0] == '#' {
			continue
		}

		if !strings.Contains(line, "=") {
			return nil, fmt.Errorf(
				"Malformed entry in environments file: %s",
				line,
			)
		}

		lines = append(lines, line)
	}

	return lines, nil
}
