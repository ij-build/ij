package options

import "github.com/efritz/ij/logging"

type AppOptions struct {
	Colorize    bool
	ConfigPath  string
	Env         []string
	EnvFiles    []string
	Quiet       bool
	Verbose     bool
	FileFactory logging.FileFactory
}
