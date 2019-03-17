IMAGE := gcr.io/msmini-item:$(VERSION)

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
	@go get github.com/golang/dep/cmd/dep
	@dep ensure -v -vendor-only
	@go test -race -coverpkg=./... -coverprofile=coverage.txt ./...
	@go get github.com/haya14busa/reviewdog/cmd/reviewdog
	@reviewdog -conf=.reviewdog.yml -diff="git diff master" -reporter=github-pr-review

.PHONY: test-local
test-local:
	@go test -race -cover ./...

.PHONY: cloudbuild
cloudbuild:
	@gcloud builds submit . \
		--project=msmini-item \
	 	--substitutions="_IMAGE=$(IMAGE),_VERSION=$(VERSION)"
