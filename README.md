<p align="center">
    <img width="100" src="https://github.com/ij-build/ij/blob/master/ij.svg" alt="IJ logo">
</p>

<p align="center">
    <a href="https://godoc.org/github.com/ij-build/ij"><img src="https://godoc.org/github.com/ij-build/ij?status.svg" alt="GoDoc"></a>
    <a href="http://travis-ci.org/ij-build/ij"><img src="https://secure.travis-ci.org/ij-build/ij.png" alt="Build Status"></a>
    <a href="https://codeclimate.com/github/ij-build/ij/maintainability"><img src="https://api.codeclimate.com/v1/badges/63b7e45a56b21d361a62/maintainability" alt="Maintainability"></a>
    <a href="https://codeclimate.com/github/ij-build/ij/test_coverage"><img src="https://api.codeclimate.com/v1/badges/63b7e45a56b21d361a62/test_coverage" alt="Test Coverage"></a>
</p>

IJ is a build tool using Docker containers.

## Concepts

A build is defined by a sequence of [tasks](https://github.com/ij-build/ij/blob/master/docs/tasks.md#user-content-tasks) which perform some unit of work (e.g. building a binary, running integration tests, pushing artifacts to a remote server), usually within a docker container. This ensures that the development host stays clean, and that the exact build dependencies are available within the container (no more development tool version hell). The build process is defined by assembling a sequence of tasks into a [plan](https://github.com/ij-build/ij/blob/master/docs/plans.md#user-content-plans). A build is simply the invocation of one or more plans. A project's [config file](https://github.com/ij-build/ij/blob/master/docs/config.md#user-content-config) declares tasks and plans which builds the project.

## Installation

Simply run `go install github.com/ij-build/ij`.

## Usage

There are currently four IJ subcommands (`run`, `login`, `logout`, `clean`, and `rotate-logs`) each discussed below. The following command line flags are applicable for all IJ commands.

| Name     | Short Flag | Description |
| -------- | ---------- | ----------- |
| config   | f          | The path to the config file. If not supplied, `ij.yaml` and `ij.yml` are attempted in the current directory. |
| env      | e          | Set an environment variable. Use `-e VAR=VAL` to set an explicit value for the variable `VAR`. Use `-e VAR` to use the host value of `$VAR`. |
| env-file |            | The path to an [environment file](https://github.com/ij-build/ij/blob/master/docs/environment.md#user-content-environment-files). |
| no-color |            | Disable colorized output. |
| verbose  | v          | Show debug-level output. |

### Run Command

This command can be invoked as `ij [run]? (plan-name)*`. The run keyword is assumed if not supplied. This command runs a series of plans or metaplans defined in the config file. If no plan-names are supplied, the plan named `default` is invoked.

| Name                 | Short Flag | Description |
| -------------------- | ---------- | ----------- |
| cpu-shares           | c          | The proc limit for run task containers. |
| force-sequential     |            | Disable running tasks in parallel. |
| healthcheck-interval |            | How frequently to check the health of service containers. |
| keep-workspace       | k          | Do not prune the scratch directory (useful for debugging failed plans). |
| login                |            | Login to registries before invoking plans and logout from registries after (useful for builds that push image artifacts). |
| memory               | m          | The memory limit for run task containers. |
| ssh-identity         |            | An additional SSH key fingerprint required to be present in the host's SSH agent. |
| ssh-agent-container  |            | Mount your `~/.ssh` directory into a container and start an ssh-agent. This is required for using SSH keys on Windows. |
| timeout              |            | The maximum time a build plan can run in total. |

### Login Command

This command can be invoked as `ij login`. Login to all [registries](https://github.com/ij-build/ij/blob/master/docs/registries.md#user-content-registries) defined in the config file.

### Logout Command

This command can be invoked as `ij logout`. Logout from all [registries](https://github.com/ij-build/ij/blob/master/docs/registries.md#user-content-registries) defined in the config file.

### Rotate Logs Command

This command can be invoked as `ij rotate-logs`. Remove all but the most recent run from the `.ij` directory in the current project.

### Clean Command

This command can be invoked as `ij clean`. Remove files exported from the workspace on previous runs.

| Name                 | Short Flag | Description |
| -------------------- | ---------- | ----------- |
| --force              |            | Do not prompt before removing files or directories. |

### Show Config Command

This command cna be invoked as `ij show-config`. This will print the effective config after resolving inheritance and extension. This output of this command, if successful, should also be another valid config file.

## License

This project is a spiritual and technical successor to [rw_grim/convey](https://bitbucket.org/rw_grim/convey). All source in this project is licensed under GPLv3+.
