.PHONY: all clean

all: clean build

build:
	go build -o bin/tugle cmd/main.go

run: all
	@echo "\n--------"
	@./bin/tugle

test:
	@go test -race -cover -coverprofile=coverage.out ./pkg/tests

clean:
	rm -rf ./bin
