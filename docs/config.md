# Config

A configuration file defines tasks (units of work) and plans (configurations of tasks). The following properties are supported.

| Name        | Default    | Description |
| ----------- | ---------- | ----------- |
| env_file    | []         | A list paths to environment files on the host. Value may be a string or a list. |
| environment | []         | A list of environment variable definitions. Value may be a string or a list. |
| export      | {}         | A [file list object](https://github.com/efritz/ij/blob/master/docs/config.md#user-content-file-list) describing the export phase. |
| extends     | ''         | The path (relative/absolute on-disk, or an HTTP(S) URL) to the parent configuration. |
| import      | {}         | A [file list object](https://github.com/efritz/ij/blob/master/docs/config.md#user-content-file-list) describing the import phase. |
| metaplans   | {}         | A name-metaplan mapping object. See [plans](https://github.com/efritz/ij/blob/master/docs/plans.md#user-content-metaplans) for the definition of these objects. |
| options     | {}         | An [options object](https://github.com/efritz/ij/blob/master/docs/config.md#user-content-options). |
| plans       | {}         | A name-plan mapping object. See [plans](https://github.com/efritz/ij/blob/master/docs/plans.md#user-content-plans) for the definition of these objects. |
| registries  | []         | A list of [docker registries](https://github.com/efritz/ij/blob/master/docs/registries.md) used for login. |
| tasks       | {}         | A name-task mapping object. See [tasks](https://github.com/efritz/ij/blob/master/docs/tasks.md#user-content-tasks) for the definition of these objects. |
| workspace   | /workspace | The default workspace to use for [run tasks](https://github.com/efritz/ij/blob/master/docs/tasks.md#user-content-run-task). |

For details on the `env_file` and `environment` properties, see the documentation on [environments](https://github.com/efritz/ij/blob/master/docs/environment.md).

## Options

The following options can be amended/overridden by [override files](https://github.com/efritz/ij/blob/master/docs/override.md) or via command line arguments.

| Name                 | Default | Description |
| -------------------- | ------- | ----------- |
| force-sequential     | false   | If true, running tasks in parallel will be disabled. |
| healthcheck-interval | 5s      | The duration to wait between health checks of a service container. |
| ssh-identities       | []      | A set of SSH key fingerprints (SHA256 or MD5). Value may be a string or a list. |

If any ssh-identities are supplied in the configuration file or on the command line, then at least one matching fingerprint must exist in the host's SSH agent. On success, the SSH auth socket will be mounted in all containers launched by a *run* task.

The following example object supplies two SHA256 SSH key fingerprints (one prefixed with the checksum type).

```yaml
options:
    ssh-identities:
        - 4aIf67ySUwltykmcwNCDEnCdpkvJ/GRweCdtGuNno9c
        - SHA256:bpCS6MPUbsg3iOmQet2bNZ3m8DAED7ym6h9IrlYzPTU
```

Alternatively, specifying the fingerprint '*' will allow any key that's bene added to the host's SSH agent (if at least one exists).

```yaml
options:
    ssh-identities: '*'
```

## File List

A file list object controls the files which move into the workspace on import and back into the project directory on export. Each direction is configured independently.

| Name    | Default | Description |
| ------- | ------- | ----------- |
| exclude | []      | Glob patterns for files to be ignored during import or export. Value may be a string or a list. |
| files   | []      | Glob patterns for files targeted for transfer during import or export. Value may be a string or a list. |

Files matching a pattern in the `files` property will be *recursively* transferred in or out of the workspace. If that file also matches a pattern in the `exclude` property, it will be skipped. All symlinks are skipped during transfer. Glob patterns support `*` for optional text and `**` for multiple directories.