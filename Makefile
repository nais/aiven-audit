test:
	go test ./...

format:
	gofmt -w -s .

lint:
	golangci-lint run

build:
	go build -o ./bin/aiven-audit ./cmd/aiven-audit/
