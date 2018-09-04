package network

//go:generate go-mockgen github.com/efritz/ij/command -i Runner -d mocks -f

import (
	"context"
	"fmt"
	"io"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"

	"github.com/efritz/ij/network/mocks"
)

type NetworkSuite struct{}

func (s *NetworkSuite) TestSetupTeardown(t sweet.T) {
	runner := mocks.NewMockRunner()

	network, err := newNetwork(
		context.Background(),
		"abcdef0",
		logging.NilLogger,
		runner,
	)

	Expect(err).To(BeNil())
	Expect(runner.RunForOutputFuncCallCount()).To(Equal(1))
	Expect(runner.RunForOutputFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker", "network", "create", "abcdef0",
	}))

	network.Teardown()
	Expect(runner.RunForOutputFuncCallCount()).To(Equal(2))
	Expect(runner.RunForOutputFuncCallParams()[1].Arg1).To(Equal([]string{
		"docker", "network", "rm", "abcdef0",
	}))
}

func (s *NetworkSuite) TestSetupError(t sweet.T) {
	runner := mocks.NewMockRunner()
	runner.RunForOutputFunc = func(_ context.Context, args []string, _ io.ReadCloser) (string, string, error) {
		return "", "", fmt.Errorf("utoh")
	}

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

	runner := mocks.NewMockRunner()
	runner.RunForOutputFunc = func(ctx context.Context, args []string, _ io.ReadCloser) (string, string, error) {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context canceled")
		default:
		}

		return "", "", nil
	}

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
	runner := mocks.NewMockRunner()

	network, err := newNetwork(
		ctx,
		"abcdef0",
		logging.NilLogger,
		runner,
	)

	Expect(err).To(BeNil())

	cancel()
	network.Teardown()
	Expect(runner.RunForOutputFuncCallCount()).To(Equal(2))

	select {
	case <-runner.RunForOutputFuncCallParams()[1].Arg0.Done():
		t.Fail()
	default:
	}
}
