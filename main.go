package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/efritz/pvc/loader"
	"github.com/efritz/pvc/logging"
	"github.com/efritz/pvc/runtime"
	"github.com/google/uuid"
)

const Version = "0.1.0"

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := parseArgs(); err != nil {
		return err
	}

	config, err := loader.LoadFile(*configPath)
	if err != nil {
		return err
	}

	processor := logging.NewProcessor()
	processor.Start()
	defer processor.Shutdown()

	raw, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	runtime := runtime.NewRuntime(
		raw.String(),
		config,
		processor,
		*env,
	)

	if err := runtime.Setup(); err != nil {
		return err
	}

	go watchSignals(runtime)
	defer runtime.Shutdown()

	if err := runtime.Run(*plans); err != nil {
		return err
	}

	return nil
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
