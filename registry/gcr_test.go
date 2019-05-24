package registry

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type GCRSuite struct{}

func (s *GCRSuite) TestLoginKey(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.GCRRegistry{
		Hostname: "gcr.io",
		Key:      `{"some": "json", "blob": "here"}`,
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

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"_json_key",
		"--password-stdin",
		"https://gcr.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal(`{"some": "json", "blob": "here"}`))
}

func (s *GCRSuite) TestLoginKeyMappedEnvironment(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.GCRRegistry{
		Hostname: "eu.gcr.io",
		Key:      `${KEY}`,
	}

	env := []string{
		`KEY={"some": "json", "blob": "here"}`,
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
	Expect(server).To(Equal("https://eu.gcr.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"_json_key",
		"--password-stdin",
		"https://eu.gcr.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal(`{"some": "json", "blob": "here"}`))
}

func (s *GCRSuite) TestLoginKeyFile(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.GCRRegistry{
		Hostname: "gcr.io",
		KeyFile:  "./test-files/gcr.key",
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

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"_json_key",
		"--password-stdin",
		"https://gcr.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal(`{"some": "json", "blob": "here"}`))
}

func (s *GCRSuite) TestLoginKeyFileMappedEnvironment(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.GCRRegistry{
		Hostname: "eu.gcr.io",
		KeyFile:  "${KEY_FILE}",
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
	Expect(server).To(Equal("https://eu.gcr.io"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"_json_key",
		"--password-stdin",
		"https://eu.gcr.io",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(strings.TrimSpace(string(content))).To(Equal(`{"some": "json", "blob": "here"}`))
}
