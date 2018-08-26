package environment

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ExpandSuite struct{}

func (s *ExpandSuite) TestExpandString(t sweet.T) {
	env := Environment{"FOO": "x", "BAR": "y", "BAZ": "z"}
	val, err := env.ExpandString("foo=${FOO}, bar=$BAR, baz=${BAZ}")
	Expect(err).To(BeNil())
	Expect(val).To(Equal("foo=x, bar=y, baz=z"))
}

func (s *ExpandSuite) TestExpandStringNoName(t sweet.T) {
	env := Environment{"FOO": "x", "BAR": "y", "BAZ": "z"}
	val, err := env.ExpandString("${FOO} $(./${BAR} $BAZ)")
	Expect(err).To(BeNil())
	Expect(val).To(Equal("x $(./y z)"))
}

func (s *ExpandSuite) TestExpandStringMissing(t sweet.T) {
	env := Environment{"FOO": "x"}
	val, err := env.ExpandString("foo=${FOO}, bar=${BAR}, baz=${BAZ}")
	Expect(err).To(BeNil())
	Expect(val).To(Equal("foo=x, bar=${BAR}, baz=${BAZ}"))
}

func (s *ExpandSuite) TestExpandStringMissingNoBrackets(t sweet.T) {
	env := Environment{"FOO": "x"}
	val, err := env.ExpandString("foo=$FOO, bar=$BAR, baz=$BAZ")
	Expect(err).To(BeNil())
	Expect(val).To(Equal("foo=x, bar=${BAR}, baz=${BAZ}"))
}
func (s *ExpandSuite) TestExpandStringIterative(t sweet.T) {
	env := Environment{"FOO": "x", "BAR": "y${FOO}", "BAZ": "z${BAR}"}
	val, err := env.ExpandString("foo=${FOO}, bar=${BAR}, baz=${BAZ}")
	Expect(err).To(BeNil())
	Expect(val).To(Equal("foo=x, bar=yx, baz=zyx"))
}

func (s *ExpandSuite) TestExpandStringRecursive(t sweet.T) {
	env := Environment{"FOO": "x${BAR}", "BAR": "y${FOO}"}
	_, err := env.ExpandString("foo=${FOO}")
	Expect(err).NotTo(BeNil())
}

func (s *ExpandSuite) TestExpandSlice(t sweet.T) {
	env := Environment{"FOO": "x", "BAR": "y", "BAZ": "z"}
	vals, err := env.ExpandSlice([]string{"foo=${FOO}", "bar=${BAR}", "baz=${BAZ}"})
	Expect(err).To(BeNil())
	Expect(vals).To(Equal([]string{"foo=x", "bar=y", "baz=z"}))
}

func (s *ExpandSuite) TestExpandSliceRecursive(t sweet.T) {
	env := Environment{"FOO": "x${BAR}", "BAR": "y${FOO}"}
	_, err := env.ExpandSlice([]string{"foo=${FOO}", "bar=${BAR}"})
	Expect(err).NotTo(BeNil())
}
