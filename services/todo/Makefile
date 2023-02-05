build:
	mkdir -p bin
	go build -o bin/todo cmd/todo.go

test:
	mkdir -p bin
	go test -coverprofile bin/cover.out ./...

cover: test
	mkdir -p bin
	go tool cover -html=bin/cover.out -o bin/cover.html

coverpkg:
	mkdir -p bin
	go test -v -coverpkg=./... -coverprofile=bin/profile.cov ./...
	go tool cover -html=bin/profile.cov -o bin/pkg-cover.html
clean:
	[ -d bin ] && rm -rf bin

all: clean build cover