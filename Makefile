
generate: # generate all mocks based on go:generate definitions
	go generate ./...

lint: # it's assumed that golangci-lint is already installed (https://golangci-lint.run/usage/install/#local-installation)
	golangci-lint run

test:
	go test -count=1 -race -v ./...

test.coverage:
	go test -coverprofile=coverage.out ./...

test.coveragehtml: test.coverage
	go tool cover -html=coverage.out