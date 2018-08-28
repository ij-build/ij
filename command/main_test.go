package command

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	if os.Getenv(TestEnvFlag) == "1" {
		// If we're invoking our own TestOutput to test with
		// os/exec, don't enable the sweet framework things.
		os.Exit(m.Run())
	}

	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&RunnerSuite{})
		s.AddSuite(&OutputSuite{})
		s.AddSuite(&BuilderSuite{})
	})
}

//
// Helpers

var testArgs = []string{
	os.Args[0],
	"-test.run=TestOutput",
	"--",
}

func TestOutput(t *testing.T) {
	if os.Getenv(TestEnvFlag) != "1" {
		return
	}

	buffer := make([]byte, 64)
	if n, _ := os.Stdin.Read(buffer); n > 0 {
		fmt.Printf("x > %s\n", string(buffer[:n]))
	}

	for i, arg := range os.Args[3:] {
		if strings.ToUpper(arg) != arg {
			fmt.Printf("%d > %s\n", i, arg)
		} else {
			fmt.Fprintf(os.Stderr, "%d > %s\n", i, arg)
		}
	}

	os.Exit(0)
}
