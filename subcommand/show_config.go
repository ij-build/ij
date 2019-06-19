package subcommand

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/options"
)

func NewShowConfigCommand(appOptions *options.AppOptions) CommandRunner {
	return func(config *config.Config) error {
		serialized, err := yaml.Marshal(config)
		if err != nil {
			return err
		}

		fmt.Println(string(serialized))
		return nil
	}
}
