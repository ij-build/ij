package logging

import (
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/mgutz/ansi"
)

type (
	ColorPicker interface {
		Colorize(val string) string
	}

	colorPicker struct {
		enabled bool
		cache   map[string]string
		mutex   sync.RWMutex
	}
)

var NilColorPicker = newColorPicker(false)

func newColorPicker(enabled bool) ColorPicker {
	return &colorPicker{
		enabled: enabled,
		cache:   map[string]string{},
	}
}

func (p *colorPicker) Colorize(val string) string {
	if !p.enabled {
		return val
	}

	return fmt.Sprintf(
		"%s%s%s",
		p.colorFor(val),
		val,
		ansi.Reset,
	)
}

func (p *colorPicker) colorFor(val string) string {
	p.mutex.RLock()
	if color, ok := p.cache[val]; ok {
		p.mutex.RUnlock()
		return color
	}

	p.mutex.RUnlock()
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if color, ok := p.cache[val]; ok {
		return color
	}

	color := chooseColor(val)
	p.cache[val] = color
	return color
}

func chooseColor(val string) string {
	hash := fnv.New32a()
	hash.Write([]byte(val))

	return colors[int(hash.Sum32())%len(colors)]
}
