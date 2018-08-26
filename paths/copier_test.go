package paths

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type CopierSuite struct{}

func (s *CopierSuite) TestCopyFile(t sweet.T) {
	srcRoot := buildTempDir(map[string]string{
		"x.txt": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	})

	defer os.RemoveAll(srcRoot)

	destRoot := buildTempDir(nil)
	defer os.RemoveAll(destRoot)

	copier := NewCopier(
		logging.NilLogger,
		srcRoot,
		map[string]struct{}{},
	)

	err := copier.Copy(
		filepath.Join(srcRoot, "x.txt"),
		filepath.Join(destRoot, "x.txt"),
	)

	Expect(err).To(BeNil())

	content, err := ioutil.ReadFile(filepath.Join(destRoot, "x.txt"))
	Expect(err).To(BeNil())
	Expect(string(content)).To(Equal("Lorem ipsum dolor sit amet, consectetur adipiscing elit."))
}

func (s *CopierSuite) TestCopyDirectory(t sweet.T) {
	srcRoot := buildTempDir(map[string]string{
		"foo/a.txt":         "Aenean euismod ultrices molestie.",
		"foo/bar/b.txt":     "Vivamus eu bibendum urna.",
		"foo/bar/baz/c.txt": "Etiam viverra ac arcu eget elementum.",
	})

	defer os.RemoveAll(srcRoot)

	destRoot := buildTempDir(nil)
	defer os.RemoveAll(destRoot)

	copier := NewCopier(
		logging.NilLogger,
		srcRoot,
		map[string]struct{}{},
	)

	err := copier.Copy(
		filepath.Join(srcRoot, "foo"),
		filepath.Join(destRoot, "foo"),
	)

	Expect(err).To(BeNil())

	content1, err := ioutil.ReadFile(filepath.Join(destRoot, "foo", "a.txt"))
	Expect(err).To(BeNil())
	Expect(string(content1)).To(Equal("Aenean euismod ultrices molestie."))

	content2, err := ioutil.ReadFile(filepath.Join(destRoot, "foo", "bar", "b.txt"))
	Expect(err).To(BeNil())
	Expect(string(content2)).To(Equal("Vivamus eu bibendum urna."))

	content3, err := ioutil.ReadFile(filepath.Join(destRoot, "foo", "bar", "baz", "c.txt"))
	Expect(err).To(BeNil())
	Expect(string(content3)).To(Equal("Etiam viverra ac arcu eget elementum."))
}

func (s *CopierSuite) TestCopyBlacklistedFile(t sweet.T) {
	srcRoot := buildTempDir(map[string]string{
		"foo/a.txt":         "Aenean euismod ultrices molestie.",
		"foo/bar/b.txt":     "Vivamus eu bibendum urna.",
		"foo/bar/baz/c.txt": "Etiam viverra ac arcu eget elementum.",
	})

	defer os.RemoveAll(srcRoot)

	destRoot := buildTempDir(nil)
	defer os.RemoveAll(destRoot)

	copier := NewCopier(
		logging.NilLogger,
		srcRoot,
		map[string]struct{}{
			filepath.Join(srcRoot, "foo", "bar", "b.txt"): struct{}{},
		},
	)

	err := copier.Copy(
		filepath.Join(srcRoot, "foo"),
		filepath.Join(destRoot, "foo"),
	)

	Expect(err).To(BeNil())

	_, err = os.Stat(filepath.Join(destRoot, "foo", "a.txt"))
	Expect(err).To(BeNil())

	_, err = os.Stat(filepath.Join(destRoot, "foo", "bar", "b.txt"))
	Expect(err).NotTo(BeNil())

	_, err = os.Stat(filepath.Join(destRoot, "foo", "bar", "baz", "c.txt"))
	Expect(err).To(BeNil())
}

func (s *CopierSuite) TestCopyBlacklistedDirectory(t sweet.T) {
	srcRoot := buildTempDir(map[string]string{
		"foo/a.txt":         "Aenean euismod ultrices molestie.",
		"foo/bar/b.txt":     "Vivamus eu bibendum urna.",
		"foo/bar/baz/c.txt": "Etiam viverra ac arcu eget elementum.",
	})

	defer os.RemoveAll(srcRoot)

	destRoot := buildTempDir(nil)
	defer os.RemoveAll(destRoot)

	copier := NewCopier(
		logging.NilLogger,
		srcRoot,
		map[string]struct{}{
			filepath.Join(srcRoot, "foo", "bar"): struct{}{},
		},
	)

	err := copier.Copy(
		filepath.Join(srcRoot, "foo"),
		filepath.Join(destRoot, "foo"),
	)

	Expect(err).To(BeNil())

	_, err = os.Stat(filepath.Join(destRoot, "foo", "a.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(destRoot, "foo", "bar", "b.txt"))
	Expect(err).NotTo(BeNil())
	_, err = os.Stat(filepath.Join(destRoot, "foo", "bar", "baz", "c.txt"))
	Expect(err).NotTo(BeNil())
}
