package jsonconfig

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
		s.AddSuite(&OverrideSuite{})
		s.AddSuite(&PlanSuite{})
		s.AddSuite(&PlanTaskSuite{})
		s.AddSuite(&RegistrySuite{})
		s.AddSuite(&RemoveTaskSuite{})
		s.AddSuite(&RunTaskSuite{})
		s.AddSuite(&StageSuite{})
		s.AddSuite(&TaskSuite{})
	})
}
