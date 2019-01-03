# Plans

A plan is a configuration of tasks organized into *stages*, described below.

| Name        | Required | Default | Description |
| ----------- | -------- | ------- | ----------- |
| disabled    |          | ''      | A flag that, if non-emptyh, will cause the plan to be skipped. |
| environment |          | []      | A list of environment variable definitions. Value may be a string or a list. |
| extend      |          | false   | Whether or not the plan is extending a plan defined in the parent config with the same name. |
| stages      |          | []      | A list of [stage](https://github.com/efritz/ij/blob/master/docs/plans.md#user-content-stage) objects. |

See the section on [extending a plan](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-plan) about the semantics of the `extends` property.

The disabled flag is notably **not** a boolean value in order for it to be parameterized via environment variable [expansion](https://github.com/efritz/ij/blob/master/docs/environment.md#user-content-environment-expansion). The disabled flag is evaluated with the environment available at the time the plan executes. Any non-empty string will cause the disable flag to be interpreted as true. Stages and stage tasks also have a disabled flag that works in the same manner.

## Stage

A stage is a direct collection of tasks.

| Name         | Required | Default    | Description |
| ------------ | -------- | ---------- | ----------- |
| after-stage  |          | ''         | A target sibling stage in the same plan (applicable only when the parent plan is extending). |
| before-stage |          | ''         | A target sibling stage in the same plan (applicable only when the parent plan is extending). |
| disabled     |          | ''         | A flag that, if non-empty, will cause the stage to be skipped. |
| environment  |          | []         | A list of environment variable definitions. Value may be a string or a list. |
| name         | yes      |            | The name of the stage. Must be unique within the plan. |
| parallel     |          | false      | Whether or not to run tasks sequentially or in parallel. |
| run-mode     |          | on-success | One of `on-success`, `on-failure`, or `always`. Determines if a stage should run in the presence of a previous stage failure. |
| tasks        |          | []         | A list of tasks to run. Values in this list can be a string (supplying only the task name), or a [stage task object](https://github.com/efritz/ij/blob/master/docs/plans.md#user-content-stage-task). |

See the section on [extending a plan](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-plan) about the semantics of the `after-stage` and `before-stage` properties.

If `parallel` is true, then each task in the list is run in a different thread. The stage will end once all tasks have ended. If `parallel` is false (the default), then each task in the stage is run to completion in sequence.

If the `run mode` property is set to `on-success` (the default), then the stage will only run if no previous failure has occurred. Use the value `on-failure` to mark a stage as an error handler (when the stage will only occur if a previous failure has occurred), and use the value `always` to mark a stage for some cleanup or *finally*-like behavior.

## Stage Task

| Name         | Required | Default    | Description |
| ------------ | -------- | ---------- | ----------- |
| disabled     |          | ''         | A flag that, if non-empty, will cause the task in this stage to be skipped. |
| environment  |          | []         | A list of environment variable definitions. Value may be a string or a list. |
| name         | yes      |            | The name of the task. |

# Metaplans

A metaplan is simply a list of plans and is semantically equivalent to running the stages of the listed plans back-to-back. A metaplan can be referenced in any place that a plan can be referenced.

# Example

This example defines a plan with two stages. The first stage installs golang dependencies and the second stage builds three golang binaries in parallel.

```yaml
plans:
  build:
    stages:
      - name: vendors
        tasks:
          - install-vendors
      - name: build
        tasks:
          - name: go-build
            environment: APP=a
          - name: go-build
            environment: APP=b
          - name: go-build
            environment: APP=c
        parallel: true

# tasks not shown
```
