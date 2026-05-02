# F-003 Flattening Contract

## Input

- `DependencyNode` root graph (already loaded upstream)
- `FlattenOptionsView{IncludeDev bool}`

## Output

- `ModuleInventoryMap` keyed by `name@version`
- one-entry-per-key via visited guard

## Compatibility notes

- Flatten stage boundary preserved: downstream consumers use flat inventory map only.
- First encounter wins for duplicate/cycle revisits.
- Excluded dev nodes are not inserted when `IncludeDev=false`.

## Open decisions

- **OD-001 (NEEDS CLARIFICATION)**: Missing `name` and/or `version` behavior for canonical key creation. Current PoC behavior: skip such nodes.
- **OD-002 (NEEDS CLARIFICATION)**: Mandatory vs optional metadata field set for strict compatibility parity.
