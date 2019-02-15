package runner

import (
	"sync"

	"github.com/efritz/ij/environment"
)

type RunContext struct {
	parent           *RunContext
	Failure          bool
	Environment      environment.Environment
	tags             []string
	tagsMutex        sync.RWMutex
	exportedEnv      []string
	exportedEnvMutex sync.RWMutex
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

func (c *RunContext) AddTags(tags []string) {
	if c.parent != nil {
		c.parent.AddTags(tags)
		return
	}

	c.tagsMutex.Lock()
	c.tags = append(c.tags, tags...)
	c.tagsMutex.Unlock()
}

func (c *RunContext) GetTags() []string {
	if c.parent != nil {
		return c.parent.GetTags()
	}

	c.tagsMutex.RLock()
	defer c.tagsMutex.RUnlock()

	tags := []string{}
	for _, tag := range c.tags {
		tags = append(tags, tag)
	}

	return tags
}

func (c *RunContext) ExportEnv(line string) {
	if c.parent != nil {
		c.parent.ExportEnv(line)
		return
	}

	c.exportedEnvMutex.Lock()
	c.exportedEnv = append(c.exportedEnv, line)
	c.exportedEnvMutex.Unlock()
}

func (c *RunContext) GetExportedEnv() []string {
	if c.parent != nil {
		return c.parent.GetExportedEnv()
	}

	c.exportedEnvMutex.RLock()
	defer c.exportedEnvMutex.RUnlock()

	env := []string{}
	for _, line := range c.exportedEnv {
		env = append(env, line)
	}

	return env
}
