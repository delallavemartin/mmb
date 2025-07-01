# Makefile
lint:
	cd mb && golangci-lint run ./...

test:
	cd mb && go test -v ./...

run:
	go run ./mb/mb.go