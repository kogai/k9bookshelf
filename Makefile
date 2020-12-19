MD_FILES := $(shell find ./ -type f -name '*.md')
GO_FILES := $(shell find ./ -type f -name '*.go' | grep -v "generated")
TK := npx bazelisk run //:theme --
LINT := yarn theme-lint
WP := yarn webpack
BZL := yarn bazelisk --
BZL_BIN := $(shell npx bazelisk info bazel-bin)
VERSION := $(shell cat content/.version | tr -d '\n')
ONIX_FILE := ""

.PHONY: deploy/theme,deploy/contents,watch,download/theme

deploy/theme:
	$(WP)
	$(TK) deploy --dir theme

deploy/contents: $(MD_FILES)
	$(BZL) run //content/cmd/content -- deploy \
		--dir $(PWD)/contents \
		--config $(PWD)/config.yml

download/contents: $(MD_FILES)
	$(BZL) run //content/cmd/content -- download \
		--dir $(PWD)/contents \
		--config $(PWD)/config.yml

import: $(GO_FILES)
	$(BZL) run //onix/cmd -- --input $(PWD)/onix/$(ONIX_FILE)

watch:
	$(TK) watch --dir theme

download/theme:
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

gqlgenc/client/client.go: gqlgenc/main.go gqlgenc/*.gql
	$(BZL) run gqlgenc
	cp -r $(BZL_BIN)/gqlgenc/gqlgenc_/gqlgenc.runfiles/k9books/gqlgenc/client $(CURDIR)/gqlgenc

k9bookshelf/content: $(GO_FILES) WORKSPACE
	mkdir -p k9bookshelf
	for target in darwin_amd64 linux_amd64 ; do \
		$(BZL) build --platforms=@io_bazel_rules_go//go/toolchain:$$target //content/cmd/content ; \
		cp -f $(BZL_BIN)/content/cmd/content/content_/content k9bookshelf/$(VERSION)-content.$$target ; \
	done

.PHONY: release
release:
	git cm -a -m "bump version"
	git tag -af "$(VERSION)" -m ""

.PHONY: setup
setup: WORKSPACE */BUILD.bazel GOPATH

.PHONY: content/BUILD.bazel gqlgenc/BUILD.bazel

*/BUILD.bazel: $(GO_FILES)
	$(BZL) run //:gazelle

WORKSPACE: go.mod content/BUILD.bazel
	$(BZL) run //:gazelle -- update-repos -from_file=go.mod

.PHONY: GOPATH

GOPATH:
	$(BZL) build //:gopath
