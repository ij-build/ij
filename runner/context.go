package runner

import (
	"sync"

	"github.com/efritz/ij/environment"
)

type RunContext struct {
	parent      *RunContext
	Failure     bool
	Environment environment.Environment
	exportedEnv []string
	envMutex    sync.RWMutex
}

func NewRunContext(parent *RunContext) *RunContext {
	context := &RunContext{
		parent: parent,
	}

	if parent != nil {
		context.Failure = parent.Failure
		context.Environment = parent.Environment
	}

	return context
}

func (c *RunContext) ExportEnv(line string) {
	if c.parent != nil {
		c.parent.ExportEnv(line)
		return
	}

	c.envMutex.Lock()
	c.exportedEnv = append(c.exportedEnv, line)
	c.envMutex.Unlock()
}

func (c *RunContext) GetExportedEnv() []string {
	if c.parent != nil {
		return c.parent.GetExportedEnv()
	}

	c.envMutex.RLock()
	defer c.envMutex.RUnlock()

	env := []string{}
	for _, line := range c.exportedEnv {
		env = append(env, line)
	}

	return c.exportedEnv
}
