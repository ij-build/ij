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
		Logger(prefix string, outfile, errfile io.WriteCloser) Logger
	}

	processor struct {
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

func NewProcessor() Processor {
	return &processor{
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

func (p *processor) Logger(prefix string, outfile, errfile io.WriteCloser) Logger {
	p.mutex.Lock()
	p.handles = append(p.handles, outfile, errfile)
	p.mutex.Unlock()

	return newLogger(
		p,
		prefix,
		p.colorPicker.next(),
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
			emergencyLog("failed to write log (%s)", err.Error())
		}

		if err := writeAll(message.file, []byte(fileText)); err != nil {
			emergencyLog("failed to write log (%s)", err.Error())
		}
	}
}
