MD_FILES := $(shell find ./ -type f -name '*.md')
GO_FILES := $(shell find ./ -type f -name '*.go' | grep -v "generated")
TK := npx bazelisk run //:theme --
LINT := yarn theme-lint
BZL := yarn bazelisk --
BZL_BIN := $(shell npx bazelisk info bazel-bin)

.PHONY: deploy/theme,deploy/products,watch,download/theme

deploy/theme:
	$(TK) deploy --dir theme

deploy/products: $(MD_FILES) bin/syncdata
	./bin/syncdata deploy

download/products: $(MD_FILES) bin/syncdata
	./bin/syncdata download

watch:
	$(TK) watch --dir theme

download/theme:
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

generated/client.go: bin/gqlgenc syncdata/*.gql
	./bin/gqlgenc

bin/%: $(GO_FILES) WORKSPACE
	mkdir -p bin
	$(BZL) build //$(@F):all
	cp -f $(BZL_BIN)/$(@F)/$(@F)_/$(@F) bin/

.PHONY: syncdata/BUILD.bazel gqlgenc/BUILD.bazel

*/BUILD.bazel:
	$(BZL) run //:gazelle

WORKSPACE: go.mod
	$(BZL) run //:gazelle -- update-repos -from_file=go.mod
