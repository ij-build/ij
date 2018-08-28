package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/loader/schema"
)

type (
	ExtendHint struct {
		Extend string `json:"extends"`
	}

	TypeHint struct {
		Type string `json:"type"`
	}

	Task interface {
		Translate(name string) (config.Task, error)
	}
)

func translateTask(
	parent *config.Config,
	name string,
	data json.RawMessage,
) (config.Task, error) {
	// Assume task is type run
	typeHint := &TypeHint{Type: "run"}

	if parent != nil {
		extendHint := &ExtendHint{}
		if err := json.Unmarshal(data, extendHint); err != nil {
			return nil, err
		}

		// Update assumption if we're extending a task defined
		// in the parent (we will still allow an explicit overwrite
		// with a check later to ensure something bad doesn't happen).
		if parentTask, ok := parent.Tasks[extendHint.Extend]; ok {
			typeHint.Type = parentTask.GetType()
		}
	}

	// See if the user provided an _explicit_ type which may contradict
	// our assumptions above. We do this to catch errors by the user in
	// the case they extend the wrong task.

	if err := json.Unmarshal(data, typeHint); err != nil {
		return nil, err
	}

	// Before validating against a schema, ensure that the type hint
	// given is something that we expect. Return an error for unknown
	// task types (instead of failing to find the schema).

	structMap := map[string]Task{
		"run":    &RunTask{},
		"build":  &BuildTask{},
		"push":   &PushTask{},
		"remove": &RemoveTask{},
	}

	task, ok := structMap[typeHint.Type]
	if !ok {
		return nil, fmt.Errorf("unknown task type '%s'", typeHint.Type)
	}

	// Now, validate the fields of the payload against the claimed type.
	// Required fields are taken care of by the schema (it will either
	// require an extends or require the fields to be filled in).

	assetName := fmt.Sprintf("schema/%s.yaml", typeHint.Type)

	if err := schema.Validate(assetName, data); err != nil {
		return nil, fmt.Errorf("failed to validate task %s: %s", name, err.Error())
	}

	// Now we can create an empty struct and populate it now that we
	// know it contains only valid fields.

	if err := json.Unmarshal(data, task); err != nil {
		return nil, err
	}

	return task.Translate(name)
}
