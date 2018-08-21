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
	"github.com/kballard/go-shellquote"
)

type (
	Runner struct {
		config       *config.Config
		logProcessor logging.Processor
		logger       logging.Logger
		env          []string
		runID        string
		ctx          context.Context
		cancel       func()
		cleanup      *Cleanup
		buildDir     *BuildDir
		workspace    *Workspace
	}

	RunStatus  int
	RunnerFunc func() bool
)

const (
	RunStatusSuccess RunStatus = iota
	RunStatusSetupFailure
	RunStatusFailure
)

func NewRunner(runID string, buildDir *BuildDir, logProcessor logging.Processor) *Runner {
	ctx, cancel := context.WithCancel(context.Background())

	return &Runner{
		runID:        runID,
		buildDir:     buildDir,
		logProcessor: logProcessor,
		ctx:          ctx,
		cancel:       cancel,
		cleanup:      NewCleanup(),
	}
}

func (r *Runner) Shutdown() {
	r.cancel()
	r.cleanup.Cleanup()
}

func (r *Runner) Run(configPath string, plans []string, env []string) RunStatus {
	if err := r.loadConfig(configPath); err != nil {
		logging.EmergencyLog(
			"error: failed to load config: %s",
			err.Error(),
		)

		return RunStatusSetupFailure
	}

	for _, name := range plans {
		if _, ok := r.config.Plans[name]; !ok {
			logging.EmergencyLog(
				"error: unknown plan %s",
				name,
			)

			return RunStatusSetupFailure
		}
	}

	_, err := command.RunForOutput(r.ctx, []string{
		"docker",
		"ps",
		"-q",
	})

	if err != nil {
		logging.EmergencyLog("error: docker is not running")
		return RunStatusSetupFailure
	}

	if err := r.setupLogger(); err != nil {
		logging.EmergencyLog(
			"error: failed to create log files: %s",
			err.Error(),
		)

		return RunStatusSetupFailure
	}

	// TODO - if failed before this point, delete build dir

	r.logger.Info(
		nil,
		"Beginning run %s",
		r.runID,
	)

	if err := r.setupWorkspace(); err != nil {
		r.logger.Error(
			nil,
			"error: failed to create workspace volume: %s",
			err.Error(),
		)

		return RunStatusFailure
	}

	r.env = env

	for _, name := range plans {
		prefix := logging.NewPrefix(name)

		if !r.runPlan(r.config.Plans[name], prefix) {
			r.logger.Error(
				prefix,
				"Plan failed",
			)

			return RunStatusFailure
		}

		r.logger.Info(
			prefix,
			"Plan completed successfully",
		)
	}

	return RunStatusSuccess
}

func (r *Runner) setupLogger() error {
	outfile, errfile, err := r.buildDir.MakeLogFiles("pvc")
	if err != nil {
		return err
	}

	r.logger = r.logProcessor.Logger(outfile, errfile)
	return nil
}

func (r *Runner) loadConfig(configPath string) error {
	config, err := loader.NewLoader().Load(configPath)
	if err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		return err
	}

	r.config = config
	return nil
}

func (r *Runner) setupWorkspace() error {
	workspace := NewWorkspace(r.runID, r.ctx, r.logger)
	if err := workspace.Setup(); err != nil {
		return err
	}

	r.workspace = workspace
	r.cleanup.Register(workspace.Teardown)
	return nil
}

func (r *Runner) runPlan(
	plan *config.Plan,
	prefix *logging.Prefix,
) bool {
	r.logger.Info(
		prefix,
		"Beginning plan",
	)

	for _, stage := range plan.Stages {
		stagePrefix := prefix.Append(stage.Name)

		if !r.runStage(plan, stage, stagePrefix) {
			r.logger.Error(
				stagePrefix,
				"Stage failed",
			)

			return false
		}

		r.logger.Info(
			stagePrefix,
			"Stage completed successfully",
		)
	}

	return true
}

func (r *Runner) runStage(
	plan *config.Plan,
	stage *config.Stage,
	prefix *logging.Prefix,
) bool {
	r.logger.Info(
		prefix,
		"Beginning stage",
	)

	runnerFuncs := []RunnerFunc{}
	for i, st := range stage.Tasks {
		var (
			index     = i
			stageTask = st
			task      = r.config.Tasks[stageTask.Name]
		)

		runner := func() bool {
			taskPrefix := prefix.Append(fmt.Sprintf(
				"%s.%d",
				task.Name,
				index,
			))

			success := r.runTask(
				plan,
				stage,
				task,
				index,
				taskPrefix,
				environment.Merge(
					environment.New(r.config.Environment),
					environment.New(task.Environment),
					environment.New(plan.Environment),
					environment.New(stage.Environment),
					environment.New(stageTask.Environment),
					environment.New(r.env),
				),
			)

			if !success {
				r.logger.Info(
					taskPrefix,
					"Task has failed",
				)
			} else {
				r.logger.Info(
					taskPrefix,
					"Task has completed successfully",
				)
			}

			return success
		}

		runnerFuncs = append(runnerFuncs, runner)
	}

	if stage.Concurrent {
		return runConcurrent(runnerFuncs)
	}

	return runSequential(runnerFuncs)
}

func (r *Runner) runTask(
	plan *config.Plan,
	stage *config.Stage,
	task *config.Task,
	index int,
	prefix *logging.Prefix,
	env environment.Environment,
) bool {
	r.logger.Info(
		prefix,
		"Beginning task",
	)

	if ok, missing := util.ContainsAll(env.Keys(), task.RequiredEnvironment); !ok {
		r.logger.Error(
			prefix,
			"Missing environment values: %s",
			strings.Join(missing, ", "),
		)

		return false
	}

	args, err := r.buildTaskCommandArgs(task, env)
	if err != nil {
		r.logger.Error(
			prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	r.logger.Debug(
		prefix,
		"Running command: %s",
		strings.Join(args, " "),
	)

	outfile, errfile, err := r.buildDir.MakeLogFiles(
		prefix.Serialize(nil),
	)

	if err != nil {
		r.logger.Error(
			prefix,
			"Failed to create task run log files: %s",
			err.Error(),
		)

		return false
	}

	commandErr := command.Run(
		context.Background(),
		prefix,
		args,
		r.logProcessor.Logger(
			outfile,
			errfile,
		),
	)

	if commandErr != nil {
		r.logger.Error(
			prefix,
			"Command failed: %s",
			commandErr.Error(),
		)

		return false
	}

	return true
}

func (r *Runner) buildTaskCommandArgs(
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
		path, err := r.buildDir.WriteScript(task.Script)
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

func runConcurrent(runnerFuncs []RunnerFunc) bool {
	okChan := make(chan bool)
	defer close(okChan)

	for _, runnerFunc := range runnerFuncs {
		go func(runnerFunc RunnerFunc) {
			okChan <- runnerFunc()
		}(runnerFunc)
	}

	success := true
	for i := 0; i < len(runnerFuncs); i++ {
		if ok := <-okChan; !ok {
			success = false
		}
	}

	return success
}

func runSequential(runnerFuncs []RunnerFunc) bool {
	for _, runnerFunc := range runnerFuncs {
		if !runnerFunc() {
			return false
		}
	}

	return true
}
