package subcommand

import (
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/options"
	"github.com/ghodss/yaml"
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
