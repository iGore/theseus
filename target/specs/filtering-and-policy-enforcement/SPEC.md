# Feature Specification: F-005 Filtering, package restriction, and policy exits

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Input artifacts: `/Users/igorbesel/Projekte/theseus/target/specs/ANALYSIS.md`, `OVERVIEW.md`, `FUNCTIONALITY-INDEX.md`. Focus slice only: F-005 — Filtering, package restriction, and policy exits (slug `filtering-and-policy-enforcement`)."

## Scope and Traceability

This specification is strictly bounded to functionality slice **F-005** (filtering behavior and policy exit behavior). It excludes dependency traversal, license detection internals, and output rendering modes except where needed as preconditions for this slice.

### Source Artifact IDs

- **SA-FI-F005**: `FUNCTIONALITY-INDEX.md` F-005 definition and open question (`29-34`).
- **SA-OV-R005**: `OVERVIEW.md` requirement task for filtering semantics (`50`).
- **SA-OV-R006**: `OVERVIEW.md` requirement task for policy fail exits (`51`).
- **SA-AN-F005**: `ANALYSIS.md` F-005 behavior summary (`110-120`).
- **SA-AN-FLOW**: `ANALYSIS.md` flow ordering for filters then policy exits (`163-165`).
- **SA-AN-EDGE**: `ANALYSIS.md` edge cases for escaped commas, BSD expansion, private handling (`199-203`).
- **SA-AN-RISK**: `ANALYSIS.md` critical failure path for `failOn`/`onlyAllow` process exits (`194-197`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enforce compliance gates in CI (Priority: P1)

As a compliance owner, I can define disallowed licenses (`failOn`) or an allowlist (`onlyAllow`) so that a dependency policy violation fails the process immediately.

**Why this priority**: This is the highest-risk compliance control and is explicitly tied to CI gate use cases and hard exit behavior. [SA-FI-F005][SA-OV-R006][SA-AN-RISK]

**Independent Test**: Can be fully tested by running a scan over a known package set with one violating license and verifying non-zero exit with violation text. [SA-AN-RISK]

**Acceptance Scenarios**:

1. **Given** a package set containing license `MIT` and policy `failOn=MIT`, **When** filtering/policy checks run, **Then** processing MUST emit a policy error message and terminate with exit code `1`. [SA-OV-R006][SA-AN-F005]
2. **Given** a package set containing a license not present in `onlyAllow`, **When** policy checks run with `onlyAllow`, **Then** processing MUST emit a policy error message and terminate with exit code `1`. [SA-OV-R006][SA-AN-F005]

---

### User Story 2 - Restrict scan results to relevant dependency subset (Priority: P2)

As an auditor, I can apply `exclude`, package include/exclude lists, and private-package exclusion so that final results contain only relevant packages.

**Why this priority**: This directly controls report precision and legal review scope before any output mode is applied. [SA-OV-R005][SA-AN-FLOW]

**Independent Test**: Can be fully tested by running filtering with a fixture that includes multiple license types, package keys, and private packages, then asserting final keys. [SA-AN-F005][SA-AN-EDGE]

**Acceptance Scenarios**:

1. **Given** an `exclude` list containing one matched license expression, **When** filtering runs, **Then** packages with satisfying licenses MUST be removed from the result map. [SA-AN-F005]
2. **Given** `packages` include-list with semicolon-delimited `name@version` keys, **When** filtering runs, **Then** only listed package keys MUST remain. [SA-AN-F005]
3. **Given** `excludePrivatePackages=true` and at least one private package, **When** filtering runs, **Then** private package entries MUST be removed from the result map. [SA-AN-F005][SA-AN-EDGE]

---

### User Story 3 - Focus unknown-license investigation (Priority: P3)

As a reviewer, I can use `unknown` and `onlyunknown` to isolate guessed or unresolved license findings.

**Why this priority**: Useful for manual review workflows but secondary to hard compliance gating and scope restriction. [SA-OV-R005][SA-AN-F005]

**Independent Test**: Can be fully tested with fixtures containing guessed (`*`) and non-guessed license values and asserting transformed/retained entries. [SA-AN-F005]

**Acceptance Scenarios**:

1. **Given** a package license value ending in `*` and option `unknown=true`, **When** filtering runs, **Then** that package license value MUST be rewritten to `UNKNOWN`. [SA-AN-F005]
2. **Given** mixed known and unknown/guessed licenses with `onlyunknown=true`, **When** filtering runs, **Then** only packages with `UNKNOWN` or `*`-suffixed licenses MUST remain. [SA-AN-F005]

### Edge Cases

- How escaped commas in `exclude` entries are handled (e.g., license strings containing commas) MUST remain supported. [SA-AN-EDGE]
- `exclude=BSD` MUST continue to match BSD variant licenses via current expansion behavior. [SA-AN-EDGE]
- Private packages may be labeled `UNLICENSED` upstream; private exclusion logic MUST still remove them when enabled. [SA-AN-EDGE]
- Policy check order is significant: filters complete before policy exits are evaluated. [SA-AN-FLOW]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST apply `unknown` transformation by converting guessed (`*`-suffixed) license values to `UNKNOWN` before subsequent restriction checks that consume license strings. **Source**: [SA-AN-F005][SA-OV-R005]
- **FR-002**: System MUST support `onlyunknown` mode that retains only entries with `UNKNOWN` or guessed (`*`-suffixed) license values and removes all others. **Source**: [SA-AN-F005][SA-OV-R005]
- **FR-003**: System MUST support `exclude` license filtering with compatibility for escaped commas and SPDX-satisfaction matching behavior used by the current slice. **Source**: [SA-AN-F005][SA-AN-EDGE]
- **FR-004**: System MUST support package restriction by semicolon-delimited package keys for both include (`packages`) and exclude (`excludePackages`) modes, using canonical `name@version` keys. **Source**: [SA-AN-F005][SA-OV-R005]
- **FR-005**: System MUST support `excludePrivatePackages` by removing entries marked private from the filtered result set. **Source**: [SA-AN-F005][SA-AN-EDGE]
- **FR-006**: System MUST enforce `failOn` as a hard policy gate: on first violating package it MUST emit a violation message and exit with status code `1`. **Source**: [SA-OV-R006][SA-AN-RISK]
- **FR-007**: System MUST enforce `onlyAllow` as a hard policy gate: if a package license is outside the allowed set under current matching semantics, it MUST emit a violation message and exit with status code `1`. **Source**: [SA-FI-F005][SA-OV-R006][SA-AN-RISK]
- **FR-008**: System MUST apply filtering stages before policy exit evaluation in the same scan lifecycle. **Source**: [SA-AN-FLOW]
- **FR-009**: CLI compatibility constraint: when both `failOn` and `onlyAllow` are provided together, invocation MUST fail with exit code `1` as an invalid policy combination. **Source**: [SA-OV-R006][SA-AN-F005]

### Inputs and Outputs (Slice I/O Contract)

- **Inputs**
  - Filter toggles: `unknown` (boolean), `onlyunknown` (boolean), `excludePrivatePackages` (boolean). [SA-AN-F005]
  - Filter lists: `exclude` (license-expression list), `packages`/`excludePackages` (semicolon-delimited package-key list). [SA-AN-F005]
  - Policy lists: `failOn`, `onlyAllow` (license list semantics as currently implemented). [SA-AN-F005]
  - Package map keyed by `name@version` with at least: `licenses` (string), `private` (boolean marker when applicable). [SA-AN-F005]
- **Outputs**
  - Filtered/restricted package map keyed by `name@version`. [SA-AN-F005]
  - On policy violation: error emission + process termination with exit code `1`. [SA-AN-RISK]

### Key Entities *(include if feature involves data)*

- **ModuleEntry**: One dependency record keyed as `name@version`, containing `licenses` and package attributes used for filtering/policy checks.
- **FilterPolicyConfig**: Runtime filter/policy options (`unknown`, `onlyunknown`, `exclude`, `packages`, `excludePackages`, `excludePrivatePackages`, `failOn`, `onlyAllow`).
- **RestrictedMap**: Post-filter result map consumed by output serialization or termination flow.

## Compatibility Notes

- Preserve semicolon-based package list parsing for include/exclude package restrictions. [SA-AN-F005]
- Preserve escaped-comma behavior for `exclude` parsing. [SA-AN-EDGE]
- Preserve BSD special-case exclusion compatibility. [SA-AN-EDGE]
- Preserve hard-exit (`code=1`) policy failure semantics for CI compatibility. [SA-OV-R006][SA-AN-RISK]
- Preserve current ordering: filtering stages precede policy exits. [SA-AN-FLOW]

## Constraints

- **C-001 (Scope constraint):** This slice MUST NOT redefine license discovery/normalization algorithms; it consumes already-computed license strings. [SA-AN-FLOW]
- **C-002 (Behavioral constraint):** This slice MUST operate on canonical `name@version` map keys and not alternate package identifiers. [SA-AN-F005]
- **C-003 (Execution constraint):** Policy checks MUST remain fail-fast within a run (process exits immediately on detected violation). [SA-AN-RISK]
- **C-004 (Input-model constraint):** This slice assumes upstream option parsing has already normalized raw CLI strings into this slice’s runtime options. [SA-OV-R005][SA-OV-R006]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For fixtures with guessed licenses, enabling `unknown` rewrites 100% of `*`-suffixed license values to `UNKNOWN` in the post-filter map. [SA-AN-F005]
- **SC-002**: For fixtures with mixed licenses, enabling `onlyunknown` removes all entries except unknown/guessed entries with zero false inclusions. [SA-AN-F005]
- **SC-003**: For fixtures with one known `failOn` or `onlyAllow` violation, process exits with code `1` and emits at least one violation message in 100% of runs. [SA-AN-RISK]
- **SC-004**: For package/private filtering fixtures, resulting key set exactly matches include/exclude/private expectations derived from policy inputs in 100% of runs. [SA-AN-F005][SA-AN-EDGE]

## Open Decisions (slice-limited)

1. **OD-001 — `onlyAllow` matching strictness**: Current behavior is documented as substring-based in Analyzer open questions; decide whether strict compatibility requires preserving that behavior exactly or tightening to exact/SPDX-aware matching in rewrite. **Status**: NEEDS CLARIFICATION. **Source**: [SA-FI-F005]

## Assumptions

- Dependency flattening has already completed successfully and produced a valid package map before this slice executes. [SA-AN-FLOW]
- License values are available as normalized strings, including guessed-marker (`*`) conventions used by this slice. [SA-AN-F005]
- Programmatic and CLI entry paths both invoke this filtering/policy stage with equivalent option semantics. [SA-AN-F005]

## Relevant Files for This Requirement

- `/Users/igorbesel/Projekte/theseus/target/specs/ANALYSIS.md`
- `/Users/igorbesel/Projekte/theseus/target/specs/OVERVIEW.md`
- `/Users/igorbesel/Projekte/theseus/target/specs/FUNCTIONALITY-INDEX.md`
