#!/bin/sh
# Mechanically rewrite this fork from the upstream module path to ours (and
# repoint any forked dependencies). Run by `make release` on a throwaway
# release branch — never on the patches branch — so rebasing on upstream stays
# conflict-free. Config in .fork/config.
set -eu

. "$(dirname "$0")/config"

# rewrite OLD->NEW everywhere it matters, precisely:
#   - quoted import strings ("github.com/OLD...) — the leading quote is what an
#     import spec always has and what comment URLs / path-literal test fixtures
#     never have, so those correctly keep pointing upstream.
#   - //go:linkname directives, which reference the path UNQUOTED in a comment
#     (a dangling linkname breaks at run time, not compile time).
rewrite() {
	old="$1"
	new="$2"
	grep -rl "$old" --include='*.go' . 2>/dev/null | while IFS= read -r f; do
		sed -i "s#\"$old#\"$new#g" "$f"
		if grep -q 'go:linkname' "$f"; then
			sed -i "/go:linkname/ s#$old#$new#g" "$f"
		fi
	done
}

rewrite "$OLD" "$NEW"
go mod edit -module "$NEW"

# Repoint forked dependencies: DEPS="oldpath=newpath@version ...".
for dep in ${DEPS:-}; do
	dold=${dep%%=*}
	drest=${dep#*=}
	dnew=${drest%@*}
	dver=${drest##*@}
	rewrite "$dold" "$dnew"
	go mod edit -dropreplace="$dold" 2>/dev/null || true
	go mod edit -droprequire="$dold" 2>/dev/null || true
	go mod edit -require="$dnew@$dver"
done

go mod tidy
