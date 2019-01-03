package config

type (
	Task interface {
		GetName() string
		GetType() string
		GetExtends() string
		GetEnvironment() []string
		GetRequiredEnvironment() []string
		Extend(parent Task) error
	}

	TaskMeta struct {
		Name                string   `json:"-"`
		Extends             string   `json:"extends,omitempty"`
		Environment         []string `json:"environment,omitempty"`
		RequiredEnvironment []string `json:"required-environment,omitempty"`
	}
)

func (t *TaskMeta) GetName() string                  { return t.Name }
func (t *TaskMeta) GetExtends() string               { return t.Extends }
func (t *TaskMeta) GetEnvironment() []string         { return t.Environment }
func (t *TaskMeta) GetRequiredEnvironment() []string { return t.RequiredEnvironment }

func (t *TaskMeta) extendMeta(parent TaskMeta) {
	t.Environment = append(parent.Environment, t.Environment...)
	t.RequiredEnvironment = append(parent.RequiredEnvironment, t.RequiredEnvironment...)
}
