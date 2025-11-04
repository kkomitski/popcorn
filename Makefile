# Default Go variables
GO=go

# Run all tests
.PHONY: test
test:
	$(GO) test ./...

# Build target depends on test
.PHONY: build
build: test
	$(GO) build -o popcorn main.go
