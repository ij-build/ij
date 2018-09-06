package paths

import (
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RemoverSuite struct{}

func (s *RemoverSuite) TestRemove(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"bar/a.txt",
		"bar/b.txt",
		"bar/c.txt",
		"baz/a.txt",
		"baz/b.txt",
		"baz/c.txt",
		"foo/a.txt",
		"foo/b.txt",
		"foo/c.txt",
	}))

	defer os.RemoveAll(name)

	patterns := []string{
		"**/*.txt",
		"bar",
	}

	blacklistPatterns := []string{
		"baz/b.txt",
		"**/a*",
	}

	err := NewRemover(name).Remove(patterns, blacklistPatterns, func(p string) (bool, error) {
		return true, nil
	})

	Expect(err).To(BeNil())

	// Safe
	_, err = os.Stat(filepath.Join(name, "baz", "a.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(name, "baz", "b.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(name, "foo", "a.txt"))
	Expect(err).To(BeNil())

	// Removed
	_, err = os.Stat(filepath.Join(name, "baz", "c.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(name, "foo", "b.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(name, "foo", "c.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())

	// Removed (entire directory)
	_, err = os.Stat(filepath.Join(name, "bar"))
	Expect(os.IsNotExist(err)).To(BeTrue())
}

func (s *RemoverSuite) TestRemoveDenied(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{
		"a.txt",
		"b.txt",
		"c.txt",
		"d.txt",
		"e.txt",
	}))

	defer os.RemoveAll(name)

	patterns := []string{"*"}
	err := NewRemover(name).Remove(patterns, nil, func(path string) (bool, error) {
		if path == filepath.Join(name, "b.txt") || path == filepath.Join(name, "d.txt") {
			return true, nil
		}

		return false, nil
	})

	Expect(err).To(BeNil())

	// Safe
	_, err = os.Stat(filepath.Join(name, "a.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(name, "c.txt"))
	Expect(err).To(BeNil())
	_, err = os.Stat(filepath.Join(name, "e.txt"))
	Expect(err).To(BeNil())

	// Removed
	_, err = os.Stat(filepath.Join(name, "b.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())
	_, err = os.Stat(filepath.Join(name, "d.txt"))
	Expect(os.IsNotExist(err)).To(BeTrue())
}

func (s *RemoverSuite) TestRemoveOutsideProjectPath(t sweet.T) {
	name := buildTempDir(nil)

	defer os.RemoveAll(name)

	patterns := []string{"../*"}
	err := NewRemover(name).Remove(patterns, nil, func(path string) (bool, error) {
		if path == filepath.Join(name, "b.txt") || path == filepath.Join(name, "d.txt") {
			return true, nil
		}

		return false, nil
	})

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("is outside of workspace directory"))
}
