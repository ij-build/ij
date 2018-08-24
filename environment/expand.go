package environment

import (
	"fmt"
	"os"
)

const ExpandMaxIterations = 50

func (e Environment) ExpandString(template string) (string, error) {
	return e.expandString(template, ExpandMaxIterations)
}

func (e Environment) ExpandSlice(templates []string) ([]string, error) {
	expanded := []string{}
	for _, template := range templates {
		val, err := e.ExpandString(template)
		if err != nil {
			return nil, err
		}
		expanded = append(expanded, val)
	}
	return expanded, nil
}

//
//

func (e Environment) expandString(template string, count int) (string, error) {
	if count == 0 {
		return "", fmt.Errorf(
			"exceeded %d iterations while expanding environment: current template is `%s`",
			ExpandMaxIterations,
			template,
		)
	}

	if expanded := os.Expand(template, e.translate); expanded != template {
		return e.expandString(expanded, count-1)
	}

	return template, nil
}

func (e Environment) translate(name string) string {
	if name == "" {
		return "$"
	}

	if value, ok := e[name]; ok {
		return value
	}

	return fmt.Sprintf("${%s}", name)
}
