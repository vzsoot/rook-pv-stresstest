FROM golang:1.22.1 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

RUN GOMAXPROCS=1 GOGC=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o application main.go

FROM rockylinux:9
WORKDIR /
COPY --from=builder /workspace/application .
USER 65532:65532

ENTRYPOINT ["/application"]
