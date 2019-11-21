# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go apps
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/gorchestrator \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor_add \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor_sub \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor_error \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/processor_print \
&& CGO_ENABLED=0 GOOS=linux go build ./cmd/producer

# final stage
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
WORKDIR /app
COPY docker/start.sh /app/
COPY --from=builder /app/gorchestrator /app/processor_add /app/processor_sub /app/processor_error /app/processor_print /app/producer /app/
ENTRYPOINT ["/app/start.sh"]