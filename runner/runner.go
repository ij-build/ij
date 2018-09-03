package runner

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/state"
)

type Runner struct {
	state *state.State
	plans []string
}

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func NewRunner(state *state.State, plans []string) *Runner {
	return &Runner{
		state: state,
		plans: plans,
	}
}

func (r *Runner) Run() bool {
	r.state.Logger.Info(
		nil,
		"Beginning run %s",
		r.state.RunID,
	)

	go r.watchSignals()
	defer r.state.Cleanup.Cleanup()

	defer func() {
		r.state.Logger.Info(
			nil,
			"Finished run %s",
			r.state.RunID,
		)
	}()

	transferer := paths.NewTransferer(
		r.state.Scratch.Project(),
		r.state.Scratch.Scratch(),
		r.state.Scratch.Workspace(),
		r.state.Logger,
	)

	r.state.Logger.Info(nil, "Importing files to workspace")

	importErr := transferer.Import(
		r.state.Config.Import.Files,
		r.state.Config.Import.Excludes,
	)

	if importErr != nil {
		r.state.Logger.Error(
			nil,
			"Failed to import files to workspace: %s",
			importErr.Error(),
		)

		return false
	}

	failure := false
	for _, name := range r.plans {
		result := NewPlanRunner(r.state).Run(name, logging.NewPrefix(), &RunContext{
			Failure:     failure,
			Environment: environment.New(nil),
		})

		if !result {
			failure = true
		}
	}

	if failure {
		return false
	}

	r.state.Logger.Info(nil, "Exporting files from workspace")

	exportErr := transferer.Export(
		r.state.Config.Export.Files,
		r.state.Config.Export.Excludes,
	)

	if exportErr != nil {
		r.state.Logger.Error(
			nil,
			"Failed to export files from workspace: %s",
			exportErr.Error(),
		)

		return false
	}

	return true
}

func (r *Runner) watchSignals() {
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
