package config

import "time"

func extendString(child, parent string) string {
	if child != "" {
		return child
	}

	return parent
}

func extendBool(child, parent bool) bool {
	return child || parent
}

func extendInt(child, parent int) int {
	if child > 0 {
		return child
	}

	return parent
}

func extendDuration(child, parent time.Duration) time.Duration {
	if child > 0 {
		return child
	}

	return parent
}
