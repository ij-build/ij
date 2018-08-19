package logging

import (
	"fmt"
	"io"
	"sync"
)

type (
	Processor interface {
		Start()
		Shutdown()
		Logger(prefix string, outfile, errfile io.WriteCloser) Logger
	}

	processor struct {
		queue   chan *message
		handles []io.Closer
		mutex   sync.Mutex
		once    sync.Once
		wg      sync.WaitGroup
	}

	message struct {
		level  LogLevel
		format string
		args   []interface{}
		prefix string
		file   io.Writer
	}
)

func NewProcessor() Processor {
	return &processor{
		queue:   make(chan *message),
		handles: []io.Closer{},
	}
}

func (p *processor) Start() {
	p.wg.Add(1)
	go p.process()
}

func (p *processor) Shutdown() {
	p.once.Do(func() {
		close(p.queue)
	})

	p.wg.Wait()

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, handle := range p.handles {
		handle.Close()
	}

	p.handles = p.handles[:0]
}

func (p *processor) Logger(prefix string, outfile, errfile io.WriteCloser) Logger {
	p.mutex.Lock()
	p.handles = append(p.handles, outfile, errfile)
	p.mutex.Unlock()

	return newLogger(p, prefix, outfile, errfile)
}

func (p *processor) enqueue(message *message) {
	p.queue <- message
}

func (p *processor) process() {
	defer p.wg.Done()

	for message := range p.queue {
		fmt.Printf(fmt.Sprintf("%s%s\n", message.level.Prefix(), message.format), message.args...)

		// TODO - probably want timestamps too
		// TODO - short writes
		message.file.Write([]byte(fmt.Sprintf(fmt.Sprintf("%s\n", message.format), message.args...)))
	}
}
