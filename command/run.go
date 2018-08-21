package command

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/efritz/pvc/logging"
	"github.com/efritz/pvc/util"
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
		logger,
	)
}

func RunForOutput(
	ctx context.Context,
	args []string,
	logger logging.Logger,
) (string, string, error) {
	outProcessor := newStringProcessor()
	errProcessor := newStringProcessor()

	err := run(
		ctx,
		args,
		outProcessor,
		errProcessor,
		logger,
	)

	return outProcessor.String(), errProcessor.String(), err
}

//
//

func run(
	ctx context.Context,
	args []string,
	outProcessor outputProcessor,
	errProcessor outputProcessor,
	logger logging.Logger,
) error {
	if logger != nil {
		logger.Debug(
			nil,
			"Running command: %s", strings.Join(args, " "),
		)
	}

	command := exec.CommandContext(
		ctx,
		args[0],
		args[1:]...,
	)

	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	outReader, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	errReader, err := command.StderrPipe()
	if err != nil {
		return err
	}

	wg := util.RunParallel(
		func() { processOutput(outReader, outProcessor) },
		func() { processOutput(errReader, errProcessor) },
	)

	if err := command.Run(); err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func processOutput(r io.Reader, p outputProcessor) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		p.Process(scanner.Text())
	}
}
