# Fork notes

Maintained fork with performance patches, published under our own module path
so consumers need **no `replace` directive** (Go ignores `replace` directives
in dependencies).

## Two branches

- **patches** (the default branch) — upstream + our patches, keeps the upstream
  module path. Rebasing on upstream only conflicts on our patches.
- **`release`** — generated from it by `.fork/rename.sh` with the module path
  renamed to ours; never rebased. Consumers depend on this via tags.

## Versioning

`vX.Y.Z-fgsi.N` — upstream version + our patch iteration.

## Workflows (all config in `.fork/config`)

```sh
make verify          # build + test the patches branch
make rebase          # fetch upstream + rebase our patches
make release         # regenerate the renamed release branch
git checkout release && make verify-release
git tag vX.Y.Z-fgsi.N release && git push origin release --tags
```

The `fork-sync` GitHub Action runs this weekly and on demand.
