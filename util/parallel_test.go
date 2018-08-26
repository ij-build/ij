package util

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ParallelSuite struct{}

func (s *ParallelSuite) TestRunParallel(t sweet.T) {
	var (
		block = make(chan struct{})
		done  = make(chan struct{})
		args  = make(chan string, 5)
	)

	wg := RunParallel(
		func() { <-block; args <- "a" },
		func() { <-block; args <- "b" },
		func() { <-block; args <- "c" },
		func() { <-block; args <- "d" },
		func() { <-block; args <- "e" },
	)

	go func() {
		defer close(done)
		defer close(args)
		wg.Wait()
	}()

	Consistently(done).ShouldNot(BeClosed())
	close(block)
	Eventually(done).Should(BeClosed())
	Expect(receive(args)).To(ConsistOf("a", "b", "c", "d", "e"))
}

func (s *ParallelSuite) TestRunParallelArgs(t sweet.T) {
	var (
		block = make(chan struct{})
		done  = make(chan struct{})
		args  = make(chan string, 5)
	)

	fn := func(arg string) {
		<-block
		args <- arg
	}

	go func() {
		defer close(done)
		defer close(args)
		RunParallelArgs(fn, "a", "b", "c", "d", "e")
	}()

	Consistently(done).ShouldNot(BeClosed())
	close(block)
	Eventually(done).Should(BeClosed())
	Expect(receive(args)).To(ConsistOf("a", "b", "c", "d", "e"))
}

//
// Helpers

func receive(ch <-chan string) []string {
	vals := []string{}
	for val := range ch {
		vals = append(vals, val)
	}

	return vals
}
