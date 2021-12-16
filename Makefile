PACKAGE=github.com/vine-io/apimachinery

all: release

vendor:
	go mod vendor

proto:
	cd $(GOPATH)/src && \
	protoc -I=$(GOPATH)/src --gogo_out=:. --validator_out=:. --deepcopy_out=:. $(PACKAGE)/apis/meta/v1/meta.proto

release:
ifeq "$(TAG)" ""
	@echo "missing tag"
	exit 1
endif
	git tag $(TAG)
	git add .
	git commit -m "$(TAG)"
	git tag -d $(TAG)
	git tag $(TAG)

.PHONY: build-tag release
