# Feature Specification: Diagnostics, Error Handling, and Exit Codes

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned slice F-007 Diagnostics, error handling, and exit codes"

## Scope

This specification covers only F-007: observable diagnostics, stderr/stdout channel selection for diagnostic paths, debug opt-in behavior, and process exit-code behavior for the `license-checker` CLI and `checker.init` error callback path.

**In scope**:

- Help/version diagnostic channels and exit codes. [SA-001]
- CLI preflight mutual-exclusion and comma-delimiter warning diagnostics. [SA-002]
- `--failOn` and `--onlyAllow` failure diagnostics and exit codes after F-005 filtering/restriction has produced the candidate package set. [SA-003]
- Callback error reporting from `checker.init` to the CLI. [SA-004]
- Missing-license-file warnings emitted by `--files` side effects owned by F-006. [SA-005]
- `DEBUG=license-checker*` namespace behavior documented by the source project. [SA-006]

**Out of scope**:

- Argument parsing/default semantics beyond the flags needed to trigger these diagnostics; covered by F-001/F-002.
- License policy matching/filtering semantics that decide which package/license is restricted; covered by F-005.
- Rendered tree/JSON/CSV/Markdown/file-copy byte contracts except for diagnostics emitted by those paths; covered by F-006.
- Byte-capturing live Node `Error` rendering; the Analyzer explicitly did not execute the source system. [SA-007]

## User Scenarios & Testing *(mandatory)*

### User Story 1 - CLI preflight diagnostics are predictable (Priority: P1)

As a CLI user, I want help/version/conflicting-flag feedback to appear on the same channel and with the same exit status as the source CLI, so scripts can treat preflight outcomes compatibly.

**Why this priority**: These paths occur before dependency scanning and are the smallest independently testable compatibility surface for F-007.

**Independent Test**: Can be fully tested by invoking the Go CLI with `--help`, `--version`, and incompatible policy flags against any working directory and asserting stderr/stdout plus exit code.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** any working directory, **When** the user runs `license-checker --help`, **Then** the process writes the source help/version lines to stderr in the order specified by the Output Contract, writes no diagnostic text to stdout, and exits with code `0`. [SA-001]
2. **Given** any working directory, **When** the user runs `license-checker --version`, **Then** the process writes `25.0.1\n` to stderr and exits with code `1`. [SA-001]
3. **Given** a command line containing both `--failOn` and `--onlyAllow`, **When** the CLI performs preflight validation, **Then** it writes exactly `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` to stderr and exits with code `1` before scanning or rendering. [SA-002]
4. **Given** a command line containing `--failOn MIT,ISC`, **When** the CLI performs preflight validation, **Then** it writes the comma-delimiter warning for `--failOn` to stderr and continues into the normal scan/policy pipeline instead of exiting because of the warning alone. [SA-002]
5. **Given** a command line containing `--onlyAllow MIT,ISC`, **When** the CLI performs preflight validation, **Then** it writes the comma-delimiter warning for `--onlyAllow` to stderr and continues into the normal scan/policy pipeline instead of exiting because of the warning alone. [SA-002]

---

### User Story 2 - License policy failures terminate with source-compatible diagnostics (Priority: P2)

As a compliance user, I want forbidden or non-allowed licenses to stop the command with the original failure messages and exit code, so CI checks that rely on `license-checker` keep their meaning.

**Why this priority**: Policy-failure exit behavior is the core operational purpose of this slice and depends on F-005's candidate package/license decisions.

**Independent Test**: Can be fully tested with fixtures whose post-filter package set contains a license matching `--failOn` and a package whose license is not accepted by `--onlyAllow`.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** the restricted package set contains a package whose license string is exactly `MIT`, **When** the user runs with `--failOn MIT`, **Then** the process writes `Found license defined by the --failOn flag: "MIT". Exiting.` to stderr and exits with code `1`. [SA-003]
2. **Given** the restricted package set contains any package whose license string exactly matches one of the semicolon-separated `--failOn` tokens, **When** the policy check reaches that package, **Then** the first matching license value is interpolated into `Found license defined by the --failOn flag: "<license>". Exiting.` and the process exits with code `1`. [SA-003]
3. **Given** the restricted package set contains package `<item>` with license `<license>` and `<license>` contains none of the semicolon-separated `--onlyAllow` tokens, **When** the user runs with `--onlyAllow`, **Then** the process writes `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` to stderr and exits with code `1`. [SA-003]

