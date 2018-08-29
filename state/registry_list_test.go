package state

//go:generate go-mockgen github.com/efritz/ij/registry -i Login -d mocks -f

import (
	"context"
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry"
	"github.com/efritz/ij/state/mocks"
)

type RegistryListSuite struct{}

func (s *RegistryListSuite) TestSetupTeardown(t sweet.T) {
	var (
		login   = mocks.NewMockLogin()
		runner  = mocks.NewMockRunner()
		args    = []config.Registry{}
		servers = make(chan string, 3)
	)

	servers <- "x"
	servers <- "y"
	servers <- "z"
	close(servers)

	registries := []config.Registry{
		&config.GCRRegistry{KeyFile: "a"},
		&config.GCRRegistry{KeyFile: "b"},
		&config.GCRRegistry{KeyFile: "c"},
	}

	factory := func(
		_ context.Context,
		_ logging.Logger,
		_ environment.Environment,
		arg config.Registry,
	) registry.Login {
		args = append(args, arg)
		return login
	}

	login.GetServerFunc = func() (string, error) {
		return <-servers, nil
	}

	registryList, err := newRegistryList(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registries,
		factory,
		runner,
	)

	Expect(err).To(BeNil())
	Expect(args).To(ConsistOf(registries[0], registries[1], registries[2]))
	Expect(login.LoginFuncCallCount()).To(Equal(3))

	registryList.Teardown()
	Expect(runner.RunFuncCallCount()).To(Equal(3))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{"docker", "logout", "x"}))
	Expect(runner.RunFuncCallParams()[1].Arg1).To(Equal([]string{"docker", "logout", "y"}))
	Expect(runner.RunFuncCallParams()[2].Arg1).To(Equal([]string{"docker", "logout", "z"}))
}

func (s *RegistryListSuite) TestSetupError(t sweet.T) {
	var (
		login   = mocks.NewMockLogin()
		runner  = mocks.NewMockRunner()
		args    = []config.Registry{}
		servers = make(chan string, 3)
		errors  = make(chan error, 3)
	)

	servers <- "x"
	servers <- "y"
	servers <- "z"
	close(servers)

	errors <- nil
	errors <- nil
	errors <- fmt.Errorf("utoh")
	close(errors)

	registries := []config.Registry{
		&config.GCRRegistry{KeyFile: "a"},
		&config.GCRRegistry{KeyFile: "b"},
		&config.GCRRegistry{KeyFile: "c"},
	}

	factory := func(
		_ context.Context,
		_ logging.Logger,
		_ environment.Environment,
		arg config.Registry,
	) registry.Login {
		args = append(args, arg)
		return login
	}

	login.GetServerFunc = func() (string, error) {
		return <-servers, nil
	}

	login.LoginFunc = func() error {
		return <-errors
	}

	_, err := newRegistryList(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registries,
		factory,
		runner,
	)

	Expect(err).To(MatchError("utoh"))
	Expect(runner.RunFuncCallCount()).To(Equal(2))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{"docker", "logout", "x"}))
	Expect(runner.RunFuncCallParams()[1].Arg1).To(Equal([]string{"docker", "logout", "y"}))
}

func (s *RegistryListSuite) TestCancelSetup(t sweet.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	registries := []config.Registry{
		&config.GCRRegistry{KeyFile: "a"},
		&config.GCRRegistry{KeyFile: "b"},
		&config.GCRRegistry{KeyFile: "c"},
	}

	_, err := newRegistryList(
		ctx,
		logging.NilLogger,
		environment.New(nil),
		registries,
		testLoginFactory,
		mocks.NewMockRunner(),
	)

	Expect(err).To(MatchError("context canceled"))
}

func (s *RegistryListSuite) TestCancelDuringTeardown(t sweet.T) {
	ctx, cancel := context.WithCancel(context.Background())
	runner := mocks.NewMockRunner()

	registries := []config.Registry{
		&config.GCRRegistry{KeyFile: "a"},
		&config.GCRRegistry{KeyFile: "b"},
		&config.GCRRegistry{KeyFile: "c"},
	}

	registryList, err := newRegistryList(
		ctx,
		logging.NilLogger,
		environment.New(nil),
		registries,
		testLoginFactory,
		runner,
	)

	Expect(err).To(BeNil())

	cancel()
	registryList.Teardown()
	Expect(runner.RunFuncCallCount()).To(Equal(3))

	select {
	case <-runner.RunFuncCallParams()[1].Arg0.Done():
		t.Fail()
	default:
	}
}

//
// Build

func testLoginFactory(
	ctx context.Context,
	_ logging.Logger,
	_ environment.Environment,
	arg config.Registry,
) registry.Login {
	login := mocks.NewMockLogin()

	login.LoginFunc = func() error {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		default:
		}

		return nil
	}

	return login
}
