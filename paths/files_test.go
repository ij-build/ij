package paths

import (
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FilesSuite struct{}

func (s *FilesSuite) TestDirContents(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"a.txt",
		"b.txt",
		"c.txt",
		"sub/a.txt",
	}))

	defer os.RemoveAll(name)

	infos, err := DirContents(name)
	Expect(err).To(BeNil())

	names := []string{}
	for _, info := range infos {
		names = append(names, info.Name())
	}

	Expect(names).To(ConsistOf("a.txt", "b.txt", "c.txt", "sub"))
}

func (s *FilesSuite) TestFileExists(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{"a.txt"}))
	defer os.RemoveAll(name)

	exists, err := FileExists(filepath.Join(name, "a.txt"))
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
}

func (s *FilesSuite) TestFileExistsNoFile(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	exists, err := FileExists(filepath.Join(name, "a.txt"))
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
}

func (s *FilesSuite) TestFileExistsDirectory(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{"sub/a.txt"}))
	defer os.RemoveAll(name)

	_, err := FileExists(filepath.Join(name, "sub"))
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("not a file"))
}

func (s *FilesSuite) TestDirExists(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{"sub/a.txt"}))
	defer os.RemoveAll(name)

	exists, err := DirExists(filepath.Join(name, "sub"))
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
}

func (s *FilesSuite) TestDirExistsNoFile(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	exists, err := DirExists(filepath.Join(name, "sub"))
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
}

func (s *FilesSuite) TestDirExistsFile(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{"sub"}))
	defer os.RemoveAll(name)

	_, err := DirExists(filepath.Join(name, "sub"))
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("not a directory"))
}

func (s *FilesSuite) TestEnsureDirExist(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	target := filepath.Join(name, "a/b/c/d/e")
	Expect(EnsureDirExists(target, os.ModePerm)).To(BeNil())
	info, err := os.Stat(target)
	Expect(err).To(BeNil())
	Expect(info.IsDir()).To(BeTrue())
}

func (s *FilesSuite) TestEnsureParentExists(t sweet.T) {
	name := buildTempDir(nil)
	defer os.RemoveAll(name)

	var (
		parent = filepath.Join(name, "a/b/c/d/e")
		target = filepath.Join(parent, "foo.txt")
	)

	err := EnsureParentExists(target, os.ModePerm)
	Expect(err).To(BeNil())

	info, err := os.Stat(parent)
	Expect(err).To(BeNil())
	Expect(info.IsDir()).To(BeTrue())
}
