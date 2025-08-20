.PHONY: test

all: test 
	go build -race .

fmt:
	go fmt ./...

test:
	go test -race -cover ./...	
