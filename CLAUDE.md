# AI Travel Planner: rules for the coding assistant

These rules apply to Claude Code and any other coding agent working in this repo.
The permissions that enforce them are in `.claude/settings.json`.

## Layout
- `api/`: Go backend (Go module `github.com/thinhduyng/ai-travel-planner/api`, standard `net/http`).
  - `cmd/server`: entry point only.
  - `internal/server`: route wiring (`server.New`).
  - `internal/api`: **generated** from the OpenAPI spec. Never edit `api.gen.go`; run `make generate`.
  - `internal/itinerary`: handlers and validation for the contract endpoints.
  - `internal/planner`: the `Planner` interface in front of the LLM, plus `Fake` for tests/CI.
  - `contract/`: HTTP contract tests.
- `docs/contract/`: endpoint contracts (`*.md` rules + `*.openapi.yaml` schema).
- `web/`: frontend (not created yet).

## Commands. The same ones run in CI (`.github/workflows/ci.yml`)
| Check | Command |
|---|---|
| Lint | `make lint` (golangci-lint, version pinned in `Makefile`) |
| Format check | `make fmt-check` (fix with `make fmt`) |
| Type check | `make typecheck` (`go vet`) |
| Generated code up to date | `make generate-check` (regenerate with `make generate`) |
| Unit tests | `make test` |
| Contract tests | `make contract` |
| All of the above | `make check` |

Run `make check` before saying a task is done. If it fails, the task is not done.

## The contract is the source of truth
- `docs/contract/*.md` defines the behaviour; each rule has a § number. Do **not** edit
  contracts. If one looks wrong, stop and say which sentence and why.
- The OpenAPI file next to it is the schema used for codegen. If it and the Markdown disagree, the Markdown wins. Report the mismatch.
- Do not edit `api/contract/**` to make a failing test pass. Fix the code.
- Contract test names start with the § they check (`§2.1 ...`).

## Code rules
- Errors always use the contract shape `{"error":{"code","message","details":[]}}`; `details` is never null.
- The LLM sits behind `planner.Planner`. Tests and CI use `planner.Fake` and never call a real
  model or need an API key.
- Never put provider errors, stack traces or secrets in a response. Log them instead.
- Standard library first; ask before adding a dependency.

## Git and secrets
- `main` is protected: work on a branch, open a PR, and let CI decide. Never force-push.
- The human makes the commits. Prepare changes and stop; don't run `git commit` unless asked.
- Never read, print or commit `.env`. Document new variables in `.env.example` with an empty value.
