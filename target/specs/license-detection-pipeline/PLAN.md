# Implementation Plan: License detection and normalization pipeline (F-004)

**Date**: 2026-05-01 | **Spec**: `target/specs/license-detection-pipeline/SPEC.md`  
**Input**: Feature specification from `target/specs/license-detection-pipeline/SPEC.md`

**Note**: This plan is scoped strictly to F-004 and derived from `SPEC.md` + shared `target/specs/ARCHITECTURE`.

## Summary

Reimplement the F-004 Go license-resolution slice that produces exactly one deterministic normalized license value per module by applying the required precedence chain: metadata (`license`/`licenses`) -> README-derived signal -> ordered license-file fallback. Preserve compatibility-significant normalization semantics (`*` for inferred values, `Custom: <value>` for URL/file references) and keep this slice side-effect free relative to filtering/policy/output stages.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md`
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Scope boundary locked to F-004 only (no F-002/F-003/F-005/F-006/F-007 behavior changes).
- [x] Go target retained, mapped to `internal/license` (`detect.go`, `normalize.go`, `files.go`) per shared architecture.
- [x] CLI behavior and process exits are excluded from this slice; F-004 returns normalized data only (FR-008).
- [x] Context7-aligned package conventions adopted: keep non-public implementation under `internal/*` and preserve export boundary via package API (`/websites/go_dev_doc`, modules/layout).
- [x] Context7-aligned execution/error convention noted for integration boundary: prefer error-returning flows (`RunE`/returned errors) over inline exits where command integration is involved (`/spf13/cobra`).

## Open Questions

- [ ] **OD-001 / NEEDS CLARIFICATION**: Exact legacy tie-break when both `license` and `licenses` are present but disagree.
- [ ] **OD-002 / NEEDS CLARIFICATION**: Exact fallback output when no recognizable signal exists in metadata, README, or candidate files.

## Technical Context

**Language/Version**: Go 1.24.x (aligned with architecture Go target)  
**Primary Dependencies**: Go stdlib for F-004 core logic; existing module may use Cobra/testify at integration/test layers but F-004 logic should remain stdlib-first  
**Storage**: Local filesystem reads for README/license candidate files (no DB)  
**Testing**: `go test` table-driven unit tests + fixture-driven compatibility tests for precedence/normalization parity  
**Target Platform**: Cross-platform Go CLI/library runtime (Linux/macOS/Windows) with deterministic CI behavior  
**Project Type**: Go CLI + reusable library, implementing bounded internal pipeline slice  
**Performance Goals**: Deterministic per-module resolution with linear scan over bounded candidate set; no separate throughput target defined in spec  
**Constraints**: Preserve legacy precedence and normalization semantics; no side effects outside license-value production; deterministic outputs for identical inputs  
**Scale/Scope**: Single pipeline slice (F-004) affecting module-level license resolution only

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/memory/constitution.md`: not present.
- `.specify/extensions.yml`: not present (no `hooks.before_plan` / `hooks.after_plan` to surface).
- Repository setup/agent-context update scripts: none found; no script execution required.
- Workflow gate result: **PASS** for planning continuation with current artifacts.

## Project Structure

### Documentation (this feature)

```text
target/specs/license-detection-pipeline/
├── PLAN.md              # This file
├── research.md          # Optional; only if further clarification research is needed
├── data-model.md        # Optional; only if expanded modeling is needed
├── quickstart.md        # Optional; validation walkthrough if requested
├── contracts/           # Optional; external interface contracts if introduced
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
└── license/
    ├── detect.go
    ├── normalize.go
    └── files.go

pkg/licensechecker/
  api.go

test/
└── integration/
```

**Structure Decision**: Use the shared Go service/CLI layout from `ARCHITECTURE`; confine F-004 work to `internal/license` plus minimal `core` model/service wiring and focused tests.

## Work Packages (F-004 only)

### WP1 — Domain contract for license-resolution inputs/outputs
- [x] Define/confirm `ModuleLicenseInput`, `NormalizedLicenseValue`, and candidate-file descriptors in core model interfaces used by F-004.
- [x] Ensure contract supports metadata forms (`license`, `licenses` string/object/array), README signal input, and ordered file candidates.
- **Traceability**: FR-001, FR-002, Key Entities.

### WP2 — Metadata parsing and normalization-first path
- [x] Implement/plan metadata resolver path that evaluates SPDX-compatible expressions before heuristic inference.
- [x] Preserve inferred marker behavior (`*`) when heuristic classification is used.
- [x] Normalize URL/file references to exact `Custom: <value>` output shape.
- **Traceability**: FR-001, FR-003, FR-004, FR-005; AC-002, AC-003.

### WP3 — Deterministic fallback orchestration
- [x] Implement/plan precedence chain: metadata -> README -> file candidates.
- [x] Implement/plan deterministic file candidate order: `LICENSE*`, `LICENCE*`, `COPYING`, then `README`.
- [x] Ensure README without recognizable license text does not terminate fallback early.
- **Traceability**: FR-002, FR-006; AC-001, AC-004; Edge Cases.

### WP4 — Compatibility lock-in for unresolved semantics
- [x] Add explicit regression test placeholders for OD-001 and OD-002 outcomes once confirmed from legacy evidence.
- [x] Keep behavior compatibility-focused with existing parser tests.
- **Traceability**: FR-007, C-002, Open Decisions; AC-005.

### WP5 — Boundary integrity and side-effect isolation
- [x] Verify F-004 stage emits normalized value only and does not perform filtering/policy exits or output formatting logic.
- [x] Ensure deterministic behavior across repeated runs with identical fixtures/inputs.
- **Traceability**: FR-008, C-001, C-003, C-004; SC-001, SC-004.

## Dependencies & Execution Order

1. WP1 (domain contract)  
2. WP2 (metadata parse/normalize path)  
3. WP3 (fallback orchestration + deterministic file selection)  
4. WP5 (boundary integrity checks integrated while implementing WP2/WP3)  
5. WP4 (compatibility lock-in tests for known/open semantics)

Blocking dependencies:
- OD-001 and OD-002 require evidence confirmation before final parity sign-off.

## Validation Approach

- [x] Unit tests (table-driven):
  - metadata form coverage (`license`, `licenses` variants),
  - SPDX vs heuristic (`*`) distinction,
  - `Custom: <value>` normalization,
  - deterministic filename precedence.
- [x] Fixture/integration tests:
  - metadata-first, README-first, file-first precedence scenarios,
  - repeat-run determinism assertions,
  - legacy parser-coverage parity for F-004.
- [x] Negative/edge scenarios:
  - conflicting `license` vs `licenses` (pending OD-001 expected result),
  - no recognizable signal fallback value (pending OD-002 expected result),
  - README present but no match continues to file fallback.

## Complexity Tracking

No constitution violations identified; section not applicable at this stage.