---

### User Story 3 - Non-fatal diagnostics preserve channels without inventing output (Priority: P3)

As a user or integrator, I want scan errors, missing copied license files, and debug output to remain source-compatible without the rewrite inventing new structured output, so existing stderr/stdout consumers are not surprised.

**Why this priority**: These paths are less common than preflight/policy failures but are observable and explicitly documented in the source artifacts.

**Independent Test**: Can be tested with fixtures or injected scanner/file-copy conditions that trigger `checker.init` errors, missing `licenseFile` records, and `DEBUG=license-checker*`.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** `checker.init` returns a non-null error to the CLI callback, **When** the CLI handles the callback, **Then** it writes `Found error` to stderr, writes the error object using the target's declared Node-compatible error-formatting decision, and does not introduce a new explicit non-zero exit solely for this callback error. [SA-004] [SA-007]
2. **Given** the user runs with `--files` and a package record lacks a usable `licenseFile`, **When** file-copy output is processed, **Then** the process writes `no license file found for: <moduleName>` to stderr as a warning and continues processing remaining packages. [SA-005]
3. **Given** the environment enables `DEBUG=license-checker*`, **When** scanning starts or an init error is debug-logged, **Then** debug output is enabled for the documented `license-checker:log` and `license-checker:error` namespaces without changing functional output contracts or exit-code decisions. [SA-006]

### Edge Cases

- `--failOn` and `--onlyAllow` are both present and one or both contain commas: the mutual-exclusion error takes precedence and exits `1`; comma warnings are not required after the fatal preflight error because the source exits immediately in the mutual-exclusion branch. [SA-002]
- `--failOn` or `--onlyAllow` contains empty semicolon-separated tokens: empty trimmed tokens are ignored by the source policy list construction; no F-007 diagnostic is emitted for those empty tokens. [SA-003]
- A comma appears in a license name supplied to `--failOn`/`--onlyAllow`: the CLI emits the historical warning preserving the misspelling `delimeters`, but it still passes the raw argument into downstream policy behavior. [SA-002]
- A scan produces no packages: source `init` sets `Error('No packages found in this path..')`, debug-logs it, returns it to the callback, and the CLI emits `Found error`; byte-exact rendering of the subsequent `Error` object is an open decision. [SA-004] [SA-007]
- File-system errors from renderer side effects (`mkdirp`, `fs.writeFileSync`, `fs.readFileSync`) are not captured as handled diagnostics in the Analyzer evidence; do not invent friendly messages for them in this slice. [SA-007]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The CLI MUST emit help output on stderr and exit with code `0` when `--help` is present. It MUST NOT require dependency scanning to satisfy help. [SA-001]
- **FR-002**: The CLI MUST emit version output on stderr and exit with code `1` when `--version` is present. The pinned source package version is `25.0.1`. [SA-001]
- **FR-003**: The CLI MUST reject simultaneous `--failOn` and `--onlyAllow` before scanning, write exactly `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` to stderr, and exit with code `1`. [SA-002]
- **FR-004**: The CLI MUST emit a non-fatal comma-delimiter warning to stderr when the active policy argument value for `--failOn` or `--onlyAllow` contains `,`. The warning MUST preserve the source misspelling `delimeters`. [SA-002]
- **FR-005**: A comma-delimiter warning MUST NOT by itself change the eventual exit code; after warning, execution MUST continue into the normal scan, filter, policy, and output pipeline. [SA-002]
- **FR-006**: During `--failOn` policy checking, the system MUST exit with code `1` immediately when a restricted package's license string exactly equals one of the non-empty trimmed semicolon-separated `--failOn` tokens. [SA-003]
- **FR-007**: On `--failOn` failure, the system MUST write exactly `Found license defined by the --failOn flag: "<license>". Exiting.` to stderr, interpolating the matched license string. [SA-003]
- **FR-008**: During `--onlyAllow` policy checking, the system MUST exit with code `1` when a restricted package's license string contains none of the non-empty trimmed semicolon-separated `--onlyAllow` tokens. [SA-003]
- **FR-009**: On `--onlyAllow` failure, the system MUST write exactly `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` to stderr, interpolating the restricted package key and license string. [SA-003]
- **FR-010**: When `checker.init` reports an error to the CLI callback, the CLI MUST write `Found error` to stderr before writing the error object, and MUST NOT add a new explicit non-zero process exit solely for that callback error. [SA-004] [SA-007]
- **FR-011**: The byte formatting of the second callback-error line produced by source `console.error(err)` MUST be treated as an explicit open decision until a maintainers-approved Node baseline is captured; the Go rewrite MUST NOT invent a structured JSON/YAML/error-code replacement for it. [SA-004] [SA-007]
- **FR-012**: The `--files` side-effect path MUST write exactly `no license file found for: <moduleName>` to stderr as a warning for each package without a usable `licenseFile`, and MUST continue with remaining package records. [SA-005]
- **FR-013**: A normal CLI invocation with no explicit fatal diagnostic path MUST exit with code `0`, as asserted by the source CLI test. [SA-008]
- **FR-014**: `DEBUG=license-checker*` MUST enable diagnostics for the source namespaces `license-checker:log` and `license-checker:error`; this debug setting MUST NOT alter policy decisions or documented exit codes. [SA-006]
- **FR-015**: The implementation MUST preserve source diagnostic spelling, punctuation, capitalization, interpolation quotes, and stderr/stdout channel choices for all exact strings listed in the Output Contract. [SA-001] [SA-002] [SA-003] [SA-004] [SA-005]

