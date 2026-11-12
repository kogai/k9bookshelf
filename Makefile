GO_FILES := $(shell find . -type f -name '*.go' -not -path './node_modules/*')
# NOTE: Generated GraphQL client is excluded from format check.
GOFMT_FILES := $(filter-out ./gqlgenc/client/%,$(GO_FILES))
LINT := yarn theme-lint
WP := yarn webpack
VERSION := $(shell cat content/.version | tr -d '\n')
ONIX_FILE := ""
DIST := k9bookshelf

THEMEKIT_VERSION := v1.0.2
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
THEMEKIT_PLATFORM := darwin-amd64
THEMEKIT_SHA256 := 4b734d08db62e001adcf8ac5607414aa9cb6cfb410faa8f53b02f9f5c8a5733c
else
THEMEKIT_PLATFORM := linux-amd64
THEMEKIT_SHA256 := 21a058bc47555ec343088c9322187bb0ad57e4055454f509c0751c092af6b12b
endif
TK := ./bin/theme

CONTENT := go run ./content/cmd/content

.PHONY: deploy/theme deploy/contents download/contents download/theme watch lint test import import/dry release gqlgenc dist

bin/theme:
	mkdir -p bin
	curl -sSfL -o $@.tmp https://shopify-themekit.s3.amazonaws.com/$(THEMEKIT_VERSION)/$(THEMEKIT_PLATFORM)/theme
	echo "$(THEMEKIT_SHA256)  $@.tmp" | shasum -a 256 -c -
	chmod +x $@.tmp
	mv $@.tmp $@

deploy/theme: bin/theme
	$(WP)
	$(TK) deploy --dir theme

watch: bin/theme
	$(TK) watch --dir theme

download/theme: bin/theme
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

test:
	test -z "$$(gofmt -l $(GOFMT_FILES))" || (gofmt -l $(GOFMT_FILES); exit 1)
	go vet ./...
	go test ./...

deploy/contents:
	$(CONTENT) deploy \
		--dir $(CURDIR)/contents \
		--config $(CURDIR)/config.yml

download/contents:
	$(CONTENT) download \
		--dir $(CURDIR)/contents \
		--config $(CURDIR)/config.yml

import:
	go run ./onix/cmd --input $(CURDIR)/onix/$(ONIX_FILE)

# dry-run writes snapshot json to onix/snapshots
import/dry:
	go run ./onix/cmd --dryRun --input $(CURDIR)/onix/$(ONIX_FILE)
	mv $(CURDIR)/onix/*.json $(CURDIR)/onix/snapshots/

gqlgenc:
	go run ./gqlgenc

# Builds content-kit binaries for release into $(DIST)
dist: $(GO_FILES)
	mkdir -p $(DIST)
	for target in darwin_amd64 darwin_arm64 linux_amd64 ; do \
		CGO_ENABLED=0 GOOS=$${target%_*} GOARCH=$${target#*_} \
			go build -o $(DIST)/$(VERSION)-content.$$target ./content/cmd/content ; \
	done

release:
	git cm -a -m "bump version"
	git tag -af "$(VERSION)" -m ""
