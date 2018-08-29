package registry

//go:generate go-mockgen github.com/efritz/ij/command -i Runner -d mocks -f

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry/mocks"
)

type ServerSuite struct{}

func (s *ServerSuite) TestLogin(t sweet.T) {
	runner := mocks.NewMockRunner()
	registry := &config.ServerRegistry{
		Server:   "docker.io",
		Username: "admin",
		Password: "secret",
	}

	server, err := newServerLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registry,
		runner,
	).Login()

	Expect(err).To(BeNil())
	Expect(server).To(Equal("docker.io"))

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"admin",
		"--password-stdin",
		"docker.io",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("secret"))
}

func (s *ServerSuite) TestLoginPasswordFile(t sweet.T) {
	runner := mocks.NewMockRunner()
	registry := &config.ServerRegistry{
		Server:       "docker.io",
		Username:     "admin",
		PasswordFile: "./test-files/secret.key",
	}

	server, err := newServerLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registry,
		runner,
	).Login()

	Expect(err).To(BeNil())
	Expect(server).To(Equal("docker.io"))

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"admin",
		"--password-stdin",
		"docker.io",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal("super secret file"))
}

func (s *ServerSuite) TestLoginMappedEnvironment(t sweet.T) {
	runner := mocks.NewMockRunner()
	registry := &config.ServerRegistry{
		Server:   "${DOCKER_HOST}",
		Username: "${DOCKER_USERNAME}",
		Password: "${DOCKER_PASSWORD}",
	}

	env := []string{
		"DOCKER_HOST=docker.io",
		"DOCKER_USERNAME=admin",
		"DOCKER_PASSWORD=secret",
	}

	server, err := newServerLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(env),
		registry,
		runner,
	).Login()

	Expect(err).To(BeNil())
	Expect(server).To(Equal("docker.io"))

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"admin",
		"--password-stdin",
		"docker.io",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("secret"))
}
