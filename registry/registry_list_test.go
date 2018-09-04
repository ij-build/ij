package registry

//go:generate go-mockgen github.com/efritz/ij/registry -i Login -d mocks -f

import (
	"context"
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry/mocks"
)

type RegistrySetSuite struct{}

func (s *RegistrySetSuite) TestLoginLogout(t sweet.T) {
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
	) Login {
		args = append(args, arg)
		return login
	}

	login.GetServerFunc = func() (string, error) {
		return <-servers, nil
	}

	registrySet, err := newRegistrySet(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registries,
		factory,
		runner,
	)

	Expect(err).To(BeNil())
	Expect(registrySet.Login()).To(BeNil())
	Expect(args).To(ConsistOf(registries[0], registries[1], registries[2]))
	Expect(login.LoginFuncCallCount()).To(Equal(3))

	registrySet.Logout()
	Expect(runner.RunFuncCallCount()).To(Equal(3))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{"docker", "logout", "x"}))
	Expect(runner.RunFuncCallParams()[1].Arg1).To(Equal([]string{"docker", "logout", "y"}))
	Expect(runner.RunFuncCallParams()[2].Arg1).To(Equal([]string{"docker", "logout", "z"}))
}

func (s *RegistrySetSuite) TestSetupError(t sweet.T) {
	var (
		login  = mocks.NewMockLogin()
		runner = mocks.NewMockRunner()
		errors = make(chan error, 3)
	)

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
	) Login {
		return login
	}

	login.GetServerFunc = func() (string, error) {
		return "", <-errors
	}

	_, err := newRegistrySet(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registries,
		factory,
		runner,
	)

	Expect(err).To(MatchError("utoh"))
}

func (s *RegistrySetSuite) TestLoginError(t sweet.T) {
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
	) Login {
		args = append(args, arg)
		return login
	}

	login.GetServerFunc = func() (string, error) {
		return <-servers, nil
	}

	login.LoginFunc = func() error {
		return <-errors
	}

	registrySet, err := newRegistrySet(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registries,
		factory,
		runner,
	)

	Expect(err).To(BeNil())
	Expect(registrySet.Login()).To(MatchError("utoh"))
	Expect(runner.RunFuncCallCount()).To(Equal(2))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{"docker", "logout", "x"}))
	Expect(runner.RunFuncCallParams()[1].Arg1).To(Equal([]string{"docker", "logout", "y"}))
}

func (s *RegistrySetSuite) TestCancelSetup(t sweet.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	registries := []config.Registry{
		&config.GCRRegistry{KeyFile: "a"},
		&config.GCRRegistry{KeyFile: "b"},
		&config.GCRRegistry{KeyFile: "c"},
	}

	registrySet, err := newRegistrySet(
		ctx,
		logging.NilLogger,
		environment.New(nil),
		registries,
		testLoginFactory,
		mocks.NewMockRunner(),
	)

	Expect(err).To(BeNil())
	Expect(registrySet.Login()).To(MatchError("context canceled"))
}

func (s *RegistrySetSuite) TestCancelDuringTeardown(t sweet.T) {
	ctx, cancel := context.WithCancel(context.Background())
	runner := mocks.NewMockRunner()

	registries := []config.Registry{
		&config.GCRRegistry{KeyFile: "a"},
		&config.GCRRegistry{KeyFile: "b"},
		&config.GCRRegistry{KeyFile: "c"},
	}

	registrySet, err := newRegistrySet(
		ctx,
		logging.NilLogger,
		environment.New(nil),
		registries,
		testLoginFactory,
		runner,
	)

	Expect(err).To(BeNil())
	Expect(registrySet.Login()).To(BeNil())

	cancel()
	registrySet.Logout()
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
) Login {
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
