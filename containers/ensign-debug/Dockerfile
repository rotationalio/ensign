# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.19-buster
ARG FINAL_IMAGE=debian:buster-slim

# Build stage
FROM ${BUILDER_IMAGE} AS builder

# Build Args
ARG GIT_REVISION=""

# Ensure ca-certificates are up to date on the image
RUN update-ca-certificates

# Use modules for dependencies
WORKDIR $GOPATH/src/github.com/rotationalio/ensign

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy package
COPY . .

# Build binary
RUN go build -v -o /go/bin/ensign-debug -ldflags="-X 'github.com/rotationalio/ensign/pkg.GitVersion=${GIT_REVISION}'" ./cmd/debug

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="Generates or consumes Ensign events as a job for debugging purposes."

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage
COPY --from=builder /go/bin/ensign-debug /usr/local/bin/ensign-debug

ENTRYPOINT [ "/usr/local/bin/ensign-debug" ]
CMD [ "consume" ]