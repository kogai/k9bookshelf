TK := npx bazelisk run //:theme --
LINT := yarn theme-lint

.PHONY: deploy,watch,download

deploy:
	$(TK) deploy

watch:
	$(TK) watch

download:
	$(TK) download

version:
	$(TK) version

lint:
	$(LINT) .

list:
	$(TK) get --list
