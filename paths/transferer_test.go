package paths

import (
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type TransfererSuite struct{}

func (s *TransfererSuite) TestImport(t sweet.T) {
	srcRoot := buildTempDir(buildEmptyFiles([]string{
		"bar/b.txt",
		"bar/junit-other.xml",
		"bar/junit.xml",
		"baz/c.txt",
		"baz/junit.xml",
		"foo/a.txt",
		"foo/baz.txt",
		"foo/junit.xml",
		"quux/unique.txt",
	}))

	defer os.RemoveAll(srcRoot)

	destRoot := buildTempDir(buildEmptyFiles([]string{
		"workspace/empty",
	}))

	defer os.RemoveAll(destRoot)

	transferer := NewTransferer(
		srcRoot,
		destRoot,
		filepath.Join(destRoot, "workspace"),
		logging.NilLogger,
	)

	patterns := []string{
		"foo",
		"b*r",
		"quux/unique.txt:moved.txt",
	}

	blacklistPatterns := []string{
		"baz",
		"**/junit*.xml",
	}

	err := transferer.Import(patterns, blacklistPatterns)
	Expect(err).To(BeNil())

	// Moved
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "bar", "b.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "foo", "a.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "foo", "baz.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "moved.txt"))
	Expect(err).To(BeNil())

	// Omitted
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "bar", "junit-other.xml"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "bar", "junit.xml"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "baz", "c.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "baz", "junit.xml"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(destRoot, "workspace", "foo", "junit.xml"))
	Expect(os.IsNotExist(err)).To(BeTrue())
}

func (s *TransfererSuite) TestImportOutsideProjectPath(t sweet.T) {
	srcRoot := buildTempDir(nil)
	defer os.RemoveAll(srcRoot)
	destRoot := buildTempDir(nil)
	defer os.RemoveAll(destRoot)

	transferer := NewTransferer(
		srcRoot,
		destRoot,
		filepath.Join(destRoot, "workspace"),
		logging.NilLogger,
	)

	patterns := []string{"../*"}
	err := transferer.Import(patterns, nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("is outside of workspace directory"))
}

func (s *TransfererSuite) TestExport(t sweet.T) {
	srcRoot := buildTempDir(nil)
	defer os.RemoveAll(srcRoot)

	destRoot := buildTempDir(buildEmptyFiles([]string{
		"workspace/foo.txt",
		"workspace/bar.txt",
		"workspace/baz.pdf",
		"workspace/sub/foo.txt",
	}))

	defer os.RemoveAll(destRoot)

	transferer := NewTransferer(
		srcRoot,
		destRoot,
		filepath.Join(destRoot, "workspace"),
		logging.NilLogger,
	)

	patterns := []string{
		"*.txt",
		"baz.pdf:bonk.pdf",
	}

	err := transferer.Export(patterns, nil)
	Expect(err).To(BeNil())

	// Moved
	_, err = os.Stat(filepath.Join(srcRoot, "foo.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(srcRoot, "bar.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(srcRoot, "bonk.pdf"))
	Expect(err).To(BeNil())

	// Omitted
	_, err = os.Stat(filepath.Join(srcRoot, "baz.pdf"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(srcRoot, "sub", "foo.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())
}

func (s *TransfererSuite) TestExportOutsideWorkspacePath(t sweet.T) {
	srcRoot := buildTempDir(nil)
	defer os.RemoveAll(srcRoot)
	destRoot := buildTempDir(nil)
	defer os.RemoveAll(destRoot)

	transferer := NewTransferer(
		srcRoot,
		destRoot,
		filepath.Join(destRoot, "workspace"),
		logging.NilLogger,
	)

	patterns := []string{"../../*"}
	err := transferer.Export(patterns, nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("is outside of workspace directory"))
}
