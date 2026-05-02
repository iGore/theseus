# Implementation Plan: Dependency Graph Loading and Traversal Control (F-002)

**Date**: 2026-05-01 | **Spec**: `target/specs/dependency-graph-loading/SPEC.md`  
**Input**: Feature specification from `target/specs/dependency-graph-loading/SPEC.md`

**Note**: Scope is strictly limited to F-002 graph-loading behavior (start path, depth, dev toggle, success/error handoff) and its Go reimplementation plan.

## Summary

Reimplement the F-002 slice as a Go graph-loading stage that:
- accepts normalized scan inputs,
- builds loader options (`depth`, `dev`, logging hook),
- invokes dependency graph acquisition from the exact `start` path,
- forwards either loaded tree (success) or surfaced error (failure) to the next stage.

The plan follows shared architecture sequencing (F-002 after CLI/normalization) and uses a loader abstraction in `internal/graph` orchestrated by `internal/core`.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md`
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Keep F-002 isolated to `internal/graph` + `internal/core` orchestration boundaries from `ARCHITECTURE` (no flatten/filter/output behavior in this plan).
- [x] Treat option normalization as upstream dependency (F-001); F-002 consumes normalized `direct` depth and scope flags only.
- [x] Use context-aware subprocess invocation conventions for loader adapter behavior (Context7 `/golang/go`: `context`, `os/exec` cancellation/error handling patterns).
- [x] Preserve compatibility-first semantics from spec/architecture: direct-only depth (`0`), recursive depth otherwise, dev-loading default enabled unless production/development flag set.

## Open Questions

- [ ] None currently blocking for F-002 planning slice.

## Technical Context

**Language/Version**: Go (target Go 1.24.x per architecture baseline; exact patch version can follow repo toolchain)  
**Primary Dependencies**: Go stdlib (`context`, `os/exec`, `errors`, `log/slog`), internal packages `internal/core`, `internal/graph`; test support via `testing`/`testify` as defined in shared architecture  
**Storage**: N/A for this slice (read-only graph acquisition from local project path)  
**Testing**: `go test` with table-driven unit tests + fixture-driven integration checks for loader success/failure and option mapping  
**Target Platform**: Cross-platform CLI runtime (Linux/macOS/Windows)  
**Project Type**: Go CLI + reusable package (slice implemented in internal pipeline packages)  
**Performance Goals**: `NEEDS CLARIFICATION` (no explicit numeric SLO in F-002 spec)  
**Constraints**: Must preserve legacy behavior parity for start path, depth semantics, dev toggle semantics, and error propagation; no code paths from out-of-scope F-003+  
**Scale/Scope**: Single bounded pipeline stage (graph load boundary only)

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/Users/igorbesel/Projekte/theseus/memory/constitution.md`: not present (no additional constitution gates discovered).
- `.specify/extensions.yml`: not present (no before/after-plan hooks discovered).
- Repository root scanned for setup/agent-context update scripts: none found; no script execution required.
- AGENTS.md workflow gates satisfied:
  - Plan artifact stored in use-case folder under `target/specs/`.
  - Go target preserved.
  - Scope bounded to one functionality slice (F-002).

## Project Structure

### Documentation (this feature)

```text
target/specs/dependency-graph-loading/
├── SPEC.md
└── PLAN.md
```

### Source Code (repository root)

```text
cmd/license-checker/
  main.go

internal/cli/
  command.go
  options.go

internal/core/
  service.go
  model.go

internal/graph/
  loader.go
  npm_ls_loader.go
```

**Structure Decision**: Use the shared Go service/CLI structure from `ARCHITECTURE`; F-002 work is constrained to `internal/graph` contracts/adapters and `internal/core` orchestration touchpoints.

## Work Packages (F-002 only)

### WP1 — Confirm F-002 boundary contracts and inputs
- [x] Map FR-001..FR-008 and SC-001..SC-004 from spec to explicit F-002 stage responsibilities.
- [x] Confirm consumed inputs from upstream normalization (start, normalized direct depth semantics, production/development flags).
- [x] Define exact stage input/output contract for `core.Service` ↔ `graph.DependencyLoader` handoff.

**Traceability**: FR-001..FR-008, Assumptions section, SRC-AN-002A/B.

### WP2 — Plan loader option mapping rules (no implementation)
- [x] Document deterministic mapping rules:
  - `depth=0` when direct mode enabled,
  - recursive depth sentinel otherwise,
  - `dev=true` by default,
  - `dev=false` when `production || development`.
- [x] Define where mapping is applied in pipeline sequence (during graph-load stage initialization, before loader call).
- [x] Define logger linkage expectations for loader options without introducing cross-slice behavior.

**Traceability**: FR-002..FR-006, SC-002..SC-003, SRC-CODE-ARGS, SRC-CODE-INIT.

### WP3 — Plan dependency loader invocation and error channel semantics
- [x] Define invocation contract ensuring exact `start` path is used for each load.
- [x] Define success handoff contract to next stage boundary (flattening input only, no flatten logic).
- [x] Define failure contract: propagate non-nil loader errors and block successful handoff.
- [x] Align subprocess/context cancellation and error interpretation behavior with Go stdlib conventions (Context7 `/golang/go`).

**Traceability**: FR-001, FR-007, FR-008, SC-001, SC-004, SRC-TEST-ERR.

### WP4 — Validation plan for F-002 compatibility
- [x] Unit cases for options mapping matrix (direct on/off × production/development combinations).
- [x] Unit cases for loader invocation argument integrity (start path passthrough).
- [x] Unit/integration cases for success and failure handoffs (tree forwarded vs error surfaced).
- [x] Fixture case for invalid/unresolvable start path to verify failure propagation.

**Traceability**: User Stories 1-3, Edge Cases, SC-001..SC-004.

## Dependencies & Execution Order

1. **Upstream prerequisite**: F-001 option normalization available (spec assumption).
2. **WP1** contract confirmation.
3. **WP2** options mapping plan.
4. **WP3** invocation/error semantics plan.
5. **WP4** validation planning and acceptance mapping.

Blocking dependency: if F-001 normalization output contract changes, F-002 mapping assumptions must be revalidated before task decomposition.

## Validation Approach

- Map each FR and SC to at least one planned verification case.
- Preserve strict stage boundary: validation must not assert flatten/filter/output semantics (out of scope).
- Ensure all failure-path checks assert both conditions: error surfaced + no success handoff.
- Record compatibility evidence against cited source anchors from SPEC (`lib/index.js`, `lib/args.js`, `tests/test.js`) during Builder phase, not in this planning artifact.

## Complexity Tracking

No constitution violations identified; section intentionally left without entries.
