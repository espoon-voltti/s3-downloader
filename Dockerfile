# SPDX-FileCopyrightText: 2018-2020 City of Espoo
#
# SPDX-License-Identifier: MIT

FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum.license ./

# Download dependencies
RUN go mod download

# Copy source code
COPY ./ ./

# Test stage
FROM builder AS test
RUN go test -v ./...

# Build stage
FROM builder AS build
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o s3-downloader .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from build stage
COPY --from=build /app/s3-downloader .

ENTRYPOINT ["./s3-downloader"]
