package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	"github.com/ij-build/ij/logging"
	. "github.com/onsi/gomega"
)

type PatternsSuite struct{}

func (s *PatternsSuite) TestConstructBlacklist(t sweet.T) {
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

func (s *PatternsSuite) TestConstructBlacklistWithDestination(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"foo/txt",
	}))

	defer os.RemoveAll(name)

	_, err := constructBlacklist(name, []string{
		"foo.txt:bar.txt",
	})

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("blacklist contains destination path"))
}

func (s *PatternsSuite) TestRunOnPatternIllegal(t sweet.T) {
	err := runOnPattern("foo/*:bar", "", logging.NilLogger, func(_ FilePair) error { return nil })
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("illegal pattern"))
}

func (s *PatternsSuite) TestRunOnSplitPattern(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"foo",
	}))

	var src, dest string
	err := runOnSplitPattern("foo:bar", name, logging.NilLogger, func(fp FilePair) error {
		src = fp.Src
		dest = fp.Dest
		return nil
	})

	Expect(err).To(BeNil())
	Expect(src).To(Equal(filepath.Join(name, "foo")))
	Expect(dest).To(Equal(filepath.Join(name, "bar")))
}

func (s *PatternsSuite) TestRunOnSplitPatternError(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"foo",
	}))

	err := runOnSplitPattern("foo:bar", name, logging.NilLogger, func(fp FilePair) error {
		return fmt.Errorf("utoh")
	})

	Expect(err).To(MatchError("utoh"))
}

func (s *PatternsSuite) TestRunOnGlobPattern(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"a/foo.txt",
		"b/foo.txt",
		"c/bar.txt",
		"c/foo.txt",
	}))

	defer os.RemoveAll(name)

	var (
		args  = []FilePair{}
		path1 = filepath.Join(name, "a", "foo.txt")
		path2 = filepath.Join(name, "b", "foo.txt")
		path3 = filepath.Join(name, "c", "foo.txt")
	)

	err := runOnGlobPattern("**/foo.txt", name, logging.NilLogger, func(fp FilePair) error {
		args = append(args, fp)
		return nil
	})

	Expect(err).To(BeNil())
	Expect(args).To(ConsistOf(
		FilePair{path1, path1},
		FilePair{path2, path2},
		FilePair{path3, path3},
	))
}

func (s *PatternsSuite) TestRunOnGlobPatternError(t sweet.T) {
	name := buildTempDir(map[string]string{"a/foo.txt": ""})
	defer os.RemoveAll(name)

	err := runOnGlobPattern("**/foo.txt", name, logging.NilLogger, func(fp FilePair) error {
		return fmt.Errorf("utoh")
	})

	Expect(err).To(MatchError("utoh"))
}

func (s *PatternsSuite) TestRunOnGlobPatternNoMatches(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	called := false
	err := runOnGlobPattern("foo.txt", name, logging.NilLogger, func(fp FilePair) error {
		called = true
		return nil
	})

	Expect(err).To(BeNil())
	Expect(called).To(BeFalse())
}

func (s *PatternsSuite) TestSplitPathImplicitDestination(t sweet.T) {
	src, dest := splitPath("a/b/c")
	Expect(src).To(Equal("a/b/c"))
	Expect(dest).To(Equal("a/b/c"))
}

func (s *PatternsSuite) TestSplitPathExplicitDestination(t sweet.T) {
	src, dest := splitPath("a/b/c:x/y/z")
	Expect(src).To(Equal("a/b/c"))
	Expect(dest).To(Equal("x/y/z"))
}
