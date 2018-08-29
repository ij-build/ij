package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

type (
	ecrLogin struct {
		ctx      context.Context
		logger   logging.Logger
		env      environment.Environment
		registry *config.ECRRegistry
		runner   command.Runner
	}

	awsCredentials struct {
		AccessKeyID     string
		SecretAccessKey string
		AccountID       string
		Region          string
		Role            string
	}
)

const ECRServerFormat = "https://%s.dkr.ecr.%s.amazonaws.com"

func NewECRLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registry *config.ECRRegistry,
) Login {
	return newECRLogin(
		ctx,
		logger,
		env,
		registry,
		command.NewRunner(logger),
	)
}

func newECRLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registry *config.ECRRegistry,
	runner command.Runner,
) Login {
	return &ecrLogin{
		ctx:      ctx,
		logger:   logger,
		env:      env,
		registry: registry,
		runner:   runner,
	}
}

func (l *ecrLogin) GetServer() (string, error) {
	credentials, err := getAWSCredentials(l.env, l.registry)
	if err != nil {
		return "", err
	}

	return l.getServer(credentials), nil
}

func (l *ecrLogin) Login() error {
	credentials, err := getAWSCredentials(l.env, l.registry)
	if err != nil {
		return err
	}

	l.logger.Info(
		nil,
		"Generating an ECR access token",
	)

	token, err := getAWSToken(
		l.ctx,
		l.runner,
		credentials,
	)

	if err != nil {
		return err
	}

	return login(
		l.ctx,
		l.runner,
		l.getServer(credentials),
		"AWS",
		token,
	)
}

func (l *ecrLogin) getServer(credentials *awsCredentials) string {
	return fmt.Sprintf(
		ECRServerFormat,
		credentials.AccountID,
		credentials.Region,
	)
}

//
// Helpers

func getAWSCredentials(
	env environment.Environment,
	registry *config.ECRRegistry,
) (*awsCredentials, error) {
	accessKeyID, err := env.ExpandString(registry.AccessKeyID)
	if err != nil {
		return nil, err
	}

	secretAccessKey, err := env.ExpandString(registry.SecretAccessKey)
	if err != nil {
		return nil, err
	}

	accountID, err := env.ExpandString(registry.AccountID)
	if err != nil {
		return nil, err
	}

	region, err := env.ExpandString(registry.Region)
	if err != nil {
		return nil, err
	}

	role, err := env.ExpandString(registry.Role)
	if err != nil {
		return nil, err
	}

	return &awsCredentials{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		AccountID:       accountID,
		Region:          region,
		Role:            role,
	}, nil
}

func getAWSToken(
	ctx context.Context,
	runner command.Runner,
	credentials *awsCredentials,
) (string, error) {
	args := []string{
		"docker",
		"run",
		"--rm",
	}

	args = append(args, credentials.Env()...)
	args = append(args, "ecr-token")

	token, stderr, err := runner.RunForOutput(
		ctx,
		args,
		nil,
	)

	if err != nil || token == "" {
		return "", fmt.Errorf("failed to generate AWS token: %s", strings.TrimSpace(stderr))
	}

	return token, nil
}

func (c *awsCredentials) Env() []string {
	return []string{
		"-e", fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", c.AccessKeyID),
		"-e", fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", c.SecretAccessKey),
		"-e", fmt.Sprintf("AWS_ACCOUNT_ID=%s", c.AccountID),
		"-e", fmt.Sprintf("AWS_REGION=%s", c.Region),
		"-e", fmt.Sprintf("AWS_ROLE=%s", c.Role),
	}
}
