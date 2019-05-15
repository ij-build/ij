package options

import "github.com/efritz/ij/logging"

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