### Output Contract *(mandatory for any slice that produces observable output)*

**Golden baseline**: The baseline is source-backed, not live-captured. The Analyzer records committed/source snippets for diagnostic strings and exit behavior; live execution was intentionally not performed. [SA-007]

| Diagnostic path | Trigger | Channel | Exact emitted diagnostic text captured or source-derived | Exit-code contract | Source evidence |
|---|---|---|---|---|---|
| Help | `--help` | stderr | First line is `license-checker@25.0.1`; usage body lines are listed below in exact source order. Final trailing whitespace from Node's two-argument `console.error(usage.join('\n'), '\n')` was not live-captured. | `0` | [SA-001] |
| Version | `--version` | stderr | `25.0.1\n` | `1` | [SA-001] |
| Mutual exclusion | both `--failOn` and `--onlyAllow` | stderr | `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` | `1` | [SA-002] |
| Comma warning for failOn | active `--failOn` value contains `,` | stderr warning | `Warning: As of v17 the --failOn argument takes semicolons as delimeters instead of commas (some license names can contain commas)` | No direct exit; continue | [SA-002] |
| Comma warning for onlyAllow | active `--onlyAllow` value contains `,` | stderr warning | `Warning: As of v17 the --onlyAllow argument takes semicolons as delimeters instead of commas (some license names can contain commas)` | No direct exit; continue | [SA-002] |
| failOn failure | restricted license exactly matches a `--failOn` token | stderr | `Found license defined by the --failOn flag: "<license>". Exiting.` | `1` | [SA-003] |
| onlyAllow failure | restricted license contains none of the `--onlyAllow` tokens | stderr | `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` | `1` | [SA-003] |
| Init callback error prefix | `checker.init` callback receives non-null error | stderr | `Found error` | No new explicit non-zero exit solely for this error | [SA-004] [SA-007] |
| Init callback error object | immediately after `Found error` | stderr | Node `console.error(err)` formatting is not byte-captured; preserve as open decision, not invented output | No new explicit non-zero exit solely for this error | [SA-004] [SA-007] |
| Missing copied license file | `--files` package has no usable `licenseFile` | stderr warning | `no license file found for: <moduleName>` | No direct exit; continue | [SA-005] |
| Debug log | `DEBUG=license-checker*` includes `license-checker:log` | debug output; source binds non-error debug logger to stdout-compatible `console.log` behavior | Debug dependency formatting is not byte-captured; namespace and enablement are required | No exit impact | [SA-006] [SA-007] |
| Debug error | `DEBUG=license-checker*` includes `license-checker:error` and init error occurs | debug output using debug's error namespace behavior | Debug dependency formatting is not byte-captured; namespace and enablement are required | No exit impact | [SA-006] [SA-007] |

Source help body line contract for `--help` after the first `license-checker@25.0.1` line:

