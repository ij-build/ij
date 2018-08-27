package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type TaskSuite struct{}

func (s *TaskSuite) TestTranslateTask(t sweet.T) {
	task, err := translateTask("foo", json.RawMessage(`{
		"type": "build",
		"extends": "parent",
		"dockerfile": "Dockerfile"
	}`))

	Expect(err).To(BeNil())
	Expect(task.GetName()).To(Equal("foo"))
	Expect(task.GetExtends()).To(Equal("parent"))
	Expect(task.(*config.BuildTask).Dockerfile).To(Equal("Dockerfile"))
}

func (s *TaskSuite) TestTranslateTaskDefaultType(t sweet.T) {
	task, err := translateTask("foo", json.RawMessage(`{
		"extends": "parent",
		"image": "alpine"
	}`))

	Expect(err).To(BeNil())
	Expect(task.GetName()).To(Equal("foo"))
	Expect(task.GetExtends()).To(Equal("parent"))
	Expect(task.(*config.RunTask).Image).To(Equal("alpine"))
}

func (s *TaskSuite) TestTranslateTaskUnknownType(t sweet.T) {
	_, err := translateTask("foo", json.RawMessage(`{
		"type": "register",
		"extends": "parent",
		"image": "alpine"
	}`))

	Expect(err).To(MatchError("unknown task type 'register'"))
}
