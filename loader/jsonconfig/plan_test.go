package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type PlanSuite struct{}

func (s *PlanSuite) TestTranslate(t sweet.T) {
	plan := &Plan{
		Extend: true,
		Stages: []*Stage{
			&Stage{
				Name: "bar",
				Tasks: []json.RawMessage{
					json.RawMessage(`"t1"`),
					json.RawMessage(`{"name":"t2", "environment": ["W=5"]}`),
				},
				Parallel: true,
			},
			&Stage{
				Name: "baz",
				Tasks: []json.RawMessage{
					json.RawMessage(`"t3"`),
				},
				Environment: []string{"Z=4"},
			},
		},
		Environment: []string{"X=1", "Y=2", "Z=3"},
	}

	translated, err := plan.Translate("foo")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Plan{
		Name:   "foo",
		Extend: true,
		Stages: []*config.Stage{
			&config.Stage{
				Name: "bar",
				Tasks: []*config.StageTask{
					&config.StageTask{Name: "t1"},
					&config.StageTask{Name: "t2", Environment: []string{"W=5"}},
				},
				RunMode:  config.RunModeOnSuccess,
				Parallel: true,
			},
			&config.Stage{
				Name: "baz",
				Tasks: []*config.StageTask{
					&config.StageTask{Name: "t3"},
				},
				RunMode:     config.RunModeOnSuccess,
				Environment: []string{"Z=4"},
			},
		},
		Environment: []string{"X=1", "Y=2", "Z=3"},
	}))
}
