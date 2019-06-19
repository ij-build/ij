package options

import "github.com/ij-build/ij/logging"

type AppOptions struct {
	ProjectDir   string
	ScratchRoot  string
	ConfigPath   string
	Env          []string
	EnvFiles     []string
	Quiet        bool
	Verbose      bool
	DisableColor bool
	FileFactory  logging.FileFactory
}
