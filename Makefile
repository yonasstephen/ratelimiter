
generate: # generate all mocks based on go:generate definitions
	go generate ./...

test:
	go test -v ./...

test.coverage:
	go test -coverprofile=coverage.out ./...

test.coveragehtml: test.coverage
	go tool cover -html=coverage.out