package registry

import (
	"context"
	"io/ioutil"

	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/environment"
	"github.com/ij-build/ij/logging"
	. "github.com/onsi/gomega"
)

type ECRSuite struct{}

func (s *ECRSuite) TestLogin(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.ECRRegistry{
		AccessKeyID:     "testAccessKeyID",
		SecretAccessKey: "testSecretAccessKey",
		AccountID:       "testAccountID",
		Region:          "testRegion",
		Role:            "testRole",
	}

	runner.RunForOutputFunc.SetDefaultReturn("somehugetoken", "", nil)

	login := newECRLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(nil),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("https://testAccountID.dkr.ecr.testRegion.amazonaws.com"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunForOutputFunc).To(BeCalledOnce())
	Expect(runner.RunForOutputFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"run",
		"--rm",
		"-e", "AWS_ACCESS_KEY_ID=testAccessKeyID",
		"-e", "AWS_SECRET_ACCESS_KEY=testSecretAccessKey",
		"-e", "AWS_ACCOUNT_ID=testAccountID",
		"-e", "AWS_REGION=testRegion",
		"-e", "AWS_ROLE=testRole",
		"efritz/ij-ecr-token:latest",
	}, BeAnything()))

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"AWS",
		"--password-stdin",
		"https://testAccountID.dkr.ecr.testRegion.amazonaws.com",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("somehugetoken"))
}

func (s *ECRSuite) TestLoginMappedEnvironment(t sweet.T) {
	runner := NewMockRunner()
	registry := &config.ECRRegistry{
		AccessKeyID:     "${AWS_ACCESS_KEY_ID}",
		SecretAccessKey: "${AWS_SECRET_ACCESS_KEY}",
		AccountID:       "${AWS_ACCOUNT_ID}",
		Region:          "${AWS_REGION}",
		Role:            "${AWS_ROLE}",
	}

	env := []string{
		"AWS_ACCESS_KEY_ID=testAccessKeyID",
		"AWS_SECRET_ACCESS_KEY=testSecretAccessKey",
		"AWS_ACCOUNT_ID=testAccountID",
		"AWS_REGION=testRegion",
		"AWS_ROLE=testRole",
	}

	runner.RunForOutputFunc.SetDefaultReturn("somehugetoken", "", nil)

	login := newECRLogin(
		context.Background(),
		logging.NilLogger,
		environment.New(env),
		registry,
		runner,
	)

	server, err := login.GetServer()
	Expect(err).To(BeNil())
	Expect(server).To(Equal("https://testAccountID.dkr.ecr.testRegion.amazonaws.com"))
	Expect(login.Login()).To(BeNil())

	Expect(runner.RunForOutputFunc).To(BeCalledOnce())
	Expect(runner.RunForOutputFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"run",
		"--rm",
		"-e", "AWS_ACCESS_KEY_ID=testAccessKeyID",
		"-e", "AWS_SECRET_ACCESS_KEY=testSecretAccessKey",
		"-e", "AWS_ACCOUNT_ID=testAccountID",
		"-e", "AWS_REGION=testRegion",
		"-e", "AWS_ROLE=testRole",
		"efritz/ij-ecr-token:latest",
	}, BeAnything()))

	Expect(runner.RunFunc).To(BeCalledOnce())
	Expect(runner.RunFunc).To(BeCalledWith(BeAnything(), []string{
		"docker",
		"login",
		"-u",
		"AWS",
		"--password-stdin",
		"https://testAccountID.dkr.ecr.testRegion.amazonaws.com",
	}, BeAnything(), BeAnything()))

	stdin := runner.RunFunc.History()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("somehugetoken"))
}
