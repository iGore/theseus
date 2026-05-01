# Implementation Plan: Output Rendering and Export (F-006)

**Date**: 2026-05-01 | **Spec**: `target/specs/output-rendering-and-export/SPEC.md`  
**Input**: Feature specification from `target/specs/output-rendering-and-export/SPEC.md`

**Note**: This plan is strictly scoped to F-006 rendering/export behavior and aligned to shared architecture in `target/specs/ARCHITECTURE`.

## Summary

Reimplement the F-006 slice in Go by adding a deterministic rendering/export stage that: (1) selects one output mode via fixed precedence (`json > csv > markdown > summary > tree`), (2) emits rendered text to stdout or `--out` with parent directory creation, and (3) optionally exports detected license files to `--files`. The plan follows the architecture’s pipeline stage boundaries (`internal/render` + `core` orchestration) and preserves compatibility-first behavior from the spec.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` (F-006 only)
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Keep F-006 implementation in Go only, under architecture-defined render/export stage (`internal/render/*`, `core.Service` integration).
- [x] Preserve renderer selection semantics exactly per FR-002 and edge-case priority ordering.
- [x] Use Cobra flag values as inputs only; enforce mode resolution in core render-selection logic for testability and parity.
- [x] Use stdlib-first file operations for output persistence and export (`filepath.Dir`, `os.MkdirAll`, `os.WriteFile`, `os.Open`/`os.Create`, `io.Copy`) per architecture package decision.
- [x] Preserve optional colorized package-key behavior only for non-structured terminal-compatible outputs (FR-006).
- [x] Keep Markdown + `--files` behavior implemented regardless of documentation-status ambiguity; track doc alignment as follow-up clarification (FR-008).

## Open Questions

- [ ] **NEEDS CLARIFICATION**: Should `--markdown` and `--files` remain compatibility-only behaviors or be promoted to first-class documented CLI options in downstream docs/release notes?
- [ ] **NEEDS CLARIFICATION**: Exact compatibility rule for “terminal interactivity” gating colorized package keys (TTY detection strategy and non-TTY fallback expectations).

## Technical Context

**Language/Version**: Go (target per architecture; repository default).  
**Primary Dependencies**: `github.com/spf13/cobra` (flag plumbing already chosen in architecture), Go stdlib (`encoding/json`, `encoding/csv`, `bytes`/`strings`, `os`, `path/filepath`, `io`, `fmt`, `sort` as needed for deterministic render output).  
**Storage**: Local filesystem only (`--out` write target and `--files` export directory).  
**Testing**: `go test` with table-driven unit tests + golden renderer tests + fixture-driven integration compatibility tests (architecture §2).  
**Target Platform**: Cross-platform Go CLI (Linux/macOS/Windows).  
**Project Type**: Go CLI + reusable library core (`pkg/licensechecker` calling into `internal/core`).  
**Performance Goals**: `NEEDS CLARIFICATION` (spec defines correctness/compatibility KPIs, not throughput/latency SLA).  
**Constraints**: Preserve strict mode precedence, stdout default behavior, parent dir creation for `--out`, combined `--out` + `--files` side effects in one run, and no scope expansion beyond F-006.  
**Scale/Scope**: Single pipeline slice (stage 6) affecting render/export only; upstream graph/flatten/license/filter inputs treated as already prepared.

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/memory/constitution.md`: not present in repository (no additional constitution gates found).
- Workflow gates from `AGENTS.md`: satisfied for this artifact:
  - Spec-first flow respected (using existing `SPEC.md` + shared `ARCHITECTURE`).
  - Plan stored in use-case folder under `target/specs/output-rendering-and-export/`.
  - Go target maintained.
  - Scope remains bounded to F-006.
- `.specify/extensions.yml`: not present; no executable `hooks.before_plan` / `hooks.after_plan` to surface.
- Repository-provided setup/agent-context scripts: none detected; no script output to capture.

## Project Structure

### Documentation (this feature)

```text
target/specs/output-rendering-and-export/
├── SPEC.md
├── PLAN.md              # This file
├── research.md          # Optional (only if clarification research is needed)
├── data-model.md        # Optional (not required for current F-006 scope)
├── quickstart.md        # Optional validation walkthrough
├── contracts/           # Optional, if external interfaces are formalized later
└── TASKS.md             # Produced by Task Decomposer
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

internal/render/
  tree.go
  json.go
  csv.go
  markdown.go
  summary.go
  files_export.go

pkg/licensechecker/
  api.go

test/
  integration/
  fixtures/
```

**Structure Decision**: Use the Go service/CLI structure defined in shared `ARCHITECTURE` and constrain F-006 work to `internal/render`, `internal/core` integration points, and render/export-focused tests.

## Work Packages (F-006 only)

### WP1 — Renderer Mode Resolution Contract (FR-001, FR-002, FR-007)
- [x] Define/confirm render mode enum + selection helper in core/render boundary.
- [x] Encode strict precedence logic for simultaneous mode flags: `json > csv > markdown > summary > tree`.
- [x] Add table-driven tests for all relevant flag combinations, including default-to-tree behavior.

### WP2 — Renderer Implementations and Determinism (FR-001, SC-001)
- [x] Implement/align five renderer outputs (tree/json/csv/markdown/summary) against source-compatible semantics.
- [x] Enforce deterministic ordering rules required by golden tests (stable key/module iteration contract).
- [x] Add/update golden tests per renderer with compatibility fixtures.

### WP3 — Output Sink Behavior (FR-003, FR-004, SC-002, SC-004)
- [x] Implement sink routing: stdout when `--out` unset; file write when set.
- [x] Ensure parent dir creation for nested output path (`os.MkdirAll(filepath.Dir(outPath), ...)`).
- [x] Add integration tests for nested path creation, non-empty output persistence, and stdout emission.

### WP4 — License File Export Side Effect (FR-005, SC-003)
- [x] Implement `--files <dir>` export flow using detected per-module license file paths.
- [x] Define target export layout policy for copied artifacts and collision handling consistent with compatibility goals.
- [x] Add tests validating: export occurs when set, does not occur when omitted.

### WP5 — Combined Behavior + Compatibility Edges (FR-006, Edge Cases)
- [x] Verify combined execution when both `--out` and `--files` are provided in a single run.
- [x] Verify colorized package-key behavior is limited to non-structured terminal-compatible contexts.
- [x] Capture unresolved doc-status behavior (`--markdown`, `--files`) in implementation notes/tests without expanding scope.

## Dependencies & Execution Order

1. WP1 (mode resolution contract)  
2. WP2 (renderer implementations + determinism)  
3. WP3 (stdout/file sink behavior)  
4. WP4 (license file export)  
5. WP5 (combined-edge compatibility checks)

Rationale: mode contract first prevents downstream ambiguity; render output stability is prerequisite for sink/export integration and compatibility assertions.

## Validation Approach

- **Unit tests**
  - Mode precedence matrix tests (FR-002).
  - Renderer-level serialization tests for each format (FR-001).
  - Export helper tests for copy and path handling (FR-005).
- **Golden tests**
  - Tree/JSON/CSV/Markdown/Summary output fixtures to lock format compatibility.
- **Integration tests**
  - `--out` nested path creation + content write checks.
  - stdout default behavior when `--out` absent.
  - `--files` export-only and `--out + --files` combined run.
- **Success criteria traceability**
  - SC-001 ↔ renderer parity test matrix
  - SC-002 ↔ output file creation/content tests
  - SC-003 ↔ exported artifact count/path assertions
  - SC-004 ↔ stdout presence assertions

## Context7 Notes (Planning Evidence)

- Context7 `/spf13/cobra`: confirms CLI validation conventions and execution model (`RunE`, hook/validation flow, flag-group helpers) used to keep mode flags as CLI inputs while centralizing deterministic resolution in core logic.
- Context7 `/golang/go` (go1.24.x): confirms stdlib-oriented file/path/copy APIs for F-006 output persistence and export behavior (`os.MkdirAll`, file writing, `io.Copy`, `filepath` cross-platform handling).

## Complexity Tracking

No constitution violations identified; table intentionally left empty.
