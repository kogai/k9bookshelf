MD_FILES := $(shell find ./ -type f -name '*.md')
GO_FILES := $(shell find ./ -type f -name '*.go')
TK := npx bazelisk run //:theme --
LINT := yarn theme-lint

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

generated/client.go: syncdata/*.gql
	gqlgenc
	# $(BZL_BIN)/external/com_github_yamashou_gqlgenc

bin/syncdata: $(GO_FILES)
	mkdir -p bin
	go build -o bin k9bookshelf/syncdata
