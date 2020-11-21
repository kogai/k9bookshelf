TK := npx bazelisk run //:theme --
LINT := yarn theme-lint

.PHONY: deploy,watch,download

deploy:
	$(TK) deploy --dir theme

watch:
	$(TK) watch --dir theme

download:
	$(TK) --dir theme download

lint:
	$(LINT) ./theme

generated/client.go: syncdata/*.gql
	gqlgenc
	# $(BZL_BIN)/external/com_github_yamashou_gqlgenc
