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

## Example

This example defines the `venv-py2` and `venv-py3` tasks. Both tasks extend the same base task but update the environment so that different binary names are used by the script.

```yaml
tasks:
  venv:
    image: ${PY_IMAGE}
    shell: /bin/bash
    script: |
      if [ ! -d ${VENV_NAME} ]; then
        virtualenv ${VENV_NAME} --python=${VENV_PYTHON}
      fi

      source ./${VENV_NAME}/bin/activate
      ${VENV_PIP} install -U pip
      ls requirements* | xargs -I {} ${VENV_PIP} install -r {}
    required_environment:
      - VENV_NAME
      - VENV_PYTHON
      - VENV_PIP

  venv-py2:
    extends: venv
    environment:
      - VENV_NAME=venv2
      - VENV_PYTHON=python2
      - VENV_PIP=pip

  venv-py3:
    extends: venv
    environment:
      - VENV_NAME=venv3
      - VENV_PYTHON=python3
      - VENV_PIP=pip3

# plans not shown
```

# Extending a Plan

A plan defined in a child config overwrites the a plan defined in a parent config with the same name. Such a plan can rather *extend* the parent plan with additional or overridden functionality. The `extends` property of a plan can be set to the name of another previously defined plan (either in the current config or a parent config, but not a child config).

First, the environment of the child plan is appended onto the environment of the parent plan. Then, each stage defined in the child plan is added to the parent plan in the following manner:

1. If the parent plan defines a stage with the same name, it is overwritten by the child stage;
2. Otherwise, if `before-stage` is set in the child stage, the child stage is inserted directly before the named stage;
3. Otherwise, if `after-stage` is set, in the child stage, the child stage is inserted directly after the named stage;
4. Otherwise, there is an ambiguity error and the stage cannot be inserted into the parent plan.

In the first case, `before-stage` and `after-stage` must not be set in the child stage. In all cases, `before-stage` and `after-stage` must not **both** be set.

## Example

The parent config file in this example declares a four-stage plan for running integration tests. The parent declares only the structure of the plan and does not reference any tasks.

```yaml
# parent.yaml

plans:
  test-integration:
    stages:
      - name: deps
      - name: services
      - name: migrations
      - name: test
```

The project config file, extending the parent config, refines the integration test plan by adding tasks to the migrations and test stages, and adds an additional stage to run fixtures.

```yaml
# ij.yaml

extends: parent.yaml

plans:
  test-integration:
    extends: test-integration
    stages:
      - name: migrations
        tasks:
          - api-migrate-postgres
          - api-migrate-cassandra
        parallel: true
      - name: fixtures
        after-stage: migrations
        tasks:
          - fixtures
      - name: test
        tasks:
          - test

# tasks not shown
```
