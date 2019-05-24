# SSH Agent Container

This Docker image starts an ssh-agent and makes its socket available in a volume. The following three actions are available within this image.

`/ij/start-agent.sh` is the entrypoint and will start an ssh-agent server at `/tmp/ij/ssh-agent/ssh-agent.sock`. The containing directory is a container volume, which can be mounted from other containers running on the same host.

`/ij/add-keys.sh` will move the contents of the `~/.ssh` directory into a scratch space, ensure that the permissions of the directory and the private key files are appropriate, then add each matching key to the running SSH agent. The movement of keys are required when the `~/.ssh` directory is a Windows volume mount (the permissions cannot be changed otherwise and are mounted with too-open permissions incompatible with the `ssh-add` command).

`/ij/ij-ensure-keys-available` is a binary wrapper around IJ's `ssh` package behavior. This ensure that the keys available in the SSH agent and the keys supplied in the config file/command line intersect.
