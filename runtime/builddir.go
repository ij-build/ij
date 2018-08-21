package runtime

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/efritz/pvc/util"
)

type BuildDir struct {
	runID  string
	path   string
	parent string
}

func NewBuildDir(runID string) *BuildDir {
	return &BuildDir{
		runID: runID,
	}
}

func (b *BuildDir) Setup() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	parent := filepath.Join(pwd, ".build")
	path, err := util.BuildPath(parent, b.runID)
	if err != nil {
		return err
	}

	b.path = path
	b.parent = filepath.Join(parent)
	return nil
}

func (b *BuildDir) Prune() error {
	return filepath.Walk(b.path, func(path string, _ os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".err.log") {

			info, err := os.Stat(path)
			if err != nil {
				return err
			}

			if info.Size() == 0 {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (b *BuildDir) Teardown() error {
	if err := os.RemoveAll(b.path); err != nil {
		return err
	}

	names, err := util.Dirnames(b.parent)
	if err != nil {
		return err
	}

	if len(names) == 0 {
		return os.RemoveAll(b.parent)
	}

	return nil
}

func (b *BuildDir) MakeLogFiles(prefix string) (io.WriteCloser, io.WriteCloser, error) {
	outpath, err := util.BuildPath(b.path, "logs", prefix+".out.log")
	if err != nil {
		return nil, nil, err
	}

	errpath, err := util.BuildPath(b.path, "logs", prefix+".err.log")
	if err != nil {
		return nil, nil, err
	}

	outfile, err := os.Create(outpath)
	if err != nil {
		return nil, nil, err
	}

	errfile, err := os.Create(errpath)
	if err != nil {
		outfile.Close()
		return nil, nil, err
	}

	return outfile, errfile, nil
}

func (b *BuildDir) WriteScript(script string) (string, error) {
	scriptID, err := util.MakeID()
	if err != nil {
		return "", err
	}

	path, err := util.BuildPath(b.path, "scripts", scriptID)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(path, []byte(script), 0700); err != nil {
		return "", err
	}

	return path, nil
}
