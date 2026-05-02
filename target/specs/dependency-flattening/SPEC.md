# Feature Specification: F-003 Flatten dependency graph into module inventory

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Input artifacts: `/Users/igorbesel/Projekte/theseus/target/specs/ANALYSIS.md`, `OVERVIEW.md`, `FUNCTIONALITY-INDEX.md`. Focus slice only: F-003 — Flatten dependency graph into module inventory (slug `dependency-flattening`). Task: Produce implementation-ready `SPEC.md` for this single functionality in `/Users/igorbesel/Projekte/theseus/target/specs/dependency-flattening/SPEC.md`. Include explicit requirements, acceptance criteria, compatibility notes, constraints, and open decisions limited to this slice. Do not implement code. Do not cover unrelated functionality."

## Scope

This specification is strictly limited to **F-003 — flattening a loaded dependency graph into a module inventory map**.

- In scope: recursive graph walk, canonical keying, duplicate/cycle guard behavior, prod/dev gating at flatten time, and module metadata extraction performed by the flattening step.
- Out of scope: graph loading (`read-installed`), license detection/normalization internals, filtering/policy exits, and output rendering/export.

## Source Traceability Registry

- **S-FI-003**: `FUNCTIONALITY-INDEX.md` F-003 definition, slug, and evidence summary (`17-21`).
- **S-OV-R003**: `OVERVIEW.md` requirement task R-003 for dedupe + metadata extraction (`48`).
- **S-AN-F003**: `ANALYSIS.md` F-003 behavior summary (`91-99`).
- **S-AN-FLAT**: `ANALYSIS.md` flatten function inventory and scope (`60`, `93`).
- **S-AN-FLOW**: `ANALYSIS.md` flow position: flatten runs after graph load and before filters (`162-164`).
- **S-AN-INV**: `ANALYSIS.md` invariants: key format `name@version` and one entry per key (`189-191`).
- **S-AN-DOM**: `ANALYSIS.md` domain map entity note for per-module info (`145-149`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Produce canonical module inventory from nested dependencies (Priority: P1)

As a migration maintainer, I need nested dependency data converted into one canonical module map so downstream stages operate on a stable inventory shape.

**Why this priority**: This is the foundational transformation for all later slices; if inventory flattening is wrong, filtering and output cannot be trusted. (Sources: S-FI-003, S-AN-F003, S-AN-FLOW)

**Independent Test**: Can be fully tested by providing a fixture dependency tree with nested dependencies and asserting a flat map keyed as `name@version`.

**Acceptance Scenarios**:

1. **Given** a dependency tree with parent and nested child modules, **When** flattening executes, **Then** the result MUST be a flat map with canonical keys in `name@version` format. (Sources: S-AN-F003, S-AN-INV)
2. **Given** a module inventory consumer stage, **When** flattening completes, **Then** downstream stage input MUST be the flattened map structure rather than the original nested tree. (Sources: S-AN-FLOW, S-AN-DOM)

---

### User Story 2 - Prevent duplicate/circular expansion from corrupting inventory (Priority: P2)

As a maintainer, I need duplicate and circular references to resolve deterministically so scans terminate and do not duplicate module entries.

**Why this priority**: Real dependency graphs can repeat modules; deterministic dedupe is required for stable inventory counts and run completion. (Sources: S-AN-F003, S-AN-INV)

**Independent Test**: Can be tested by using fixtures where the same `name@version` appears through multiple paths and/or cycles, then validating one inventory entry per key.

**Acceptance Scenarios**:

1. **Given** a graph where the same module key is reachable via two paths, **When** flattening executes, **Then** the inventory MUST contain exactly one entry for that key. (Sources: S-AN-F003, S-AN-INV)
2. **Given** a cyclic dependency path, **When** flattening revisits an already-seen key, **Then** recursion MUST short-circuit for that branch and continue without infinite traversal. (Sources: S-AN-F003)

---

### User Story 3 - Preserve module metadata needed by later stages (Priority: P3)

As a compliance/reporting pipeline maintainer, I need flattening to carry forward module metadata fields so downstream license/filter/output stages have required context.

**Why this priority**: Metadata continuity is necessary for compatibility, but secondary to producing a correct deduplicated map. (Sources: S-OV-R003, S-AN-F003, S-AN-DOM)

**Independent Test**: Can be tested by scanning fixture modules with repository/author/path/private-related fields and verifying those fields appear in flattened entries.

**Acceptance Scenarios**:

1. **Given** module nodes containing repository/author/url/path/private-related metadata, **When** flattening executes, **Then** flattened entries MUST include extracted metadata fields used by downstream stages. (Sources: S-AN-F003, S-AN-DOM)
2. **Given** production/development gating inputs, **When** flattening executes, **Then** entries excluded by the prod/dev gate at flatten stage MUST not be added to the inventory map. (Sources: S-AN-F003)

### Edge Cases

- Revisited module keys encountered during recursion MUST not create additional entries (duplicate/cycle guard behavior). (Sources: S-AN-F003, S-AN-INV)
- Dependencies nested multiple levels deep MUST still be traversed recursively via `dependencies` relationships. (Sources: S-AN-F003)
- Modules omitted by prod/dev gate at flatten stage MUST be skipped consistently. (Sources: S-AN-F003)
- If module identity fields required for canonical key creation are missing, behavior is not fully specified in Analyzer artifacts and requires explicit clarification. (Sources: S-AN-F003)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST transform a nested dependency graph into a flat module inventory map keyed by canonical `name@version`. (Sources: S-FI-003, S-AN-F003, S-AN-INV)
- **FR-002**: System MUST traverse module dependencies recursively through dependency relationships to populate the flat inventory. (Sources: S-AN-F003)
- **FR-003**: System MUST enforce one-entry-per-key behavior using a duplicate/circular guard so revisiting an existing key does not re-expand that node. (Sources: S-AN-F003, S-AN-INV)
- **FR-004**: System MUST apply prod/dev inclusion gating during flattening so excluded nodes are not inserted into the inventory map. (Sources: S-AN-F003)
- **FR-005**: System MUST extract and persist per-module metadata required by downstream pipeline stages, including repository/author/url/path/private-related fields observed in this slice. (Sources: S-OV-R003, S-AN-F003, S-AN-DOM)
- **FR-006**: System MUST produce inventory output as the flatten stage result contract consumed by subsequent filtering and serialization stages. (Sources: S-AN-FLOW, S-AN-DOM)
- **FR-007**: System MUST preserve deterministic behavior such that identical input trees produce identical key sets and per-key uniqueness outcomes. (Sources: S-AN-INV)

### Acceptance Criteria (Requirement-Level)

- **AC-001 (FR-001/FR-002)**: For a fixture tree with depth >= 3, flattened output contains entries for all reachable included modules and no nested structure in result shape. (Sources: S-AN-F003)
- **AC-002 (FR-003/FR-007)**: For a fixture where one module key is reachable via multiple paths (including a cycle), output contains exactly one entry for that key and traversal completes. (Sources: S-AN-F003, S-AN-INV)
- **AC-003 (FR-004)**: For fixtures with mixed prod/dev nodes under gating options, output membership matches flatten-stage inclusion rules with zero unexpected entries. (Sources: S-AN-F003)
- **AC-004 (FR-005)**: For fixtures containing repository/author/path/private metadata, flattened entries expose those metadata values for downstream consumption. (Sources: S-OV-R003, S-AN-F003)
- **AC-005 (FR-006)**: Downstream stage harness consuming flatten output can process the returned module map without requiring original nested graph input. (Sources: S-AN-FLOW, S-AN-DOM)

### Inputs and Outputs (Slice I/O Contract)

- **Inputs**
  - `dependencyTree`: nested dependency-node object graph rooted at the selected start module; each node may include `name`, `version`, `dependencies`, and package metadata fields. (Sources: S-AN-F003)
  - `flattenOptions`: runtime options relevant to this slice, including prod/dev inclusion controls already normalized by upstream option handling. (Sources: S-AN-F003, S-AN-FLOW)
- **Outputs**
  - `moduleInventoryMap`: flat map keyed by `name@version`, where each value is module metadata extracted during flattening and intended for downstream filtering/output stages. (Sources: S-AN-INV, S-AN-DOM)
  - **Behavioral output guarantee**: one entry per canonical key with duplicate/cycle guard semantics applied. (Sources: S-AN-F003, S-AN-INV)

### Key Entities *(include if feature involves data)*

- **DependencyNode**: One node from the loaded dependency tree containing module identity, dependency links, and package metadata consumed during flattening. (Sources: S-AN-F003)
- **ModuleInventoryMap**: Flat map keyed as `name@version` with one entry per key and extracted metadata payload. (Sources: S-AN-INV, S-AN-DOM)
- **FlattenOptionsView**: Flatten-relevant option subset controlling inclusion gates (not full CLI model). (Sources: S-AN-F003)

## Compatibility Notes

- Canonical module key format MUST remain `name@version` for parity with downstream expectations. (Sources: S-AN-INV)
- Duplicate/circular guard behavior MUST remain first-encounter-wins (no second entry for same key). (Sources: S-AN-F003, S-AN-INV)
- Flattening MUST remain the stage boundary between graph loading and filter/policy processing. (Sources: S-AN-FLOW)
- Metadata extraction surface from flattening MUST remain available so later slices retain behavior compatibility. (Sources: S-OV-R003, S-AN-F003)

## Constraints

- **C-001 (Scope constraint)**: This slice MUST NOT redefine dependency loading strategy (`read-installed`) or graph acquisition behavior. (Sources: S-AN-FLOW)
- **C-002 (Scope constraint)**: This slice MUST NOT include license parsing/normalization logic beyond carrying required metadata into downstream stages. (Sources: S-AN-FLOW)
- **C-003 (Behavior constraint)**: This slice MUST operate on canonical key identity and not alternate identity schemes unless separately approved. (Sources: S-AN-INV)
- **C-004 (Evidence constraint)**: Any behavior not explicit in Analyzer artifacts MUST be recorded as an open decision, not inferred. (Sources: S-FI-003)

## Open Decisions (slice-limited)

1. **OD-001 — Missing identity fields**: If a dependency node lacks `name` and/or `version`, exact legacy flatten behavior for key creation is not explicit in provided Analyzer artifacts and MUST be clarified before compatibility lock. **Status**: NEEDS CLARIFICATION. (Sources: S-AN-F003)
2. **OD-002 — Metadata completeness contract**: Analyzer confirms extraction of repository/author/url/path/private-related data, but does not fully enumerate mandatory-vs-optional field set for compatibility assertions; final contract granularity MUST be confirmed during architecture/planning handoff. **Status**: NEEDS CLARIFICATION. (Sources: S-OV-R003, S-AN-F003, S-AN-DOM)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For compatibility fixtures, 100% of included reachable modules are represented in flattened output key set using `name@version` keys. (Sources: S-AN-F003, S-AN-INV)
- **SC-002**: For fixtures with repeated/cyclic references, flattened output has zero duplicate keys and traversal completes without recursion failure in 100% of runs. (Sources: S-AN-F003, S-AN-INV)
- **SC-003**: For fixtures with prod/dev mixed nodes, output inclusion/exclusion matches flatten-stage gate behavior with zero mismatches. (Sources: S-AN-F003)
- **SC-004**: For fixtures containing metadata fields used downstream, flattened entries preserve required metadata values for all represented modules. (Sources: S-OV-R003, S-AN-DOM)

## Assumptions

- Upstream stage has already produced a readable dependency tree input before this slice starts. (Sources: S-AN-FLOW)
- Option normalization (CLI/programmatic parsing) has already occurred before flattening receives options. (Sources: S-AN-FLOW)
- This slice’s output is a data contract for downstream stages, not final user-facing rendering. (Sources: S-AN-FLOW, S-AN-DOM)

## Relevant Files for This Requirement

- `/Users/igorbesel/Projekte/theseus/target/specs/ANALYSIS.md`
- `/Users/igorbesel/Projekte/theseus/target/specs/OVERVIEW.md`
- `/Users/igorbesel/Projekte/theseus/target/specs/FUNCTIONALITY-INDEX.md`
