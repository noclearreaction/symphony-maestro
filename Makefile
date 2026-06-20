.PHONY: build test clean

build:
	go build -o bin/rubato ./cmd/rubato

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

test-coverage:
	go test -cover ./...
