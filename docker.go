package main

import (
	"context"
	"time"

	"github.com/efritz/ij/command"
)

func ensureDocker() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := []string{
		"docker",
		"ps",
		"-q",
	}

	_, _, err := command.NewRunner(nil).RunForOutput(
		ctx,
		args,
		nil,
	)

	return err == nil
}
