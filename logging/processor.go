package logging

import (
	"fmt"
	"io"
	"sync"

	"github.com/mgutz/ansi"
)

type (
	Processor interface {
		Start()
		Shutdown()
		BaseLogger(prefix string, outfile, errfile io.WriteCloser) Logger
		TaskLogger(prefix string, outfile, errfile io.WriteCloser) Logger
	}

	processor struct {
		verbose     bool
		colorPicker *colorPicker
		queue       chan *message
		handles     []io.Closer
		mutex       sync.Mutex
		once        sync.Once
		wg          sync.WaitGroup
	}
)

const (
	ShortTimestampFormat = "15:04:05"
	LongTimestampFormat  = "2006-01-02 15:04:05.000"
)

func NewProcessor(verbose bool) Processor {
	return &processor{
		verbose:     verbose,
		colorPicker: newColorPicker(),
		queue:       make(chan *message),
		handles:     []io.Closer{},
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

func (p *processor) BaseLogger(prefix string, outfile, errfile io.WriteCloser) Logger {
	return p.loggerWithColor(prefix, outfile, errfile, "")
}

func (p *processor) TaskLogger(prefix string, outfile, errfile io.WriteCloser) Logger {
	p.mutex.Lock()
	colorCode := p.colorPicker.next()
	p.mutex.Unlock()

	return p.loggerWithColor(prefix, outfile, errfile, colorCode)
}

func (p *processor) loggerWithColor(prefix string, outfile, errfile io.WriteCloser, colorCode string) Logger {
	p.mutex.Lock()
	p.handles = append(p.handles, outfile, errfile)
	p.mutex.Unlock()

	return newLogger(
		p,
		prefix,
		colorCode,
		outfile,
		errfile,
	)
}

func (p *processor) enqueue(message *message) {
	p.queue <- message
}

func (p *processor) process() {
	defer p.wg.Done()

	for message := range p.queue {
		streamText := fmt.Sprintf(
			"%s%s [%s] %s%s\n",
			message.colorCode,
			message.timestamp.Format(ShortTimestampFormat),
			message.prefix,
			message.Text(),
			ansi.Reset,
		)

		fileText := fmt.Sprintf(
			"%s | %s\n",
			message.timestamp.Format(LongTimestampFormat),
			message.Text(),
		)

		if err := writeAll(message.stream, []byte(streamText)); err != nil {
			EmergencyLog("error: failed to write log: %s", err.Error())
		}

		if err := writeAll(message.file, []byte(fileText)); err != nil {
			EmergencyLog("error: failed to write log: %s", err.Error())
		}
	}
}
