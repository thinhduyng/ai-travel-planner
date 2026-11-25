# Self-assessment: PA#2

Submitted by:
    23127122 - Nguyễn Duy Thịnh
    23127357 - Lê Anh Duy
    23127487 - Thôi Đặng Trường Thọ

Repo: https://github.com/thinhduyng/ai-travel-planner (public) · tag [`pa2`](https://github.com/thinhduyng/ai-travel-planner/tree/pa2) · links: [`PA2.md`](PA2.md)

Total I claim: 92 / 100

| Criterion | Max | I claim | Evidence |
|---|---|---|---|
| Rules file | 10 | 9 | [`CLAUDE.md`](CLAUDE.md). -1 AI wrote it, no team edit |
| Permissions, no secrets | 15 | 14 | [`.claude/settings.json`](.claude/settings.json) allow/ask/deny; `.env*` ignored; gitleaks full history in CI ([run](https://github.com/thinhduyng/ai-travel-planner/actions/runs/45625811717)). -1 no local scan |
| Same command local + CI | 15 | 14 | [`Makefile`](Makefile), each [CI](.github/workflows/ci.yml) job = 1 `make` target. -1 `make check` skips `generate-check` |
| Protected main | 15 | 13 | 7 required checks, admins included, no force push ([pic](docs/pa2/branch-protection.png), [blocked](docs/pa2/pr2-merge-blocked.png)). -2 no required review |
| Contract | 15 | 14 | [`create-itinerary.md`](docs/contract/create-itinerary.md) §-numbered + [OpenAPI](docs/contract/create-itinerary.openapi.yaml) → codegen. Contract before code ([`696bdf2`](https://github.com/thinhduyng/ai-travel-planner/commit/696bdf2)). -1 md + yaml sync by hand |
| Contract tests in CI | 15 | 14 | [`api/contract/`](api/contract/) 36 tests over HTTP, named by §, job `contract`. -1 fake planner only |
| Gate bites | 15 | 14 | Red: [PR #3](https://github.com/thinhduyng/ai-travel-planner/pull/3) [run](https://github.com/thinhduyng/ai-travel-planner/actions/runs/45514941965); green: [`5e8c5cb`](https://github.com/thinhduyng/ai-travel-planner/commit/5e8c5cbe389dd8b21db5641aba7e48846c6d05d7); break: [PR #2](https://github.com/thinhduyng/ai-travel-planner/pull/2) [run](https://github.com/thinhduyng/ai-travel-planner/actions/runs/45513934539). -1 PR #1 red run cancelled, redo |

## What I did not manage

- No required review, 1 person can self-merge.
- PR #1 red run cancelled (pushed fix too fast), redo in PR #3.
- No local secret scan, CI first.

## What I would do differently

Protect main before everything + wait CI to finish before pushing the fix in PR #1.
