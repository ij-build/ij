package paths

import (
	"fmt"
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
	Expect(err.Error()).To(ContainSubstring("import file is outside of workspace directory"))
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
	Expect(err.Error()).To(ContainSubstring("export file is outside of workspace directory"))
}

func (s *TransfererSuite) TestConstructBlacklist(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"a.txt",
		"b.txt",
		"c.txt",
		"x/a.txt",
		"x/junit-other.xml",
		"x/junit.xml",
		"y/b.txt",
		"y/junit.xml",
		"z/c.txt",
		"z/junit.xml",
	}))

	defer os.RemoveAll(name)

	blacklist, err := constructBlacklist(name, []string{
		"**/b.txt",
		"**/junit*.xml",
		"a.txt",
		"z",
	})

	Expect(err).To(BeNil())
	Expect(blacklist).To(HaveKey(filepath.Join(name, "a.txt")))
	Expect(blacklist).To(HaveKey(filepath.Join(name, "b.txt")))
	Expect(blacklist).To(HaveKey(filepath.Join(name, "x", "junit-other.xml")))
	Expect(blacklist).To(HaveKey(filepath.Join(name, "y", "b.txt")))
	Expect(blacklist).To(HaveKey(filepath.Join(name, "y", "junit.xml")))
	Expect(blacklist).To(HaveKey(filepath.Join(name, "z")))
}

func (s *TransfererSuite) TestRunOnPatternIllegal(t sweet.T) {
	err := runOnPattern("foo/*:bar", "", false, false, func(_ filePair) error { return nil })
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("illegal pattern"))
}

func (s *TransfererSuite) TestRunOnSplitPattern(t sweet.T) {
	var src, dest string
	err := runOnSplitPattern("foo:bar", "root", func(fp filePair) error {
		src = fp.src
		dest = fp.dest
		return nil
	})

	Expect(err).To(BeNil())
	Expect(src).To(Equal("root/foo"))
	Expect(dest).To(Equal("root/bar"))
}

func (s *TransfererSuite) TestRunOnSplitPatternError(t sweet.T) {
	err := runOnSplitPattern("foo:bar", "root", func(fp filePair) error {
		return fmt.Errorf("utoh")
	})

	Expect(err).To(MatchError("utoh"))
}

func (s *TransfererSuite) TestRunOnGlobPattern(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"a/foo.txt",
		"b/foo.txt",
		"c/bar.txt",
		"c/foo.txt",
	}))

	defer os.RemoveAll(name)

	var (
		args  = []filePair{}
		path1 = filepath.Join(name, "a", "foo.txt")
		path2 = filepath.Join(name, "b", "foo.txt")
		path3 = filepath.Join(name, "c", "foo.txt")
	)

	err := runOnGlobPattern("**/foo.txt", name, false, func(fp filePair) error {
		args = append(args, fp)
		return nil
	})

	Expect(err).To(BeNil())
	Expect(args).To(ConsistOf(
		filePair{path1, path1},
		filePair{path2, path2},
		filePair{path3, path3},
	))
}

func (s *TransfererSuite) TestRunOnGlobPatternError(t sweet.T) {
	name := buildTempDir(map[string]string{"a/foo.txt": ""})
	defer os.RemoveAll(name)

	err := runOnGlobPattern("**/foo.txt", name, false, func(fp filePair) error {
		return fmt.Errorf("utoh")
	})

	Expect(err).To(MatchError("utoh"))
}

func (s *TransfererSuite) TestRunOnGlobPatternNoMatches(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	called := false
	err := runOnGlobPattern("foo.txt", name, false, func(fp filePair) error {
		called = true
		return nil
	})

	Expect(err).To(BeNil())
	Expect(called).To(BeFalse())
}

func (s *TransfererSuite) TestRunOnGlobPatternNoMatchesStrict(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	err := runOnGlobPattern("foo.txt", name, true, func(fp filePair) error {
		return nil
	})

	Expect(err).NotTo(BeNil())
}

func (s *TransfererSuite) TestSplitPathImplicitDestination(t sweet.T) {
	src, dest := splitPath("a/b/c")
	Expect(src).To(Equal("a/b/c"))
	Expect(dest).To(Equal("a/b/c"))
}

func (s *TransfererSuite) TestSplitPathExplicitDestination(t sweet.T) {
	src, dest := splitPath("a/b/c:x/y/z")
	Expect(src).To(Equal("a/b/c"))
	Expect(dest).To(Equal("x/y/z"))
}
