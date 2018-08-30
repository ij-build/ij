package logging

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/efritz/glock"
	"github.com/mgutz/ansi"
)

type (
	Processor interface {
		Start()
		Shutdown()
		Logger(outFile, errFile io.WriteCloser, writePrefix bool) Logger
	}

	processor struct {
		verbose     bool
		colorize    bool
		clock       glock.Clock
		outStream   io.Writer
		errStream   io.Writer
		colorPicker ColorPicker
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
	return newProcessor(
		verbose,
		colorize,
		glock.NewRealClock(),
		os.Stdout,
		os.Stderr,
	)
}

func newProcessor(
	verbose bool,
	colorize bool,
	clock glock.Clock,
	outStream io.Writer,
	errStream io.Writer,
) Processor {
	return &processor{
		verbose:     verbose,
		colorize:    colorize,
		clock:       clock,
		outStream:   outStream,
		errStream:   errStream,
		colorPicker: newColorPicker(colorize),
		queue:       make(chan *message),
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

func (p *processor) Logger(outFile, errFile io.WriteCloser, writePrefix bool) Logger {
	p.mutex.Lock()
	p.handles = append(p.handles, outFile, errFile)
	p.mutex.Unlock()

	return newLogger(
		p,
		p.outStream,
		outFile,
		p.errStream,
		errFile,
		writePrefix,
	)
}

func (p *processor) enqueue(message *message) {
	p.queue <- message
}

func (p *processor) process() {
	defer p.wg.Done()

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
			"%s | %s%s\n",
			message.timestamp.Format(LongTimestampFormat),
			p.buildPrefixForFile(message.prefix, message.writePrefix),
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

func (p *processor) buildPrefixForFile(prefix *Prefix, writePrefix bool) string {
	if prefix == nil || !writePrefix {
		return ""
	}

	return fmt.Sprintf("%s: ", prefix.Serialize(NilColorPicker))
}
