test:
	go test ./...

format:
	gofmt -w -s .

build:
	go build -o ./bin/aiven-audit ./cmd/aiven-audit/
