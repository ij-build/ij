package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/loader"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/runtime"
	"github.com/efritz/ij/util"
)

const Version = "0.1.0"

var (
	app                    = kingpin.New("ij", "").Version(Version)
	plans                  = app.Arg("plans", "").Default("default").Strings()
	configPath             = app.Flag("filename", "").Short('f').String()
	env                    = app.Flag("env", "").Short('e').Strings()
	verbose                = app.Flag("verbose", "").Short('v').Default("false").Bool()
	colorize               = app.Flag("color", "").Default("true").Bool()
	forceSequential        = app.Flag("force-sequential", "").Default("false").Bool()
	rawHealthcheckInterval = app.Flag("healthcheck-interval", "").Default("5s").String()
	healthcheckInterval    time.Duration

	defaultConfigPaths = []string{
		"ij.yaml",
		"ij.yml",
	}
)

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

	config, ok := loadConfig()
	if !ok {
		return false
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)

	defer cancel()

	if !ensureDocker(ctx) {
		logging.EmergencyLog("error: docker is not running")
		return false
	}

	state, err := runtime.NewState(
		config,
		*plans,
		*env,
		*verbose,
		*colorize,
		*forceSequential,
		healthcheckInterval,
	)

	if err != nil {
		return false
	}

	return runtime.NewPlanRunner(state).Run()
}

func parseArgs() error {
	args := os.Args[1:]

	if _, err := app.Parse(args); err != nil {
		return err
	}

	if *configPath == "" {
		for _, path := range defaultConfigPaths {
			ok, err := util.FileExists(path)
			if err != nil {
				return err
			}

			if ok {
				*configPath = path
				break
			}
		}
	}

	parsed, err := time.ParseDuration(*rawHealthcheckInterval)
	if err != nil {
		return fmt.Errorf(
			"illegal healthcheck interval '%s'",
			*rawHealthcheckInterval,
		)
	}

	healthcheckInterval = parsed
	return nil
}

func loadConfig() (*config.Config, bool) {
	config, err := loader.NewLoader().Load(*configPath)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to load config: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := config.Validate(); err != nil {
		logging.EmergencyLog(
			"error: failed to validate config: %s",
			err.Error(),
		)

		return nil, false
	}

	for _, name := range *plans {
		if _, ok := config.Plans[name]; !ok {
			logging.EmergencyLog(
				"error: unknown plan %s",
				name,
			)

			return nil, false
		}
	}

	return config, true
}

func ensureDocker(ctx context.Context) bool {
	_, _, err := command.RunForOutput(
		ctx,
		[]string{
			"docker",
			"ps",
			"-q",
		},
		nil,
	)

	return err == nil
}
