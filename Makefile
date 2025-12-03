BINARY_NAME=hyprtask
input_dir=cmd/hyprtask/main.go

.PHONY: all build install clean run

all: build

# Builds the binary into a 'bin' folder
build:
	@echo "Building..."
	@mkdir -p bin
	@go build -o bin/$(BINARY_NAME) $(input_dir)
	@echo "Built $(BINARY_NAME) in ./bin"

# Installs the binary to your $GOPATH/bin (usually ~/go/bin)
# This allows you to run 'hyprtask' from anywhere
install:
	@echo "Installing..."
	@go install ./cmd/hyprtask
	@echo "Done! You can now run '$(BINARY_NAME)' from anywhere."

# Runs the project directly (good for dev)
run:
	@go run $(input_dir)

# Cleans up build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin
	@echo "Cleaned."