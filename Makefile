TK := npx bazelisk run //:theme --
LINT := yarn theme-lint

.PHONY: deploy,watch,download

deploy:
	$(TK) deploy --dir theme

watch:
	$(TK) watch --dir theme

download:
	$(TK) --dir theme download

version:
	$(TK) version

lint:
	$(LINT) ./theme

list:
	$(TK) get --list
