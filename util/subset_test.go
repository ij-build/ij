package util

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type SubsetSuite struct{}

func (s *SubsetSuite) TestContainsAllSuperset(t sweet.T) {
	sup := []string{"a", "b", "c", "d"}
	sub := []string{"a", "b", "c"}

	subset, missing := ContainsAll(sup, sub)
	Expect(subset).To(BeTrue())
	Expect(missing).To(BeEmpty())
}

func (s *SubsetSuite) TestContainsAllEqual(t sweet.T) {
	sup := []string{"a", "b", "c"}
	sub := []string{"a", "b", "c"}

	subset, missing := ContainsAll(sup, sub)
	Expect(subset).To(BeTrue())
	Expect(missing).To(BeEmpty())
}

func (s *SubsetSuite) TestContainsAllSubset(t sweet.T) {
	sup := []string{"a", "b"}
	sub := []string{"a", "b", "c"}

	subset, missing := ContainsAll(sup, sub)
	Expect(subset).To(BeFalse())
	Expect(missing).To(ConsistOf("c"))
}

func (s *SubsetSuite) TestContainsAllEmpty(t sweet.T) {
	sup := []string{"a", "b", "c"}
	subset, missing := ContainsAll(sup, nil)
	Expect(subset).To(BeTrue())
	Expect(missing).To(BeEmpty())
}
