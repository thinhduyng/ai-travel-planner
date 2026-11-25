# AI-LOG.md

## 2026-11-24: choosing the stack (PA#2)
Tool: Claude Code.
Asked for: a plan for the PA#2 checkpoint and a stack to build the AI Travel Planner on.
Kept: the plan (rules file, permissions, one command per check, protected main, contract, contract tests).
Changed: nothing.
Rejected: its suggested TypeScript/Node backend (Fastify + Zod). The team prefers working in Go, so we switched the backend to Go and kept the TypeScript idea only for the future frontend.

## 2026-11-25: initial code and the API contract (PA#2)
Tool: Claude Code.
Asked for: the Go API skeleton, CI, assistant rules and permissions, the contract for `POST /v1/itineraries` (Markdown rules + OpenAPI schema for codegen), its contract tests, and the implementation.
Kept: all of it: `CLAUDE.md`, `.claude/settings.json`, `Makefile`, `ci.yml`, `docs/contract/create-itinerary.*`, the 36 contract tests and the handler. Checked with `make check` locally and the 7 CI jobs. Lint (gosec, errcheck) caught two issues in its first draft (an unvalidated `PORT` written to the log, an unchecked `Body.Close`), and it fixed them before we committed.
Changed: nothing.
Rejected: nothing.
