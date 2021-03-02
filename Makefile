build: *.go
	go build -o bin/decisive-engine -v

.PHONY: test

test:
	go test ./...
