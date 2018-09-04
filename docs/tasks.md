# Tasks

A task is an object representing a unit of work in the run of a build plan. These objects use the type property to determine the other available properties. Each of the following task types are discussed in the sections below. Every task defines the following properties.

| Name    | Required | Default | Description |
| ------- | -------- | ------- | ----------- |
| extends |          | ''      | The name of the task this task extends (if any). |
| type    |          | run     | The type of task. May also be one of `build`, `push`, `remove`, or `plan`. |

See the section on [extending a task](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-task) about the semantics of the `extends` property. It may be of note that the `extends` property does **not** support environment substitution.

## Run Task

A run task runs a Docker container.

| Name                    | Required | Default    | Description |
| ----------------------- | -------- | ---------- | ----------- |
| command                 |          | ''         | The command to run. If this value contains shell-specific tokens (e.g. chaining, pipes, or redirection), then `script` property should be used instead. |
| detach                  |          | false      | If true, this container is run in the background until container exit or the end of the build plan. |
| entrypoint              |          | ''         | The entrypoint of the container. |
| environment             |          | []         | A list additional environment variable definitions. Value may be a string or a list. |
| export_environment_file |          | ''         | The path (relative to the working directory) to the file where exported environment variables are written. |
| healthcheck             |          | {}         | A [healthcheck configuration object](https://github.com/efritz/ij/blob/master/docs/tasks.md#user-content-healthcheck-configuration). |
| hostname                |          | ''         | The container's network alias. |
| image                   | yes      |            | The name of the image to run. |
| required_environment    |          | []         | A list of environment variable names which MUST be defined as non-empty for this task to run. |
| script                  |          | ''         | Lke the `command` property, but supports multi-line strings and shell features. |
| shell                   |          | /bin/sh    | The shell used to invoke the supplied script. |
| user                    |          | ''         | The username to invoke the command or script under. |
| workspace               |          | /workspace | The working directory within the container. If a global value is set, that is used as a fallback before using the default. |

Some of these properties work only in tandem (or mutually exclusively) with other properties:

- `shell` is useful only when `script` is supplied
- `entrypoint` is useful only when `script` is absent
- `healthcheck` parameters are only useful *in convey* when `detach` is true (but will still affect external `docker inspect` commands)
- `export_environment_file` is useful only when `detach` is false

This task will run containers in the foreground (blocking until the container exits) unless `detach` is set to true. The task succeeds if the container exits with a zero status. When `detach` is set to true, the container will be run in the background. If the container defines a healthcheck (either via Dockerfile or the task healthcheck configuration defined below), the task will block until the container becomes healthy. The task succeeds if the container becomes healthy.

The file referenced by `export_environment_file` should contain lines of the form `VAR=VAL`. Whitespace and `#` comments ignored. Each of these lines will be added to the working environment set made available to tasks in future stages in the same run.

### Healthcheck Configuration

Supplying any of the following parameters will overwrite any healthcheck defined in the Dockerfile used to build the running image. For details on how these properties affect a running container, see [the Docker documentation](https://docs.docker.com/engine/reference/builder/#healthcheck).

| Name         | Required | Default | Description |
| ------------ | -------- | ------- | ----------- |
| command      |          |         | The command to exec in the container. |
| interval     |          |         | The duration between health checks. |
| retries      |          | 0       | The number of times to check an unhealthy container before failing. |
| start_period |          |         | The duration after container startup in which failed health checks are not counted against the retry count. |
| timeout      |          |         | The maximum runtime of a single health check. |

## Build Task

A build task builds a Docker image from a Dockerfile.

| Name       | Required | Default    | Description |
| ---------- | -------- | ---------- | ----------- |
| dockerfile |          | Dockerfile | The path to the Dockerfile on the host. |
| labels     |          | []         | Metadata for the resulting image. Value may be a string or a list. |
| tags       |          | []         | A list of tags for the resulting image. Value may be a string or a list. |

## Push Task

A push task pushes image tags to a remote registry. For this task to succeed, the target registry must be writable by the current host and user. This may require previously running `ij login` or invoking this plan with the `--login` option.

| Name   | Required | Default | Description |
| ------ | -------- | ------- | ----------- |
| images |          | []      | A list of image tags to push to a remote registry. Value may be a string or a list. |

## Remove Task

A remove task removes image from the host.

| Name   | Required | Default | Description |
| ------ | -------- | ------- | ----------- |
| images |          | []      | A list of image tags to remove from the host. Value may be a string or a list. |

## Plan Task

A plan task (recursively) invokes a plan or a metaplan defined in the same configuration.

| Name        | Required | Default | Description |
| ----------- | -------- | ------- | ----------- |
| environment |          | []      | A list additional environment variable definitions. Value may be a string or a list. |
| name        | yes      |         | The name of the plan or metaplan to invoke. |

It may be of note that the `name` property does **not** support environment substitution.
