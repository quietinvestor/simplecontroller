# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o simplecontroller .

# Runtime stage
FROM scratch

COPY --from=builder /build/simplecontroller /simplecontroller

USER 65532:65532 # nobody:nogroup

ENTRYPOINT ["/simplecontroller"]
