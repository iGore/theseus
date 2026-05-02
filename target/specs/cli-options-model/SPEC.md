# Feature Specification: CLI options model and normalization

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Hard gate retry: artifact still missing. Focus ONLY F-001 (slug `cli-options-model`)."

## Source Traceability

- **SA-001**: `FUNCTIONALITY-INDEX.md:5-10` (F-001 scope, slug, open question).  
- **SA-002**: `ANALYSIS.md:75-83` (F-001 behavior summary).  
- **SA-003**: `ANALYSIS.md:56-59` (functions mapped to F-001: `parse`, `setDefaults`, `raw/clean/has`, CLI guards).  
- **SA-004**: `ANALYSIS.md:158-160` (flow: parse defaults, then CLI validation).  
- **SA-005**: `OVERVIEW.md:46` (R-001 requirement anchor for compatibility).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Parse and normalize options for a scan run (Priority: P1)

As a CLI user, I can provide supported flags and receive a normalized option set with deterministic defaults so the scan can start without extra setup.

**Why this priority**: This is the entry gate for all downstream behaviors; without stable normalization, no other feature slice can execute reliably.

**Independent Test**: Can be fully tested by passing representative arg combinations into parsing and verifying normalized output values (including defaults) without running dependency scanning.

**Acceptance Scenarios**:

1. **Given** no explicit scan path, **When** options are parsed, **Then** `start` defaults to the current working directory. [SA-002]
2. **Given** `--direct` is enabled, **When** options are parsed, **Then** traversal depth intent is represented as direct-only (`0`) instead of full traversal (`Infinity`). [SA-002]
3. **Given** structured output mode (`json`/`csv`/`markdown`) is selected, **When** defaults are applied, **Then** color output is disabled. [SA-002]

---

### User Story 2 - Enforce early CLI guardrails before scanning (Priority: P2)

As a CLI user, I receive immediate feedback for invalid flag combinations and delimiter misuse so incorrect invocations fail fast.

**Why this priority**: Fast failure protects users from long-running scans with invalid policy flags and preserves compatibility with observed CLI behavior.

**Independent Test**: Can be tested by invoking CLI parsing/guard layer only and asserting warning/error/exit behavior for specific flag combinations.

**Acceptance Scenarios**:

1. **Given** both `--failOn` and `--onlyAllow` are set, **When** CLI guardrails run, **Then** execution exits with a non-zero status before scan execution. [SA-003][SA-004]
2. **Given** `--failOn` or `--onlyAllow` contains commas, **When** CLI guardrails run, **Then** user receives warning guidance to use semicolon delimiters. [SA-002]

---

### User Story 3 - Access utility CLI meta commands (Priority: P3)

As a CLI user, I can invoke help/version command paths with deterministic process behavior independent of scanning.

**Why this priority**: Lower business impact than scan-path parsing, but needed for compatibility and predictable automation behavior.

**Independent Test**: Can be tested by invoking `--help` and `--version` in isolation and validating output/exit semantics.

**Acceptance Scenarios**:

1. **Given** `--help` is passed, **When** CLI command runs, **Then** usage text is printed and process exits `0`. [SA-003]
2. **Given** `--version` is passed, **When** CLI command runs, **Then** version is printed and process exits according to legacy behavior (currently observed as `1`, pending confirmation). [SA-001]

### Edge Cases

- What happens when no args are passed at all? System defaults `start` to CWD and proceeds with normalized defaults. [SA-002]
- How does system handle policy delimiter mistakes? It warns when commas are used for `failOn/onlyAllow`, expecting semicolons. [SA-002]
- How does system handle incompatible policy flags? `--failOn` + `--onlyAllow` is rejected before scanning. [SA-003][SA-004]
- How should legacy `--version` exit code behavior be treated in compatibility mode? **NEEDS CLARIFICATION** (observed `1`). [SA-001]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST parse the established CLI option model for F-001-compatible flags and expose a normalized options object for downstream flow. [SA-001][SA-003][SA-005]
- **FR-002**: System MUST default `start` to the current working directory when not provided by caller input. [SA-002]
- **FR-003**: System MUST normalize direct-traversal intent such that direct mode maps to depth `0` and non-direct mode maps to full traversal intent (`Infinity`). [SA-002]
- **FR-004**: System MUST disable color output by default when structured output flags (`json`, `csv`, `markdown`) are selected. [SA-002]
- **FR-005**: System MUST reject incompatible simultaneous usage of `failOn` and `onlyAllow` before scan execution and return a non-zero process outcome. [SA-003][SA-004]
- **FR-006**: System MUST emit guidance warning when `failOn` or `onlyAllow` uses commas, indicating semicolon-delimited input is expected. [SA-002]
- **FR-007**: System MUST support meta command paths for `help` and `version` with deterministic process exits. [SA-003]
- **FR-008**: System MUST preserve `--version` exit semantics exactly as source-compatible behavior, or explicitly document a migration deviation. [NEEDS CLARIFICATION: preserve legacy exit `1` vs normalize to `0` not yet decided] [SA-001]

### Key Entities *(include if feature involves data)*

- **RawArgs**: User-provided argument vector before parsing; may contain aliases, booleans, strings, and absent values. [SA-003]
- **NormalizedOptions**: Parsed/typed option set after defaults and normalization (`start`, color behavior, direct-depth intent, policy-related flags). [SA-002][SA-003]
- **GuardrailResult**: Early validation outcome that determines whether execution continues, warns, or exits. [SA-003][SA-004]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of P1 acceptance scenarios pass using parser/normalizer-only tests with no dependency graph execution required.
- **SC-002**: In compatibility tests, invoking conflicting `failOn` + `onlyAllow` always produces pre-scan non-zero termination.
- **SC-003**: In structured-output option tests (`json`, `csv`, `markdown`), normalized option state always disables color.
- **SC-004**: Help/version command-path tests produce deterministic output and exit behavior across repeated runs (including explicit handling of unresolved version-exit decision).

## Assumptions

- This specification is strictly bounded to F-001 and does not define dependency loading, filtering, or output rendering internals from other functionality slices. [SA-001]
- Programmatic callers and CLI flow both depend on the same normalized option semantics from the parser layer. [SA-003][SA-004]
- Semicolon-delimited policy token expectations remain part of compatibility behavior for this slice. [SA-002]
- Any change to legacy `--version` exit code is a product decision outside this spec and must be explicitly resolved before implementation lock. [SA-001]

## Open Questions

1. Should `--version` continue exiting with code `1` for strict compatibility, or be changed to `0` for conventional CLI semantics? [SA-001]
