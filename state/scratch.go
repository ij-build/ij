package state

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/util"
)

type ScratchSpace struct {
	runID          string
	keepWorkspace  bool
	project        string
	rootDirFactory func() (string, error)
	scratch        string
	runpath        string
	workspace      string
}

const (
	ScratchDir   = ".ij"
	WorkspaceDir = "workspace"
	ScriptsDir   = "scripts"
	LogsDir      = "logs"
	OutLogSuffix = ".out.log"
	ErrLogSuffix = ".err.log"
)

func NewScratchSpace(runID string, keepWorkspace bool) *ScratchSpace {
	return newScratchSpace(
		runID,
		keepWorkspace,
		os.Getwd,
	)
}

func newScratchSpace(
	runID string,
	keepWorkspace bool,
	rootDirFactory func() (string, error),
) *ScratchSpace {
	return &ScratchSpace{
		runID:          runID,
		keepWorkspace:  keepWorkspace,
		rootDirFactory: rootDirFactory,
	}
}

func (s *ScratchSpace) Setup() error {
	pwd, err := s.rootDirFactory()
	if err != nil {
		return err
	}

	s.project = pwd
	s.scratch = filepath.Join(s.project, ScratchDir)
	s.runpath = filepath.Join(s.scratch, s.runID)
	s.workspace = filepath.Join(s.runpath, WorkspaceDir)

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

	path, err := buildPath(filepath.Join(s.runpath, ScriptsDir, scriptID))
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(path, []byte(script), 0700); err != nil {
		return "", err
	}

	return path, nil
}

func (s *ScratchSpace) MakeLogFiles(prefix string) (*os.File, *os.File, error) {
	outpath, err := buildPath(filepath.Join(s.runpath, LogsDir, prefix+OutLogSuffix))
	if err != nil {
		return nil, nil, err
	}

	errpath, err := buildPath(filepath.Join(s.runpath, LogsDir, prefix+ErrLogSuffix))
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
	if !s.keepWorkspace {
		if err := os.RemoveAll(s.workspace); err != nil {
			return err
		}

		if err := os.RemoveAll(filepath.Join(s.runpath, ScriptsDir)); err != nil {
			return err
		}
	}

	return filepath.Walk(s.runpath, func(path string, _ os.FileInfo, err error) error {
		if strings.HasSuffix(path, ErrLogSuffix) {
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
