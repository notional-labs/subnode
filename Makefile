all: install

LD_FLAGS = -w -s


BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build:
	@echo "Building subnode"
	@go build -mod readonly $(BUILD_FLAGS) -o build/subnode main.go

test:
	@echo "Testing basic"
	@go test -mod readonly --timeout=10m $(BUILD_FLAGS) `go list ./... |grep -v github.com/notional-labs/subnode/test`

test-osmosis:
	@echo "Testing subnode with default osmosis config"
	@go test -mod readonly --timeout=10m-ldflags '$(LD_FLAGS) -X github.com/notional-labs/subnode/test.Chain=osmosis' ./test

test-evmos:
	@echo "Testing subnode with evmos config"
	@go test -mod readonly --timeout=10m -ldflags '$(LD_FLAGS) -X github.com/notional-labs/subnode/test.Chain=evmos' ./test

lint:
	@echo "Running golangci-lint"
	golangci-lint run --timeout=10m

install:
	@echo "Installing subnode"
	@go install -mod readonly $(BUILD_FLAGS) ./...

clean:
	rm -rf build

.PHONY: all lint test race msan tools clean build
