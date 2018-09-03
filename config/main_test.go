package config

import (
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&BuildTaskSuite{})
		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&PlanSuite{})
		s.AddSuite(&PlanTaskSuite{})
		s.AddSuite(&PushTaskSuite{})
		s.AddSuite(&RemoveTaskSuite{})
		s.AddSuite(&ResolverSuite{})
		s.AddSuite(&RunTaskSuite{})
		s.AddSuite(&StageSuite{})
		s.AddSuite(&UtilSuite{})
	})
}
