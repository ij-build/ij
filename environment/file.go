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

		name, value := split(line)

		lines = append(lines, fmt.Sprintf(
			"%s=%s",
			strings.ToUpper(name),
			value,
		))
	}

	return lines, nil
}
