package network

//go:generate go-mockgen -f github.com/efritz/ij/command -i Runner -o mock_runner_test.go

import (
	"context"
	"fmt"
	"io"

	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type NetworkSuite struct{}

func (s *NetworkSuite) TestSetupTeardown(t sweet.T) {
	runner := NewMockRunner()

	network, err := newNetwork(
		context.Background(),
		"abcdef0",
		logging.NilLogger,
		runner,
	)

	Expect(err).To(BeNil())
	Expect(runner.RunForOutputFunc).To(BeCalledOnce())
	Expect(runner.RunForOutputFunc).To(BeCalledWith(BeAnything(), []string{
		"docker", "network", "create", "abcdef0",
	}, BeAnything()))

	network.Teardown()
	Expect(runner.RunForOutputFunc).To(BeCalledN(2))
	Expect(runner.RunForOutputFunc).To(BeCalledWith(BeAnything(), []string{
		"docker", "network", "rm", "abcdef0",
	}, BeAnything()))
}

func (s *NetworkSuite) TestSetupError(t sweet.T) {
	runner := NewMockRunner()
	runner.RunForOutputFunc.SetDefaultReturn("", "", fmt.Errorf("utoh"))

	_, err := newNetwork(
		context.Background(),
		"abcdef0",
		logging.NilLogger,
		runner,
	)

	Expect(err).To(MatchError("utoh"))
}

func (s *NetworkSuite) TestCancelSetup(t sweet.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runner := NewMockRunner()
	runner.RunForOutputFunc.SetDefaultHook(func(ctx context.Context, args []string, _ io.ReadCloser) (string, string, error) {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context canceled")
		default:
		}

		return "", "", nil
	})

	network, err := newNetwork(
		ctx,
		"abcdef0",
		logging.NilLogger,
		runner,
	)

	Expect(err).To(MatchError("context canceled"))
	Expect(network).To(BeNil())
}

func (s *NetworkSuite) TestCancelDuringTeardown(t sweet.T) {
	ctx, cancel := context.WithCancel(context.Background())
	runner := NewMockRunner()

	network, err := newNetwork(
		ctx,
		"abcdef0",
		logging.NilLogger,
		runner,
	)

	Expect(err).To(BeNil())

	cancel()
	network.Teardown()
	Expect(runner.RunForOutputFunc).To(BeCalledN(2))

	select {
	case <-runner.RunForOutputFunc.History()[1].Arg0.Done():
		t.Fail()
	default:
	}
}
