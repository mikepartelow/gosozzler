.PHONY: test

all: test 
	go build .

fmt:
	go fmt ./...

test:
	go test -race -cover ./...	
