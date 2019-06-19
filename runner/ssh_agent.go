package runner

import (
	"context"
	"fmt"
	"os/user"

	"github.com/ij-build/ij/command"
	"github.com/ij-build/ij/logging"
	"github.com/ij-build/ij/scratch"
)

const (
	SSHAgentImage    = "efritz/ij-ssh-agent:latest"
	SocketVolumePath = "/tmp/ij/ssh-agent"
	SocketPath       = "/tmp/ij/ssh-agent/ssh-agent.sock"
)

func startSSHAgent(
	runID string,
	identities []string,
	scratch *scratch.ScratchSpace,
	containerLists *ContainerLists,
	logger logging.Logger,
) error {
	containerName := fmt.Sprintf("%s-ssh-agent", runID)

	if err := startContainer(runID, containerName, scratch, containerLists, logger); err != nil {
		return err
	}

	if err := addKeys(containerName, logger); err != nil {
		command.NewRunner(logger).RunForOutput(
			context.Background(),
			[]string{"docker", "kill", containerName},
			nil,
		)

		return err
	}

	if err := ensureKeys(containerName, identities, logger); err != nil {
		command.NewRunner(logger).RunForOutput(
			context.Background(),
			[]string{"docker", "kill", containerName},
			nil,
		)

		return err
	}

	return nil
}

func startContainer(
	runID string,
	containerName string,
	scratch *scratch.ScratchSpace,
	containerLists *ContainerLists,
	logger logging.Logger,
) error {
	builder, err := sshAgentCommandBuilderFactory(
		runID,
		scratch,
		containerName,
	)

	if err != nil {
		return fmt.Errorf("failed to build command args: %s", err.Error())
	}

	args, _, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build command args: %s", err.Error())
	}

	containerLists.ContainerStopper.Add(containerName)

	_, errOutput, err := command.NewRunner(logger).RunForOutput(
		context.Background(),
		args,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to start ssh-agent container: %s, %s", err.Error(), errOutput)
	}

	return nil
}

func addKeys(containerName string, logger logging.Logger) error {
	_, errOutput, err := command.NewRunner(logger).RunForOutput(
		context.Background(),
		[]string{
			"docker",
			"exec",
			containerName,
			"/ij/add-keys.sh",
		},
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to add ssh-keys: %s, %s", err.Error(), errOutput)
	}

	return nil
}

func ensureKeys(containerName string, identities []string, logger logging.Logger) error {
	identityArgs := []string{}
	for _, identity := range identities {
		identityArgs = append(identityArgs, "--ssh-identity")
		identityArgs = append(identityArgs, identity)
	}

	_, errOutput, err := command.NewRunner(logger).RunForOutput(
		context.Background(),
		append([]string{
			"docker",
			"exec",
			containerName,
			"/ij/ij-ensure-keys-available",
		}, identityArgs...),
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to validate ssh keys: %s, %s", err.Error(), errOutput)
	}

	return nil
}

func sshAgentCommandBuilderFactory(
	runID string,
	scratch *scratch.ScratchSpace,
	containerName string,
) (*command.Builder, error) {
	current, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user (%s)", err.Error())
	}

	builder := command.NewBuilder([]string{
		"docker",
		"run",
		"--rm",
	}, nil)

	builder.AddArgs(SSHAgentImage)
	builder.AddFlagValue("-v", fmt.Sprintf("%s/.ssh:/root/.ssh", current.HomeDir))
	builder.AddFlagValue("--name", containerName)
	builder.AddFlag("-d")
	builder.AddFlagValue("--network", runID)

	return builder, nil
}
