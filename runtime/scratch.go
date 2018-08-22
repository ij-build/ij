package runtime

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/util"
)

type ScratchSpace struct {
	runID     string
	project   string
	scratch   string
	runpath   string
	workspace string
}

func NewScratchSpace(runID string) *ScratchSpace {
	return &ScratchSpace{
		runID: runID,
	}
}

func (s *ScratchSpace) Setup() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	s.project = pwd
	s.scratch = filepath.Join(s.project, ".ij")
	s.runpath = filepath.Join(s.scratch, s.runID)
	s.workspace = filepath.Join(s.runpath, "workspace")

	return paths.EnsureDirExists(s.workspace, os.ModePerm)
}

func (s *ScratchSpace) Teardown() error {
	if err := os.RemoveAll(s.runpath); err != nil {
		return err
	}

	entries, err := paths.DirContents(s.scratch)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return os.RemoveAll(s.scratch)
	}

	return nil
}

func (s *ScratchSpace) Project() string {
	return s.project
}

func (s *ScratchSpace) Scratch() string {
	return s.scratch
}

func (s *ScratchSpace) Runpath() string {
	return s.runpath
}

func (s *ScratchSpace) Workspace() string {
	return s.workspace
}

func (s *ScratchSpace) WriteScript(script string) (string, error) {
	scriptID, err := util.MakeID()
	if err != nil {
		return "", err
	}

	path, err := buildPath(filepath.Join(s.runpath, "scripts", scriptID))
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(path, []byte(script), 0700); err != nil {
		return "", err
	}

	return path, nil
}

func (s *ScratchSpace) MakeLogFiles(prefix string) (io.WriteCloser, io.WriteCloser, error) {
	outpath, err := buildPath(filepath.Join(s.runpath, "logs", prefix+".out.log"))
	if err != nil {
		return nil, nil, err
	}

	errpath, err := buildPath(filepath.Join(s.runpath, "logs", prefix+".err.log"))
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

func (s *ScratchSpace) Prune() error {
	return filepath.Walk(s.runpath, func(path string, _ os.FileInfo, err error) error {
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

//
// Helpers

func buildPath(path string) (string, error) {
	if err := paths.EnsureParentExists(path, os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}
