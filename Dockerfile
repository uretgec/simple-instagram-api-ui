FROM golang:1.17-alpine AS builder

ARG SERVICE_BUILD
ARG SERVICE_COMMIT_ID

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /usr/local/go/src/poster/

WORKDIR /usr/local/go/src/poster

COPY . .

RUN ls .

RUN go install -v -ldflags="-w -s -X main.ServiceBuild=${SERVICE_BUILD} -X main.ServiceCommitId=${SERVICE_COMMIT_ID}" ./...