MD_FILES := $(shell find ./ -type f -name '*.md')
GO_FILES := $(shell find ./ -type f -name '*.go' | grep -v "generated")
TK := npx bazelisk run //:theme --
LINT := yarn theme-lint
BZL := yarn bazelisk --
BZL_BIN := $(shell npx bazelisk info bazel-bin)

.PHONY: deploy/theme,deploy/contents,watch,download/theme

deploy/theme:
	$(TK) deploy --dir theme

deploy/contents: $(MD_FILES) bin/syncdata
	./bin/syncdata deploy --input $(PWD)/contents

download/contents: $(MD_FILES) bin/syncdata
	./bin/syncdata download --output $(PWD)/contents

watch:
	$(TK) watch --dir theme

download/theme:
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

syncdata/generated/client.go: gqlgenc/main.go syncdata/*.gql
	$(BZL) run gqlgenc
	cp -r $(BZL_BIN)/gqlgenc/gqlgenc_/gqlgenc.runfiles/k9books/syncdata/generated $(CURDIR)/syncdata

bin/*: $(GO_FILES) WORKSPACE
	mkdir -p bin
	$(BZL) build //$(@F):all
	cp -f $(BZL_BIN)/$(@F)/$(@F)_/$(@F) bin/

bin/syncdata: $(GO_FILES) WORKSPACE
	mkdir -p bin
	$(BZL) build //syncdata/cmd:all
	cp -f $(BZL_BIN)/syncdata/cmd/cmd_/cmd bin/syncdata

.PHONY: setup
setup: WORKSPACE */BUILD.bazel GOPATH

.PHONY: syncdata/BUILD.bazel gqlgenc/BUILD.bazel

*/BUILD.bazel: $(GO_FILES)
	$(BZL) run //:gazelle

WORKSPACE: go.mod syncdata/BUILD.bazel
	$(BZL) run //:gazelle -- update-repos -from_file=go.mod

.PHONY: GOPATH

GOPATH:
	$(BZL) build //:gopath
