# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.19-buster
ARG FINAL_IMAGE=debian:buster-slim

# Build stage
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS builder

# Build Args
ARG GIT_REVISION=""

# Ensure ca-certificates are up to date on the image
RUN update-ca-certificates

# Use modules for dependencies
WORKDIR $GOPATH/src/github.com/rotationalio/ensign

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=1
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy package
COPY . .

# Build binary
ARG TARGETOS TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /go/bin/uptime -ldflags="-X 'github.com/rotationalio/ensign/pkg.GitVersion=${GIT_REVISION}'" ./cmd/uptime

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="Uptime monitor that provides Rotational service status updates"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage
COPY --from=builder /go/bin/uptime /usr/local/bin/uptime

CMD [ "/usr/local/bin/uptime", "serve" ]