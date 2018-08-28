package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/util"
)

type (
	Runner interface {
		Run(ctx context.Context, args []string, stdin io.ReadCloser, prefix *logging.Prefix) error
		RunForOutput(ctx context.Context, args []string, stdin io.ReadCloser) (string, string, error)
	}

	runner struct {
		logger  logging.Logger
		testing bool
	}
)

const TestEnvFlag = "TEST_OS_EXEC"

func NewRunner(logger logging.Logger) Runner {
	return newRunner(logger, false)
}

func newRunner(logger logging.Logger, testing bool) *runner {
	return &runner{
		logger:  logger,
		testing: testing,
	}
}

func (r *runner) Run(
	ctx context.Context,
	args []string,
	stdin io.ReadCloser,
	prefix *logging.Prefix,
) error {
	return r.run(
		ctx,
		args,
		stdin,
		newLogProcessor(prefix, r.logger.Info),
		newLogProcessor(prefix, r.logger.Error),
	)
}

func (r *runner) RunForOutput(
	ctx context.Context,
	args []string,
	stdin io.ReadCloser,
) (string, string, error) {
	outProcessor := newStringProcessor()
	errProcessor := newStringProcessor()

	err := r.run(
		ctx,
		args,
		stdin,
		outProcessor,
		errProcessor,
	)

	return outProcessor.String(), errProcessor.String(), err
}

func (r *runner) run(
	ctx context.Context,
	args []string,
	stdin io.ReadCloser,
	outProcessor outputProcessor,
	errProcessor outputProcessor,
) error {
	if r.logger != nil {
		r.logger.Debug(
			nil,
			"Running command: %s", strings.Join(args, " "),
		)
	}

	command := exec.CommandContext(
		ctx,
		args[0],
		args[1:]...,
	)

	if stdin != nil {
		defer stdin.Close()
		command.Stdin = stdin
	}

	if r.testing {
		command.Env = []string{
			fmt.Sprintf("%s=1", TestEnvFlag),
		}
	}

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
