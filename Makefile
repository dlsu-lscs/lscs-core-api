build:
	@go build -o bin/lscs-central-auth ./cmd/api/main.go

run: build
	@./bin/lscs-central-auth

test:
	@go test -v ./...
