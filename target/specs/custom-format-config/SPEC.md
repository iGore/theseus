# Feature Specification: Custom Format and JSON Config Loading

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Gate retry: required artifact missing. Input artifacts: `/Users/igorbesel/Projekte/theseus/target/specs/ANALYSIS.md`, `OVERVIEW.md`, `FUNCTIONALITY-INDEX.md`. Focus slice only: F-007 (slug `custom-format-config`). Required output: MUST create `/Users/igorbesel/Projekte/theseus/target/specs/custom-format-config/SPEC.md`. If folder does not exist, create it. Do not implement code. Keep scope strictly F-007."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add custom fields via inline format object (Priority: P1)

As a user running a dependency license scan, I can provide an inline `customFormat` object so each module result includes additional named fields, using module values when present and configured defaults when missing.

**Why this priority**: This is the core customization behavior of F-007 and directly affects the resulting schema users consume downstream. [Source: F-007, R-009, UC-06]

**Independent Test**: Can be fully tested by running a scan with inline `customFormat` keys (including one nonexistent source key) and verifying returned module entries contain those keys with expected value fallback behavior. [Source: F-007; tests/test.js:431-447]

**Acceptance Scenarios**:

1. **Given** a scan run with `customFormat` containing keys that exist in package metadata (for example `name`, `description`), **When** module output is generated, **Then** those fields appear in each module entry with package-provided string values where available. [Source: F-007; lib/index.js:95-106]
2. **Given** a scan run with `customFormat` containing a key not present in package metadata, **When** module output is generated, **Then** the configured default value is emitted for that field. [Source: F-007; lib/index.js:102-104; tests/test.js:437-445]

---

### User Story 2 - Load custom format from JSON file path (Priority: P2)

As a user, I can provide `customPath` to load `customFormat` from a JSON file instead of inline options.

**Why this priority**: This is the supported configuration entrypoint for repeatable CLI usage and is explicitly listed for F-007. [Source: F-007, R-009, UC-06]

**Independent Test**: Can be fully tested by running with a valid JSON config file and confirming all declared config keys are present in scan output fields. [Source: F-007; tests/test.js:450-475]

**Acceptance Scenarios**:

1. **Given** a valid `customPath` JSON file, **When** scan initialization runs, **Then** JSON content is parsed and used as `customFormat` for output shaping. [Source: F-007; lib/index.js:264-266, 593-604]
2. **Given** a valid `customPath` JSON file with multiple keys, **When** scan output is produced, **Then** each module entry includes each configured key. [Source: F-007; tests/test.js:468-473]

---

### User Story 3 - Preserve error and opt-out semantics in config handling (Priority: P3)

As a user and integrator, I get predictable behavior when custom config input is invalid, and I can explicitly disable selected output properties.

**Why this priority**: This preserves compatibility-critical edge behavior (error object return path and property exclusion via `false`) required for faithful migration of F-007. [Source: F-007; Open Question on invalid customPath]

**Independent Test**: Can be fully tested by calling parse behavior with malformed/missing inputs and by setting selected properties to `false` in custom format, then verifying error-return and field omission behavior. [Source: F-007; lib/index.js:55-58, 593-608; tests/test.js:701-716, 535-548]

**Acceptance Scenarios**:

1. **Given** `parseJson` receives non-string, missing-file, or malformed-JSON input, **When** parsing executes, **Then** the function returns an `Error` object rather than throwing to caller. [Source: F-007; lib/index.js:593-608; tests/test.js:701-716]
2. **Given** `customFormat` sets a known property to `false`, **When** module output is assembled, **Then** that property is excluded from emitted module fields. [Source: F-007; lib/index.js:55-58]

### Edge Cases

