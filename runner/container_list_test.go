package runner

import (
	"github.com/aphistic/sweet"
	"github.com/ij-build/ij/logging"
	. "github.com/onsi/gomega"
)

type ContainerListSuite struct{}

func (s *ContainerListSuite) TestExecute(t sweet.T) {
	args := make(chan string, 3)

	containers := NewContainerList(
		"test",
		func(name string) { args <- name },
		logging.NilLogger,
	)

	containers.Add("foo")
	containers.Add("bar")
	containers.Add("baz")
	containers.Execute()
	close(args)

	Expect(receive(args)).To(ConsistOf("foo", "bar", "baz"))
}

func (s *ContainerListSuite) TestRemove(t sweet.T) {
	args := make(chan string, 3)

	containers := NewContainerList(
		"test",
		func(name string) { args <- name },
		logging.NilLogger,
	)

	containers.Add("foo")
	containers.Add("bar")
	containers.Remove("bar")
	containers.Add("baz")
	containers.Execute()
	close(args)

	Expect(receive(args)).To(ConsistOf("foo", "baz"))
}

//
// Helpers

func receive(ch <-chan string) []string {
	vals := []string{}
	for val := range ch {
		vals = append(vals, val)
	}

	return vals
}
