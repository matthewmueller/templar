precommit: test

test:
	@ go mod tidy
	@ go vet ./...
	@ go run honnef.co/go/tools/cmd/staticcheck@latest ./...
	@ go test -race ./...

release: VERSION := $(shell awk '/[0-9]+\.[0-9]+\.[0-9]+/ {print $$2; exit}' Changelog.md)
release: test build-dist install-dist
	@ go mod tidy
	@ test -n "$(VERSION)" || (echo "Unable to read the version." && false)
	@ test -z "`git tag -l $(VERSION)`" || (echo "Aborting because the $(VERSION) tag already exists." && false)
	@ test -z "`git status --porcelain | grep -vE 'Changelog\.md'`" || (echo "Aborting from uncommitted changes." && false)
	@ git add Changelog.md
	@ git commit -m "Release $(VERSION)"
	@ git tag "$(VERSION)"
	@ git push origin main "$(VERSION)"
	@ go run github.com/cli/cli/v2/cmd/gh@latest release create --generate-notes "$(VERSION)" dist/*.tar.gz dist/*.txt
