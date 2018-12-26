package config

type (
	Task interface {
		GetName() string
		GetType() string
		GetExtends() string
		GetEnvironment() []string
		Extend(parent Task) error
	}

	TaskMeta struct {
		Name    string
		Extends string
	}
)

func (t *TaskMeta) GetName() string    { return t.Name }
func (t *TaskMeta) GetExtends() string { return t.Extends }
