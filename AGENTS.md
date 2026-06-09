# AGENTS.md

## Repository purpose

`garethpaul/fsq-go-explore` is a Go project. Foursquare GoLang Explore Sample

## Project structure

- `Makefile` - repository verification targets
- `scripts` - baseline checks and helper scripts
- `docs` - plans, notes, and generated README assets
- `templates` - server-rendered templates
- `go.mod` - Go module definition

## Development commands

- Install dependencies: `go mod download`
- Full baseline: `make check`
- Go test all packages: `go test ./...`
- Go vet all packages: `go vet ./...`
- Go build all packages: `go build ./...`
- If a command above skips because a platform toolchain is missing, verify on a machine with that SDK before claiming platform behavior is tested.

## Coding conventions

- Language mix noted in the README: Go (14).
- Keep imports compatible with module path `github.com/garethpaul/fsq-go-explore`.
- Run gofmt on changed Go files and keep table-driven tests close to the package under change.

## Testing guidance

- Test-related files detected: `auth_test.go`, `cache_test.go`, `edit_test.go`, `fsq/api_test.go`, `fsq/keys_test.go`, `search_test.go`
- Start with the narrowest relevant test or Make target, then run `make check` before handing off if the change is not documentation-only.
- Keep README verification notes in sync when commands, fixtures, or supported toolchains change.

## PR / change guidance

- Keep diffs focused on the requested repository and avoid unrelated modernization or formatting churn.
- Preserve public APIs, sample behavior, file formats, and documented environment variables unless the task explicitly changes them.
- Update tests, README notes, or docs/plans when behavior, security posture, or validation commands change.
- Call out skipped platform validation, legacy toolchain assumptions, and any risky files touched in the final summary.

## Safety and gotchas

- Required Foursquare settings: `FSQ_CLIENT_ID`, `FSQ_CLIENT_SECRET`, and `FSQ_VERSION`.
- Keep API keys, OAuth credentials, access tokens, `.env` files, and deployment-specific config out of source control.
- Cache keys are deterministic SHA-256 digests and should not expose raw query, token, or user fields.
- OAuth login uses per-request state values and HTTP-only cookies for callback validation.
- OAuth callbacks with matching state still fail before token exchange; missing OAuth authorization codes are rejected.
- Auth cookie values are validated as generated user cache keys before memcache lookup, so malformed cookie values do not reach access-token cache work.

## Agent workflow

1. Inspect the README, Makefile, manifests, and the files directly related to the request.
2. Make the smallest source or docs change that satisfies the task; avoid generated, vendored, or local-environment files unless required.
3. Run the narrowest useful validation first, then `make check` or the documented package/platform gate when available.
4. If a required SDK, service credential, or external runtime is unavailable, record the skipped command and why.
5. Summarize changed files, commands run, and remaining risks or follow-up validation.
