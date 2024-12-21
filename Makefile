# Go parameters
BINARY_NAME=logd
MAIN_PACKAGE=./cmd/logd

# Build flags
GOOS?=linux
GOARCH?=amd64
BUILD_FLAGS=-ldflags="-s -w"

# Test flags
TEST_FLAGS=-v -race
INTEGRATION_TEST_FLAGS=-v -race -tags=integration
E2E_TEST_FLAGS=-v -race -tags=e2e

# Kubernetes
K3D_CLUSTER_NAME=ergo-logs-cluster
K3D_CONFIG=deploy/k8s/k3d-config.yaml

.PHONY: all build clean test test-integration test-e2e k3d-create k3d-delete

all: test build

# Build the application
build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run unit tests
test:
	go test $(TEST_FLAGS) ./...

# Run integration tests
test-integration:
	go test $(INTEGRATION_TEST_FLAGS) ./test/integration/...

# Run end-to-end tests
test-e2e:
	go test $(E2E_TEST_FLAGS) ./test/e2e/...

# Create K3D cluster
k3d-create:
	k3d cluster create $(K3D_CLUSTER_NAME) --config $(K3D_CONFIG)

# Delete K3D cluster
k3d-delete:
	k3d cluster delete $(K3D_CLUSTER_NAME)

# Watch tests during development (requires watchexec)
watch-test:
	watchexec -c -e go -- go test $(TEST_FLAGS) ./...

# Run all tests
test-all: test test-integration test-e2e

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...
