package runner

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ContextSuite struct{}

func (s *ContextSuite) TestExportEnv(t sweet.T) {
	a := NewRunContext(nil)
	b := NewRunContext(a)
	c := NewRunContext(b)

	c.ExportEnv("X=1")
	b.ExportEnv("Y=2")
	a.ExportEnv("Z=3")

	Expect(a.GetExportedEnv()).To(Equal([]string{"X=1", "Y=2", "Z=3"}))
	Expect(b.GetExportedEnv()).To(Equal([]string{"X=1", "Y=2", "Z=3"}))
	Expect(c.GetExportedEnv()).To(Equal([]string{"X=1", "Y=2", "Z=3"}))
}
