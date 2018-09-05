# Extending a Config

A child config (with the `extends` property set) extends its parent config. Tasks, plans, and metaplans defined in the child config will be merged with the parent config. If the same name is defined in both config files, then the child definition overwrites the parent (except for plans defined in the child where the `extends` property is true, then see the [extending a plan](https://github.com/efritz/ij/blob/master/docs/extend.md#user-content-extending-a-plan) section below).

The values of all remaining properties *except* for `ssh-identities` are merged in the following manner:

1. If the property is a list type, then the child value is appended onto the parent value;
2. Otherwise, if the property is defined in the child, the value of the child is used;
3. Otherwise, the value of the parent is used (which may be a zero-value for that type).

Like [override files](https://github.com/efritz/ij/blob/master/docs/override.md#user-content-override-files), any ssh-identities specified in a child config file will *replace* ssh-identities defined in the parent.

# Extending a Task

A child task (with the `extends` property set) extends its parent task. If the child also sets a value for its `type` property, it must match the value of the same property defined in the parent. The values of all remaining properties are merged in the manner described in the section above.

It is legal to form a task extend chain (*a* extends *b*, *b* extends *c*, etc), but this chain may not contain cycles.

# Extending a Plan

This section applies when a child and parent config both define a plan with the same name. If the plan defined in the child config does not have its `extend` property set to true, then the plan defined in the child config overwrites the plan defined in the parent config. In this section we describe the other case.

First, the environment of the child plan is appended onto the environment of the parent plan. Then, each stage defined in the child plan is added to the parent plan in the following manner:

1. If the parent plan defines a stage with the same name, it is overwritten by the child stage;
2. Otherwise, if `before_stage` is set in the child stage, the child stage is inserted directly before the named stage;
3. Otherwise, if `after_stage` is set, in the child stage, the child stage is inserted directly after the named stage;
4. Otherwise, there is an ambiguity error and the stage cannot be inserted into the parent plan.

In the first case, `before_stage` and `after_stage` must not be set in the child stage. In all cases, `before_stage` and `after_stage` must not **both** be set.