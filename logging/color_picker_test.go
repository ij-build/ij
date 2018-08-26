package logging

import (
	"strings"

	"github.com/aphistic/sweet"
	"github.com/mgutz/ansi"
	. "github.com/onsi/gomega"
)

type ColorPickerSuite struct{}

func (s *ColorPickerSuite) TestDisabled(t sweet.T) {
	Expect(NilColorPicker.Colorize("foo")).To(Equal("foo"))
}

func (s *ColorPickerSuite) TestEnabled(t sweet.T) {
	picker := newColorPicker(true)

	c1 := picker.Colorize("foo")
	c2 := picker.Colorize("bar")
	c3 := picker.Colorize("baz")
	c4 := picker.Colorize("foo")

	// Test hashing + caching
	Expect(c1).NotTo(Equal(c2))
	Expect(c1).NotTo(Equal(c3))
	Expect(c2).NotTo(Equal(c3))
	Expect(c1).To(Equal(c4))

	Expect(c1).To(ContainSubstring("foo"))
	Expect(c2).To(ContainSubstring("bar"))
	Expect(c3).To(ContainSubstring("baz"))
	Expect(c1).To(HaveSuffix(ansi.Reset))
	Expect(c2).To(HaveSuffix(ansi.Reset))
	Expect(c3).To(HaveSuffix(ansi.Reset))

	found := false
	for _, color := range colors {
		if strings.HasPrefix(c1, color) {
			found = true
		}
	}

	Expect(found).To(BeTrue())
}
