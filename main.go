package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/pvc/logging"
	"github.com/efritz/pvc/runtime"
	"github.com/efritz/pvc/util"
)

const Version = "0.1.0"

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	if !run() {
		os.Exit(1)
	}
}

func run() bool {
	if err := parseArgs(); err != nil {
		logging.EmergencyLog("error: %s", err.Error())
		return false
	}

	runID, err := util.MakeID()
	if err != nil {
		logging.EmergencyLog(
			"error: failed to generate run id: %s",
			err.Error(),
		)

		return false
	}

	buildDir := runtime.NewBuildDir(runID)
	if err := buildDir.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create build directory: %s",
			err.Error(),
		)

		return false
	}

	processor := logging.NewProcessor(*verbose, *colorize)
	processor.Start()
	defer processor.Shutdown()

	runner := runtime.NewRunner(runID, buildDir, processor)

	go watchSignals(runner)
	defer runner.Shutdown()

	status := runner.Run(*configPath, *plans, *env)

	if status == runtime.RunStatusSetupFailure {
		if err := buildDir.Teardown(); err != nil {
			logging.EmergencyLog(
				"error: failed to teardown build directory: %s",
				err.Error(),
			)
		}

		return false
	}

	if err := buildDir.Prune(); err != nil {
		logging.EmergencyLog(
			"error: failed to prune build directory: %s",
			err.Error(),
		)
	}

	return status == runtime.RunStatusSuccess
}

func watchSignals(runtime *runtime.Runner) {
	signals := make(chan os.Signal, 1)

	for _, s := range shutdownSignals {
		signal.Notify(signals, s)
	}

	for range signals {
		runtime.Shutdown()
		os.Exit(1)
	}
}
