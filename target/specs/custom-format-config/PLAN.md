# Implementation Plan: Custom Format and JSON Config Loading (F-007)

**Date**: 2026-05-01 | **Spec**: `target/specs/custom-format-config/SPEC.md`  
**Input**: Feature specification from `target/specs/custom-format-config/SPEC.md`

**Note**: This plan is scoped strictly to F-007 and the shared architecture artifact at `target/specs/ARCHITECTURE`.

## Summary

Reimplement F-007 in Go by adding a bounded `internal/config` slice that supports: (1) inline `customFormat` field mapping with default fallback and opt-out semantics, and (2) `customPath` JSON loading with explicit legacy-compatible parse-error behavior. The plan follows architecture order (custom format at slice 8), uses stdlib-first config parsing (`os.ReadFile` + `encoding/json`), and preserves compatibility-sensitive behaviors identified in SPEC FR-001..FR-010.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` (F-007 only)
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Keep F-007 isolated to config/custom-field behavior; do not absorb unrelated filtering/policy/output logic.
- [x] Implement config parsing in `internal/config/custom_format.go` per architecture structure.
- [x] Use Go stdlib (`os.ReadFile`, `encoding/json`) for JSON file loading/parsing to stay aligned with stdlib-first architecture decision.
- [x] Preserve parser contract as **value-or-error-object return behavior** for compatibility until FR-010 decision is resolved.
- [x] Treat `customFormat[property] == false` as explicit field exclusion during module field assembly.
- [x] Context7 references used for planning conventions:
  - `/golang/go` and `/websites/go_dev_doc` for Go error handling and JSON parsing conventions.
  - `/spf13/cobra` for CLI error propagation via `RunE` where config validation/initialization failures are surfaced.

## Open Questions

- [x] **FR-010 resolved for bounded PoC slice**: keep compatibility-mode default from `ARCHITECTURE` and preserve legacy-style parse behavior by returning parser errors as values from parsing boundary helpers (no panic/throw behavior).
- [x] Compatibility mode default is used for this repository slice (strict mode deferred).

## Technical Context

**Language/Version**: Go (target per architecture; exact toolchain version **NEEDS CLARIFICATION**, recommend Go 1.24.x baseline from architecture context)  
**Primary Dependencies**: Go stdlib (`os`, `encoding/json`, `errors`, `fmt`, `context` where needed); Cobra at CLI boundary only (already architecture-chosen)  
**Storage**: Local filesystem (read JSON config file via `customPath`)  
**Testing**: `go test` with table-driven unit tests + fixture/integration checks + golden compatibility assertions where output shaping is involved  
**Target Platform**: Cross-platform Go CLI (Linux/macOS/Windows)  
**Project Type**: Go CLI + reusable package pipeline slice  
**Performance Goals**: No dedicated throughput target for F-007; must keep per-module custom-field mapping linear in configured field count  
**Constraints**: Preserve F-007 legacy semantics (fallback defaults, `false` opt-out, parse returns Error object compatibility), avoid scope creep beyond F-007  
**Scale/Scope**: One bounded migration slice: custom format ingestion + mapping + parse/error semantics

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/Users/igorbesel/Projekte/theseus/memory/constitution.md`: **not present**.
- Workflow/rule check from `AGENTS.md`: **PASS**
  - Plan stored in correct use-case folder under `target/specs/custom-format-config/`.
  - Go target preserved.
  - Scope bounded to F-007.
  - No code implementation performed.
- `.specify/extensions.yml`: **not present** (no before/after-plan hooks to surface).
- Repository setup/agent-context scripts: none discovered; proceed without script execution.

## Project Structure

### Documentation (this feature)

```text
target/specs/custom-format-config/
├── PLAN.md              # This file
├── research.md          # Optional; create only if FR-010 resolution research is needed
├── data-model.md        # Optional; create if additional model formalization is needed
├── quickstart.md        # Optional; execution/validation walkthrough for F-007
├── contracts/           # Optional; only if external interfaces are introduced
└── TASKS.md             # Produced later by Task Decomposer
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

internal/config/
  custom_format.go      # F-007 primary implementation surface

internal/flatten/
  flatten.go            # caller/integration point for module field shaping

pkg/licensechecker/
  api.go

test/
  integration/
  fixtures/
```

**Structure Decision**: Use the architecture’s Option 1 Go service/CLI layout; implement F-007 behavior primarily in `internal/config/custom_format.go` with integration touchpoints in core/flatten pipeline and CLI option mapping.

## Work Packages (F-007 Only)

### WP1 — Define F-007 data contracts and boundaries
- [x] Map SPEC entities to Go types/interfaces:
  - `CustomFormatConfig`
  - `ParsedCustomConfigResult` (object or error-compatible outcome)
  - module custom field projection contract
- [x] Record exact ownership boundaries between `internal/cli`, `internal/config`, and `internal/core`.
- **Traceability**: FR-001, FR-004, Entities section.

### WP2 — Inline custom format mapping behavior
- [x] Plan deterministic mapping routine: for each configured key, use module string value else configured default.
- [x] Plan explicit opt-out behavior for `value == false` (exclude property).
- [x] Ensure plan includes `licenseText` and `copyright` conditional enrichment hooks.
- **Traceability**: FR-001, FR-002, FR-003, FR-008, FR-009; SC-001, SC-004.

### WP3 — `customPath` JSON loading and parser behavior
- [x] Plan JSON-file load sequence at initialization stage before output shaping.
- [x] Plan parse behavior matrix for valid path, non-string path, missing file, malformed JSON.
- [x] Preserve/flag compatibility behavior for `Error` object return path pending FR-010 decision.
- **Traceability**: FR-004, FR-005, FR-006, FR-007, FR-010; SC-002, SC-003.

### WP4 — Error semantics and compatibility toggle decision
- [x] Resolve FR-010 with explicit repository decision (or document blocker if unresolved).
- [x] If unresolved at build start, require guarded implementation path with clearly marked compatibility mode.
- [x] Define failure/reporting behavior at CLI boundary using Cobra `RunE` propagation.
- **Traceability**: FR-010, Edge Cases, Architecture risk #4.

### WP5 — Validation and regression coverage plan
- [x] Unit tests: table-driven coverage for parser outcomes and mapping fallback semantics.
- [x] Integration tests: inline `customFormat`, file-backed `customPath`, false-property exclusion.
- [x] Compatibility assertions for CSV-related normalization behavior when custom fields include `licenseText`.
- [x] Confirm deterministic behavior for malformed/missing path scenarios (return vs fail based on FR-010 resolution).
- **Traceability**: SC-001..SC-004; User Stories 1..3; Edge Cases.

## Dependencies & Execution Order

1. **WP1** (contracts/boundaries) must finish before all other WPs.
2. **WP2** and **WP3** can proceed in parallel after WP1.
3. **WP4** (FR-010 decision) must be resolved before finalizing builder tasks for error paths.
4. **WP5** runs after WP2/WP3 behavior is fixed and WP4 decision is explicit.

## Validation Approach

- Validate every FR with at least one mapped test case:
  - FR-001/002/003: inline mapping + default fallback + false exclusion tables.
  - FR-004/005/006/007: parser input matrix (valid, missing, malformed, non-string).
  - FR-008/009: targeted field population checks when keys are configured.
  - FR-010: explicit expected behavior test split by resolved decision.
- Maintain parity-style fixtures where possible to confirm F-007 compatibility outcomes.
- Capture pass/fail evidence in downstream `TASKS.md`/`BUILD.md` (not in this planning artifact).

## Complexity Tracking

No constitution violations identified; table not required for this plan.
