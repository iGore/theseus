# Implementation Plan: F-003 Dependency Flattening

**Date**: 2026-05-01 | **Spec**: `target/specs/dependency-flattening/SPEC.md`  
**Input**: Feature specification from `target/specs/dependency-flattening/SPEC.md`

**Note**: This plan is scoped only to F-003 (flatten dependency graph into module inventory map) and uses the shared architecture baseline in `target/specs/ARCHITECTURE`.

## Summary

Reimplement the flattening stage in Go as `internal/flatten/flatten.go`, converting a nested dependency tree into a deterministic, deduplicated `ModuleInventoryMap` keyed by `name@version`, while preserving required metadata and applying prod/dev inclusion gating during traversal. The plan keeps strict stage boundaries (graph-load input -> flatten output contract) and defers unresolved legacy ambiguities (missing identity fields, metadata mandatory set) as explicit blockers/clarifications rather than inferred behavior.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` (F-003 only)
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Flatten implementation stays isolated in `internal/flatten` and consumes upstream-loaded graph only (no loader logic in scope). (Trace: FR-001, FR-002, C-001)
- [x] Canonical identity remains `name@version` and drives dedupe/cycle guards via visited-key semantics. (Trace: FR-001, FR-003, C-003)
- [x] Deterministic outcomes are defined as stable key set + uniqueness, with deterministic presentation achieved by sorted key materialization in tests/consumers where ordering is needed. (Trace: FR-007, SC-001/SC-002; Context7 `/golang/go`: map iteration order unspecified)
- [x] Metadata extraction remains in flatten stage and is carried in inventory values for downstream stages; exact mandatory-vs-optional field contract remains an open compatibility decision. (Trace: FR-005, OD-002)
- [x] Validation strategy uses table-driven unit tests with fixture graphs plus focused cycle/duplicate/gating cases. (Trace: AC-001..AC-004; Context7 `/stretchr/testify` testing patterns)

## Open Questions

- [ ] **OD-001 (NEEDS CLARIFICATION)**: Legacy-compatible behavior when `name` and/or `version` is missing during key construction.
- [ ] **OD-002 (NEEDS CLARIFICATION)**: Final mandatory metadata field contract for flattened entries (vs optional passthrough fields).

## Technical Context

**Language/Version**: Go (targeting repo architecture baseline; Go toolchain version to be pinned during build stage)  
**Primary Dependencies**: Go stdlib (`context`, maps/slices/sort patterns), `github.com/stretchr/testify` for tests  
**Storage**: N/A for this slice (in-memory graph -> in-memory map transform)  
**Testing**: `go test` with table-driven tests; testify `require/assert` for preconditions and result checks  
**Target Platform**: Cross-platform Go CLI/library runtime (Linux/macOS/Windows)  
**Project Type**: Go CLI + reusable package architecture; this slice is internal pipeline stage  
**Performance Goals**: Complete traversal of all reachable included nodes with O(V+E) walk behavior and single-entry-per-key dedupe semantics  
**Constraints**: Preserve key format `name@version`; no graph-loading or license-resolution logic in this slice; recursion must short-circuit revisits/cycles without non-termination  
**Scale/Scope**: Bounded to F-003 only (flatten traversal, dedupe/cycle guard, metadata extraction, prod/dev gate)

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/memory/constitution.md`: **not present**.
- Repository workflow rules (`AGENTS.md`) checked:
  - Go target preserved ✅
  - Spec-first flow preserved (planning from SPEC + ARCHITECTURE only) ✅
  - Artifact location under `target/specs/<use-case-slug>/` preserved ✅
  - Scope bounded to a single functionality slice (F-003) ✅
- `.specify/extensions.yml`: **not present** (no before_plan/after_plan hooks to surface).
- Repository setup/agent-context scripts: none discovered via pattern scan (continued without script execution).

## Work Packages (Checklist-Oriented)

