# Override Files

Developers can define user and project-level overrides which are applied on top of a config. A user-level override is read from `~/.ij/override.yaml`, and a project-level override file is read from `ij.override.yaml` or `ij.override.yml` in the project directory. Project-overrides take precedence over user-level overrides.

An override file is basically a [config](https://github.com/efritz/ij/blob/master/docs/config.md#user-content-config) with some properties stripped. The following properties are supported.

| Name        | Notes |
| ----------- | ----- |
| env_file    | |
| environment | |
| export      | Only exclude patterns can be supplied. |
| import      | Only exclude patterns can be supplied. |
| options     | |
| registries  | |

The values of all properties *except* for `ssh-identities` will be appended to the parent. In the case of environment variables, the values defined in the override are given precedence in the case of name collision. Any ssh-identities specified in an override file will *replace* ssh-identities defined in the parent.
