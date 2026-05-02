---
name: orchestrator-completion
description: >
  Use this skill before declaring any phase or the overall workflow complete.
  It defines explicit exit criteria for each phase gate and for the full
  pipeline. The Orchestrator runs in a continuous verification loop: emit a
  completion promise only after all checklist items pass, then run Oracle
  verification before stopping. Do not rely on a subjective sense of "done".
---

# Orchestrator Completion

## Purpose

The Orchestrator operates in a self-referential completion loop. Believing
the work is done is not the same as the work being done. This skill provides
the Verifier checklist that must pass before the Orchestrator may stop.

A subjective "I think this phase is complete" is a **completion promise** —
not a completion fact. This checklist (run by the Verifier) must confirm the promise
before the Orchestrator proceeds or stops.

## Completion Loop Protocol

1. When you believe a phase or the workflow is complete, do **not** stop.
2. Emit an internal completion promise: note that you believe the phase is done.
3. Immediately run the Verifier checklist for that phase (checklist below).
4. Only if **every item passes** may you treat the phase as complete.
5. If any item fails, continue working rather than stopping or skipping.
6. Repeat until the Verifier checklist passes for all phases.

The loop only ends after Phase 3 Verifier checklist confirms every item.
The `steps` parameter in the agent profile ensures enough iterations to
complete this loop without premature termination.

## Phase 1 Checklist — Rekonstruktion

Invoke the Verifier for Phase 1. The Verifier checks each step granularly:

- ANALYSIS.md exists and has all required sections (Scope, Entry Points, Function Inventory, Risk Map)
- OVERVIEW.md exists with Use Cases and Requirements Task List
- FUNCTIONALITY-INDEX.md exists, every item has an ID, slug, and explicit status
- Every `ready` item has a SPEC.md with Gherkin scenarios, FR entries with source references, and Success Criteria
- SPEC-INDEX.md lists all requirement folders

**Checklist result**: Verifier returns PASS → Phase 1 complete. Any FAIL → continue Phase 1.

## Phase 2 Checklist — Transformationsplanung

Invoke the Verifier for Phase 2. The Verifier checks each step granularly:

- ARCHITECTURE exists with Technical Context, Structure Decision, Package Decisions (each with documented reason), and Implementation Order
- Every use-case folder has a PLAN.md with work packages and explicit dependencies

**Checklist result**: Verifier returns PASS → Phase 2 complete. Any FAIL → continue Phase 2.

## Phase 3 Checklist — Neuimplementierung

Invoke the Verifier for Phase 3. The Verifier checks each step granularly:

- Every use-case folder has a TASKS.md with correct format, story labels, and no unchecked boxes
- Every use-case folder has a BUILD.md with verification evidence
- `go build ./...` and `go vet ./...` pass from `target/`

**Checklist result**: Verifier returns PASS → workflow complete. Any FAIL → continue Phase 3.

## Guardrails

- Do not declare a phase complete if any checklist item is unverified
- Do not assume a background run finished — confirm the artifact exists on disk
- If a required artifact is missing, route the work back to the responsible role rather than skipping the gate
- Do not substitute a summary or status note for an actual artifact check
- A completion promise without a passing Verifier checklist has no effect — keep working
