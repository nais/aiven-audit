FROM golang:1.22 AS builder

WORKDIR /app
ENV CGO_ENABLED=0

COPY . .

RUN go mod download

RUN go test -v ./...

RUN go build -o aiven-audit ./cmd/aiven-audit/

FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /app

COPY --from=builder /app/aiven-audit /app/aiven-audit

CMD ["/app/aiven-audit"]
