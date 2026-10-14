HOSTNAME    = registry.terraform.io
NAMESPACE   = lannuttia
NAME        = quay
BINARY      = terraform-provider-$(NAME)
VERSION     = 0.1.0

OS_ARCH     = $(shell go env GOOS)_$(shell go env GOARCH)

default: build

.PHONY: build
build:
	go build -o $(BINARY)

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)
	mv $(BINARY) ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./... -v -count=1 -timeout 120s

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./internal/provider/... -v -count=1 -timeout 600s $(TESTARGS)

.PHONY: clean
clean:
	rm -f $(BINARY)
