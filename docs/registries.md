# Registries

A registry is an object containing authentication metadata for logging into a remote Docker registry. These objects use the `type` property (shown below) to determine the other available properties. Each of the following registry types are discussed in the sections below.

| Name | Required | Default | Description |
| ---- | -------- | ------- | ----------- |
| type |          | server  | The type of registry. May also be one of `ecr` or `gcr`. |

## Server

A Docker container registry which can be logged in via username and password.

| Name          | Required | Default | Description |
| ------------- | -------- | ------- | ----------- |
| password      |          |         | The password used for login. |
| password_file |          |         | The path to a file on the host containing the password used for login. |
| server        | yes      |         | The hostname of the registry. |
| username      | yes      |         | The username used for login. |

One of `password` and `password_file` variables must be supplied.

## AWS / ECR

An AWS Elastic Container Registry which can be logged in via AWS account credentials and an optional IAM role.

| Name              | Required | Default   | Description |
| ----------------- | -------- | --------- | ----------- |
| access_key_id     | yes      |           | The user's AWS credentials. |
| account_id        |          |           | The identifier of the account owning the registry. |
| region            |          | us-east-1 | The region where the registry is available. |
| role              |          |           | The target assumed role of the provided account. |
| secret_access_key | yes      |           | The user's AWS credentials. |

If `role` and `account_id` are not supplied, then the registry is assumed to belong to the same account as the authenticated user. For cross-account use, the `role` and `account_id` variables can be supplied to force a role to be assumed on a secondary account. For more details, see the ecr-token [readme](https://github.com/efritz/ij/blob/master/images/ecr-token/README.md#user-content-ecr-token).

## Google Cloud / GCR

A Google Container Registry which can be logged in via the [JSON Key File](https://cloud.google.com/container-registry/docs/advanced-authentication#json_key_file) authentication mechanism.

| Name     | Required | Default | Description |
| -------- | -------- | ------- | ----------- |
| hostname |          | gcr.io  | The GCR hostname. May also be one of `us.gcr.io`, `eu.gcr.io`, or `asia.gcr.io`. |
| key_file | yes      |         | The path to a [service account JSON key file](https://support.google.com/cloud/answer/6158849#serviceaccounts) on the host. |

# Example

This example illustrates a registry list with all three registry types. It is suggested to store server registry credentials (when it is not possible to use password_file) as well as AWS credentials in the user's global override file.

```yaml
registries:
  - server: registry.example.io
    username: admin
    password: secret

  - type: ecr
    access_key_id: ${AWS_ACCESS_KEY_ID}
    secret_access_key: ${AWS_SECRET_ACCESS_KEY}
    account_id: 641844361036
    role: Developer

  - type: gcr
    key_file: /etc/docker-agent-keyfile.json
```
