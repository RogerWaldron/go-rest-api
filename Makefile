build:
	@go build -o bin/api

run: build
	@./bin.api

test: @go test -v ./server/...

test_coverage: 
@go test ./server/... -coverprofile=coverage.out