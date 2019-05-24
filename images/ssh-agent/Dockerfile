FROM golang:1.12 AS builder
WORKDIR /ij
COPY go.mod go.sum ./
COPY main.go ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ij-ensure-keys-available

FROM alpine:3.8

RUN set -ex && \
    apk -v --update add openssh-client

ENV SSH_AUTH_VOLUME="/tmp/ij/ssh-agent"
ENV SSH_AUTH_SOCK="${SSH_AUTH_VOLUME}/ssh-agent.sock"

WORKDIR /ij
CMD ["/ij/start-agent.sh"]
VOLUME ["${SSH_AUTH_VOLUME}"]
COPY add-keys.sh /ij
COPY start-agent.sh /ij
COPY --from=builder /ij/ij-ensure-keys-available .
