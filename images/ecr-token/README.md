# ECR Token Generator

This Docker image runs commands via AWS CLI and outputs an AWS authorization
token capable of logging into an ECR registry.

The following environment variables tell the entrypoint how to authenticate.

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_REGION` (defaults to *us-east-1*)

Additionally, if the following variables are supplied, the entrypoint will try
to a assume a role in a different account.

- `AWS_ACCOUNT_ID`
- `AWS_ROLE`
