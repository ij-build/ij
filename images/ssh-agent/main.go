package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/efritz/ij/consts"
	"github.com/efritz/ij/ssh"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func doMain() error {
	identities := []string{}
	app := kingpin.New("ij-ensure-keys-available", "Ensure ssh keys are in the containerized ssh-agent.").Version(consts.Version)
	app.Flag("ssh-identity", "Enable ssh-agent for the given identities.").StringsVar(&identities)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		return err
	}

	_, err := ssh.EnsureKeysAvailable(identities)
	return err
}
