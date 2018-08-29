package registry

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry/mocks"
	. "github.com/onsi/gomega"
)

type GCRSuite struct{}

func (s *GCRSuite) TestLogin(t sweet.T) {
	runner := mocks.NewMockRunner()
	registry := &config.GCRRegistry{
		KeyFile: "./test-files/gcr.key",
	}

	login := newGCRLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("https://gcr.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"_json_key",
		"--password-stdin",
		"https://gcr.io",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal(`{"some": "json", "blob": "here"}`))
}

func (s *GCRSuite) TestLoginMappedEnvironment(t sweet.T) {
	runner := mocks.NewMockRunner()
	registry := &config.GCRRegistry{
		KeyFile: "${KEY_FILE}",
	}

	env := []string{
		"KEY_FILE=./test-files/gcr.key",
	}

	login := newGCRLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(env),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("https://gcr.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"_json_key",
		"--password-stdin",
		"https://gcr.io",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal(`{"some": "json", "blob": "here"}`))
}
