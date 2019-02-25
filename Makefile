LINT_TOOLS=\
	golang.org/x/lint/golint \
	github.com/client9/misspell \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck

.PHONY: test-coverage-reviewdog
test-coverage-reviewdog:
	@for tool in $(LINT_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool; \
	done
	@dep ensure -v -vendor-only
	@go test -race -coverpkg=./... -coverprofile=coverage.txt ./...
	@go get github.com/haya14busa/reviewdog/cmd/reviewdog
	@reviewdog -conf=.reviewdog.yml -diff="git diff master" -reporter=github-pr-review
