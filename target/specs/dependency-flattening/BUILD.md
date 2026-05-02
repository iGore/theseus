# BUILD.md — F-003 Dependency Flattening

## 1. Implementation Scope

- **Files changed (code):**
  - `target/internal/core/model.go`
  - `target/internal/core/service.go`
  - `target/internal/flatten/flatten.go`
  - `target/internal/flatten/flatten_test.go`
  - `target/test/integration/flatten_contract_test.go`
  - `target/go.mod`
- **Files changed (stage-owned artifacts):**
  - `target/specs/dependency-flattening/PLAN.md`
  - `target/specs/dependency-flattening/TASKS.md`
  - `target/specs/dependency-flattening/quickstart.md`
  - `target/specs/dependency-flattening/contracts/flattening-contract.md`

**Spec requirements covered (F-003 only):** FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007; AC-001..AC-005.

## 2. Build Steps

1. Added/extended F-003 domain contracts in `internal/core/model.go`:
   - `DependencyNode`
   - `ModuleInventoryMap`
   - `FlattenOptionsView`
   - metadata-bearing `ModuleEntry`
2. Added flatten stage interface boundary in `internal/core/service.go` (`FlattenStage`).
3. Implemented `internal/flatten/flatten.go`:
   - recursive walk over `dependencies`
   - canonical `name@version` key generation
   - visited-key dedupe/cycle short-circuit
   - prod/dev inclusion gate (`IncludeDev`)
   - metadata extraction into inventory entries
4. Added unit tests in `internal/flatten/flatten_test.go` for:
   - deep nested traversal to flat map
   - duplicate path + cycle behavior
   - metadata propagation
   - prod/dev gating
   - deterministic key-set comparison via sorted keys
5. Added integration contract test in `test/integration/flatten_contract_test.go` ensuring downstream consumption from flattened map only.
6. Updated execution artifacts and checklists:
   - checked completed items in `PLAN.md` and `TASKS.md`
   - recorded contract/open decisions and verification walkthrough.

## 3. Test Evidence

### Tests added/updated

- `target/internal/flatten/flatten_test.go` (new)
- `target/test/integration/flatten_contract_test.go` (new)

### Commands run

From `target/`:

```bash
gofmt -w internal/core/model.go internal/core/service.go internal/flatten/flatten.go internal/flatten/flatten_test.go test/integration/flatten_contract_test.go
go test ./internal/flatten ./test/integration -run Flatten
```

### Observed results

```text
ok   theseus/target/internal/flatten      0.644s
ok   theseus/target/test/integration      0.398s
```

## 4. Residual Risks

- **OD-001 unresolved:** behavior for nodes missing `name` and/or `version` is still `NEEDS CLARIFICATION`; current implementation skips such nodes.
- **OD-002 unresolved:** mandatory vs optional metadata contract remains `NEEDS CLARIFICATION`; current implementation carries repository/author/url/path/private/dev fields.
- Full-repo test suite is not part of this bounded F-003 slice evidence; scope-limited tests were executed for flattening slice behavior.

## 5. Traceability

- **FR-001 / FR-002 / AC-001 / SC-001** → `internal/flatten/flatten.go` recursive walk + canonical key map population; nested fixture in `internal/flatten/flatten_test.go`.
- **FR-003 / AC-002 / SC-002** → visited-key guard in `walk()`; duplicate/cycle test in `internal/flatten/flatten_test.go`.
- **FR-004 / AC-003 / SC-003** → `includeNode()` gate before insertion; dev/prod fixture test.
- **FR-005 / AC-004 / SC-004** → `toModuleEntry()` metadata propagation + metadata assertions in tests.
- **FR-006 / AC-005** → explicit flatten stage contract (`FlattenStage` in `internal/core/service.go`) and downstream map-only integration harness (`test/integration/flatten_contract_test.go`).
- **FR-007** → deterministic uniqueness asserted via sorted key helper in `internal/flatten/flatten_test.go`.
