# Feature Specification: Output Rendering and Export

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Gate retry for missing artifact. Focus slice: F-006 (`output-rendering-and-export`)."

## Source Evidence Index

- **SE-001**: F-006 scope and evidence summary (`FUNCTIONALITY-INDEX.md:35-40`).
- **SE-002**: Output modes and persistence behavior summary (`ANALYSIS.md:121-129`).
- **SE-003**: CLI mode selection and write flow (`ANALYSIS.md:37-45`, `ANALYSIS.md:122-128`).
- **SE-004**: Serializer functions (tree/csv/markdown/summary) (`ANALYSIS.md:63-68`, `ANALYSIS.md:222-223`).
- **SE-005**: Required migration tasks for F-006 (`OVERVIEW.md:52-54`, `OVERVIEW.md:66-67`).
- **SE-006**: Use-case expectation for machine-readable output + `--out` (`OVERVIEW.md:19-23`).
- **SE-007**: Use-case expectation for `--files` license export (`OVERVIEW.md:34-38`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate the correct output format (Priority: P1)

As a user, I can request one output mode (JSON, CSV, Markdown, Summary, or Tree) and get the correctly rendered result for the scanned dependency set. [SE-001][SE-002][SE-004]

**Why this priority**: Rendering a usable output is the core value of F-006; without this, the scan result cannot be consumed. [SE-001][SE-006]

**Independent Test**: Run scan commands that toggle one mode at a time and validate that each response format matches the selected renderer.

**Acceptance Scenarios**:

1. **Given** scan results are available, **When** user sets `--json`, **Then** output is JSON rendering and not CSV/Markdown/Summary/Tree. [SE-003]
2. **Given** scan results are available, **When** user sets no structured mode flags, **Then** output defaults to tree rendering. [SE-002][SE-004]
3. **Given** scan results are available and multiple mode flags are provided, **When** rendering is selected, **Then** mode priority is `json > csv > markdown > summary > tree`. [SE-002][SE-003]

---

### User Story 2 - Persist rendered output to file (Priority: P2)

As a user, I can write the rendered output to an output file path for CI artifacts and offline review. [SE-002][SE-003][SE-006]

**Why this priority**: File persistence is required for automated pipelines and report archival, but depends on rendering behavior from Story 1. [SE-005][SE-006]

**Independent Test**: Run command with `--out <path>` and verify file creation plus content equality with expected renderer output.

**Acceptance Scenarios**:

1. **Given** rendered output text is available, **When** user sets `--out` to a nested path, **Then** parent directories are created and output is written to that file path. [SE-002][SE-003]
2. **Given** `--out` is not set, **When** rendering completes, **Then** output is emitted to standard output. [SE-003]

---

### User Story 3 - Export detected license files per package (Priority: P3)

As a user, I can request export of detected license files into a target directory for legal evidence collection. [SE-002][SE-007]

**Why this priority**: This is valuable for legal/compliance workflows but secondary to primary rendering and report persistence. [SE-007]

**Independent Test**: Run with `--files <dir>` and verify per-module license artifacts are copied into the output folder structure.

**Acceptance Scenarios**:

1. **Given** modules include detected license file paths, **When** user sets `--files <dir>`, **Then** license files are copied to the specified export directory. [SE-002][SE-007]
2. **Given** `--files` is omitted, **When** output rendering completes, **Then** no license-file export side effect occurs. [SE-002]

### Edge Cases

- If multiple output mode flags are provided simultaneously, renderer selection MUST follow fixed priority ordering (`json > csv > markdown > summary > tree`). [SE-002][SE-003]
- If `--out` and `--files` are both set, system MUST perform both actions in one run (write rendered output and export license files). [SE-002][SE-003]
- Under-documented modes (`--markdown`, `--files`) are implemented in source behavior; documentation visibility for rewrite remains open. [SE-001]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support rendering in five modes: tree, JSON, CSV, Markdown, and Summary. [SE-001][SE-002][SE-004]
- **FR-002**: System MUST select output mode using strict precedence: JSON first, then CSV, then Markdown, then Summary, else Tree. [SE-002][SE-003]
- **FR-003**: System MUST emit rendered output to stdout when `--out` is not provided. [SE-003]
- **FR-004**: System MUST create parent directories for the `--out` path and persist rendered output to that file. [SE-002][SE-003]
- **FR-005**: System MUST support license-file export via `--files <dir>` by copying detected per-module license files into the target directory. [SE-002][SE-007]
- **FR-006**: System MUST preserve optional colorized package-key behavior only for non-structured output contexts compatible with terminal interactivity. [SE-002][SE-003]
- **FR-007**: System MUST keep F-006 scoped to rendering/export behavior and treat custom-format ingestion rules as external dependency on F-007. [SE-001][SE-005]
- **FR-008**: System MUST retain support for Markdown and Files output behaviors, while documentation-status alignment is **NEEDS CLARIFICATION**. [SE-001]

### Key Entities *(include if feature involves data)*

- **Module Result Map**: Flattened package result collection keyed by `name@version`; input payload for all renderers and export actions. [SE-002][SE-004]
- **Rendered Output Artifact**: Final textual representation produced by selected mode (tree/json/csv/markdown/summary). [SE-002][SE-004]
- **License Export Artifact Set**: File collection copied from detected module license-file paths into `--files` target directory. [SE-002][SE-007]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In compatibility tests, 100% of F-006 runs choose the same renderer as source behavior for equivalent flag combinations. [SE-001][SE-003]
- **SC-002**: For `--out` test cases, 100% of runs create the target file (including nested parent dirs) and persist non-empty rendered content. [SE-002][SE-003]
- **SC-003**: For `--files` test cases with available license-file sources, 100% of expected license artifacts are exported to the target directory. [SE-002][SE-007]
- **SC-004**: For non-`--out` runs, 100% of cases emit output to stdout with no missing primary report artifact. [SE-003]

## Assumptions

- Scan/flatten/filter phases already produced a valid module result map before F-006 rendering starts (handled outside this slice). [SE-005]
- Filesystem permissions allow directory creation and file writes for `--out` and `--files` destinations.
- Structured output compatibility checks compare behavior, not byte-for-byte ordering beyond source-defined serializer semantics.
- Scope is restricted to F-006 only; policy filtering, traversal, and custom-format parsing logic are specified in other functionality slices. [SE-001][SE-005]

## Open Questions

1. Should under-documented but implemented flags (`--markdown`, `--files`) remain behavior-only compatibility features or become first-class documented options in downstream artifacts? [SE-001]
