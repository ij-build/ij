// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package runner

import (
	"fmt"
	"os/user"

	"github.com/ij-build/ij/command"
)

const FlashPermissionsImage = "alpine:3.8"

func (r *Runner) tryFlashPermissions() {
	r.logger.Info(
		nil,
		"Flashing workspace permissions",
	)

	if err := r.flashPermissions(); err != nil {
		r.logger.Error(
			nil,
			"Failed to flash workspace permissions: %s",
			err.Error(),
		)
	}
}

func (r *Runner) flashPermissions() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	builder := command.NewBuilder([]string{
		"docker",
		"run",
		"--rm",
	}, nil)

	builder.AddArgs(FlashPermissionsImage)
	builder.AddArgs("chown", fmt.Sprintf("%s:%s", user.Uid, user.Gid), "-R", ".")
	builder.AddFlagValue("-w", "/workspace")
	builder.AddFlagValue("-v", fmt.Sprintf("%s:/workspace", r.scratch.Workspace()))

	args, _, err := builder.Build()
	if err != nil {
		return err
	}

	return command.NewRunner(r.logger).Run(
		r.ctx,
		args,
		nil,
		nil,
	)
}
