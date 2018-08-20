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
		Logger(outfile, errfile io.WriteCloser) Logger
	}

	processor struct {
		verbose     bool
		colorize    bool
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

func NewProcessor(verbose, colorize bool) Processor {
	return &processor{
		verbose:     verbose,
		colorize:    colorize,
		colorPicker: newColorPicker(colorize),
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

func (p *processor) Logger(outfile, errfile io.WriteCloser) Logger {
	p.mutex.Lock()
	p.handles = append(p.handles, outfile, errfile)
	p.mutex.Unlock()

	return newLogger(
		p,
		outfile,
		errfile,
	)
}

func (p *processor) enqueue(message *message) {
	p.queue <- message
}

func (p *processor) process() {
	defer p.wg.Done()

	// TODO - also need colors based on level

	for message := range p.queue {
		text := message.Text()

		if p.colorize {
			text = fmt.Sprintf(
				"%s%s%s",
				levelColors[message.level],
				text,
				ansi.Reset,
			)
		}

		streamText := fmt.Sprintf(
			"%s %s%s\n",
			message.timestamp.Format(ShortTimestampFormat),
			p.buildPrefix(message.prefix),
			text,
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

func (p *processor) buildPrefix(prefix *Prefix) string {
	if prefix == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s: ",
		prefix.Serialize(p.colorPicker),
	)
}
