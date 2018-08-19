package runtime

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/efritz/pvc/paths"
	"github.com/efritz/pvc/util"
)

type Builddir struct {
	runID string
	path  string
}

func NewBuilddir(runID string) *Builddir {
	return &Builddir{
		runID: runID,
	}
}

func (b *Builddir) Setup() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := makePath(pwd, ".build", b.runID)
	if err != nil {
		return err
	}

	b.path = path
	return nil
}

func (b *Builddir) Teardown() error {
	// TODO - nothing right now!
	return nil
}

func (b *Builddir) MakeLogFiles(prefix string) (io.WriteCloser, io.WriteCloser, error) {
	outpath, err := makePath(b.path, "logs", prefix+".out.log")
	if err != nil {
		return nil, nil, err
	}

	errpath, err := makePath(b.path, "logs", prefix+".err.log")
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

func (b *Builddir) WriteScript(script string) (string, error) {
	scriptID, err := util.MakeID()
	if err != nil {
		return "", err
	}

	path, err := makePath(b.path, "scripts", scriptID)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(path, []byte(script), 0700); err != nil {
		return "", err
	}

	return path, nil
}

func makePath(parts ...string) (string, error) {
	fullPath := filepath.Join(parts...)

	if err := paths.EnsureDirExists(filepath.Dir(fullPath)); err != nil {
		return "", err
	}

	return fullPath, nil
}
