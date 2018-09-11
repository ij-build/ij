package options

import "github.com/efritz/ij/logging"

type AppOptions struct {
	ProjectDir  string
	ScratchRoot string
	Colorize    bool
	ConfigPath  string
	Env         []string
	EnvFiles    []string
	Quiet       bool
	Verbose     bool
	FileFactory logging.FileFactory
}
