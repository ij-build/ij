# Plans

A plan is a configuration of tasks organized into *stages*, described below.

| Name        | Required | Default | Description |
| ----------- | -------- | ------- | ----------- |
| environment |          | []      | A list of environment variable definitions. Value may be a string or a list. |
| extend      |          | false   | Whether or not the plan is extending a plan defined in the parent config with the same name. |
| stages      |          | []      | A list of [stage](https://github.com/efritz/ij/blob/master/docs/plans.md#user-content-stage) objects. |

See the section on [extending a plan](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-plan) about the semantics of the `extends` property.

## Stage

A stage is a direct collection of tasks.

| Name         | Required | Default    | Description |
| ------------ | -------- | ---------- | ----------- |
| after_stage  |          | ''         | A target sibling stage in the same plan (applicable only when the parent plan is extending). |
| before_stage |          | ''         | A target sibling stage in the same plan (applicable only when the parent plan is extending). |
| environment  |          | []         | A list of environment variable definitions. Value may be a string or a list. |
| name         | yes      |            | The name of the stage. Must be unique within the plan. |
| parallel     |          | false      | Whether or not to run tasks sequentially or in parallel. |
| run-mode     |          | on-success | One of `on-success`, `on-failure`, or `always`. Determines if a stage should run in the presence of a previous stage failure. |
| tasks        |          | []         | A list of tasks to run. Values in this list can be a string (the task name), or a object with a `name` and `environment` property. |

See the section on [extending a plan](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-plan) about the semantics of the `after_stage` and `before_stage` properties.

If `parallel` is true, then each task in the list is run in a different thread. The stage will end once all tasks have ended. If `parallel` is false (the default), then each task in the stage is run to completion in sequence.

If the `run mode` property is set to `on-success` (the default), then the stage will only run if no previous failure has occurred. Use the value `on-failure` to mark a stage as an error handler (when the stage will only occur if a previous failure has occurred), and use the value `always` to mark a stage for some cleanup or *finally*-like behavior.

# Metaplans

A metaplan is simply a list of plans and is semantically equivalent to running the stages of the listed plans back-to-back. A metaplan can be referenced in any place that a plan can be referenced.