```text

   --production only show production dependencies.
   --development only show development dependencies.
   --unknown report guessed licenses as unknown licenses.
   --start [path of the initial json to look for]
   --onlyunknown only list packages with unknown or guessed licenses.
   --json output in json format.
   --csv output in csv format.
   --csvComponentPrefix column prefix for components in csv file
   --out [filepath] write the data to a specific file.
   --customPath to add a custom Format file in JSON
   --exclude [list] exclude modules which licenses are in the comma-separated list from the output
   --relativeLicensePath output the location of the license files as relative paths
   --summary output a summary of the license usage
   --failOn [list] fail (exit with code 1) on the first occurrence of the licenses of the semicolon-separated list
   --onlyAllow [list] fail (exit with code 1) on the first occurrence of the licenses not in the semicolon-seperated list
   --direct look for direct dependencies only
   --packages [list] restrict output to the packages (package@version) in the semicolon-seperated list
   --excludePackages [list] restrict output to the packages (package@version) not in the semicolon-seperated list
   --excludePrivatePackages restrict output to not include any package marked as private

   --version The current version
   --help  The text you are reading right now :)

```

**Field order & presence**:

1. Preflight order is help first, version second, mutual-exclusion third, comma warning fourth, then scan initialization. [SA-001] [SA-002]
2. If help is present, help output is emitted and the process exits `0`; later preflight checks are not reached. [SA-001]
3. If version is present and help is absent, version output is emitted and the process exits `1`; later preflight checks are not reached. [SA-001]
4. If both policy flags are present and help/version are absent, only the mutual-exclusion diagnostic is required before exit `1`. [SA-002]
5. Policy-failure diagnostics occur after F-005 filtering/restriction has produced the restricted package set and before any CLI renderer callback output is emitted for that run. [SA-003]
6. `Found error` precedes the implementation's declared Node-compatible error-object formatting. [SA-004] [SA-007]
7. Missing-license-file warnings are emitted once per affected module during `--files` processing. [SA-005]

**Sentinels & literals**:

- Preserve exact literals `--failOn`, `--onlyAllow`, `Warning:`, `delimeters`, `can not`, `Found error`, `Exiting.`, and `no license file found for:` where they appear above. [SA-002] [SA-003] [SA-004] [SA-005]
- Interpolation placeholders are raw source values: `<license>` is the restricted record's license string, `<item>` is the package key such as `name@version`, and `<moduleName>` is the package key used by the file-copy path. [SA-003] [SA-005]

**Sort/aggregation order**:

- This slice does not introduce new sorted diagnostic collections.
- Policy-failure iteration follows the restricted package key order supplied by F-005/F-006 pipeline output; the source builds this from lexically sorted data before later restrictions. [SA-003]
- Missing-license-file warning order follows the package key order of the JSON object passed into `asFiles`, which is owned by F-006. [SA-005]

**Whitespace contract**:

- Exact diagnostic strings above are line-oriented stderr messages as emitted by source `console.error` or `console.warn`; preserve one line per diagnostic message in baseline tests unless the open Node formatting decision says otherwise for `console.error(err)`. [SA-002] [SA-003] [SA-004] [SA-005]
- Do not add stdout copies of stderr diagnostics. [SA-001] [SA-002] [SA-003] [SA-004] [SA-005]
- Full byte-exact help trailing whitespace after the listed lines and Node `Error` object formatting were not captured by Analyzer live execution and remain open decisions. [SA-007]

### Algorithm Fidelity *(mandatory when the slice transforms values)*

This slice does not own license classification or output rendering transformations, but it owns ordered diagnostic/exit control flow:

1. Parse/default CLI arguments according to F-001/F-002. [SA-001]
2. If `help` is true: emit help/version text to stderr and exit `0`. [SA-001]
3. Else if `version` is true: emit version to stderr and exit `1`. [SA-001]
4. Else if both `failOn` and `onlyAllow` are truthy: emit mutual-exclusion diagnostic to stderr and exit `1`. [SA-002]
5. Else if the active policy argument (`failOn` or `onlyAllow`) contains a comma: emit the corresponding warning to stderr and continue. [SA-002]
6. Initialize scanner/pipeline from F-003 through F-006. [SA-004]
7. Build non-empty trimmed semicolon token lists for `failOn` or `onlyAllow`; ignore empty trimmed tokens. [SA-003]
8. For each restricted package: if `failOn` list is non-empty and `restricted[item].licenses` exactly equals one token, emit the failOn diagnostic and exit `1`. [SA-003]
9. For each restricted package: if `onlyAllow` list is non-empty and no token is contained as a substring of `restricted[item].licenses`, emit the onlyAllow diagnostic and exit `1`. [SA-003]
10. If an init error is present after scan/filter processing, debug-log it, return it to the callback, and have the CLI emit `Found error` plus the error object without introducing a new explicit exit code. [SA-004] [SA-006] [SA-007]
11. During `--files`, for each missing/unusable `licenseFile`, emit the missing-file warning and continue. [SA-005]

