FROM golang:1.17-alpine AS builder

WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go build -o aiven-audit ./cmd/aiven-audit/

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/aiven-audit /app/aiven-audit

CMD ["/app/aiven-audit"]
