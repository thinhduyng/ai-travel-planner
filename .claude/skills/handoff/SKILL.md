---
name: handoff
description: >
  Save current task state and prepare a clean context
  handoff for the next Claude Code session.
---

# Session Handoff

Prepare a handoff for the current task.

## Inspect

- git status
- git diff
- Completed implementation
- Remaining implementation
- Tests and results
- Known bugs
- Technical decisions

## Save

Write or update:

.ai/tasks/<task-name>/status.md
.ai/tasks/<task-name>/handoff.md

## Handoff format

# Task Handoff

## Objective
What are we building?

## Completed
What has been implemented?

## Current state
What is actually working?

## Decisions
What technical decisions must be preserved?

## Relevant files
List only files needed by the next phase.

## Blockers
What remains unresolved?

## Next actions
Ordered, actionable steps.

## Validation
Commands and tests to run.

## Context rules

- Never copy the entire conversation.
- Do not duplicate CLAUDE.md.
- Do not repeat information available directly
  from source code.
- Verify facts against the actual repository.
- Keep the handoff concise and actionable.

Finish by telling the user how to start
the next session.