### Key Entities *(include if feature involves data)*

- **Diagnostic event**: An observable line-oriented message emitted to stderr or debug output, with a trigger, exact literal/interpolation contract, and exit-code effect.
- **Policy token list**: Non-empty trimmed tokens derived from semicolon-separated `--failOn` or `--onlyAllow`; used only after F-005 restriction/filter decisions.
- **Restricted package item**: Package key and record produced by earlier pipeline slices; F-007 consumes its key and `licenses` value for policy diagnostics.
- **Exit outcome**: Process status code resulting from explicit source exits or normal process completion.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Automated CLI tests assert `stderr`, `stdout`, and process status for `--help`, `--version`, mutual `--failOn`/`--onlyAllow`, comma warning, failOn failure, and onlyAllow failure.
- **SC-002**: For all exact diagnostics listed in the Output Contract, tests compare exact spelling, punctuation, interpolation quotes, and channel selection; misspellings from the source are preserved.
- **SC-003**: A policy-failure fixture verifies exit code `1` for `--failOn MIT` and semicolon multi-token matching, matching source tests for fail behavior.
- **SC-004**: A non-fatal warning fixture verifies comma warnings and missing-license-file warnings do not by themselves force exit code `1`.
- **SC-005**: Differential-parity baseline: for the committed/source-backed diagnostic fixtures in this slice, target stderr/stdout and exit codes are byte-identical where the Analyzer captured exact strings, with declared exceptions only for full help body whitespace, debug dependency formatting, and Node `Error` object formatting gaps.

## Assumptions

- The Go target is a CLI-compatible rewrite of `license-checker` at pinned source version `25.0.1` / commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`. [SA-001]
- F-001/F-002 provide parsed flags and aliases compatible enough for this slice to observe `help`, `version`, `failOn`, `onlyAllow`, and `files`. [SA-001] [SA-002]
- F-005 provides the restricted package set and license strings used by `--failOn` and `--onlyAllow`; this slice does not redefine SPDX/filter matching. [SA-003]
- F-006 owns normal renderers and file-copy outputs; this slice only specifies diagnostics and exit behavior that those paths emit. [SA-005]
- Node `console.error(err)` byte formatting is a known open decision because Analyzer did not run the live source. A maintainers-approved fixture or compatibility decision is required before claiming byte-identical parity for that second error line. [SA-007]
- Debug package timestamp/color/prefix formatting was not byte-captured; this spec requires namespace enablement and no exit-code impact, not an invented byte-exact debug format. [SA-006] [SA-007]

## Source Evidence

- **SA-001**: `ANALYSIS.md:11-49`, `ANALYSIS.md:61-67`; source references `bin/license-checker:17-52`, `tests/bin-test.js:7-17`; pinned version from `ANALYSIS.md:5-8`.
- **SA-002**: `ANALYSIS.md:121-127`; source reference `bin/license-checker:54-63`; comma warning asserted in `tests/failOn-test.js:27-40`.
- **SA-003**: `ANALYSIS.md:36-37`, `ANALYSIS.md:121-127`, `ANALYSIS.md:342-350`; source reference `lib/index.js:437-456`; fail exit asserted in `tests/failOn-test.js:7-42`.
- **SA-004**: `ANALYSIS.md:61-67`, `ANALYSIS.md:254-266`; source references `bin/license-checker:66-72`, `lib/index.js:461-467`.
- **SA-005**: `ANALYSIS.md:230-232`; source reference `lib/index.js:610-627`.
- **SA-006**: `ANALYSIS.md:49`, `ANALYSIS.md:270-291`, `ANALYSIS.md:350`; source references `README.md:158-170`, `lib/index.js:21-25`.
- **SA-007**: `ANALYSIS.md:68-70`, `ANALYSIS.md:304-309`; `USE-CASES.md:85-89`.
- **SA-008**: `ANALYSIS.md:65-66`; source reference `tests/bin-test.js:7-17`.
