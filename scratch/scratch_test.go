package scratch

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/paths"
	. "github.com/onsi/gomega"
)

type ScratchSuite struct{}

func (s *ScratchSuite) TestPaths(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, true)
	Expect(scratch.Setup()).To(BeNil())
	Expect(scratch.Project()).To(Equal(name))
	Expect(scratch.Scratch()).To(Equal(filepath.Join(name, ".ij")))
	Expect(scratch.Runpath()).To(Equal(filepath.Join(name, ".ij", "abcdef0")))
	Expect(scratch.Workspace()).To(Equal(filepath.Join(name, ".ij", "abcdef0", "workspace")))

	info, err := os.Stat(scratch.Workspace())
	Expect(err).To(BeNil())
	Expect(info.IsDir()).To(BeTrue())
}

func (s *ScratchSuite) TestWriteScript(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, true)
	scratch.Setup()

	scriptPath, err := scratch.WriteScript("foo\nbar\nbaz\n")
	Expect(err).To(BeNil())
	Expect(scriptPath).To(HavePrefix(filepath.Join(name, ".ij", "abcdef0")))

	content, err := ioutil.ReadFile(scriptPath)
	Expect(err).To(BeNil())
	Expect(string(content)).To(Equal("foo\nbar\nbaz\n"))
}

func (s *ScratchSuite) TestMakeLogFiles(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, true)
	scratch.Setup()

	outFile, errFile, err := scratch.MakeLogFiles("foo/bar/baz")
	Expect(err).To(BeNil())

	outFile.Write([]byte("log1\n"))
	errFile.Write([]byte("log2\n"))
	outFile.Write([]byte("log3\n"))
	errFile.Write([]byte("log4\n"))

	outFile.Close()
	errFile.Close()

	Expect(outFile.Name()).To(Equal(filepath.Join(name, ".ij", "abcdef0", "logs", "foo", "bar", "baz.out.log")))
	Expect(errFile.Name()).To(Equal(filepath.Join(name, ".ij", "abcdef0", "logs", "foo", "bar", "baz.err.log")))

	outContent, err := ioutil.ReadFile(outFile.Name())
	Expect(err).To(BeNil())
	Expect(string(outContent)).To(Equal("log1\nlog3\n"))

	errContent, err := ioutil.ReadFile(errFile.Name())
	Expect(err).To(BeNil())
	Expect(string(errContent)).To(Equal("log2\nlog4\n"))
}

func (s *ScratchSuite) TestPrune(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, true)
	scratch.Setup()

	paths := []string{
		"a1",
		"a2",
		"p1/b",
		"p1/c/d",
		"p2/b",
		"p2/c/d",
	}

	// Populate workspace and scripts
	err := ioutil.WriteFile(filepath.Join(scratch.Workspace(), "f1"), []byte("a"), os.ModePerm)
	Expect(err).To(BeNil())
	_, err = scratch.WriteScript("foo")
	Expect(err).To(BeNil())

	// Populate log dir

	for _, path := range paths {
		outFile, errFile, err := scratch.MakeLogFiles(path)
		Expect(err).To(BeNil())
		outFile.Close()
		errFile.Close()
	}

	Expect(scratch.Prune()).To(BeNil())

	for _, path := range paths {
		_, outErr := os.Stat(filepath.Join(name, ".ij", "abcdef0", "logs", path+".out.log"))
		_, errErr := os.Stat(filepath.Join(name, ".ij", "abcdef0", "logs", path+".err.log"))
		Expect(outErr).To(BeNil())
		Expect(os.IsNotExist(errErr)).To(BeTrue())
	}

	// Workspace is not removed
	_, err = os.Stat(scratch.Workspace())
	Expect(err).To(BeNil())

	// Scripts dir is not removed
	_, err = os.Stat(scratch.Runpath() + "/scripts")
	Expect(err).To(BeNil())
}

func (s *ScratchSuite) TestPruneDiscardWorkspace(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, false)
	scratch.Setup()

	// Populate workspace and scripts
	err := ioutil.WriteFile(filepath.Join(scratch.Workspace(), "f1"), []byte("a"), os.ModePerm)
	Expect(err).To(BeNil())
	_, err = scratch.WriteScript("foo")
	Expect(err).To(BeNil())

	// Populate log dir
	outFile, errFile, err := scratch.MakeLogFiles("build")
	Expect(err).To(BeNil())
	outFile.Close()
	errFile.Close()

	Expect(scratch.Prune()).To(BeNil())

	// Output logs are not removed
	_, err = os.Stat(filepath.Join(name, ".ij", "abcdef0", "logs", "build.out.log"))
	Expect(err).To(BeNil())

	// Workspace is removed
	_, err = os.Stat(scratch.Workspace())
	Expect(os.IsNotExist(err)).To(BeTrue())

	// Scripts dir is removed
	_, err = os.Stat(scratch.Runpath() + "/scripts")
	Expect(os.IsNotExist(err)).To(BeTrue())
}

func (s *ScratchSuite) TestTeardown(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, true)
	scratch.Setup()
	Expect(scratch.Teardown()).To(BeNil())
	_, err := os.Stat(scratch.Scratch())
	Expect(os.IsNotExist(err)).To(BeTrue())
}

func (s *ScratchSuite) TestTeardownMultipleProjects(t sweet.T) {
	name, _ := ioutil.TempDir("", "ij-test")
	defer os.RemoveAll(name)

	scratch := NewScratchSpace("abcdef0", name, name, true)

	// Make another project first
	paths.EnsureDirExists(filepath.Join(name, ".ij", "abcdef1", "workspace", name), os.ModePerm)

	scratch.Setup()
	Expect(scratch.Teardown()).To(BeNil())
	_, err := os.Stat(scratch.Runpath())
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(scratch.Scratch())
	Expect(err).To(BeNil())
}
