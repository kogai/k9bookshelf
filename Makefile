MD_FILES := $(shell find ./ -type f -name '*.md')
TK := npx bazelisk run //:theme --
LINT := yarn theme-lint

.PHONY: deploy/theme,deploy/products,watch,download

deploy/theme:
	$(TK) deploy --dir theme

deploy/products: $(MD_FILES)
	go run syncdata/main.go --name deploy

watch:
	$(TK) watch --dir theme

download:
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

generated/client.go: syncdata/*.gql
	gqlgenc
	# $(BZL_BIN)/external/com_github_yamashou_gqlgenc
