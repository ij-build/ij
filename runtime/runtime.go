package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/efritz/pvc/command"
	"github.com/efritz/pvc/config"
	"github.com/efritz/pvc/environment"
	"github.com/efritz/pvc/logging"
	shellquote "github.com/kballard/go-shellquote"
)

type (
	Runtime struct {
		config       *config.Config
		logProcessor logging.Processor
		logger       logging.Logger
		env          []string
		id           string
		ctx          context.Context
		cancel       func()
		cleanup      *Cleanup
		builddir     *Builddir
		workspace    *Workspace
	}

	Runner func() error
)

func NewRuntime(id string, config *config.Config, logProcessor logging.Processor, env []string) *Runtime {
	ctx, cancel := context.WithCancel(context.Background())

	return &Runtime{
		config:       config,
		logProcessor: logProcessor,
		env:          env,
		id:           id,
		ctx:          ctx,
		cancel:       cancel,
		cleanup:      NewCleanup(),
	}
}

func (r *Runtime) Setup() error {
	setupFuncs := []func() error{
		r.setupBuilddir,
		r.setupWorkspace,
		r.setupLogger,
	}

	for _, f := range setupFuncs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) setupBuilddir() error {
	builddir := NewBuilddir()

	if err := builddir.Setup(); err != nil {
		return err
	}

	r.builddir = builddir
	r.cleanup.Register(builddir.Teardown)
	return nil
}

func (r *Runtime) setupWorkspace() error {
	workspace := NewWorkspace(r)
	if err := workspace.Setup(); err != nil {
		return err
	}

	r.workspace = workspace
	r.cleanup.Register(workspace.Teardown)
	return nil
}

func (r *Runtime) setupLogger() error {
	outfile, errfile, err := r.builddir.LogFiles("pvc")
	if err != nil {
		return err
	}

	r.logger = r.logProcessor.Logger("pvc", outfile, errfile)
	return nil
}

func (r *Runtime) Shutdown() {
	r.cancel()
	r.cleanup.Cleanup()
}

func (r *Runtime) Run(plans []string) error {
	for _, name := range plans {
		plan, ok := r.config.Plans[name]
		if !ok {
			return fmt.Errorf("plan %s not found", name)
		}

		if err := r.runPlan(plan); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) runPlan(plan *config.Plan) error {
	for _, stage := range plan.Stages {
		if err := r.runStage(plan, stage); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) runStage(plan *config.Plan, stage *config.Stage) error {
	runners := []Runner{}
	for i, taskInstance := range stage.Tasks {
		task, ok := r.config.Tasks[taskInstance.Name]
		if !ok {
			return fmt.Errorf("task %s not found", taskInstance.Name)
		}

		env := environment.Merge(
			environment.New(r.env),
			environment.New(r.config.Environment),
			environment.New(task.Environment),
			environment.New(plan.Environment),
			environment.New(stage.Environment),
			environment.New(taskInstance.Environment),
		)

		index := i
		runner := func() error {
			return r.runTask(plan, stage, task, index, env)
		}

		runners = append(runners, runner)
	}

	if stage.Concurrent {
		return runConcurrent(runners)
	}

	return runSequential(runners)
}

func (r *Runtime) runTask(
	plan *config.Plan,
	stage *config.Stage,
	task *config.Task,
	index int,
	env environment.Environment,
) error {
	if ok, missing := containsAll(env.Keys(), task.RequiredEnvironment); !ok {
		panic(fmt.Sprintf("missing environment values: %#v", missing))
	}

	args, err := r.buildTaskCommandArgs(task, env)
	if err != nil {
		return err
	}

	outfile, errfile, err := r.builddir.LogFiles(strings.Join([]string{
		plan.Name,
		stage.Name,
		task.Name,
		fmt.Sprintf("%d", index),
	}, "."))

	if err != nil {
		return err
	}

	return command.Run(
		context.Background(),
		args,
		r.logProcessor.Logger(
			fmt.Sprintf(
				"%s/%s/%s (%d)",
				plan.Name,
				stage.Name,
				task.Name,
				index,
			),
			outfile,
			errfile,
		),
	)
}

func (r *Runtime) buildTaskCommandArgs(
	task *config.Task,
	env environment.Environment,
) ([]string, error) {
	args := []string{
		"docker",
		"run",
		"--rm",
		"-v",
		fmt.Sprintf("%s:/workspace", r.workspace.VolumePath),
		"-w",
		"/workspace",
	}

	for _, line := range env.Serialize() {
		args = append(args, "-e", line)
	}

	args = append(args, task.Image)

	commandArgs, err := shellquote.Split(task.Command)
	if err != nil {
		return nil, err
	}

	args = append(args, commandArgs...)
	return args, nil
}

func runConcurrent(runners []Runner) error {
	errChan := make(chan error)
	defer close(errChan)

	for _, runner := range runners {
		go func(runner func() error) {
			errChan <- runner()
		}(runner)
	}

	var firstErr error
	for i := 0; i < len(runners); i++ {
		if err := <-errChan; firstErr == nil {
			firstErr = err
		}
	}

	if firstErr != nil {
		return firstErr
	}

	return nil
}

func runSequential(runners []Runner) error {
	for _, runner := range runners {
		if err := runner(); err != nil {
			return err
		}
	}

	return nil
}
