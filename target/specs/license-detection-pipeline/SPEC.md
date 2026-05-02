# Feature Specification: License detection and normalization pipeline

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Input artifacts: `/Users/igorbesel/Projekte/theseus/target/specs/ANALYSIS.md`, `OVERVIEW.md`, `FUNCTIONALITY-INDEX.md`. Focus slice only: F-004 — License detection and normalization pipeline (slug `license-detection-pipeline`). Task: Produce implementation-ready `SPEC.md` for this single functionality in `/Users/igorbesel/Projekte/theseus/target/specs/license-detection-pipeline/SPEC.md`. Include explicit requirements, acceptance criteria, compatibility notes, constraints, and open decisions limited to this slice. Do not implement code. Do not cover unrelated functionality."  

## Scope

This specification covers only **F-004 — License detection and normalization pipeline**.

- In scope: extraction, precedence, parsing, normalization, and marking of license values per dependency record.
- Out of scope: dependency traversal (F-002/F-003), filtering and policy exits (F-005), output rendering/export (F-006), custom format loading (F-007).

## Source Traceability Registry

- **S1**: F-004 index entry and evidence pointers (`FUNCTIONALITY-INDEX.md:23-27`).
- **S2**: Requirement linkage R-004: preserve precedence SPDX -> heuristics -> files (`OVERVIEW.md:49`).
- **S3**: F-004 behavior summary (`ANALYSIS.md:100-109`).
- **S4**: `flatten` performs license collection and fallback chain (`ANALYSIS.md:60`, `ANALYSIS.md:102-106`, `ANALYSIS.md:162-163`).
- **S5**: Parser supports `license` and `licenses` shapes, README fallback, and inferred marker `*` (`ANALYSIS.md:104-108`).
- **S6**: License file precedence order is deterministic (`ANALYSIS.md:191`, `ANALYSIS.md:224`).
- **S7**: URL/file refs normalized as `Custom: ...` (`ANALYSIS.md:108`, `ANALYSIS.md:203`).
- **S8**: Parser test coverage for license heuristics exists in `tests/license.js` (`FUNCTIONALITY-INDEX.md:26`, `ANALYSIS.md:225`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Resolve a license value for each module (Priority: P1)

As a compliance user, I need each scanned module to receive a normalized license value so inventory and policy stages can operate on consistent data.

**Why this priority**: Without normalized license resolution, downstream filtering/policy and reporting cannot function correctly. This is the core purpose of F-004. (Sources: S1, S3)

**Independent Test**: Can be fully tested by feeding module metadata with known `license`/`licenses` variants and verifying one normalized license output per module.

**Acceptance Scenarios**:

1. **Given** a module with `license` as an SPDX-valid string, **When** the pipeline resolves license data, **Then** that SPDX value is returned as the module license without file-based fallback. (Sources: S3, S5)
2. **Given** a module with `licenses` as array/object forms, **When** the pipeline resolves license data, **Then** the module receives a normalized license derived from those forms. (Sources: S3, S5)

---

### User Story 2 - Fall back deterministically when metadata is incomplete (Priority: P2)

As a compliance user, I need deterministic fallback from metadata to README and then to license files, so repeated scans produce stable, explainable license outcomes.

**Why this priority**: Real package metadata is often incomplete; deterministic fallback preserves reproducibility and legal auditability. (Sources: S2, S3, S6)

**Independent Test**: Can be tested with fixture packages missing metadata but containing README/license files and verifying precedence order and stable outputs across runs.

**Acceptance Scenarios**:

1. **Given** a module with no usable `license`/`licenses` fields but README containing recognizable license text, **When** resolution runs, **Then** README-derived result is used before scanning license files. (Sources: S3, S5)
2. **Given** a module requiring file fallback, **When** directory files include multiple candidates, **Then** selection follows deterministic precedence (`LICENSE*`, `LICENCE*`, `COPYING`, `README`). (Sources: S3, S6)

---

### User Story 3 - Preserve heuristic and custom-reference normalization semantics (Priority: P3)

As a compatibility-focused maintainer, I need inferred and non-standard license signals normalized exactly as legacy behavior expects.

**Why this priority**: Compatibility drift in edge normalization breaks historical outputs and downstream expectations. (Sources: S3, S7, S8)

**Independent Test**: Can be tested by passing non-SPDX text and URL/file-style references and verifying `*` inference markers and `Custom: ...` formatting.

**Acceptance Scenarios**:

1. **Given** license text that is inferred heuristically (not exact SPDX), **When** parsing completes, **Then** the normalized value is marked with `*` (for example `MIT*`). (Sources: S3, S5)
2. **Given** a license field containing URL/file-style custom reference, **When** parsing completes, **Then** the output is normalized as `Custom: <reference>`. (Sources: S3, S7)

### Edge Cases

- What happens when both `license` and `licenses` are present but conflict? Resolve using existing source behavior precedence and document final tie-break outcome in tests for compatibility lock-in. (Sources: S3, S8)  
- What happens when README exists but contains no recognizable license phrase? Pipeline must continue to file-based fallback rather than stopping at README presence. (Sources: S3, S5)  
- What happens when no license signal is found in metadata, README, or candidate files? Output must remain explicit and stable (exact legacy value requires compatibility verification). (Sources: S3, S8)  
- How does the system handle uncommon filename casing/variants beyond listed precedence buckets? Preserve current deterministic matcher boundaries unless explicitly changed. (Sources: S6)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST resolve module license information from package metadata fields supporting `license` and `licenses` (string/object/array forms) before fallback mechanisms. (Sources: S3, S5)
- **FR-002**: System MUST apply fallback precedence in this order: metadata fields -> README-derived signal -> license-file-derived signal. (Sources: S2, S3)
- **FR-003**: System MUST parse and normalize SPDX-compatible expressions through the dedicated license parser path prior to heuristic inference. (Sources: S1, S3)
- **FR-004**: System MUST mark heuristically inferred license values with a trailing `*` to distinguish inferred vs direct recognition. (Sources: S3, S5)
- **FR-005**: System MUST normalize URL/file-style non-standard license references to `Custom: <value>`. (Sources: S3, S7)
- **FR-006**: System MUST use deterministic license-file candidate precedence: `LICENSE*`, `LICENCE*`, `COPYING`, then `README` when file fallback is required. (Sources: S3, S6)
- **FR-007**: System MUST keep license resolution behavior compatible with parser expectations covered by existing license parser tests. (Sources: S1, S8)
- **FR-008**: System MUST keep this slice side-effect free with respect to process exits and output formatting decisions; it only produces normalized license values for later stages. (Sources: S3, S4)

### Acceptance Criteria (Requirement-Level)

- **AC-001 (FR-001/FR-002)**: For fixtures with complete metadata, README-only, and file-only license signals, outputs prove precedence order by selecting metadata first, README second, files third. (Sources: S2, S3)
- **AC-002 (FR-003/FR-004)**: SPDX-valid expressions pass through normalized parser path; heuristic matches are distinguishable by `*`. (Sources: S3, S5)
- **AC-003 (FR-005)**: URL/file license references produce `Custom: ...` exactly. (Sources: S7)
- **AC-004 (FR-006)**: In a fixture containing multiple candidate files, selected source follows deterministic sequence defined in FR-006. (Sources: S6)
- **AC-005 (FR-007)**: Legacy parser test scenarios for license normalization pass without behavior regression for this slice. (Sources: S8)

### Key Entities *(include if feature involves data)*

- **ModuleLicenseInput**: Per-module inputs consumed by the pipeline (package metadata license fields, README content signal, candidate file list/content). (Sources: S3, S4)
- **NormalizedLicenseValue**: Final resolved string used by downstream filtering/reporting (may be SPDX form, inferred form with `*`, or `Custom: ...`). (Sources: S3, S5, S7)
- **LicenseCandidateSet**: Ordered candidate files considered during file fallback (`LICENSE*`, `LICENCE*`, `COPYING`, `README`). (Sources: S6)

## Compatibility Notes

- Must preserve legacy precedence and normalization semantics for F-004 to avoid downstream drift in filtering/policy and output layers. (Sources: S2, S3)
- Inferred-license marker `*` is compatibility-significant and must remain observable in resolved values. (Sources: S3, S5)
- `Custom: ...` formatting for URL/file references is compatibility-significant and must remain unchanged unless explicitly approved in a separate decision. (Sources: S7)
- Deterministic license-file precedence is an invariant and must not be reordered in this slice. (Sources: S6)

## Constraints

- **C-001 (Scope constraint)**: This slice MUST NOT introduce or change behavior in filtering/policy exits, CLI options, or rendering layers. (Sources: S1, S3)
- **C-002 (Evidence constraint)**: Only source-observed behaviors are in scope; ambiguous behavior remains an open decision and cannot be guessed. (Sources: S8)
- **C-003 (Determinism constraint)**: For identical module inputs, license resolution output MUST be deterministic. (Sources: S6)
- **C-004 (Pipeline boundary constraint)**: Output of this slice is normalized license data for later pipeline stages, not final user-facing formatting. (Sources: S4)

## Open Decisions (slice-limited)

1. **OD-001**: When `license` and `licenses` both exist and disagree, confirm exact legacy tie-break order from direct code/tests and lock with explicit regression tests for rewrite parity. (Sources: S3, S8)
2. **OD-002**: Confirm exact fallback value when no recognizable signal exists across metadata, README, and files (value shape is not explicit in provided Analyzer summary for this slice). (Sources: S3, S8)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of modules in compatibility fixtures receive one deterministic normalized license value via metadata/README/file fallback path. (Sources: S2, S3, S6)
- **SC-002**: Existing parser-focused test coverage for license normalization (including heuristic and custom-reference cases) passes with zero regressions for F-004 behavior. (Sources: S8)
- **SC-003**: Inference-vs-direct distinction is preserved: all heuristic outcomes are marked with `*`, and direct SPDX-normalized outcomes are not. (Sources: S3, S5)
- **SC-004**: File-based fallback selects the same candidate file source across repeated runs for identical inputs (deterministic precedence invariant). (Sources: S6)

## Assumptions

- Analyzer evidence for F-004 is sufficient to define precedence and normalization obligations without requiring target-stack design choices. (Sources: S1, S3)
- Input module records supplied by upstream flattening stage already contain access to required metadata and readable filesystem context. (Sources: S4)
- This specification intentionally defers unresolved edge semantics to explicit open decisions rather than inferring undocumented behavior. (Sources: S8)
