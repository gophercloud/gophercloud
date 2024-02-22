lint:
	go fmt ./...
	go vet ./...
.PHONY: lint

unit:
	go test ./openstack/...
.PHONY: unit

acceptance: acceptance-baremetal acceptance-blockstorage acceptance-clustering acceptance-compute acceptance-container acceptance-containerinfra acceptance-db acceptance-dns acceptance-identity acceptance-imageservice acceptance-keymanager acceptance-loadbalancer acceptance-messaging acceptance-networking acceptance-objectstorage acceptance-orchestration acceptance-placement acceptance-sharedfilesystems acceptance-workflow
.PHONY: acceptance

acceptance-baremetal:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/baremetal/httpbasic
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/baremetal/noauth
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/baremetal/v2
.PHONY: acceptance-baremetal

acceptance-blockstorage:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/blockstorage/extensions
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/blockstorage/noauth
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/blockstorage/v1
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/blockstorage/v2
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/blockstorage/v3
.PHONY: acceptance-blockstorage

acceptance-clustering:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/clustering/v1
.PHONY: acceptance-clustering

acceptance-compute:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/compute/v2
.PHONY: acceptance-compute

acceptance-container:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/container/v1
.PHONY: acceptance-container

acceptance-containerinfra:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/containerinfra/v1
.PHONY: acceptance-containerinfra

acceptance-db:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/db/v1
.PHONY: acceptance-db

acceptance-dns:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/dns/v2
.PHONY: acceptance-dns

acceptance-identity:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/identity/v2
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/identity/v3
.PHONY: acceptance-identity

acceptance-image:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/imageservice/v2
.PHONY: acceptance-image

acceptance-keymanager:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/keymanager/v1
.PHONY: acceptance-keymanager

acceptance-loadbalancer:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/loadbalancer/v2
.PHONY: acceptance-loadbalancer

acceptance-messaging:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/messaging/v2
.PHONY: acceptance-messaging

acceptance-networking:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/networking/v2
.PHONY: acceptance-networking

acceptance-objectstorage:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/objectstorage/v1
.PHONY: acceptance-objectstorage

acceptance-orchestration:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/orchestration/v1
.PHONY: acceptance-orchestration

acceptance-placement:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/placement/v1
.PHONY: acceptance-placement

acceptance-sharedfilesystems:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/sharedfilesystems/v2
.PHONY: acceptance-sharefilesystems

acceptance-workflow:
	go test -v -tags "fixtures acceptance" ./internal/acceptance/openstack/workflow/v2
.PHONY: acceptance-workflow
