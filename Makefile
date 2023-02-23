# the target used when no target is specified
.PHONY: default
default: clean protocol build

# the target used to clean up the directory
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin internal/protocol

# the target used to build the protocol buffer
.PHONY: protocol
protocol:
	@echo "Building the protocol..."
	mkdir -p internal/protocol
	protoc --go_out=paths=source_relative:internal/protocol protocol.proto

# the target used to build the program
.PHONY: build
build:
	@echo "Building main program..."
	mkdir -p bin
	go mod tidy
	go build -o bin/gpbs ./
