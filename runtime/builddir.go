package runtime

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/efritz/pvc/paths"
	"github.com/efritz/pvc/util"
)

type BuildDir struct {
	runID string
	path  string
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

	path, err := makePath(pwd, ".build", b.runID)
	if err != nil {
		return err
	}

	b.path = path
	return nil
}

func (b *BuildDir) MakeLogFiles(prefix string) (io.WriteCloser, io.WriteCloser, error) {
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

func (b *BuildDir) WriteScript(script string) (string, error) {
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
