MD_FILES := $(shell find ./ -type f -name '*.md')
GO_FILES := $(shell find ./ -type f -name '*.go')
TK := npx bazelisk run //:theme --
LINT := yarn theme-lint
BZL := yarn bazelisk --
BZL_BIN := $(shell npx bazelisk info bazel-bin)

.PHONY: deploy/theme,deploy/products,watch,download

deploy/theme:
	$(TK) deploy --dir theme

deploy/products: $(MD_FILES) bin/syncdata
	./bin/syncdata --name deploy

watch:
	$(TK) watch --dir theme

download:
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

generated/client.go: gqlgenc syncdata/*.gql
	$(BZL_BIN)/external/com_github_yamashou_gqlgenc/gqlgenc_/gqlgenc

bin/syncdata: $(GO_FILES)
	mkdir -p bin
	go build -o bin k9bookshelf/syncdata

WORKSPACE: go.mod
	$(BZL) run //:gazelle -- update-repos -from_file=go.mod

gqlgenc: WORKSPACE
	$(BZL) build @com_github_99designs_gqlgen//:...
	$(BZL) build @com_github_yamashou_gqlgenc//:gqlgenc

