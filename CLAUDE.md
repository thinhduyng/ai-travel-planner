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

## AI-LOG.md (course requirement, CSC13114)
`AI-LOG.md` at the root records every task where an assistant did part of the work. It is part
of every milestone submission (PA#1–PA#5). When you finish a task, append an entry (or draft one
for the human to confirm). Newest entries go at the bottom.

```md
## <YYYY-MM-DD>: <what was worked on> (<milestone, e.g. PA#2>)
Tool: <Claude Code / Copilot / Cursor / none>
Asked for: <the task, in one sentence>
Kept: <what was taken as it came>
Changed: <what the human rewrote, and why>
Rejected: <what was thrown away, and why>
```

- Three honest lines are enough: an account, not a transcript. Write it as you go, not the night before.
- `Changed` and `Rejected` are read most closely: say what the assistant got wrong and **which
  check caught it** (lint, a contract test, a review). "Nothing went wrong" is fine if it says what was run.
- Never overstate the human's part or understate the assistant's. Heavy use declared honestly is fine;
  the team must be able to explain everything in it at the oral.
- Never put secrets, API keys, real user data or a classmate's work in a prompt or in the log.
- For the product's own LLM feature, the log (or the docs) must say what data leaves the system and why that is acceptable.
- Bad entry (don't write this): "Week 3. Used AI for the frontend. Fixed some things."

## Git and secrets
- `main` is protected: work on a branch, open a PR, and let CI decide. Never force-push.
- The human makes the commits. Prepare changes and stop; don't run `git commit` unless asked.
- Never read, print or commit `.env`. Document new variables in `.env.example` with an empty value.
