# Fork maintenance (see FORK.md). Two branches:
#   main/master — upstream + our patches, upstream module path (easy rebase).
#   release     — generated from it with the module path renamed to ours; never
#                 rebased. Consumers pin tags vX.Y.Z-fgsi.N.
include .fork/config

RELEASE_BRANCH ?= release

.PHONY: verify rebase release verify-release

verify:
	go build ./...
	go test ./...

rebase:
	@git remote get-url upstream >/dev/null 2>&1 || git remote add upstream $(UPSTREAM_URL)
	git fetch upstream
	git rebase upstream/$(UPSTREAM_BRANCH)

release:
	git checkout -B $(RELEASE_BRANCH)
	sh .fork/rename.sh
	git add -A
	git commit -q -m "release: module $(NEW) (generated)"
	git checkout -

verify-release:
	go build ./...
	go test ./...
