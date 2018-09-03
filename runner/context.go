package runner

import "github.com/efritz/ij/environment"

type RunContext struct {
	Failure     bool
	Environment environment.Environment
}
