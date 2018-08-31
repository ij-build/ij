package environment

import (
	"os"
	"strings"
)

func split(value string) (string, string) {
	if parts := strings.SplitN(value, "=", 2); len(parts) == 2 {
		return parts[0], parts[1]
	}

	return value, os.Getenv(value)
}
