package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/efritz/pvc/command"
	"github.com/efritz/pvc/config"
	"github.com/efritz/pvc/environment"
	"github.com/efritz/pvc/loader"
	"github.com/efritz/pvc/logging"
	"github.com/efritz/pvc/util"
	shellquote "github.com/kballard/go-shellquote"
)

type (
	Runtime struct {
		config       *config.Config
		logProcessor logging.Processor
		logger       logging.Logger
		env          []string
		runID        string
		ctx          context.Context
		cancel       func()
		cleanup      *Cleanup
		builddir     *Builddir
		workspace    *Workspace
	}

	Runner func() error
)

func NewRuntime(logProcessor logging.Processor) *Runtime {
	ctx, cancel := context.WithCancel(context.Background())

	return &Runtime{
		logProcessor: logProcessor,
		ctx:          ctx,
		cancel:       cancel,
		cleanup:      NewCleanup(),
	}
}

func (r *Runtime) Setup() error {
	setupFuncs := []func() error{
		r.setupRunID,
		r.setupBuilddir,
		r.setupLogger,
	}

	for _, f := range setupFuncs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) setupRunID() error {
	runID, err := util.MakeID()
	if err != nil {
		return err
	}

	r.runID = runID
	return nil
}

func (r *Runtime) setupBuilddir() error {
	builddir := NewBuilddir(r.runID)

	if err := builddir.Setup(); err != nil {
		return err
	}

	r.builddir = builddir
	r.cleanup.Register(builddir.Teardown)
	return nil
}

func (r *Runtime) setupLogger() error {
	outfile, errfile, err := r.builddir.MakeLogFiles("pvc")
	if err != nil {
		return err
	}

	r.logger = r.logProcessor.Logger("pvc", outfile, errfile)
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

func (r *Runtime) Shutdown() {
	r.cancel()
	r.cleanup.Cleanup()
}

func (r *Runtime) Run(configPath string, plans []string, env []string) bool {
	config, err := loader.LoadFile(configPath)
	if err != nil {
		r.logger.Error("error: %s", err.Error())
		return false
	}

	r.env = env
	r.config = config

	if err := r.setupWorkspace(); err != nil {
		r.logger.Error("error: %s", err.Error())
		return false
	}

	for _, name := range plans {
		plan, ok := r.config.Plans[name]
		if !ok {
			r.logger.Error("plan %s not found", name)
			return false
		}

		if err := r.runPlan(plan); err != nil {
			r.logger.Error("error: %s", err.Error())
			return false
		}
	}

	return true
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
	for index, taskInstance := range stage.Tasks {
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

		stableIndex := index

		runner := func() error {
			return r.runTask(
				plan,
				stage,
				task,
				stableIndex,
				env,
			)
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
	if ok, missing := util.ContainsAll(env.Keys(), task.RequiredEnvironment); !ok {
		panic(fmt.Sprintf("missing environment values: %#v", missing))
	}

	args, err := r.buildTaskCommandArgs(task, env)
	if err != nil {
		return err
	}

	parts := []string{
		plan.Name,
		stage.Name,
		fmt.Sprintf("%d.%s", index, task.Name),
	}

	prefix := strings.Join(parts, "/")

	outfile, errfile, err := r.builddir.MakeLogFiles(prefix)
	if err != nil {
		return err
	}

	return command.Run(
		context.Background(),
		args,
		r.logProcessor.Logger(
			prefix,
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

	if task.Script != "" {
		path, err := r.builddir.WriteScript(task.Script)
		if err != nil {
			return nil, err
		}

		shell := task.Shell
		if shell == "" {
			shell = "/bin/sh"
		}

		args = append(
			args,
			"-v",
			path+":/workspace/script", // TODO - something more unique
			"--entrypoint",
			shell,
			task.Image,
			"/workspace/script",
		)
	} else {
		commandArgs, err := shellquote.Split(task.Command)
		if err != nil {
			return nil, err
		}

		args = append(args, task.Image)
		args = append(args, commandArgs...)
	}

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
