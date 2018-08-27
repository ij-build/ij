package logging

import (
	"errors"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type UtilSuite struct{}

func (s *UtilSuite) TestWriteAll(t sweet.T) {
	var (
		w    = &slowWriter{}
		data = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	)

	err := writeAll(w, data)
	Expect(err).To(BeNil())
	Expect(w.data).To(Equal(data))
	Expect(w.numCalls).To(Equal(5))
}

func (s *UtilSuite) TestWriteAllError(t sweet.T) {
	var (
		w    = &failingSlowWriter{}
		data = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	)

	Expect(writeAll(w, data)).To(MatchError("utoh"))
	Expect(w.data).To(Equal([]byte{1, 2, 3, 4, 5, 6}))
	Expect(w.numCalls).To(Equal(4))
}

//
//

type slowWriter struct {
	numCalls int
	data     []byte
}

func (w *slowWriter) Write(p []byte) (int, error) {
	w.numCalls++
	w.data = append(w.data, p[:2]...)
	return 2, nil
}

//
//

type failingSlowWriter struct {
	numCalls int
	data     []byte
}

func (w *failingSlowWriter) Write(p []byte) (int, error) {
	w.numCalls++

	if len(w.data) > 5 {
		return 0, errors.New("utoh")
	}

	w.data = append(w.data, p[:2]...)
	return 2, nil
}
