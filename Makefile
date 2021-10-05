test:
	go test ./...

format:
	gofmt -w .

build:
	go build -o ./bin/aiven-audit ./main.go