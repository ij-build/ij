package environment

import (
	"fmt"
	"sort"
)

type Environment map[string]string

func New(values []string) Environment {
	env := Environment{}
	for _, line := range values {
		k, v := split(line)
		env[k] = v
	}

	return env
}

func (e Environment) Keys() []string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (e Environment) Serialize() []string {
	lines := []string{}
	for _, k := range e.Keys() {
		lines = append(lines, fmt.Sprintf("%s=%s", k, e[k]))
	}

	return lines
}
