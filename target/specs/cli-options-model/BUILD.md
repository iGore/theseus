# BUILD: CLI options model and normalization (F-001)

## 1. Implementation Scope

- **Files changed (F-001 slice only):**
  - `target/internal/core/model.go`
  - `target/internal/cli/options.go`
  - `target/internal/cli/command.go`
  - `target/internal/cli/options_test.go`
  - `target/internal/cli/command_test.go`
  - `target/go.mod` (+ `target/go.sum` via `go mod tidy`)
  - `target/specs/cli-options-model/PLAN.md`
  - `target/specs/cli-options-model/TASKS.md`
  - `target/specs/cli-options-model/BUILD.md`
- **Spec requirements covered:** FR-001 through FR-008 (with FR-008 still explicitly gated by unresolved product decision for version exit code compatibility mode).

## 2. Build Steps

1. Read `SPEC.md`, `PLAN.md`, `TASKS.md`, and shared `ARCHITECTURE` for F-001 scope boundaries.
2. Implemented F-001 entities and contracts in `internal/core/model.go` (`RawArgs`, `Options`, `GuardrailResult`, `VersionExitMode`).
3. Implemented option normalization in `internal/cli/options.go`:
   - default `start` from `getwd` when empty,
   - direct traversal depth mapping,
   - structured output forcing `color=false`,
   - incompatible `failOn` + `onlyAllow` validation helper.
4. Implemented Cobra command execution path in `internal/cli/command.go` using `RunE` for deterministic error propagation and explicit version compatibility mode handling.
5. Added table-driven normalization tests in `internal/cli/options_test.go`.
6. Added command-path tests in `internal/cli/command_test.go` for guardrail errors, delimiter warnings, help determinism, and version dual-mode behavior.
7. Ran dependency resolution and tests; then synchronized `PLAN.md` and `TASKS.md` checklists.

## 3. Test Evidence

- **Tests added/updated**
  - `target/internal/cli/options_test.go`
  - `target/internal/cli/command_test.go`

- **Commands run + observed results**
  - `go mod tidy` (workdir: `target/`) -> downloaded and resolved Cobra + pflag dependencies successfully.
  - `go test ./internal/cli` -> `ok   theseus/target/internal/cli`
  - `go test ./internal/core ./internal/filter ./internal/cli` -> all packages `ok`.

## 4. Residual Risks

- **FR-008 open decision remains**: legacy `--version` exit code compatibility (`1`) vs conventional success (`0`) is unresolved at product level; implementation keeps an explicit mode toggle (`VersionExitMode`) and test coverage for both branches.
- `FullTraversalDepth` sentinel currently uses existing project sentinel semantics in `core` and should remain aligned with future global compatibility decisions.

## 5. Traceability

- **FR-001** -> `core.RawArgs`, `core.Options`, CLI parse/normalize path in `internal/cli/options.go` and `internal/cli/command.go`.
- **FR-002** -> `NormalizeOptions` defaulting `start` via `getwd`; tested in `TestNormalizeOptions/defaults_start_to_cwd_when_empty`.
- **FR-003** -> `direct=true -> depth=0`, otherwise full traversal sentinel; tested in `TestNormalizeOptions/direct_maps_depth_zero`.
- **FR-004** -> structured output disables color; tested in `TestNormalizeOptions/structured_output_disables_color`.
- **FR-005** -> conflicting `failOn` + `onlyAllow` returns non-zero error path; tested in `TestExecute_GuardrailsAndWarnings/conflicting_failOn_and_onlyAllow_returns_non-zero`.
- **FR-006** -> comma delimiter warning messages for `failOn`/`onlyAllow`; tested in `TestExecute_GuardrailsAndWarnings/comma_delimiter_emits_warning`.
- **FR-007** -> deterministic help/version command paths via Cobra `RunE`; tested with repeated executions in `TestExecute_HelpAndVersionDeterministic`.
- **FR-008** -> compatibility gate retained through `VersionExitMode`; tested in `TestExecute_VersionLegacyCompatibilityBranch` and success branch in `TestExecute_HelpAndVersionDeterministic`.

- **Context7 usage evidence**
  - Used Context7 `/spf13/cobra` documentation to confirm `RunE` error propagation behavior and help/version customization patterns before implementing command lifecycle behavior.
