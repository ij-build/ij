package config

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type UtilSuite struct{}

func (s *UtilSuite) TestExtendString(t sweet.T) {
	Expect(extendString("foo", "bar")).To(Equal("foo"))
	Expect(extendString("", "bar")).To(Equal("bar"))
}

func (s *UtilSuite) TestExtendBool(t sweet.T) {
	Expect(extendBool(true, true)).To(BeTrue())
	Expect(extendBool(true, false)).To(BeTrue())
	Expect(extendBool(false, true)).To(BeTrue())
	Expect(extendBool(false, false)).To(BeFalse())
}

func (s *UtilSuite) TestExtendInt(t sweet.T) {
	Expect(extendInt(5, 10)).To(Equal(5))
	Expect(extendInt(0, 10)).To(Equal(10))
}

func (s *UtilSuite) TestExtendDuration(t sweet.T) {
	Expect(extendDuration(time.Second, time.Minute)).To(Equal(time.Second))
	Expect(extendDuration(0, time.Minute)).To(Equal(time.Minute))
}
