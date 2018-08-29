package runner

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/state"
)

type PlanRunner struct {
	state *state.State
}

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func NewPlanRunner(state *state.State) *PlanRunner {
	return &PlanRunner{
		state: state,
	}
}

func (r *PlanRunner) Run() bool {
	r.state.Logger.Info(
		nil,
		"Beginning run %s",
		r.state.RunID,
	)

	go r.watchSignals()
	defer r.state.Cleanup.Cleanup()

	transferer, err := paths.NewTransferer(
		r.state.Scratch.Project(),
		r.state.Scratch.Scratch(),
		r.state.Scratch.Workspace(),
		r.state.Config.Excludes,
		r.state.Logger,
	)

	if err != nil {
		r.state.Logger.Error(
			nil,
			"Failed to prepare file import blacklist: %s",
			err.Error(),
		)

		return false
	}

	if err := transferer.Import(r.state.Config.Imports); err != nil {
		r.state.Logger.Error(
			nil,
			"Failed to import files to workspace: %s",
			err.Error(),
		)

		return false
	}

	failure := false
	for _, name := range r.state.Plans {
		if !r.runPlanOrMetaplan(name, logging.NewPrefix(), failure) {
			failure = true
		}
	}

	if failure {
		return false
	}

	if err := transferer.Export(r.state.Config.Exports); err != nil {
		r.state.Logger.Error(
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
		r.state.Logger.Error(
			nil,
			"Received signal",
		)

		r.state.Cancel()
		return
	}
}

func (r *PlanRunner) runPlanOrMetaplan(
	name string,
	prefix *logging.Prefix,
	failure bool,
) bool {
	prefix = prefix.Append(name)

	r.state.Logger.Info(
		prefix,
		"Beginning plan",
	)

	// Stash this so we know if we failed due to a new
	// error or to a failure of a previously failed plan.
	previousFailure := failure

	if plans, ok := r.state.Config.Metaplans[name]; ok {
		for _, plan := range plans {
			if !r.runPlanOrMetaplan(plan, prefix.Append(plan), failure) {
				failure = true
			}
		}
	} else {
		failure = !r.runPlan(name, prefix, failure)
	}

	if failure {
		suffix := ""
		if previousFailure {
			suffix = " (due to previous failure)"
		}

		r.state.Logger.Error(
			prefix,
			"Plan failed%s",
			suffix,
		)
	} else {
		r.state.Logger.Info(
			prefix,
			"Plan completed successfully",
		)
	}

	return !failure
}

func (r *PlanRunner) runPlan(
	name string,
	prefix *logging.Prefix,
	failure bool,
) bool {
	plan := r.state.Config.Plans[name]

	for _, stage := range plan.Stages {
		stagePrefix := prefix.Append(stage.Name)

		if !stage.ShouldRun(failure) {
			r.state.Logger.Info(
				stagePrefix,
				"Skipping stage",
			)

			continue
		}

		if !NewStageRunner(r.state, plan, stage, stagePrefix).Run() {
			r.state.Logger.Error(
				stagePrefix,
				"Stage failed",
			)

			failure = true
			continue
		}

		r.state.Logger.Info(
			stagePrefix,
			"Stage completed successfully",
		)
	}

	return !failure
}
