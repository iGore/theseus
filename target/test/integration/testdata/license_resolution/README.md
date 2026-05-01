# F-004 License Resolution Fixture Index

This fixture index maps SPEC acceptance criteria to deterministic F-004 scenarios.

| Fixture Name | Purpose | SPEC Mapping |
|---|---|---|
| `metadata_first` | Metadata wins over README/files | AC-001, FR-001, FR-002 |
| `readme_first` | README used when metadata missing | AC-001, User Story 2 Scenario 1 |
| `file_first` | File fallback used when metadata/README unresolved | AC-001 |
| `ordered_files` | LICENSE* > LICENCE* > COPYING > README | AC-004, FR-006 |
| `spdx_vs_heuristic` | Direct SPDX vs inferred `*` marker | AC-002, FR-003, FR-004 |
| `custom_reference` | `Custom: <value>` normalization | AC-003, FR-005 |
| `compat_od_placeholders` | Explicit unresolved OD-001/OD-002 markers | AC-005, Open Decisions |
