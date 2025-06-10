# Makefile
lint:
    golangci-lint run ./...

test:
    go test -v ./...

run:
    go run mb.go