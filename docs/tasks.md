# Tasks

A task is an object representing a unit of work in the run of a build plan. These objects use the type property to determine the other available properties. Each of the following task types are discussed in the sections below. Every task defines the following properties.

| Name                 | Required | Default | Description |
| -------------------- | -------- | ------- | ----------- |
| extends              |          | ''      | The name of the task this task extends (if any). |
| type                 |          | run     | The type of task. May also be one of `build`, `push`, `remove`, or `plan`. |
| environment          |          | []      | A list of environment variable definitions. Value may be a string or a list. |
| required_environment |          | []      | A list of environment variable names which MUST be defined as non-empty for this task to run. |

See the section on [extending a task](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-task) about the semantics of the `extends` property. It may be of note that the `extends` property does **not** support environment expansion.

The `type` property is not always required when not the default value -- if the `extends` property is set and the `type` property is not, then the value of the `type` property is inferred by type of the parent task. It is an error to supply both the `type` and `extends` property in an inconsistent manner (it is not possible to extend a task of a different type).

## Run Task

A run task runs a Docker container.

| Name                    | Required | Default    | Description |
| ----------------------- | -------- | ---------- | ----------- |
| command                 |          | ''         | The command to run. If this value contains shell-specific tokens (e.g. chaining, pipes, or redirection), then `script` property should be used instead. |
| detach                  |          | false      | If true, this container is run in the background until container exit or the end of the build plan. |
| entrypoint              |          | ''         | The entrypoint of the container. |
| export_environment_file |          | ''         | The path (relative to the working directory) to the file where exported environment variables are written. |
| healthcheck             |          | {}         | A [healthcheck configuration object](https://github.com/efritz/ij/blob/master/docs/tasks.md#user-content-healthcheck-configuration). |
| hostname                |          | ''         | The container's network alias. |
| image                   | yes      |            | The name of the image to run. |
| script                  |          | ''         | Lke the `command` property, but supports multi-line strings and shell features. |
| shell                   |          | /bin/sh    | The shell used to invoke the supplied script. |
| user                    |          | ''         | The username to invoke the command or script under. |
| workspace               |          |            | The working directory within the container. If a global value is set, that is used as a fallback. |

Some of these properties work only in tandem (or mutually exclusively) with other properties:

- `shell` is useful only when `script` is supplied
- `entrypoint` is useful only when `script` is absent
- `healthcheck` parameters are only useful *in convey* when `detach` is true (but will still affect external `docker inspect` commands)
- `export_environment_file` is useful only when `detach` is false

This task will run containers in the foreground (blocking until the container exits) unless `detach` is set to true. The task succeeds if the container exits with a zero status. When `detach` is set to true, the container will be run in the background. If the container defines a healthcheck (either via Dockerfile or the task healthcheck configuration defined below), the task will block until the container becomes healthy. The task succeeds if the container becomes healthy.

The file referenced by `export_environment_file` should be formatted like an env file as discussed in the documentation on [environments](https://github.com/efritz/ij/blob/master/docs/environment.md#user-content-environment). Each relevant line of the file will be added to the working environment set made available to tasks in future stages in the same run.

### Healthcheck Configuration

Supplying any of the following parameters will overwrite any healthcheck defined in the Dockerfile used to build the running image. For details on how these properties affect a running container, see [the Docker documentation](https://docs.docker.com/engine/reference/builder/#healthcheck).

| Name         | Required | Default | Description |
| ------------ | -------- | ------- | ----------- |
| command      |          |         | The command to exec in the container. |
| interval     |          |         | The duration between health checks. |
| retries      |          | 0       | The number of times to check an unhealthy container before failing. |
| start_period |          |         | The duration after container startup in which failed health checks are not counted against the retry count. |
| timeout      |          |         | The maximum runtime of a single health check. |

### Example

This first example runs the image `${GO_IMAGE}`, defined in the global environment section. It defines a single `script` which adds a GitHub public key to the known hosts file and installs vendor dependencies via, but only if the `vendor` directory was not imported from the host.

```yaml
environment:
  - GO_IMAGE=registry.example.io/devops/go-build:master-latest

tasks:
  glide-install:
    image: ${GO_IMAGE}
    script: |
      if [ -d vendor ]; then
        # Skip if vendor was imported from host
        exit 0;
      fi

      # Install deps from glide.yaml (may include private repos)
      ssh-keyscan -H github.com 2> /dev/null 1> ~/.ssh/known_hosts
      glide install

# plans not shown
```

The second example declares a task to run `redis` in the background, and another task to run the image `api` with an environment pointed to the redis hostname. This example shows building blocks useful for end-to-end integration testing with a live (locally-hosted) database.

```yaml
environment:
  - REDIS_HOST=redis.ij

tasks:
  redis:
    image: redis
    detach: true
    hostname: ${REDIS_HOST}
    healthcheck:
      command: redis-cli ping
      interval: 1s

  api:
    image: api
    environment:
      - API_REDIS_HOST=${REDIS_HOST}

# plans not shown
```

## Build Task

A build task builds a Docker image from a Dockerfile.

| Name       | Required | Default    | Description |
| ---------- | -------- | ---------- | ----------- |
| dockerfile |          | Dockerfile | The path to the Dockerfile on the host. |
| labels     |          | []         | Metadata for the resulting image. Value may be a string or a list. |
| tags       |          | []         | A list of tags for the resulting image. Value may be a string or a list. |
| target     |          |            | The target stage to build in a multi-stage dockerfile. |

### Example

This example tags an image with the project's current git status, and adds the same information plus the time of the build to the image labels.

```yaml
tasks:
  build-api:
    type: build
    dockerfile: Dockerfile.api
    tags:
      - registry.example.io/devops/api:${GIT_BRANCH_NORMALIZED}-latest
      - registry.example.io/devops/api:${GIT_BRANCH_NORMALIZED}-${GIT_COMMIT_SHORT}
    labels:
      - BUILD_TIME=${BUILD_TIME}
      - GIT_BRANCH=${GIT_BRANCH}
      - GIT_COMMIT=${GIT_COMMIT}

# plans not shown
```

## Push Task

A push task pushes image tags to a remote registry. For this task to succeed, the target registry must be writable by the current host and user. This may require previously running `ij login` or invoking this plan with the `--login` option.

| Name   | Required | Default | Description |
| ------ | -------- | ------- | ----------- |
| images |          | []      | A list of image tags to push to a remote registry. Value may be a string or a list. |

### Example

This example speaks for itself.

```yaml
tasks:
  push-images:
    type: push
    images:
      - registry.example.io/devops/api:${GIT_BRANCH_NORMALIZED}-latest
      - registry.example.io/devops/api:${GIT_BRANCH_NORMALIZED}-${GIT_COMMIT_SHORT}

# plans not shown
```

## Remove Task

A remove task removes image from the host.

| Name   | Required | Default | Description |
| ------ | -------- | ------- | ----------- |
| images |          | []      | A list of image tags to remove from the host. Value may be a string or a list. |

### Example

This example speaks for itself.

```yaml
tasks:
  remove-images:
    type: remove
    images:
      - registry.example.io/devops/api:${GIT_BRANCH_NORMALIZED}-latest
      - registry.example.io/devops/api:${GIT_BRANCH_NORMALIZED}-${GIT_COMMIT_SHORT}

# plans not shown
```

## Plan Task

A plan task (recursively) invokes a plan or a metaplan defined in the same configuration.

| Name | Required | Default | Description |
| ---- | -------- | ------- | ----------- |
| name | yes      |         | The name of the plan or metaplan to invoke. |

It may be of note that the `name` property does **not** support environment expansion.

### Plan Task Environment

The [environment](https://github.com/efritz/ij/blob/master/docs/environment.md#user-content-environment) built to execute a plan task is merged into the environment of the target plan. The environment active at the time of this task invocation is inserted after the override environment, but before the environment of a task referenced by the target plan.

### Example

The following example defines the plans `build` and `test`, both of which have a set of tasks that must be run first. This sequence of dependencies are expressed as a separate plan which is called recursively via the `go-deps` task. This example is a bit contrived, but this basic strategy is beneficial for common dependencies, long sequences of tasks, and when extending a plan defined in a parent config.

```yaml
tasks:
  go-deps:
    type: plan
    name: go-deps

plans:
  build:
    stages:
      - name: deps
        tasks:
          - go-deps
      - name: build
        tasks:
          - go-build

  test:
    stages:
      - name: deps
        tasks:
          - go-deps
      - name: test
        tasks:
          - go-test

  go-deps:
    stages:
      - name: vendors
        tasks:
          - install-vendors
      - name: deps
        tasks:
          - generate-protobuf
          - generate-mocks
        parallel: true

# additional tasks not shown
```
