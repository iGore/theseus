# BUILD: Dependency Graph Loading and Traversal Control (F-002)

## 1. Implementation Scope

### Files changed
- `target/internal/core/model.go`
- `target/internal/core/service.go`
- `target/internal/core/service_test.go`
- `target/internal/graph/loader.go`
- `target/internal/graph/npm_ls_loader.go`
- `target/internal/graph/npm_ls_loader_test.go`
- `target/specs/dependency-graph-loading/PLAN.md`
- `target/specs/dependency-graph-loading/TASKS.md`

### Spec requirements covered
- FR-001..FR-008
- SC-001..SC-004

Out-of-scope guard: no F-003 flattening semantics, F-004 license detection, F-005 policy/filtering, or F-006 rendering/export behavior was implemented in this slice.

## 2. Build Steps

1. Added F-002 input and stage contract fields in `core`:
   - `Options.Production` / `Options.Development`
   - `LoaderOptions` and `DependencyLoader` interface.
2. Implemented F-002 orchestration in `core.Service`:
   - `BuildLoaderOptions` for depth/dev/logger mapping.
   - `LoadDependencyGraph` for exact `start` passthrough + loader error propagation.
   - `LoadAndFlatten` boundary method to hand loaded tree to flatten stage (handoff only).
3. Implemented graph adapter contract layer in `internal/graph`:
   - `loader.go` boundary aliasing to core contracts.
   - `npm_ls_loader.go` with `exec.CommandContext` runner, depth/dev argument mapping, JSON parse to `core.DependencyNode`, and context-aware cancellation/error behavior.
4. Implemented tests for F-002 behavior:
   - `core/service_test.go` for start-path passthrough, success handoff, mapping matrix, and failure propagation.
   - `graph/npm_ls_loader_test.go` for loader argument mapping, fixture-style JSON parsing, cancellation propagation, and wrapped command errors.
5. Updated `PLAN.md` WP1..WP4 checkboxes and `TASKS.md` T001..T017 execution checklist.

Context7 usage:
- Queried Context7 `/golang/go` for `os/exec` + `context` cancellation/error handling patterns before finalizing loader adapter error behavior.

## 3. Test Evidence

### Tests added/updated
- Added: `target/internal/core/service_test.go`
- Added: `target/internal/graph/npm_ls_loader_test.go`
- Updated (indirectly exercised): existing package suites under `target/internal/...` and integration tests through full `go test ./...`.

### Commands run
- `gofmt -w internal/core/model.go internal/core/service.go internal/core/service_test.go internal/graph/loader.go internal/graph/npm_ls_loader.go internal/graph/npm_ls_loader_test.go`
- `go test ./...`

### Observed results
- `go test ./...` passed.
- Key package results:
  - `ok theseus/target/internal/core`
  - `ok theseus/target/internal/graph`
  - full suite also passed for `internal/cli`, `internal/filter`, `internal/flatten`, `internal/legacy`, `internal/license`, and integration test packages.

## 4. Residual Risks

- `npm ls` output/flags can vary slightly by npm version; adapter currently targets stable `--json`, `--depth=0`, `--omit=dev` behavior and is covered with runner-based tests, but real-world npm parity may need additional fixture expansion.
- Recursive traversal sentinel is represented as `core.FullTraversalDepth` (`-1`) in service mapping; if upstream normalization contract changes, this mapping should be revalidated.

## 5. Traceability

| Spec ID | Implementation evidence |
|---|---|
| FR-001 / SC-001 | `Service.LoadDependencyGraph` forwards exact `opts.Start`; verified in `TestServiceLoadDependencyGraph_PassesExactStartPath`; adapter uses provided `start` as command dir in `NpmLSLoader.Load`. |
| FR-002 | `Service.BuildLoaderOptions` constructs loader options from normalized inputs. |
| FR-003 | `BuildLoaderOptions` maps direct mode (`Depth==0`) to `depth=0`; tested in `TestServiceBuildLoaderOptions_Mapping`. |
| FR-004 | `BuildLoaderOptions` maps non-direct to recursive sentinel (`FullTraversalDepth`); tested in `TestServiceBuildLoaderOptions_Mapping`. |
| FR-005 | `BuildLoaderOptions` defaults `Dev=true`; tested in `TestServiceBuildLoaderOptions_Mapping`. |
| FR-006 / SC-003 | `BuildLoaderOptions` sets `Dev=false` when `Production || Development`; tested in `TestServiceBuildLoaderOptions_Mapping`. |
| FR-007 / SC-004 | `Service.LoadDependencyGraph` propagates loader errors; `NpmLSLoader.Load` propagates context cancellation and wraps command failures; tested in `TestServiceLoadDependencyGraph_PropagatesLoaderError` and loader error-path tests. |
| FR-008 | `Service.LoadAndFlatten` forwards loaded dependency tree directly to flatten boundary; tested in `TestServiceLoadAndFlatten_ForwardsLoadedTree`. |
