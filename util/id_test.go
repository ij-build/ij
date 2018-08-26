package util

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type IDSuite struct{}

const RandomTrials = 50

func (s *IDSuite) TestMakeID(t sweet.T) {
	ids := map[string]struct{}{}

	for i := 0; i < RandomTrials; i++ {
		id, err := MakeID()
		Expect(err).To(BeNil())
		Expect(id).To(HaveLen(6))

		Expect(ids).NotTo(HaveKey(id))
		ids[id] = struct{}{}
	}
}
