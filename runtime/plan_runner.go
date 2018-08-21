package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/logging"
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

	for _, name := range r.state.plans {
		prefix := logging.NewPrefix(name)

		if !r.runPlan(r.state.config.Plans[name], prefix) {
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

func (r *PlanRunner) runPlan(
	plan *config.Plan,
	prefix *logging.Prefix,
) bool {
	r.state.logger.Info(
		prefix,
		"Beginning plan",
	)

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
