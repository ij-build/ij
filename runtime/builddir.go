package runtime

import (
	"io"
	"os"
	"path/filepath"
)

type Builddir struct {
	path string
}

func NewBuilddir() *Builddir {
	return &Builddir{}
}

func (b *Builddir) Setup() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(pwd, ".build")

	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	b.path = path
	return nil
}

func (b *Builddir) Teardown() error {
	// TODO - nothing right now!
	return nil
}

func (b *Builddir) LogFiles(prefix string) (io.WriteCloser, io.WriteCloser, error) {
	outpath := filepath.Join(b.path, prefix+".out.log")
	errpath := filepath.Join(b.path, prefix+".err.log")

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
