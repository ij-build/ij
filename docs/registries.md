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
| password-file |          |         | The path to a file on the host containing the password used for login. |
| server        | yes      |         | The hostname of the registry. |
| username      | yes      |         | The username used for login. |

One of `password` and `password-file` variables must be supplied.

## AWS / ECR

An AWS Elastic Container Registry which can be logged in via AWS account credentials and an optional IAM role.

| Name              | Required | Default   | Description |
| ----------------- | -------- | --------- | ----------- |
| access-key-id     | yes      |           | The user's AWS credentials. |
| account-id        |          |           | The identifier of the account owning the registry. |
| region            |          | us-east-1 | The region where the registry is available. |
| role              |          |           | The target assumed role of the provided account. |
| secret-access-key | yes      |           | The user's AWS credentials. |

If `role` and `account-id` are not supplied, then the registry is assumed to belong to the same account as the authenticated user. For cross-account use, the `role` and `account-id` variables can be supplied to force a role to be assumed on a secondary account. For more details, see the ecr-token [readme](https://github.com/efritz/ij/blob/master/images/ecr-token/README.md#user-content-ecr-token).

## Google Cloud / GCR

A Google Container Registry which can be logged in via the [JSON Key File](https://cloud.google.com/container-registry/docs/advanced-authentication#json_key_file) authentication mechanism.

| Name     | Required | Default | Description |
| -------- | -------- | ------- | ----------- |
| hostname |          | gcr.io  | The GCR hostname. May also be one of `us.gcr.io`, `eu.gcr.io`, or `asia.gcr.io`. |
| key      |          |         | A [service account JSON key](https://support.google.com/cloud/answer/6158849#serviceaccounts). |
| key-file |          |         | The path to a JSON key file on the host. |

One of `key` or `key-file` must be supplied.

# Example

This example illustrates a registry list with all three registry types. It is suggested to store server registry credentials (when it is not possible to use password-file) as well as AWS credentials in the user's global override file.

```yaml
registries:
  - server: registry.example.io
    username: admin
    password: secret

  - type: ecr
    access-key-id: ${AWS_ACCESS_KEY_ID}
    secret-access-key: ${AWS_SECRET_ACCESS_KEY}
    account-id: 641844361036
    role: Developer

  - type: gcr
    key-file: /etc/docker-agent-keyfile.json
```
