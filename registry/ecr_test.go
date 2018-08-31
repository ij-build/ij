package registry

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry/mocks"
	. "github.com/onsi/gomega"
)

type ECRSuite struct{}

func (s *ECRSuite) TestLogin(t sweet.T) {
	runner := mocks.NewMockRunner()
	registry := &config.ECRRegistry{
		AccessKeyID:     "testAccessKeyID",
		SecretAccessKey: "testSecretAccessKey",
		AccountID:       "testAccountID",
		Region:          "testRegion",
		Role:            "testRole",
	}

	runner.RunForOutputFunc = func(_ context.Context, _ []string, _ io.ReadCloser) (string, string, error) {
		return "somehugetoken", "", nil
	}

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

	Expect(runner.RunForOutputFuncCallCount()).To(Equal(1))
	Expect(runner.RunForOutputFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"run",
		"--rm",
		"-e", "AWS_ACCESS_KEY_ID=testAccessKeyID",
		"-e", "AWS_SECRET_ACCESS_KEY=testSecretAccessKey",
		"-e", "AWS_ACCOUNT_ID=testAccountID",
		"-e", "AWS_REGION=testRegion",
		"-e", "AWS_ROLE=testRole",
		"efritz/ij-ecr-token:latest",
	}))

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"AWS",
		"--password-stdin",
		"https://testAccountID.dkr.ecr.testRegion.amazonaws.com",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("somehugetoken"))
}

func (s *ECRSuite) TestLoginMappedEnvironment(t sweet.T) {
	runner := mocks.NewMockRunner()
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

	runner.RunForOutputFunc = func(_ context.Context, _ []string, _ io.ReadCloser) (string, string, error) {
		return "somehugetoken", "", nil
	}

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

	Expect(runner.RunForOutputFuncCallCount()).To(Equal(1))
	Expect(runner.RunForOutputFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"run",
		"--rm",
		"-e", "AWS_ACCESS_KEY_ID=testAccessKeyID",
		"-e", "AWS_SECRET_ACCESS_KEY=testSecretAccessKey",
		"-e", "AWS_ACCOUNT_ID=testAccountID",
		"-e", "AWS_REGION=testRegion",
		"-e", "AWS_ROLE=testRole",
		"efritz/ij-ecr-token:latest",
	}))

	Expect(runner.RunFuncCallCount()).To(Equal(1))
	Expect(runner.RunFuncCallParams()[0].Arg1).To(Equal([]string{
		"docker",
		"login",
		"-u",
		"AWS",
		"--password-stdin",
		"https://testAccountID.dkr.ecr.testRegion.amazonaws.com",
	}))

	stdin := runner.RunFuncCallParams()[0].Arg2
	content, _ := ioutil.ReadAll(stdin)
	Expect(string(content)).To(Equal("somehugetoken"))
}
