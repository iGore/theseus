# Implementation Plan: CLI options model and normalization (F-001)

**Date**: 2026-05-01 | **Spec**: `target/specs/cli-options-model/SPEC.md`  
**Input**: Feature specification from `target/specs/cli-options-model/SPEC.md`

**Note**: This plan is scoped strictly to F-001 and derived from `SPEC.md` + shared `target/specs/ARCHITECTURE`.

## Summary

Reimplement the F-001 CLI options slice in Go by building a Cobra-based option parsing and normalization layer that produces deterministic `NormalizedOptions`, applies early guardrails, and supports help/version meta paths with compatibility-controlled exit behavior. The implementation order follows shared architecture sequencing (CLI slice after core model scaffolding) and keeps all behavior bounded to FR-001..FR-008.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` (F-001 only)
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Use Cobra command lifecycle (`RunE` + validation hook path) for deterministic parse/guard error propagation, aligned with architecture package decision and Context7 guidance.
- [x] Keep CLI/business separation: `internal/cli` handles parsing + guardrails only; normalized options are passed to core-facing contracts.
- [x] Default `start` via `os.Getwd()` at normalization stage when caller omits path (FR-002).
- [x] Represent direct traversal intent as depth `0`, non-direct as full traversal (`math.MaxInt` sentinel documented as Go replacement for JS `Infinity`) (FR-003).
- [x] Keep `--version` exit semantics as a gated compatibility decision; do not lock implementation behavior until product decision is made (FR-008).

## Open Questions

- [ ] **NEEDS CLARIFICATION**: Should `--version` exit with `1` (strict legacy compatibility) or `0` (conventional CLI success)? This blocks final acceptance for FR-008 and SC-004.

## Technical Context

**Language/Version**: Go (project default; target current stable toolchain used by repo CI)  
**Primary Dependencies**: `github.com/spf13/cobra`; Go stdlib (`os`, `context`, `errors`, `fmt`, `strings`)  
**Storage**: N/A for F-001 (in-memory option parse/normalize; reads CWD only)  
**Testing**: `go test` table-driven unit tests (`internal/cli`), plus focused command-path tests for help/version and guard exits  
**Target Platform**: Cross-platform CLI (Linux/macOS/Windows)  
**Project Type**: Go CLI slice in a pipeline-oriented binary + reusable package architecture  
**Performance Goals**: Negligible parse latency; deterministic completion for local CLI invocation  
**Constraints**: Preserve spec compatibility semantics for flags/guards; no code beyond F-001 boundaries  
**Scale/Scope**: Single slice covering parse/normalize/guard/meta-command paths only

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/Users/igorbesel/Projekte/theseus/memory/constitution.md`: not present (no additional constitution gate file detected).
- Workflow gates from `AGENTS.md`: satisfied for Planner stage (spec-first flow preserved, Go target preserved, artifact placed in use-case folder).
- `.specify/extensions.yml`: not present, therefore no `hooks.before_plan`/`hooks.after_plan` entries to execute or report.
- Repository setup/agent-context scripts: none detected (`setup.sh`, `agent-context.sh`, `update-agent-context.sh` absent), so no script output to capture.

## Project Structure

### Documentation (this feature)

```text
target/specs/cli-options-model/
├── SPEC.md
├── PLAN.md
└── TASKS.md             # Produced later by Task Decomposer
```

### Source Code (repository root)

```text
cmd/license-checker/
└── main.go                      # Execute root command; map returned errors to process exit strategy

internal/cli/
├── command.go                   # Cobra root command, help/version paths, guardrail invocation order
└── options.go                   # RawArgs -> NormalizedOptions defaults + normalization

internal/core/
└── model.go                     # Shared options/entity contracts used by CLI slice

pkg/licensechecker/
└── api.go                       # Programmatic entrypoint reusing normalization semantics
```

**Structure Decision**: Use the architecture’s Go service/CLI layout and implement only files needed for F-001 behavior. No graph/filter/render slice files are in scope.

## Work Packages (F-001 only)

### WP-1: F-001 contract alignment and option model boundary

- [x] Map `RawArgs`, `NormalizedOptions`, and `GuardrailResult` from spec entities into Go type contracts.
- [x] Define which fields are required for FR-001..FR-008 and mark non-F-001 fields as out-of-scope placeholders (no behavior).
- [ ] Traceability: FR-001, Entities section.

### WP-2: Cobra command skeleton and deterministic execution path

- [x] Establish root command execution path using `RunE`/validation hooks so parse and guard failures return errors (non-zero exit mapping in caller).
- [x] Ensure help path prints usage and exits success semantics per FR-007/SC-004.
- [x] Gate version path exit code behind compatibility decision toggle/documentation marker until open question is resolved.
- [ ] Traceability: FR-005, FR-007, FR-008; US2/US3.

### WP-3: Normalization defaults and semantic transforms

- [x] Implement default `start` behavior from current working directory when absent (FR-002).
- [x] Implement direct traversal normalization (`direct=true -> depth=0`; `direct=false -> full traversal sentinel`) (FR-003).
- [x] Implement structured output color default override (`json|csv|markdown` disables color) (FR-004).
- [ ] Traceability: US1 scenarios 1-3; SC-001, SC-003.

### WP-4: Early guardrails and policy input warnings

- [x] Reject simultaneous `failOn` + `onlyAllow` before scan stage and return non-zero outcome contract (FR-005).
- [x] Emit guidance warning when commas are used in `failOn`/`onlyAllow` values, instructing semicolon delimiter usage (FR-006).
- [x] Keep behavior in CLI guard layer only; no downstream policy engine implementation in this slice.
- [ ] Traceability: US2 scenarios 1-2; SC-002.

### WP-5: Test matrix and compatibility evidence for F-001

- [x] Build table-driven tests covering parse defaults, normalization branches, structured-output color behavior, and no-args behavior.
- [x] Add command-path tests for help/version and guardrail outcomes, including explicit pending assertion branch for unresolved version exit semantics.
- [x] Produce traceability table linking tests to FR-001..FR-008 and SC-001..SC-004.
- [ ] Traceability: all success criteria.

## Dependencies and Execution Order

1. WP-1 contract alignment (foundation for all F-001 behavior).
2. WP-2 command skeleton (execution lifecycle and error path).
3. WP-3 normalization behavior (core P1 value).
4. WP-4 guardrails/warnings (pre-scan safety behavior).
5. WP-5 tests and compatibility evidence (validation closure).

Architecture alignment: this preserves the shared implementation order entry **“CLI slice (F-001)”** after base model scaffolding and before F-002+ slices.

## Validation Approach

- [x] **FR-level validation**: each FR has at least one direct unit/command-path test.
- [x] **Scenario validation**: execute US1/US2/US3 acceptance scenarios as explicit named test cases.
- [x] **Determinism check**: repeated help/version/guard invocations yield stable output + exit behavior.
- [x] **Compatibility check**: for FR-008, run both expected-exit variants in tests with one marked pending behind `NEEDS CLARIFICATION` decision.

## Context7 Notes (planning evidence)

- [x] Context7 `/spf13/cobra`: confirmed `RunE`-based error propagation and command execution lifecycle as the recommended pattern for deterministic CLI failure handling.
- [x] Context7 Go stdlib docs: confirmed `os.Getwd()` semantics for deriving default working directory and informed normalization fallback behavior.

## Complexity Tracking

No constitution violations identified; section intentionally empty.
