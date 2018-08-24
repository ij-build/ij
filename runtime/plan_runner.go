package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
)

type PlanRunner struct {
	state *State
}

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func NewPlanRunner(state *State) *PlanRunner {
	return &PlanRunner{
		state: state,
	}
}

func (r *PlanRunner) Run() bool {
	r.state.logger.Info(
		nil,
		"Beginning run %s",
		r.state.runID,
	)

	go r.watchSignals()
	defer r.state.cleanup.Cleanup()

	transferer, err := paths.NewTransferer(
		r.state.scratch.Project(),
		r.state.scratch.Scratch(),
		r.state.scratch.Workspace(),
		r.state.config.Excludes,
		r.state.logger,
	)

	if err != nil {
		r.state.logger.Error(
			nil,
			"Failed to prepare file import blacklist: %s",
			err.Error(),
		)

		return false
	}

	if err := transferer.Import(r.state.config.Imports); err != nil {
		r.state.logger.Error(
			nil,
			"Failed to import files to workspace: %s",
			err.Error(),
		)

		return false
	}

	for _, name := range r.state.plans {
		if !r.runPlanOrMetaplan(name, logging.NewPrefix()) {
			return false
		}
	}

	if err := transferer.Export(r.state.config.Exports); err != nil {
		r.state.logger.Error(
			nil,
			"Failed to export files from workspace: %s",
			err.Error(),
		)

		return false
	}

	return true
}

func (r *PlanRunner) watchSignals() {
	signals := make(chan os.Signal, 1)

	for _, s := range shutdownSignals {
		signal.Notify(signals, s)
	}

	for range signals {
		r.state.logger.Error(
			nil,
			"Received signal",
		)

		r.state.cancel()
		return
	}
}

func (r *PlanRunner) runPlanOrMetaplan(
	name string,
	prefix *logging.Prefix,
) bool {
	prefix = prefix.Append(name)

	r.state.logger.Info(
		prefix,
		"Beginning plan",
	)

	ok := true

	if plans, ok := r.state.config.Metaplans[name]; ok {
		for _, plan := range plans {
			if !r.runPlanOrMetaplan(plan, prefix.Append(plan)) {
				ok = false
				break
			}
		}
	} else {
		ok = r.runPlan(name, prefix)
	}

	if !ok {
		r.state.logger.Error(
			prefix,
			"Plan failed",
		)

		return false
	}

	r.state.logger.Info(
		prefix,
		"Plan completed successfully",
	)

	return true
}

func (r *PlanRunner) runPlan(
	name string,
	prefix *logging.Prefix,
) bool {
	plan := r.state.config.Plans[name]

	for _, stage := range plan.Stages {
		stagePrefix := prefix.Append(stage.Name)

		if !NewStageRunner(r.state, plan, stage, stagePrefix).Run() {
			r.state.logger.Error(
				stagePrefix,
				"Stage failed",
			)

			return false
		}

		r.state.logger.Info(
			stagePrefix,
			"Stage completed successfully",
		)
	}

	return true
}
