# Run integration tests in handler with full internal coverage
test-integration:
	go test ./internal/handler \
		-coverpkg=./internal/handler,./internal/service \
		-coverprofile=/tmp/gotest.o \
	&& go tool cover -html=/tmp/gotest.o
