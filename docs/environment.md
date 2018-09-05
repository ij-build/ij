# Environment

An environment is an ordered collection of values of the form `VAR` or `VAR=VAL`. In the former, the value is implicitly the value of `$VAR` on the host, and in the latter the value is the literal value `$VAL`.

If a variable is defined multiple times in a single environment, the last value takes precedence. The working environment of a build plan is populated by the following (ordered by increasing precedence):

- the config `environment` and `env_file` properties
- the override `environment` and `env_file` properties
- the active [plan task](https://github.com/efritz/ij/blob/master/docs/tasks.md#user-content-plan-task-environment)'s environment
- the active task's `environment` property
- the active plan's `environment` property
- the active stage's `environment` property
- the active stage task's `environment` property
- the environment exported by a previous run task
- the command line `--environment` and `--env-file` flags

# Environment Files

Contents of an environment file can be interpreted as environment assignments using the `env_file` property of the config and override files, the `--env-file` command line argument, or from an `exported_environment_file` property of a run task.

An environment file consists of lines of the form `VAR=VAL`, where `VAL` contains new newlines (but may contain additional equal signs). Whitespace and `#`-style comments are ignored in environment files.

Any environment file containing a non-empty, non-comment line with no `=` symbol is malformed (and will cause the currently executing build plan to fail).

# Environment Expansion

Almost all tasks and registry object properties (unless otherwise noted) will replace names references of the form `$VAR` and `${VAR}` with the value defined in the environment at the time the task is invoked. If `VAR` is not defined in the environment, the reference will remain untouched. Expansion happens recursively, so if a substituted value also contains a name reference, it will also be substituted.
