# BUILD: License detection and normalization pipeline (F-004)

## 1. Implementation Scope

- **Files/modules changed (F-004 only)**
  - `target/internal/core/model.go`
  - `target/internal/core/service.go`
  - `target/internal/core/service_license_test.go`
  - `target/internal/license/detect.go`
  - `target/internal/license/normalize.go`
  - `target/internal/license/files.go`
  - `target/internal/license/detect_test.go`
  - `target/internal/license/normalize_test.go`
  - `target/internal/license/files_test.go`
  - `target/test/integration/license_resolution_fallback_test.go`
  - `target/test/integration/license_resolution_compat_test.go`
  - `target/test/integration/testdata/license_resolution/README.md`
  - `target/specs/license-detection-pipeline/PLAN.md`
  - `target/specs/license-detection-pipeline/TASKS.md`

- **Spec requirements covered**
  - FR-001/FR-002: metadata-first then README then file fallback implemented in `detect.go`.
  - FR-003: SPDX-first path in `normalize.go` (`isSPDXExpression` before heuristic inference).
  - FR-004: inferred values marked with trailing `*`.
  - FR-005: URL/file refs normalized to `Custom: <value>`.
  - FR-006: deterministic file precedence in `files.go` (`LICENSE*`, `LICENCE*`, `COPYING`, `README`).
  - FR-007: compatibility placeholders created for unresolved legacy semantics (OD-001/OD-002).
  - FR-008: stage boundary introduced as pure resolution (`LicenseStage` + `ResolveLicenses`) with no filter/policy/output side effects.

## 2. Build Steps

1. Extended F-004 domain contracts in core model (`ModuleLicenseInput`, `LicenseCandidateFile`, `NormalizedLicenseValue`).
2. Added service-stage seam (`LicenseStage`) and per-module single-value assignment (`ResolveLicenses`).
3. Implemented metadata parser supporting `license` / `licenses` string-object-array forms.
4. Implemented normalization path: SPDX first, then heuristic (`*`), plus custom-reference format.
5. Implemented deterministic ordered file fallback matcher.
6. Added unit tests and integration tests for precedence, deterministic behavior, normalization semantics, and unresolved OD placeholders.
7. Updated `PLAN.md` and `TASKS.md` checklists to reflect completed build tasks.

## 3. Test Evidence

- **Tests added/updated**
  - `target/internal/core/service_license_test.go`
  - `target/internal/license/detect_test.go`
  - `target/internal/license/normalize_test.go`
  - `target/internal/license/files_test.go`
  - `target/test/integration/license_resolution_fallback_test.go`
  - `target/test/integration/license_resolution_compat_test.go` (explicit skip placeholders for OD-001/OD-002)

- **Command run**
  - `go test ./...` (run in `target/`)

- **Observed result**
  - PASS:
    - `ok theseus/target/internal/core`
    - `ok theseus/target/internal/license`
    - `ok theseus/target/test/integration`
    - other existing packages also passed
  - No failing tests; OD placeholders are explicit `t.Skip(...)` and documented as unresolved semantics.

## 4. Residual Risks

- **OD-001 unresolved**: exact legacy tie-break when `license` and `licenses` disagree is still `NEEDS CLARIFICATION`; current PoC behavior prefers `license` first.
- **OD-002 unresolved**: exact legacy fallback when no signal exists is still `NEEDS CLARIFICATION`; current PoC returns explicit sentinel `UNKNOWN`.
- SPDX detection is intentionally bounded for PoC compatibility and may require expansion after legacy parity confirmation.

## 5. Traceability

- **WP1 / FR-001 / Key Entities**
  - `internal/core/model.go`: `ModuleLicenseInput`, `NormalizedLicenseValue`, `LicenseCandidateFile`.
- **WP2 / FR-003 FR-004 FR-005 / AC-002 AC-003**
  - `internal/license/normalize.go` + `internal/license/normalize_test.go`.
- **WP3 / FR-002 FR-006 / AC-001 AC-004 / SC-004**
  - `internal/license/detect.go`, `internal/license/files.go`, `internal/license/detect_test.go`, `test/integration/license_resolution_fallback_test.go`.
- **WP4 / FR-007 / AC-005 / Open Decisions**
  - `test/integration/license_resolution_compat_test.go` skip placeholders with `NEEDS CLARIFICATION`.
- **WP5 / FR-008 / C-001 C-004**
  - `internal/core/service.go`: `LicenseStage` and `ResolveLicenses` remain side-effect free and pipeline-local.
