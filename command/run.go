package command

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"

	"github.com/efritz/pvc/logging"
)

func Run(
	ctx context.Context,
	prefix *logging.Prefix,
	args []string,
	logger logging.Logger,
) error {
	return run(
		ctx,
		args,
		newLogProcessor(prefix, logger.Info),
		newLogProcessor(prefix, logger.Error),
	)
}

func RunForOutput(ctx context.Context, args []string) (string, error) {
	processor := newStringProcessor()
	if err := run(ctx, args, processor, nilProcessor); err != nil {
		return "", err
	}

	return processor.String(), nil
}

//
//

func run(ctx context.Context, args []string, outProcessor, errProcessor outputProcessor) error {
	command := exec.CommandContext(
		ctx,
		args[0],
		args[1:]...,
	)

	outReader, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	errReader, err := command.StderrPipe()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go processOutput(outReader, outProcessor, wg)
	go processOutput(errReader, errProcessor, wg)

	if err := command.Start(); err != nil {
		return err
	}

	if err := command.Wait(); err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func processOutput(r io.Reader, p outputProcessor, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		p.Process(scanner.Text())
	}
}
