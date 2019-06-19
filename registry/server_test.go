package registry

//go:generate go-mockgen -f github.com/ij-build/ij/command -i Runner -o mock_runner_test.go

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	. "github.com/efritz/go-mockgen/matchers"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/environment"
	"github.com/ij-build/ij/logging"
)

type ServerSuite struct{}

func (s *ServerSuite) TestLogin(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.ServerRegistry{
		Server:   "docker.io",
		Username: "admin",
		Password: "secret",
	}

	login := newServerLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("docker.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"admin",
		"--password-stdin",
		"docker.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("secret"))
}

func (s *ServerSuite) TestLoginPasswordFile(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.ServerRegistry{
		Server:       "docker.io",
		Username:     "admin",
		PasswordFile: "./test-files/secret.key",
	}

	login := newServerLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("docker.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"admin",
		"--password-stdin",
		"docker.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal("super secret file"))
}

func (s *ServerSuite) TestLoginMappedEnvironment(t sweet.T) {
	runner := NewMockRunner()
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

	login := newServerLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(env),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("docker.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"admin",
		"--password-stdin",
		"docker.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("secret"))
}
