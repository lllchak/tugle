.PHONY: all clean

all: clean build

build:
	go build -o bin/altum cmd/main.go

run: all
	@echo "\n--------"
	@./bin/altum

test:
	@go test -v ./...

clean:
	rm -rf ./bin