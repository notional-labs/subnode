all: install

LD_FLAGS = -w -s


BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build:
	@echo "Building subnode"
	@go build -mod readonly $(BUILD_FLAGS) -o build/subnode main.go

test:
	@echo "Testing subnode"
	@go test -mod readonly --timeout=10m $(BUILD_FLAGS) ./...

lint:
	@echo "Running golangci-lint"
	golangci-lint run --timeout=10m

install:
	@echo "Installing subnode"
	@go install -mod readonly $(BUILD_FLAGS) ./...

clean:
	rm -rf build

.PHONY: all lint test race msan tools clean build
