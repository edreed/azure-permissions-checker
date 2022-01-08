BUILD_ARCH?=amd64
BUILD_OS?=linux

IMAGE_REGISTRY?=ejreed
IMAGE_NAME?=azure-permissions-checker
IMAGE_VERSION?=v1.0.0-alpha.1
IMAGE_TAG=$(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: azpermissions
azpermissions:
	CGO_ENABLED=0 GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build -o ./bin/$(BUILD_ARCH)/azpermissions ./cmd/azpermissions

.PHONY: azpermissions-container
azpermissions-container:
	docker build --tag "$(IMAGE_TAG)" --file ./cmd/azpermissions/Dockerfile . --build-arg "BUILD_ARCH=$(BUILD_ARCH)" --build-arg "BUILD_OS=$(BUILD_OS)"
	docker push "$(IMAGE_TAG)"

.PHONY: azcheckperms
azcheckperms:
	CGO_ENABLED=0 GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build -o ./bin/$(BUILD_ARCH)/azcheckperms ./cmd/azcheckperms

.PHONY: all
all: azpermissions-container azcheckperms
