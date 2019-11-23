# Start from the latest golang base image
FROM golang:latest as builder

# install protobuf from source
RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool
RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install gRPC stuff
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Create the gRPC sources
RUN protoc -I pkg/messaging/process --go_out=plugins=grpc:pkg/messaging/process pkg/messaging/process/*.proto

# Build the Go apps
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/orchestrator \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor-add \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor-sub \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor-print \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/producer

# final stage
FROM alpine:latest
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*
RUN adduser -S -D -H -h /app appuser
USER appuser
WORKDIR /app
COPY docker/start.sh /app/
COPY --from=builder /app/orchestrator /app/processor-add /app/processor-sub /app/processor-print /app/producer /app/

ENTRYPOINT ["/app/start.sh"]