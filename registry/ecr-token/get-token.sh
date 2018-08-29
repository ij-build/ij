#!/bin/sh

cat << EOF > ~/.aws/credentials
[default]
aws_access_key_id=${AWS_ACCESS_KEY_ID}
aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}
EOF

cat << EOF > ~/.aws/config
[profile default]
output=text
source_profile=default
region=${AWS_REGION:-us-east-1}
EOF

if [ x"${AWS_ROLE}" != x"" ]; then
cat << EOF >> ~/.aws/config
role_arn=arn:aws:iam::${AWS_ACCOUNT_ID}:role/${AWS_ROLE}
EOF
fi

aws ecr get-authorization-token --profile default --query 'authorizationData[].authorizationToken' | base64 -d | cut -d: -f2
