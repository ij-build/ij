package state

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CleanupSuite struct{}

func (s *CleanupSuite) TestCleanup(t sweet.T) {
	called := []string{}

	cleanup := NewCleanup()
	cleanup.Register(func() { called = append(called, "a") })
	cleanup.Register(func() { called = append(called, "b") })
	cleanup.Register(func() { called = append(called, "c") })
	cleanup.Cleanup()

	Expect(called).To(Equal([]string{"c", "b", "a"}))
}
