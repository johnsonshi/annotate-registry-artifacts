.PHONY: build-cli
build-cli:
	go build -v -o ./bin/annotation ./cmd/cli
	chmod +x ./bin/annotation
