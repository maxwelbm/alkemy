# This command will run the tests for the project
.PHONY tests:
tests:
	go test -v ./...

.PHONY int-tests:
int-tests:
	@go run cmd/main.go & \
	SERVER_PID=$$!; \
	sleep 1;  \
	go test tests/integration_test.go -tags special; \
	if ps -p $$SERVER_PID > /dev/null; then kill $$SERVER_PID; fi

# This command will generate a coverage report for the project
.PHONY: coverage
coverage:
	go test -cover -coverprofile=coverage.out ./...

# This command will display the coverage report in an HTML file
.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html="coverage.out"

# This command will display the coverage report in the terminal
.PHONY: coverage-total
coverage-total: coverage
	go tool cover -func="coverage.out"

# This command will run the linter for the project
.PHONY: linter
linter:
	staticcheck ./...