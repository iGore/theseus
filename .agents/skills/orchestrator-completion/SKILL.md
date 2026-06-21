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
  - A pure library build without a `cmd/` entry point counts as FAIL — passing `go build ./...` on packages alone is not sufficient when the migration targets a CLI tool

### Phase 3 — Golden-Output Parity Gate (test-driven, in-flow)

When the migration goal is to reproduce a source system's behavior, semantic
"close enough" is **not** an acceptable gate — it is the documented root cause
of silent output drift. The in-flow gate compares the target against the
**captured example/golden outputs**, byte-for-byte, via committed tests:

- The golden outputs are the example outputs recorded **once** by the Analyzer
  (Golden Output Snippets, captured from the pinned source version named in
  ANALYSIS.md) and committed as fixture files. They are the test oracle inside
  the flow — the flow does **not** install or run the live source system.
- The Builder writes **golden-file tests** (table-driven, committed under the
  target test tree) that run the target binary across the recorded mode × flag
  matrix and edge set and assert **byte equality** against the golden files.
- Only an explicit, documented **normalization allowlist** may be applied before
  comparison — limited to genuinely volatile values (e.g. absolute paths under a
  temp root). Each normalization must be named in BUILD.md. Ordering of
  collections, whitespace, sentinels, and trailing newlines are **in scope** and
  may not be hand-waved away.
- The Verifier confirms these golden tests **exist, are byte-level, cover the
  recorded matrix, and pass** (`go test ./...`). It checks the recorded golden
  outputs and the target — it does **not** run the live source system.
- **Any failing golden test is FAIL**, unless the affected mode/flag is listed as
  a deliberate exclusion in a `## Conformance Matrix` (see escape hatch below).

### Out-of-flow final parity run (manual, not a gate)

A live differential run against the actually-installed source system
(a maintainer-run differential script: source oracle vs. target across the full
matrix) is a **manual verification performed by the maintainers at the end**,
outside the agent flow. Agents MUST NOT install or execute the live source system as part
of any phase gate. The committed golden outputs are what makes that final manual
run cheap and likely to pass on the first try.

### Conformance Matrix (escape hatch when 100% parity is not the goal)

If full byte-parity is explicitly out of scope, the project MUST declare this
up front rather than discovering it at the gate. ARCHITECTURE.md (or SPECS.md)
carries a `## Conformance Matrix` listing, per mode/flag: `identical`,
`compatible` (documented intentional difference), or `out-of-scope`. The golden
gate then FAILs only on failing rows marked `identical`. Claiming "produces the
same results" while the matrix contains undeclared `compatible`/`out-of-scope`
rows is itself a FAIL.

**Checklist result**: Verifier returns PASS → workflow complete. Any FAIL → continue Phase 3.

## Guardrails

- Do not declare a phase complete if any checklist item is unverified
- Do not assume a background run finished — confirm the artifact exists on disk
- If a required artifact is missing, route the work back to the responsible role rather than skipping the gate
- Do not substitute a summary or status note for an actual artifact check
- A completion promise without a passing Verifier checklist has no effect — keep working