- Invalid `customPath` handling in `init`: current flow assigns parser result directly to `customFormat`, including `Error` objects; whether migration should hard-fail is unresolved. [Source: F-007 Open Question; lib/index.js:264-266, 593-608]
- `parseJson(null)` and non-string path input return `Error('did not specify a path')`. [Source: F-007; lib/index.js:594-596; tests/test.js:713-716]
- Malformed JSON file returns `Error` object instead of partial config. [Source: F-007; lib/index.js:601-606; tests/test.js:701-705]
- `licenseText` formatting differs for CSV vs non-CSV output (newline flattening and quote replacement in CSV mode). [Source: F-007; lib/index.js:181-190]
- `copyright` extraction includes copyright blocks and excludes "copyright notice" and "copyright and related rights" blocks. [Source: F-007; lib/index.js:193-205; tests/test.js:535-548]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a `customFormat` object input and evaluate configured keys during module-info assembly. [Source: F-007; lib/index.js:95-106]
- **FR-002**: For each configured custom key, system MUST use module metadata string value when present; otherwise it MUST apply configured default value. [Source: F-007; lib/index.js:97-104; tests/test.js:441-445]
- **FR-003**: System MUST treat `customFormat[property] === false` as an explicit exclusion flag and MUST NOT emit that property in output. [Source: F-007; lib/index.js:55-58]
- **FR-004**: System MUST support `customPath` input and MUST load it by invoking JSON parsing before dependency traversal output filtering completes. [Source: F-007; lib/index.js:264-266]
- **FR-005**: `parseJson` MUST return parsed JSON object for valid file path and valid JSON content. [Source: F-007; lib/index.js:601-604; tests/test.js:692-699]
- **FR-006**: `parseJson` MUST return an `Error` object for non-string path input. [Source: F-007; lib/index.js:594-596; tests/test.js:713-716]
- **FR-007**: `parseJson` MUST return an `Error` object when file cannot be read or content cannot be parsed as JSON. [Source: F-007; lib/index.js:602-606; tests/test.js:701-711]
- **FR-008**: When custom format includes `licenseText`, system MUST populate it from detected license file content and MUST apply CSV-safe normalization in CSV mode. [Source: F-007; lib/index.js:181-190]
- **FR-009**: When custom format includes `copyright`, system MUST derive it from selected license text blocks using existing inclusion/exclusion pattern rules. [Source: F-007; lib/index.js:193-205; tests/test.js:535-548]
- **FR-010**: Rewrite decision on invalid `customPath` behavior MUST be explicitly resolved because current behavior allows `Error` value to propagate as `customFormat`. [NEEDS CLARIFICATION] [Source: F-007 Open Question; OVERVIEW Open Question #3]

### Key Entities *(include if feature involves data)*

- **CustomFormatConfig**: User-provided key/value map (inline or from JSON file) that defines additional output fields, default values, and opt-out flags (`false`). [Source: F-007]
- **ParsedCustomConfigResult**: Result of `parseJson`, typed as either JSON object or `Error` object. [Source: F-007; lib/index.js:593-608]
- **ModuleCustomFields**: Per-module fields added or suppressed based on `CustomFormatConfig` during flatten output assembly. [Source: F-007; lib/index.js:95-106]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: With a valid inline `customFormat`, 100% of emitted module entries contain all configured keys, each populated by module value or configured default according to FR-002. [Source: F-007]
- **SC-002**: With a valid `customPath` JSON file, scan output includes every configured custom key across all module entries for the same run. [Source: F-007]
- **SC-003**: For malformed, missing, and null-path JSON parse inputs, parser behavior is deterministic: returns `Error` object in all tested cases (no uncaught exception). [Source: F-007]
- **SC-004**: For configuration entries set to `false`, corresponding properties are absent from emitted module fields in all tested modules. [Source: F-007]

## Assumptions

- Scope is strictly limited to F-007 behavior (custom format/config ingestion and related field shaping), excluding unrelated filtering/policy/output mode logic except where directly needed for this slice. [Source: F-007 scope assignment]
- `customPath` points to local filesystem paths accessible at runtime with read permissions. [Source: F-007; lib/index.js:602]
- JSON config keys are treated as flat top-level fields for output shaping; nested object semantics are not guaranteed by current behavior. [Source: F-007; lib/index.js:98-99]
- Existing behavior of returning (not throwing) `Error` from `parseJson` is compatibility-significant unless explicitly changed by later decision. [Source: F-007; lib/index.js:593-608; OVERVIEW Open Question #3]
