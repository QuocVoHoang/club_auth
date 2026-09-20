# Stage 1: Build the Go application
FROM golang:alpine AS builder

# Install required dependencies
RUN apk update && apk add --no-cache make git

# Install golang-migrate@v4.18.3
RUN set -ex && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

# BuildKit/buildx will provide TARGETOS and TARGETARCH at build time for cross-platform builds.
# Use build args and set GOOS/GOARCH from them so the Go binary is built for the target platform.
ARG TARGETOS
ARG TARGETARCH

# Set up Go environment
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app with maximum optimizations
RUN  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
		-ldflags '-w -s' \
		-tags=jsoniter \
		-o product-service cmd/main.go

#########################################
# Build a small image
# Stage 2: Set up Go
FROM alpine:latest

# Install only essential dependencies
RUN apk update && apk add --no-cache ca-certificates tzdata make

# Copy files from the builder stage
COPY --from=builder /app/product-service /
COPY --from=builder /app/database /database
COPY --from=builder /app/Makefile /
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Make sure migrate is executable
RUN chmod +x /usr/local/bin/migrate

# Command to start services
CMD ["/product-service"]
