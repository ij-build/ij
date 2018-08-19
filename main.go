package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/pvc/logging"
	"github.com/efritz/pvc/runtime"
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
		fmt.Printf("error: %s\n", err.Error())
		return false
	}

	processor := logging.NewProcessor()
	processor.Start()
	defer processor.Shutdown()

	runtime := runtime.NewRuntime(processor)

	if err := runtime.Setup(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return false
	}

	go watchSignals(runtime)
	defer runtime.Shutdown()

	return runtime.Run(*configPath, *plans, *env)
}

func watchSignals(runtime *runtime.Runtime) {
	signals := make(chan os.Signal, 1)

	for _, s := range shutdownSignals {
		signal.Notify(signals, s)
	}

	for range signals {
		runtime.Shutdown()
		os.Exit(1)
	}
}
