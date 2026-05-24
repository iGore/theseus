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
- ANALYSIS.md contains `## Use Cases` and `## Requirements Task List` sections (non-empty)
- For source tools that consume external data beyond CLI arguments, ANALYSIS.md contains an `## Input-Field Inventory` section (see `cli-analyzer` skill) listing, for every externally-supplied input field, its default path, every conditional override with source-line reference and triggering condition, and an explicit `unconditional` marker for overrides that no CLI flag, env var, or config key controls. Every unconditional override MUST have a corresponding probe fixture and Golden Output Snippet recorded in ANALYSIS.md
- Every unconditional override identified in the Input-Field Inventory is surfaced as an explicit FR in the SPEC.md of the slice that owns the affected output field, not only in the slice that owns a flag of the same name
- USE-CASES.md exists, every item has an ID, slug, and explicit status
- Every `ready` item has a SPEC.md with Gherkin scenarios, FR entries with source references, and Success Criteria
- SPECS.md lists all requirement folders

**Checklist result**: Verifier returns PASS → Phase 1 complete. Any FAIL → continue Phase 1.

## Phase 2 Checklist — Transformationsplanung

Invoke the Verifier for Phase 2. The Verifier checks each step granularly:

- ARCHITECTURE.md exists with Technical Context, Structure Decision, Package Decisions (each with documented reason), and Implementation Order
- For CLI-oriented target systems, ARCHITECTURE.md additionally contains a `## Composition Root` section with entry-point path, ordered stage sequence, concrete stage constructors, a flag-to-stage propagation table, and an acceptance reference to the Golden Output Snippets in ANALYSIS.md (see `go-architect` skill)
- Every use-case folder has a PLAN.md with work packages and explicit dependencies

**Checklist result**: Verifier returns PASS → Phase 2 complete. Any FAIL → continue Phase 2.

## Phase 3 Checklist — Neuimplementierung

Invoke the Verifier for Phase 3. The Verifier checks each step granularly:

- Every use-case folder has a TASKS.md with correct format, story labels, and no unchecked boxes
- Every use-case folder has a BUILD.md with verification evidence
- `go build ./...` and `go vet ./...` pass from `target/`
- For CLI-oriented target systems, an end-to-end binary exists and is verifiable:
  - `target/cmd/<tool>/main.go` is present, derived from the Composition Root section of ARCHITECTURE.md
  - `go build -o /tmp/cli-bin ./cmd/...` from `target/` produces an executable file
  - The binary responds to the standard discovery flags it declares (typically `--help` and `--version`) with the exit codes and output channels defined in the corresponding SPEC.md
  - The binary, executed on each of the minimal fixtures referenced in ANALYSIS.md, produces output that semantically matches the corresponding Golden Output Snippet. Format-level divergences such as ordering of structurally unordered collections, insignificant whitespace, or platform-dependent absolute paths are tolerated; missing entities, divergent output schemas, or empty outputs are FAIL conditions
  - A pure library build without a `cmd/` entry point counts as FAIL — passing `go build ./...` on packages alone is not sufficient when the migration targets a CLI tool

**Checklist result**: Verifier returns PASS → workflow complete. Any FAIL → continue Phase 3.

## Guardrails

- Do not declare a phase complete if any checklist item is unverified
- Do not assume a background run finished — confirm the artifact exists on disk
- If a required artifact is missing, route the work back to the responsible role rather than skipping the gate
- Do not substitute a summary or status note for an actual artifact check
- A completion promise without a passing Verifier checklist has no effect — keep working
