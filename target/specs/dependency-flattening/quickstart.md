# F-003 Quickstart Verification Notes

## Fixture coverage plan

- Nested traversal depth >= 3 (AC-001)
- Duplicate path + cycle fixture with one-entry-per-key guarantee (AC-002)
- Mixed prod/dev inclusion fixture (AC-003)
- Metadata propagation fixture for repository/author/url/path/private fields (AC-004)
- Downstream contract fixture consuming flattened map only (AC-005)

## Execution walkthrough

1. Run `go test ./...` from `target/`.
2. Confirm `internal/flatten` tests pass for traversal, dedupe/cycle, metadata, and gating.
3. Confirm `test/integration/flatten_contract_test.go` passes for downstream map-only contract.

## Evidence checklist

- [x] Unit tests executed
- [x] Integration contract test executed
- [x] Deterministic key comparison uses sorted key sets
- [x] OD-001 and OD-002 kept explicit as clarification items
