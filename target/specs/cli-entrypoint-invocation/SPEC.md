# Feature Specification: CLI entry point and invocation model

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned slice F-001 CLI entry point and invocation model"

## Scope

This specification covers only F-001: the top-level command-line executable surface for `license-checker`, its single-command invocation model, pre-scan help/version/license-policy preflight behavior, and the observable stderr/exit-code contract owned by that entry point. It is source-grounded in the Analyzer artifacts for pinned source commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`, package version `25.0.1`. [SA-001]

In scope:

- The executable name `license-checker` and signature `license-checker [flags]` with no subcommands. [SA-002]
- Top-level preflight order for `--help`, `--version`, `--failOn`/`--onlyAllow`, and comma warnings before the scan pipeline runs. [SA-003] [SA-004] [SA-005]
- Stderr/stdout routing and exit codes for those preflight outcomes. [SA-006]
- Handoff from the CLI entry point to the downstream scan/render pipeline after preflight succeeds. [SA-007]

Out of scope for this slice:

- Full flag type coercion/default semantics and custom format loading, owned by F-002.
- Dependency graph scanning and package record construction, owned by F-003/F-004.
- Filtering semantics beyond the entry-point conflict/warning checks for `--failOn` and `--onlyAllow`, owned by F-005.
- Renderer byte contracts after `checker.init` succeeds, owned by F-006.
- General runtime scan errors and debug logging beyond entry-point stderr routing, owned by F-007.

### Source Evidence

- **SA-001**: `ANALYSIS.md` Scope lines 3-9; `USE-CASES.md` F-001 lines 5-12. Pinned module, Go target, no live source execution, F-001 ready with no blocking open questions.
- **SA-002**: `ANALYSIS.md` Entry Points lines 11-19; `package.json:91-94`; `bin/license-checker:1-15`. The bin map exposes `license-checker`, and the Analyzer records a single command with no subcommands.
- **SA-003**: `bin/license-checker:17-47`; `ANALYSIS.md` lines 46, 65-66. `--help` writes version/usage to stderr and exits `0`.
- **SA-004**: `bin/license-checker:49-52`; `ANALYSIS.md` lines 47, 65-66. `--version` writes the package version to stderr and exits `1`.
- **SA-005**: `bin/license-checker:54-63`; `ANALYSIS.md` lines 36-37 and 121-127. `--failOn` and `--onlyAllow` conflict exits `1`; comma-containing values warn with exact text.
- **SA-006**: `ANALYSIS.md` Inputs and Outputs lines 61-66; `tests/bin-test.js:7-14`; `tests/failOn-test.js:7-42`. Normal CLI exits `0`; preflight/fail conditions use the recorded exit codes.
- **SA-007**: `ANALYSIS.md` Flow Map lines 254-262 and End-to-End Pipeline lines 264-266. After preflight, the CLI calls the scan pipeline and then renderer/write behavior.
- **SA-008**: `README.md:69-90`; `ANALYSIS.md` Evidence Pointers lines 313-315. README documents the option surface exposed to CLI users.
- **SA-009**: `ANALYSIS.md` Source-Dependency Contracts line 285. `nopt` parses known option types and maps `-v` to `version`, `-h` to `help`; the detailed parser contract is owned by F-002.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Invoke the license-checker executable (Priority: P1)

As a CLI user, I can run `license-checker` as a single executable without choosing a subcommand, and the entry point starts the existing scan/render pipeline when no preflight condition stops it.

**Why this priority**: This is the minimum viable entry surface for every other CLI feature; all scan, filtering, and rendering slices depend on reaching the pipeline from this command.

**Independent Test**: Can be tested with a stubbed scan pipeline by invoking `license-checker` with no preflight flags and verifying the stub receives parsed CLI options exactly once and that no subcommand dispatch is required.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** the target executable is installed on PATH, **When** a user runs `license-checker` without `--help`, `--version`, or a `--failOn`/`--onlyAllow` conflict, **Then** the entry point MUST proceed to the scan/render pipeline exactly once under the single-command signature `license-checker [flags]`. [SA-002] [SA-007]
2. **Given** a user supplies an arbitrary subcommand token such as `license-checker scan`, **When** the command is parsed, **Then** the entry point MUST NOT dispatch to a named subcommand because the source CLI has no subcommand layer; any treatment of that token belongs to the F-002 parser compatibility contract. [SA-002] [SA-009]

---

### User Story 2 - Read help from the CLI (Priority: P2)

As a CLI user, I can request help and receive the source-compatible version banner and usage text on stderr with a successful help exit.

**Why this priority**: Help is the primary discoverability path and is explicitly handled before all other preflight checks in the source entry point.

**Independent Test**: Can be fully tested by invoking `license-checker --help` against the target binary with the downstream scanner stubbed to fail if called; the expected value is stderr output plus exit code `0` and no stdout.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** the package version is `25.0.1`, **When** a user runs `license-checker --help`, **Then** stderr MUST begin with `license-checker@25.0.1`, MUST include the usage lines in the source order listed in the Output Contract, stdout MUST be empty, the scan pipeline MUST NOT be called, and the process MUST exit `0`. [SA-001] [SA-003]
2. **Given** a user combines `--help` with `--version` or license-policy flags, **When** the command runs, **Then** the `--help` branch MUST win because it is evaluated first, producing help output and exit `0` without version, conflict, warning, or scan behavior. [SA-003] [SA-004] [SA-005]

---

### User Story 3 - Observe version and entry-point policy diagnostics (Priority: P3)

As a CLI user or automation script, I can rely on source-compatible stderr and exit behavior for version requests and mutually exclusive license-policy flags before any scan starts.

**Why this priority**: Automation depends on stable exit codes and diagnostics; the source intentionally exits `1` for `--version` and for `--failOn`/`--onlyAllow` conflicts.

**Independent Test**: Can be tested with direct invocations of `license-checker --version`, `license-checker --failOn MIT --onlyAllow MIT`, and comma-containing policy inputs while asserting stderr, stdout, exit code, and whether the scan stub was called.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** the package version is `25.0.1`, **When** a user runs `license-checker --version`, **Then** stderr MUST be exactly `25.0.1\n`, stdout MUST be empty, the scan pipeline MUST NOT be called, and the process MUST exit `1`. [SA-004] [SA-006]
2. **Given** both `--failOn` and `--onlyAllow` are present, **When** a user runs the command without `--help` or `--version`, **Then** stderr MUST be exactly `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.\n`, stdout MUST be empty, the scan pipeline MUST NOT be called, and the process MUST exit `1`. [SA-005] [SA-006]
3. **Given** exactly one of `--failOn` or `--onlyAllow` contains a comma, **When** a user runs the command without a higher-priority preflight exit, **Then** stderr MUST include the exact comma-warning string for that flag, and the entry point MUST continue into the scan pipeline rather than exiting solely because of the warning. [SA-005] [SA-007]

### Edge Cases

- `--help` combined with any other flag MUST take precedence over all other preflight branches and exit `0`. [SA-003]
- `--version` combined with `--failOn` and `--onlyAllow`, but without `--help`, MUST take precedence over the conflict branch and exit `1` after printing only the version. [SA-004] [SA-005]
- If both `--failOn` and `--onlyAllow` are present and either contains a comma, the mutual-exclusion error MUST take precedence; the comma warning MUST NOT be emitted. [SA-005]
- If exactly one of `--failOn` or `--onlyAllow` contains a comma, the warning literal MUST preserve source typos: `delimeters` and, in help text, `seperated`. [SA-003] [SA-005]
- Full parser edge behavior for unknown operands, path coercion, defaults, and `--no-` forms is not owned here; F-002 MUST pin it to the `nopt` contract. [SA-009]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The target MUST expose a CLI executable named exactly `license-checker`. [SA-002]
- **FR-002**: The CLI MUST implement a single-command invocation model with no subcommand dispatch layer; its user-facing signature is `license-checker [flags]`. [SA-002]
- **FR-003**: The entry point MUST evaluate preflight branches in this exact order before starting the scan pipeline: `--help`; else `--version`; else simultaneous `--failOn` and `--onlyAllow`; else comma warning for the one present policy flag; else scan pipeline. [SA-003] [SA-004] [SA-005] [SA-007]
- **FR-004**: When `--help` is present, the entry point MUST write the help contract to stderr, MUST write nothing to stdout, MUST NOT start the scan pipeline, and MUST exit with code `0`. [SA-003] [SA-006]
- **FR-005**: When `--version` is present and `--help` is absent, the entry point MUST write exactly `25.0.1\n` to stderr, MUST write nothing to stdout, MUST NOT start the scan pipeline, and MUST exit with code `1`. [SA-001] [SA-004] [SA-006]
- **FR-006**: When `--failOn` and `--onlyAllow` are both present and neither `--help` nor `--version` is present, the entry point MUST write exactly `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.\n` to stderr, MUST write nothing to stdout, MUST NOT start the scan pipeline, and MUST exit with code `1`. [SA-005] [SA-006]
- **FR-007**: When exactly one of `--failOn` or `--onlyAllow` is present and its value contains a comma, the entry point MUST write the exact warning literal from the Output Contract to stderr using that flag name and MUST continue into the scan pipeline unless a downstream slice causes a later exit. [SA-005] [SA-007]
- **FR-008**: When no preflight exit applies, the entry point MUST call the downstream scan pipeline with the parsed arguments and then delegate result rendering or side-effect output to the renderer/write behavior owned by F-006. [SA-007]
- **FR-009**: The entry point MUST preserve the source stderr/stdout routing: preflight diagnostics go to stderr; no preflight branch writes normal scan output to stdout. [SA-003] [SA-004] [SA-005] [SA-006]
- **FR-010**: The help output MUST list the documented option names from the source help block, including `--production`, `--development`, `--unknown`, `--start`, `--onlyunknown`, `--json`, `--csv`, `--csvComponentPrefix`, `--out`, `--customPath`, `--exclude`, `--relativeLicensePath`, `--summary`, `--failOn`, `--onlyAllow`, `--direct`, `--packages`, `--excludePackages`, `--excludePrivatePackages`, `--version`, and `--help`, in the exact order shown in the Output Contract. [SA-003] [SA-008]
- **FR-011**: The implementation MUST treat `-v` as the short alias for `--version` and `-h` as the short alias for `--help` at the entry surface; full `nopt` compatibility remains owned by F-002. [SA-002] [SA-009]

### Output Contract *(mandatory for any slice that produces observable output)*

- **Golden baseline**: Source-derived stderr baselines come from `bin/license-checker:17-63`, transcribed through `ANALYSIS.md` Golden Output Snippets lines 121-127 and Entry Point/Inputs-and-Outputs lines 46-47, 65-66. No live CLI execution was performed; these contracts are source-statement-derived, not live-output captures. [SA-001] [SA-003] [SA-004] [SA-005]
- **Field order & presence**:
  1. `--help`: stderr version banner first, then usage block; stdout absent; scanner absent; exit `0`. [SA-003]
  2. `--version`: stderr package version only; stdout absent; scanner absent; exit `1`. [SA-004]
  3. `--failOn` + `--onlyAllow`: stderr mutual-exclusion error only; stdout absent; scanner absent; exit `1`. [SA-005]
  4. Comma warning for exactly one policy flag: stderr warning before scan; scan/render outputs after this point are owned by downstream slices. [SA-005] [SA-007]
- **Exact stderr strings and literals**:
  - Help banner: `license-checker@25.0.1\n`. [SA-001] [SA-003]
  - Help usage block body, in exact source order and preserving typos:

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

  - Version stderr: `25.0.1\n`. [SA-004]
  - Mutual exclusion stderr: `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.\n`. [SA-005]
  - `--failOn` comma warning stderr: `Warning: As of v17 the --failOn argument takes semicolons as delimeters instead of commas (some license names can contain commas)\n`. [SA-005]
  - `--onlyAllow` comma warning stderr: `Warning: As of v17 the --onlyAllow argument takes semicolons as delimeters instead of commas (some license names can contain commas)\n`. [SA-005]
- **Sort/aggregation order**: Not applicable for this entry-point slice; package ordering and renderer aggregation are owned by F-003/F-006. [SA-007]
- **Whitespace contract**:
  - All preflight diagnostics are emitted to stderr with newline termination equivalent to the source `console.error`/`console.warn` calls. [SA-003] [SA-004] [SA-005]
  - The help branch emits a blank line before the first usage option and a blank line before `--version`, matching the source `usage` array shape. [SA-003]
  - The source help implementation calls `console.error(usage.join('\n'), '\n')`; no live byte capture exists for runtime console spacing after the final help line, so final trailing-space byte identity is an open parity limitation to be confirmed in the manual differential run. [SA-001] [SA-003]
  - Preflight branches write no stdout bytes. [SA-006]

### Algorithm Fidelity *(mandatory when the slice transforms values)*

This slice does not classify package values or render package records. It does, however, own the ordered preflight decision table:

| Order | Condition | Action | Exit |
|---|---|---|---|
| 1 | `args.help` truthy | Emit help contract to stderr; do not scan | `0` |
| 2 | `args.version` truthy | Emit `25.0.1\n` to stderr; do not scan | `1` |
| 3 | `args.failOn && args.onlyAllow` truthy | Emit mutual-exclusion stderr; do not scan | `1` |
| 4 | Exactly one policy argument truthy and contains `,` | Emit warning for `failOn` or `onlyAllow`; continue | no immediate exit |
| 5 | No preflight exit | Call scan/render pipeline | downstream |

The parser source dependency is `nopt` (`^4.0.1`) for recognizing known options and aliases; full parser compatibility is specified by F-002, but F-001 requires that the entry point observe the resulting `help`, `version`, `failOn`, and `onlyAllow` fields with the order above. [SA-009]

### Key Entities *(include if feature involves data)*

- **CLI executable**: The installed command named `license-checker`, mapped from package metadata to the entry script. [SA-002]
- **Parsed arguments**: The top-level option object consumed by the entry point, including `help`, `version`, `failOn`, and `onlyAllow` for this slice. [SA-003] [SA-004] [SA-005] [SA-009]
- **Preflight result**: An entry-point decision that either exits before scanning or allows the downstream scan pipeline to run. [SA-007]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For source-derived preflight fixtures `--help`, `--version`, and `--failOn MIT --onlyAllow MIT`, the target exit code, stdout emptiness, and stderr contract match this SPEC for every byte except the explicitly noted final help trailing-spacing limitation. [SA-003] [SA-004] [SA-005]
- **SC-002**: `license-checker` with no preflight exit condition invokes the scan pipeline exactly once in an integration test with a stubbed downstream scanner. [SA-007]
- **SC-003**: `--help` and `-h` both take precedence over all other preflight flags and exit `0`; `--version` and `-v` both exit `1` when help is absent. [SA-003] [SA-004] [SA-009]
- **SC-004**: For comma-containing `--failOn` and `--onlyAllow` values supplied separately, the target emits the exact warning literal for the selected flag and does not exit solely because of the warning. [SA-005]
- **SC-005**: The target command exposes no subcommand-specific behavior in tests; all functionality is reachable under `license-checker [flags]`. [SA-002]

## Assumptions

- The target package version for parity in this migration slice is fixed to the pinned source `package.json` version `25.0.1`; future version stamping is outside this F-001 spec. [SA-001]
- No live source CLI output was captured during analysis. The help/version/diagnostic contracts above are derived from committed source statements and Analyzer snippets; any live differential run against an installed source system is a manual out-of-flow maintainer step, not a phase gate. [SA-001]
- The final trailing whitespace of the help usage block may depend on Node `console.error` formatting for the source call with two arguments; preserve the source structure and verify manually if byte-exact live parity is required. [SA-003]
- Full parser behavior, including unknown operands and `nopt` edge cases, is assumed to be specified by F-002; this slice only requires the entry point to honor the parsed fields needed for preflight. [SA-009]
- Normal scan output and non-preflight runtime errors are assumed to be specified by F-006/F-007; this slice only ensures the entry point reaches or bypasses those paths correctly. [SA-007]
