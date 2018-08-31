package environment

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FilesSuite struct{}

func (s *FilesSuite) TestNormalizeEnvironmentFile(t sweet.T) {
	lines, err := NormalizeEnvironmentFile(`
	X=1
	# a comment
	Y=2

	Z=3
	lower_case_guy=4
	`)

	Expect(err).To(BeNil())
	Expect(lines).To(Equal([]string{
		"X=1",
		"Y=2",
		"Z=3",
		"LOWER_CASE_GUY=4",
	}))
}

func (s *FilesSuite) TestNormalizeEnvironmentFileError(t sweet.T) {
	_, err := NormalizeEnvironmentFile(`
	X=1

	NO_VALUE
	`)

	Expect(err).To(MatchError("Malformed entry in environments file: NO_VALUE"))
}
