package runtime

import (
	"sync"
)

type Cleanup struct {
	funcs []func()
	mutex sync.Mutex
}

func NewCleanup() *Cleanup {
	return &Cleanup{}
}

func (c *Cleanup) Register(f func()) {
	c.mutex.Lock()
	c.funcs = append(c.funcs, f)
	c.mutex.Unlock()
}

func (c *Cleanup) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := len(c.funcs) - 1; i >= 0; i-- {
		c.funcs[i]()
	}
}
