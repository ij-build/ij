package main

import (
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	app        = kingpin.New("pvc", "").Version(Version)
	plans      = app.Arg("plans", "").Required().Strings()
	configPath = app.Flag("filename", "").Short('o').String()
	env        = app.Flag("env", "").Short('e').Strings()

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
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			*configPath = path
			break
		}
	}

	return nil
}
