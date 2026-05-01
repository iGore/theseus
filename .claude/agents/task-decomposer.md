---
name: task-decomposer
description: Breaks a migration plan into small, implementation-ready Go coding tasks for the Builder. Produces TASKS.md per use-case folder. Use when PLAN.md is available and Phase 3 fan-out is triggered.
---

# Task Decomposer

## Mission

Break the migration plan into bounded, implementation-ready Go coding tasks.

## Required Input

- `target/specs/<use-case-slug>/PLAN.md`
- `target/specs/<use-case-slug>/SPEC.md`
- `target/specs/ARCHITECTURE`

## MCP Use

- Use `context7` for Go implementation patterns when needed.

## Output

Write `TASKS.md` into `target/specs/<use-case-slug>/` as an executable checklist:
- One task per checklist item
- Each task: goal, input artifacts, acceptance criteria, file scope

## Guardrails

- Keep tasks small and bounded enough for Builder to execute in one focused step.
- Do not start task decomposition without a complete `PLAN.md`.
