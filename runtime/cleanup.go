package runtime

import (
	"fmt"
	"sync"
)

type Cleanup struct {
	funcs []func() error
	mutex sync.Mutex
}

func NewCleanup() *Cleanup {
	return &Cleanup{}
}

func (c *Cleanup) Register(f func() error) {
	c.mutex.Lock()
	c.funcs = append(c.funcs, f)
	c.mutex.Unlock()
}

func (c *Cleanup) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := len(c.funcs) - 1; i >= 0; i-- {
		if err := c.funcs[i](); err != nil {
			// TODO
			fmt.Printf("Error: %#v\n", err.Error())
		}
	}
}
