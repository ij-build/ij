package environment

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type MergeSuite struct{}

func (s *MergeSuite) TestMerge(t sweet.T) {
	env1 := Environment{"FOO": "x", "BAR": "y", "BAZ": "z"}
	env2 := Environment{"FOO": "a", "BNK": "c"}
	env3 := Environment{"FOO": "m", "BAR": "n", "QUX": "l"}

	Expect(Merge(env1, env2, env3)).To(Equal(Environment{
		"FOO": "m",
		"BAR": "n",
		"BAZ": "z",
		"BNK": "c",
		"QUX": "l",
	}))
}
