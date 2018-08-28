package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type StageSuite struct{}

func (s *StageSuite) TestShouldRun(t sweet.T) {
	var (
		s1 = &Stage{RunMode: RunModeAlways}
		s2 = &Stage{RunMode: RunModeOnSuccess}
		s3 = &Stage{RunMode: RunModeOnFailure}
	)

	Expect(s1.ShouldRun(true)).To(BeTrue())
	Expect(s1.ShouldRun(false)).To(BeTrue())
	Expect(s2.ShouldRun(true)).To(BeFalse())
	Expect(s2.ShouldRun(false)).To(BeTrue())
	Expect(s3.ShouldRun(true)).To(BeTrue())
	Expect(s3.ShouldRun(false)).To(BeFalse())
}
