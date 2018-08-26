package environment

import (
	"os"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type EnvironmentSuite struct{}

func (s *EnvironmentSuite) SetUpTest(t sweet.T) {
	os.Clearenv()
}

func (s *EnvironmentSuite) TestNew(t sweet.T) {
	Expect(New([]string{
		"FOO=x",
		"BAR=y",
		"BAZ=z",
	})).To(Equal(Environment{
		"FOO": "x",
		"BAR": "y",
		"BAZ": "z",
	}))
}

func (s *EnvironmentSuite) TestNewDuplicates(t sweet.T) {
	Expect(New([]string{
		"FOO=x",
		"BAR=y",
		"BAZ=z",
		"FOO=w",
	})).To(Equal(Environment{
		"FOO": "w",
		"BAR": "y",
		"BAZ": "z",
	}))
}

func (s *EnvironmentSuite) TestNewNoValue(t sweet.T) {
	os.Setenv("BAR", "y")
	os.Setenv("BAZ", "z")

	Expect(New([]string{
		"FOO=x",
		"BAR",
		"BAZ",
	})).To(Equal(Environment{
		"FOO": "x",
		"BAR": "y",
		"BAZ": "z",
	}))
}

func (s *EnvironmentSuite) TestKeys(t sweet.T) {
	env := New([]string{
		"FOO=x",
		"BAR=y",
		"BAZ=z",
	})

	Expect(env.Keys()).To(Equal([]string{"BAR", "BAZ", "FOO"}))
}

func (s *EnvironmentSuite) TestSerialize(t sweet.T) {
	env := New([]string{
		"FOO=x",
		"BAR=y",
		"BAZ=z",
	})

	Expect(env.Serialize()).To(Equal([]string{
		"BAR=y",
		"BAZ=z",
		"FOO=x",
	}))
}
