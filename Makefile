PACKAGE=github.com/vine-io/apimachinery

all: release

vendor:
	go mod vendor

proto:
	goproto-gen -p github.com/vine-io/apimachinery/apis/meta/v1
	deepcopy-gen -i github.com/vine-io/apimachinery/apis/meta/v1

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
