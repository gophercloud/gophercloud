undefine GOFLAGS

TIMEOUT := "60m"
GOLANGCI_LINT_VERSION?=v1.57.1
ifeq ($(shell command -v podman 2> /dev/null),)
	RUNNER=docker
else
	RUNNER=podman
endif

# if the golangci-lint steps fails with one of the following error messages:
#
#   directory prefix . does not contain main module or its selected dependencies
#
#   failed to initialize build cache at /root/.cache/golangci-lint: mkdir /root/.cache/golangci-lint: permission denied
#
# you probably have to fix the SELinux security context for root directory plus your cache
#
#   chcon -Rt svirt_sandbox_file_t .
#   chcon -Rt svirt_sandbox_file_t ~/.cache/golangci-lint
lint:
	mkdir -p ~/.cache/golangci-lint/$(GOLANGCI_LINT_VERSION)
	$(RUNNER) run -t --rm \
		-v $(shell pwd):/app \
		-v ~/.cache/golangci-lint/$(GOLANGCI_LINT_VERSION):/root/.cache \
		-w /app \
		-e GOFLAGS="-tags=acceptance" \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run
.PHONY: lint

unit:
	go test ./...
.PHONY: unit

coverage:
	go test -covermode count -coverprofile cover.out -coverpkg=./... ./...
.PHONY: coverage

acceptance: acceptance-basic acceptance-baremetal acceptance-blockstorage acceptance-compute acceptance-container acceptance-containerinfra acceptance-db acceptance-dns acceptance-identity acceptance-image acceptance-keymanager acceptance-loadbalancer acceptance-messaging acceptance-networking acceptance-objectstorage acceptance-orchestration acceptance-placement acceptance-sharedfilesystems acceptance-workflow
.PHONY: acceptance

acceptance-basic:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack
.PHONY: acceptance-basic

acceptance-baremetal:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/baremetal/...
.PHONY: acceptance-baremetal

acceptance-blockstorage:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/blockstorage/...
.PHONY: acceptance-blockstorage

acceptance-compute:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/compute/...
.PHONY: acceptance-compute

acceptance-container:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/container/...
.PHONY: acceptance-container

acceptance-containerinfra:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/containerinfra/...
.PHONY: acceptance-containerinfra

acceptance-db:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/db/...
.PHONY: acceptance-db

acceptance-dns:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/dns/...
.PHONY: acceptance-dns

acceptance-identity:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/identity/...
.PHONY: acceptance-identity

acceptance-image:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/image/...
.PHONY: acceptance-image

acceptance-keymanager:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/keymanager/...
.PHONY: acceptance-keymanager

acceptance-loadbalancer:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/loadbalancer/...
.PHONY: acceptance-loadbalancer

acceptance-messaging:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/messaging/...
.PHONY: acceptance-messaging

acceptance-networking:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/networking/...
.PHONY: acceptance-networking

acceptance-objectstorage:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/objectstorage/...
.PHONY: acceptance-objectstorage

acceptance-orchestration:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/orchestration/...
.PHONY: acceptance-orchestration

acceptance-placement:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/placement/...
.PHONY: acceptance-placement

acceptance-sharedfilesystems:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/sharedfilesystems/...
.PHONY: acceptance-sharefilesystems

acceptance-workflow:
	go test -timeout $(TIMEOUT) -tags "fixtures acceptance" ./internal/acceptance/openstack/workflow/...
.PHONY: acceptance-workflow
