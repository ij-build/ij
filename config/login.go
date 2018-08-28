package config

import "fmt"

type LoginTask struct {
	TaskMeta
	Server       string
	Username     string
	Password     string
	PasswordFile string
}

func (t *LoginTask) GetType() string {
	return "login"
}

func (t *LoginTask) Extend(task Task) error {
	parent, ok := task.(*LoginTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.Server = extendString(t.Server, parent.Server)
	t.Username = extendString(t.Username, parent.Username)
	t.Password = extendString(t.Password, parent.Password)
	t.PasswordFile = extendString(t.PasswordFile, parent.PasswordFile)
	return nil
}
