.PHONY: build

XAPPNAME=ransim-api-demo
VERSION=v0.0.1

build:
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/$(XAPPNAME) ./cmd/$(XAPPNAME)

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

docker:
	@go mod vendor
	docker build --network host -f build/Dockerfile -t dhanifudin/$(XAPPNAME):$(VERSION) .
	docker push dhanifudin/$(XAPPNAME):$(VERSION)
	@rm -rf vendor

install-xapp:
	helm install -n riab $(XAPPNAME) ./helm-chart/$(XAPPNAME) --values ./helm-chart/$(XAPPNAME)/values.yaml

delete-xapp:
	-helm uninstall -n riab $(XAPPNAME)

dev: delete-xapp docker install-xapp
