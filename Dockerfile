# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-alpine AS build

RUN apk add curl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x ./kubectl

WORKDIR /app

COPY go.mod ./
# RUN go mod download
COPY *.go ./

RUN go build -o /K8sLogChecker

## Deploy
FROM alpine

ENV PATH=/
# ENV KUBECONFIG=/root/.kube/config
WORKDIR /

COPY --from=build /go/kubectl /kubectl
COPY --from=build /K8sLogChecker /K8sLogChecker

ENTRYPOINT ["/K8sLogChecker"]