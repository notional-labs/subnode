all: install

LD_FLAGS = -w -s


BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build:
	@echo "Building subnode"
	@go build -mod readonly $(BUILD_FLAGS) -o build/subnode main.go

test:
	@echo "Testing subnode"
	@go test -mod readonly $(BUILD_FLAGS) ./...

lint:
	@echo "Run golangci-lint"
	golangci-lint run

install:
	@echo "Installing subnode"
	@go install -mod readonly $(BUILD_FLAGS) ./...

clean:
	rm -rf build

.PHONY: all lint test race msan tools clean build
