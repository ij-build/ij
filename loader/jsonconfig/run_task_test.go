package jsonconfig

import (
	"encoding/json"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/util"
	. "github.com/onsi/gomega"
)

type RunTaskSuite struct{}

func (s *RunTaskSuite) TestTranslate(t sweet.T) {
	task := &RunTask{
		Extends:             "parent",
		Environment:         json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
		RequiredEnvironment: []string{"X"},
		Image:               "image",
		Command:             "command",
		Shell:               "shell",
		Script:              "script",
		Entrypoint:          "entrypoint",
		User:                "user",
		Workspace:           "workspace",
		Hostname:            "hostname",
		Detach:              true,
		Healthcheck:         nil,

		ExportEnvironmentFiles: json.RawMessage(`["e1","e2"]`),
	}

	translated, err := task.Translate("run")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RunTask{
		TaskMeta: config.TaskMeta{
			Name:                "run",
			Extends:             "parent",
			Environment:         []string{"X=1", "Y=2", "Z=3"},
			RequiredEnvironment: []string{"X"},
		},
		Image:                  "image",
		Command:                "command",
		Shell:                  "shell",
		Script:                 "script",
		Entrypoint:             "entrypoint",
		User:                   "user",
		Workspace:              "workspace",
		Hostname:               "hostname",
		Detach:                 true,
		Healthcheck:            &config.Healthcheck{},
		ExportEnvironmentFiles: []string{"e1", "e2"},
	}))
}

func (s *RunTaskSuite) TestTranslateStringLists(t sweet.T) {
	task := &RunTask{
		Extends:                "parent",
		Environment:            json.RawMessage(`"X=1"`),
		ExportEnvironmentFiles: json.RawMessage(`"env"`),
	}

	translated, err := task.Translate("run")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RunTask{
		TaskMeta: config.TaskMeta{
			Name:        "run",
			Extends:     "parent",
			Environment: []string{"X=1"},
		},
		Healthcheck:            &config.Healthcheck{},
		ExportEnvironmentFiles: []string{"env"},
	}))
}

func (s *RunTaskSuite) TestTranslateHealthcheck(t sweet.T) {
	task := &RunTask{
		Extends: "parent",
		Healthcheck: &Healthcheck{
			Command:     "command",
			Interval:    util.Duration{time.Second},
			Retries:     10,
			StartPeriod: util.Duration{time.Second},
			Timeout:     util.Duration{time.Second},
		},
	}

	translated, err := task.Translate("run")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RunTask{
		TaskMeta: config.TaskMeta{
			Name:    "run",
			Extends: "parent",
		},
		Healthcheck: &config.Healthcheck{
			Command:     "command",
			Interval:    time.Second,
			Retries:     10,
			StartPeriod: time.Second,
			Timeout:     time.Second,
		},
	}))
}
