# Contributing

`linctl` is a schema-aligned Go CLI for Linear. Keep each change small and typed, and back it
with the generated GraphQL schema.

## Local checks

`go tool task ci` is the one local review gate. Run it before you open a PR.

```bash
go tool task ci
```

It runs these steps in order: `deps-check`, `fmt-check`, `generate-check`,
`domain-language-check`, `browser-login-smoke-check`, `vet`, `test`, a build, `smoke-run`,
`lint`, `shellcheck`, `actionlint`, and `vuln`. `deps-check` verifies the module checksums and
the tidy state. `generate-check` verifies the generated GraphQL client, the generated command
references, and the upstream schema and coverage-ledger checks. None of these steps modifies a
source file.

`go tool task coverage` is separate from `ci`, and it enforces 100% statement coverage on
hand-written code. Run it on its own after a change to product code.

```bash
go tool task coverage
```

### Network and offline

`go tool task ci` validates the local GraphQL operations and
`docs/internal/linear-api-coverage.md` against the upstream `linear/linear` schema. The shared
source contract is:

- Remote: `https://github.com/linear/linear.git`
- Default checkout: a run-local temporary checkout that the Taskfile manages
- Default ref: the commit SHA in `.linear-sdk-ref`. This file is the single source of the pin,
  and CI reads the same file.
- Reusable checkout: `LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear`
- Override the ref: `LINCTL_LINEAR_SDK_REF=<branch-or-tag>`
- Skip the refresh fetch against a reusable checkout: `LINCTL_LINEAR_SDK_OFFLINE=1`. This needs
  `LINCTL_LINEAR_SDK_UPSTREAM` to point at a checkout already, and it fails loud otherwise.

Prepare or refresh a reusable checkout:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear \
go tool task linear-sdk-upstream-checkout
```

Run the gate offline after that checkout exists:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear \
LINCTL_LINEAR_SDK_OFFLINE=1 \
go tool task ci
```

`scripts/upstream-check.sh` runs the upstream checks behind one checkout. Run one check on its
own with `operations`, `ledger`, `drift`, or `all`:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear bash scripts/upstream-check.sh operations
```

GitHub runs three checks outside `go tool task ci`. Dependency review runs only on a pull
request. Coverage stays in `go tool task coverage`. The live OAuth check and the integration
check need disposable fixture credentials.

Run a live integration test only with a disposable OAuth app fixture:

```bash
LINCTL_OAUTH_CLIENT_ID=<client-id> \
LINCTL_OAUTH_CLIENT_SECRET=<client-secret> \
LINCTL_OAUTH_REDIRECT_URI=http://127.0.0.1:8765/callback \
LINCTL_OAUTH_SCOPES=read,write,issues:create,comments:create \
LINCTL_OAUTH_EXPECTED_ACTOR=app \
go test -count=1 -tags=integration ./internal/client
```

The full live smoke harness is:

```bash
go tool task live-smoke
```

For the Infisical setup of this project, the fixture secrets are under `/linctl`, not under the
root secret path. Use the pinned aliases, so nobody forgets the folder:

```bash
go tool task live-oauth-infisical
go tool task live-smoke-infisical
```

Never run a write test against real project data. A test resource must use the name prefix
`linctl-it-<runid>`, and the cleanup must archive it.

## Schema changes

A nightly `schema-drift` job (`.github/workflows/integration.yml`) reports the semantic drift
between the vendored `internal/client/schema.graphql` and the upstream `linear/linear` SDK
schema. It reports each type, field, and enum value added or removed, and each field type change
on a type in both schemas. It is read-only, and it is not part of `go tool task ci`, so it never
blocks a PR. It exists to show when the vendored snapshot falls behind. Run it locally:

```bash
go tool task schema-drift-check
```

A companion nightly `schema-refresh` workflow (`.github/workflows/schema-refresh.yml`) opens or
updates a `chore/schema-refresh` PR when it detects drift. That workflow copies the upstream SDK
schema, bumps `.linear-sdk-ref`, regenerates the client and the coverage ledger, and leaves the
PR for a human review. It never merges by itself. Prefer that PR over a hand-rolled refresh when
it is already open.

The repo secret `SCHEMA_REFRESH_GITHUB_TOKEN` is optional. It is a fine-grained PAT with
contents and pull-requests write access on this repo. When it is set, the workflow uses it, and
the opened PR then triggers `ci.yml`. With `GITHUB_TOKEN` alone, GitHub does not trigger a
workflow on the bot PR.

To refresh by hand without OAuth, use the same path as the bot:

```bash
LINCTL_LINEAR_SDK_REF=master go tool task schema-refresh
go tool task ci
```

To refresh through live API introspection, you need a token and the managed Node dependencies:

```bash
npm ci
LINCTL_OAUTH_ACCESS_TOKEN=<token> bash scripts/refresh-schema.sh
go generate ./...
go tool task ci
```

Review the generated diff before you commit it. Drift in a schema area that linctl does not use
is safe in a routine refresh. Drift on a field or a type that a current operation depends on is
urgent, but `go tool task graphql-operation-check` already fails in that case.

Refresh the vendored Linear schema before you add or change a GraphQL operation:

```bash
npm ci
./scripts/refresh-schema.sh
go generate ./...
```

`scripts/refresh-schema.sh` uses the repo-managed `graphql` dependency from `package-lock.json`,
and it needs Node 22 or later. Set `LINCTL_OAUTH_ACCESS_TOKEN` for the command, and never print
or paste the token value into a log. Commit the generated code together with the operation that
needs it.

## Releases

A push of a `v*` tag starts a release. The release workflow then runs GoReleaser, which publishes
the GitHub artifacts (archives, SBOMs, `checksums.txt`, and a keyless cosign sigstore bundle) and
updates the `KyaniteHQ/homebrew-linctl` tap cask. The tap token must be
`HOMEBREW_TAP_GITHUB_TOKEN`.

Run the local non-publishing release preflight before you create and push the tag:

```bash
go tool task release-preflight
```

The preflight runs the local CI gate, the statement coverage, and `goreleaser check`. It creates
no tag, it pushes nothing to git, it publishes no release, and it needs no release secret. For a
heavier local artifact build as a final manual check, run the snapshot task. It also publishes
nothing:

```bash
go tool task release-snapshot
```

Create and push the release tag only after the preflight passes:

```bash
git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z
```
