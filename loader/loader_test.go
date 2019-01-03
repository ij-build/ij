package loader

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type LoaderSuite struct{}

func (s *LoaderSuite) TestLoad(t sweet.T) {
	loaded, err := NewLoader().Load("./test-configs/basic.yaml")
	Expect(err).To(BeNil())
	Expect(loaded).To(Equal(&config.Config{
		Options:     &config.Options{},
		Registries:  []config.Registry{},
		Import:      &config.ImportFileList{},
		Export:      &config.ExportFileList{},
		Environment: []string{"X=1", "Y=2", "Z=3"},
		Tasks: map[string]config.Task{
			"x": &config.BuildTask{TaskMeta: config.TaskMeta{Name: "x"}, Dockerfile: "Dockerfile.x"},
			"y": &config.BuildTask{TaskMeta: config.TaskMeta{Name: "y"}, Dockerfile: "Dockerfile.y"},
			"z": &config.BuildTask{TaskMeta: config.TaskMeta{Name: "z"}, Dockerfile: "Dockerfile.z"},
		},
		Plans: map[string]*config.Plan{
			"a": &config.Plan{
				Name: "a",
				Stages: []*config.Stage{
					&config.Stage{
						Name: "w",
						Tasks: []*config.StageTask{
							&config.StageTask{Name: "x"},
							&config.StageTask{Name: "y"},
							&config.StageTask{Name: "z"},
						},
						RunMode: config.RunModeOnSuccess,
					},
				},
			},
			"b": &config.Plan{
				Name: "b",
				Stages: []*config.Stage{
					&config.Stage{
						Name: "q",
						Tasks: []*config.StageTask{
							&config.StageTask{Name: "x", Environment: []string{"Z=4"}},
						},
						RunMode: config.RunModeOnSuccess,
					},
				},
			},
		},
		Metaplans: map[string][]string{
			"default": []string{"a", "b"},
		},
	}))
}

func (s *LoaderSuite) TestLoadFromURL(t sweet.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, _ := ioutil.ReadFile("./test-configs/parent.yaml")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}))

	defer ts.Close()

	loaded, err := NewLoader().Load(ts.URL)
	Expect(err).To(BeNil())
	Expect(loaded).To(Equal(&config.Config{
		Options:     &config.Options{},
		Registries:  []config.Registry{},
		Import:      &config.ImportFileList{},
		Export:      &config.ExportFileList{},
		Environment: []string{"X=1", "Y=2", "Z=3"},
		Tasks:       map[string]config.Task{},
		Plans:       map[string]*config.Plan{},
		Metaplans:   map[string][]string{},
	}))
}

func (s *LoaderSuite) TestLoadExtends(t sweet.T) {
	loaded, err := NewLoader().Load("./test-configs/child.yaml")
	Expect(err).To(BeNil())
	Expect(loaded).To(Equal(&config.Config{
		Options:     &config.Options{},
		Registries:  []config.Registry{},
		Import:      &config.ImportFileList{},
		Export:      &config.ExportFileList{},
		Environment: []string{"X=1", "Y=2", "Z=3", "X=10", "W=20"},
		Tasks:       map[string]config.Task{},
		Plans:       map[string]*config.Plan{},
		Metaplans:   map[string][]string{},
	}))
}

func (s *LoaderSuite) TestOverride(t sweet.T) {
	config := &config.Config{
		Options:     &config.Options{},
		Import:      &config.ImportFileList{},
		Export:      &config.ExportFileList{},
		Environment: []string{"X=1"},
	}

	err := NewLoader().ApplyOverrides(config, []string{
		"./test-configs/override1.yaml",
		"./test-configs/override2.yaml",
	})

	Expect(err).To(BeNil())
	Expect(config.Options.ForceSequential).To(BeTrue())
	Expect(config.Options.HealthcheckInterval).To(Equal(time.Second * 10))
	Expect(config.Environment).To(Equal([]string{
		"X=1", "X=2", "X=3",
	}))
}

func (s *LoaderSuite) TestLoadInvalidSchemaTask(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/invalid-task.yaml")
	Expect(err).To(MatchError("failed to validate task foo: Additional property image is not allowed"))
}

func (s *LoaderSuite) TestLoadInvalidSchemaPlan(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/invalid-plan.yaml")
	Expect(err).To(MatchError("failed to validate plan foo: Invalid type. Expected: object, given: string"))
}

func (s *LoaderSuite) TestLoadInvalidSchemaMetaplan(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/invalid-metaplan.yaml")
	Expect(err).To(MatchError("failed to validate metaplan foo: Invalid type. Expected: string, given: object"))
}

func (s *LoaderSuite) TestLoadExtendsCycle(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/a.yaml")
	Expect(err).To(MatchError("failed to extend cyclic config (test-configs/b.yaml -> test-configs/a.yaml -> test-configs/b.yaml)"))
}

func (s *LoaderSuite) TestOverrideInvalidSchema(t sweet.T) {
	err := NewLoader().ApplyOverrides(&config.Config{}, []string{"./test-configs/invalid-override.yaml"})
	Expect(err).To(MatchError("failed to validate override file: Invalid type. Expected: object, given: string"))
}
