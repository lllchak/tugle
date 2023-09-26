.PHONY: all clean

all: clean build

build:
	go build -o bin/tugle cmd/main.go

run: all
	@echo "\n--------"
	@./bin/altum

test:
	@go test -race -cover -coverprofile=coverage.out ./...

clean:
	rm -rf ./bin
