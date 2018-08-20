package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/efritz/pvc/paths"
)

var (
	app        = kingpin.New("pvc", "").Version(Version)
	plans      = app.Arg("plans", "").Required().Strings()
	configPath = app.Flag("filename", "").Short('o').String()
	env        = app.Flag("env", "").Short('e').Strings()
	verbose    = app.Flag("verbose", "").Short('v').Default("false").Bool()
	colorize   = app.Flag("colorize", "").Default("true").Bool()

	defaultConfigPaths = []string{
		"pvc.yaml",
		"pvc.yml",
	}
)

func parseArgs() error {
	args := os.Args[1:]

	if _, err := app.Parse(args); err != nil {
		return err
	}

	if *configPath == "" {
		for _, path := range defaultConfigPaths {
			ok, err := paths.FileExists(path)
			if err != nil {
				return err
			}

			if ok {
				*configPath = path
				break
			}
		}
	}

	return nil
}
