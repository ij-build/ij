package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/ij-build/ij/config"
	. "github.com/onsi/gomega"
)

type TaskSuite struct{}

func (s *TaskSuite) TestTranslateTask(t sweet.T) {
	task, err := translateTask(nil, "foo", json.RawMessage(`{
		"type": "build",
		"extends": "parent",
		"dockerfile": "Dockerfile"
	}`))

	Expect(err).To(BeNil())
	Expect(task.GetName()).To(Equal("foo"))
	Expect(task.GetExtends()).To(Equal("parent"))
	Expect(task.GetType()).To(Equal("build"))
	Expect(task.(*config.BuildTask).Dockerfile).To(Equal("Dockerfile"))
}

func (s *TaskSuite) TestTranslateTaskWithParentHint(t sweet.T) {
	parent := &config.Config{
		Tasks: map[string]config.Task{
			"bar": &config.BuildTask{},
		},
	}

	task, err := translateTask(parent, "foo", json.RawMessage(`{
		"extends": "bar",
		"dockerfile": "Dockerfile"
	}`))

	Expect(err).To(BeNil())
	Expect(task.GetName()).To(Equal("foo"))
	Expect(task.GetExtends()).To(Equal("bar"))
	Expect(task.GetType()).To(Equal("build"))
	Expect(task.(*config.BuildTask).Dockerfile).To(Equal("Dockerfile"))
}

func (s *TaskSuite) TestTranslateMalformedTaskWithParentHint(t sweet.T) {
	parent := &config.Config{
		Tasks: map[string]config.Task{
			"foo": &config.RunTask{},
		},
	}

	_, err := translateTask(parent, "foo", json.RawMessage(`{
		"extends": "foo",
		"dockerfile": "Dockerfile"
	}`))

	Expect(err).To(MatchError("failed to validate task foo: Additional property dockerfile is not allowed"))
}

func (s *TaskSuite) TestTranslateTaskWithOverriddenParentHint(t sweet.T) {
	parent := &config.Config{
		Tasks: map[string]config.Task{
			"foo": &config.PushTask{},
		},
	}

	task, err := translateTask(parent, "foo", json.RawMessage(`{
		"type": "build",
		"extends": "parent",
		"dockerfile": "Dockerfile"
	}`))

	Expect(err).To(BeNil())
	Expect(task.GetName()).To(Equal("foo"))
	Expect(task.GetExtends()).To(Equal("parent"))
	Expect(task.GetType()).To(Equal("build"))
	Expect(task.(*config.BuildTask).Dockerfile).To(Equal("Dockerfile"))
}

func (s *TaskSuite) TestTranslateTaskDefaultType(t sweet.T) {
	task, err := translateTask(nil, "foo", json.RawMessage(`{
		"extends": "parent",
		"image": "alpine"
	}`))

	Expect(err).To(BeNil())
	Expect(task.GetName()).To(Equal("foo"))
	Expect(task.GetExtends()).To(Equal("parent"))
	Expect(task.GetType()).To(Equal("run"))
	Expect(task.(*config.RunTask).Image).To(Equal("alpine"))
}

func (s *TaskSuite) TestTranslateTaskUnknownType(t sweet.T) {
	_, err := translateTask(nil, "foo", json.RawMessage(`{
		"type": "register",
		"extends": "parent",
		"image": "alpine"
	}`))

	Expect(err).To(MatchError("unknown task type 'register'"))
}

func (s *TaskSuite) TestTranslateTaskMalformedSchema(t sweet.T) {
	_, err := translateTask(nil, "foo", json.RawMessage(`{
		"type": "run",
		"dockerfile": "alpine"
	}`))

	Expect(err).To(MatchError("failed to validate task foo: Additional property dockerfile is not allowed"))
}
