package paths

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&FilesSuite{})
		s.AddSuite(&CopierSuite{})
		s.AddSuite(&TransfererSuite{})
	})
}

//
// Helpers

func buildTempDir(files map[string]string) string {
	name, _ := ioutil.TempDir("", "ij-test")

	for path, content := range files {
		path = filepath.Join(name, path)
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		ioutil.WriteFile(path, []byte(content), os.ModePerm)
	}

	return name
}

func buildEmptyFiles(keys []string) map[string]string {
	files := map[string]string{}
	for _, key := range keys {
		files[key] = ""
	}

	return files
}
