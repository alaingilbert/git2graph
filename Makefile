EXAMPLES_FILE=tools/renderer/examples.js
GITHUB_PAGES_FOLDER=tools/renderer
GITHUB_PAGES_BRANCH=gh-pages

deploy:
	@echo "// This file is autogenerated, do not modify it directly." > $(EXAMPLES_FILE); \
	echo "var examples = {}" >> $(EXAMPLES_FILE); \
	for f in data/*.json; do \
		c=`go run main.go -f $$f`; \
		echo "examples['$$f'] = '$$c';" >> $(EXAMPLES_FILE); \
	done; \
	echo "var examples_rows = {}" >> $(EXAMPLES_FILE); \
	for f in data/*.json; do \
		c=`go run main.go -f $$f --rows`; \
		echo "examples_rows['$$f'] = '$$c';" >> $(EXAMPLES_FILE); \
	done

github: deploy
	ghp-import -b $(GITHUB_PAGES_BRANCH) -p $(GITHUB_PAGES_FOLDER)

test:
	go test ./...

cover:
	go test -coverprofile cover.out ./git2graph/
	go tool cover -html=cover.out

.PHONY: deploy github test
