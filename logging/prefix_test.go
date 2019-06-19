package logging

//go:generate go-mockgen -f github.com/ij-build/ij/logging -i ColorPicker -o mock_color_picker_test.go

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PrefixSuite struct{}

func (s *PrefixSuite) TestSerialize(t sweet.T) {
	prefix := NewPrefix("foo", "bar", "baz")
	Expect(prefix.Serialize(NilColorPicker)).To(Equal("foo/bar/baz"))
}

func (s *PrefixSuite) TestSerializeWithPicker(t sweet.T) {
	picker := NewMockColorPicker()
	picker.ColorizeFunc.SetDefaultHook(func(val string) string {
		return fmt.Sprintf("[%s]", val)
	})

	prefix := NewPrefix("foo", "bar", "baz")
	Expect(prefix.Serialize(picker)).To(Equal("[foo]/[bar]/[baz]"))
}

func (s *PrefixSuite) TestAppend(t sweet.T) {
	var (
		prefix1 = NewPrefix("foo")
		prefix2 = prefix1.Append("bar")
		prefix3 = prefix2.Append("baz")
		prefix4 = prefix2.Append("bonk")
	)

	Expect(prefix1.Serialize(NilColorPicker)).To(Equal("foo"))
	Expect(prefix2.Serialize(NilColorPicker)).To(Equal("foo/bar"))
	Expect(prefix3.Serialize(NilColorPicker)).To(Equal("foo/bar/baz"))
	Expect(prefix4.Serialize(NilColorPicker)).To(Equal("foo/bar/bonk"))
}
