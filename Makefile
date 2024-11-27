run-tests:
	go clean -cache
	go test -v ./...
coverage:
	go test -covermode=count -coverpkg=./... -coverprofile coverage.out -v ./...
	go tool cover -html coverage.out
app-run-docker:
	docker-compose up
mock:
	go generate ./...
lint:
	golangci-lint run