### WP1 — Flatten stage contract alignment
- [x] Confirm F-003 input/output contract mapping into Go domain types (`DependencyNode` -> `ModuleInventoryMap`) without widening scope.
- [x] Record explicit compatibility notes for stage boundary (flatten output consumed by downstream filter/render stages).
- [x] Mark OD-001/OD-002 as gating clarifications for strict parity assertions.

**Traceability**: FR-001, FR-006, C-004, AC-005

### WP2 — Recursive traversal and canonical key strategy
- [x] Define recursive traversal algorithm over `dependencies` edges (depth-unbounded walk of reachable nodes).
- [x] Define canonical key construction path (`name@version`) and visited-key guard behavior.
- [x] Define revisit semantics: first encounter wins, skip re-expansion for existing key.

**Traceability**: FR-001, FR-002, FR-003, FR-007, AC-001, AC-002, SC-001, SC-002

### WP3 — Inclusion gating at flatten time
- [x] Define where prod/dev gate evaluation occurs in recursion flow (before insertion/expansion for excluded nodes).
- [x] Ensure exclusion behavior is testable independently from upstream option parsing.

**Traceability**: FR-004, AC-003, SC-003

### WP4 — Metadata extraction contract for inventory entries
- [x] Enumerate currently evidenced metadata fields (repository/author/url/path/private-related) into flatten output mapping.
- [x] Tag field-level certainty: confirmed-by-spec vs needs-clarification.

**Traceability**: FR-005, AC-004, SC-004, OD-002

### WP5 — Determinism and validation design
- [x] Define deterministic validation approach that avoids relying on Go map iteration order.
- [x] Plan fixture suite for: deep nesting, duplicate paths, cycles, mixed prod/dev inclusion, metadata preservation.
- [x] Plan downstream contract harness check that accepts flattened map without original tree dependency.

**Traceability**: FR-007, AC-001..AC-005, SC-001..SC-004

## Dependencies and Execution Order

1. **WP1** (contract baseline) must complete before algorithm details are locked.
2. **WP2** depends on WP1 contract decisions.
3. **WP3** depends on WP2 traversal checkpoints.
4. **WP4** can proceed in parallel with WP3 after WP1, but parity sign-off depends on OD-002 clarification.
5. **WP5** depends on WP2/WP3/WP4 decisions to build complete validation matrix.

## Validation Approach

- Unit-level table-driven tests for flatten logic with subtests per behavior axis (traversal depth, dedupe/cycle, gating, metadata).
- Fixture-driven compatibility cases for reachable-node coverage and one-entry-per-key guarantees.
- Determinism checks compare key sets and value expectations independent of map iteration order (keys sorted before comparison where ordering matters).
- Downstream contract test ensures flatten result shape is sufficient input for next stage harness.

Context7-backed planning notes:
- `/golang/go`: map iteration order is unspecified; tests and any order-sensitive output checks must sort keys.
- `/stretchr/testify`: use `require` for preconditions and `assert` for non-fatal multi-field verification in table-driven subtests.

## Project Structure

### Documentation (this feature)

```text
target/specs/dependency-flattening/
├── PLAN.md              # This file
├── research.md          # Optional: extra planning research if clarifications expand
├── data-model.md        # Optional: field-level model notes if needed
├── quickstart.md        # Optional: validation walkthrough for this slice
├── contracts/           # Optional: explicit stage contract docs for F-003
└── TASKS.md             # Produced later by Task Decomposer
```

### Source Code (repository root)

```text
cmd/license-checker/
  main.go

internal/
├── core/
│   ├── model.go
│   └── service.go
├── graph/
│   └── loader.go
└── flatten/
    └── flatten.go

pkg/licensechecker/
  api.go

test/
└── integration/
```

**Structure Decision**: Use the Go service/CLI architecture defined in `target/specs/ARCHITECTURE`, with F-003 implemented strictly in `internal/flatten` and validated through unit + integration fixtures.

## Complexity Tracking

No constitution violations identified; section not applicable for this plan revision.
