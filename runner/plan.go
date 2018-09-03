package runner

import (
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type PlanRunner struct {
	state *state.State
}

func NewPlanRunner(state *state.State) *PlanRunner {
	return &PlanRunner{
		state: state,
	}
}

func (r *PlanRunner) Run(
	name string,
	prefix *logging.Prefix,
	context *RunContext,
) bool {
	prefix = prefix.Append(name)

	r.state.Logger.Info(
		prefix,
		"Beginning plan",
	)

	failure := context.Failure

	if plans, ok := r.state.Config.Metaplans[name]; ok {
		for _, plan := range plans {
			result := r.Run(plan, prefix, &RunContext{
				Failure:     failure,
				Environment: context.Environment,
			})

			if !result {
				failure = true
			}
		}
	} else {
		failure = !r.runPlan(name, prefix, &RunContext{
			Failure:     failure,
			Environment: context.Environment,
		})
	}

	if failure {
		suffix := ""
		if context.Failure {
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
	context *RunContext,
) bool {
	var (
		plan    = r.state.Config.Plans[name]
		failure = context.Failure
	)

	for _, stage := range plan.Stages {
		stagePrefix := prefix.Append(stage.Name)

		if !stage.ShouldRun(context.Failure) {
			r.state.Logger.Info(
				stagePrefix,
				"Skipping stage",
			)

			continue
		}

		runner := NewStageRunner(r.state, plan, stage, stagePrefix)
		newContext := &RunContext{
			Failure:     failure,
			Environment: context.Environment,
		}

		if !runner.Run(newContext) {
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
