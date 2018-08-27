package jsonconfig

import (
	"encoding/json"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type RunSuite struct{}

func (s *RunSuite) TestTranslate(t sweet.T) {
	task := &RunTask{
		Extends:                "parent",
		Image:                  "image",
		Command:                "command",
		Shell:                  "shell",
		Script:                 "script",
		Entrypoint:             "entrypoint",
		User:                   "user",
		Workspace:              "workspace",
		Hostname:               "hostname",
		Detach:                 true,
		Healthcheck:            nil,
		Environment:            []string{"X=1", "Y=2", "Z=3"},
		RequiredEnvironment:    []string{"X"},
		ExportEnvironmentFiles: json.RawMessage(`["e1","e2"]`),
	}

	translated, err := task.Translate("run")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RunTask{
		TaskMeta: config.TaskMeta{
			Name:    "run",
			Extends: "parent",
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
		Environment:            []string{"X=1", "Y=2", "Z=3"},
		RequiredEnvironment:    []string{"X"},
		ExportEnvironmentFiles: []string{"e1", "e2"},
	}))
}

func (s *RunSuite) TestTranslateStringLists(t sweet.T) {
	task := &RunTask{
		Extends:                "parent",
		ExportEnvironmentFiles: json.RawMessage(`"env"`),
	}

	translated, err := task.Translate("run")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RunTask{
		TaskMeta: config.TaskMeta{
			Name:    "run",
			Extends: "parent",
		},
		Healthcheck:            &config.Healthcheck{},
		ExportEnvironmentFiles: []string{"env"},
	}))
}

func (s *RunSuite) TestTranslateHealthcheck(t sweet.T) {
	task := &RunTask{
		Extends: "parent",
		Healthcheck: &Healthcheck{
			Command:     "command",
			Interval:    Duration{time.Second},
			Retries:     10,
			StartPeriod: Duration{time.Second},
			Timeout:     Duration{time.Second},
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
