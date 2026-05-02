# Feature Specification: Dependency Graph Loading and Traversal Control (F-002)

**Created**: 2026-05-01  
**Status**: Draft  
**Input**: User description: "Role: Spec-Writer. Hard gate retry: artifact still missing. Focus ONLY F-002 (slug `dependency-graph-loading`)."

## Source Traceability *(mandatory)*

- **SRC-FI-002**: `FUNCTIONALITY-INDEX.md` defines F-002 as ready, scoped to `read-installed` orchestration + depth/dev toggles (`FUNCTIONALITY-INDEX.md:11-15`).
- **SRC-AN-002A**: `ANALYSIS.md` defines F-002 behavior: loader called with `start`, `depth` from normalized `direct`, and `opts.dev` toggled by `production/development` (`ANALYSIS.md:84-90`).
- **SRC-AN-002B**: `ANALYSIS.md` flow confirms `init` builds loader options and invokes dependency-tree load before flattening (`ANALYSIS.md:160-163`).
- **SRC-OV-R002**: `OVERVIEW.md` requirement R-002 requires recreating start path + depth/dev behavior (`OVERVIEW.md:47`).
- **SRC-CODE-INIT**: `lib/index.js` sets loader opts (`dev`, `depth`), flips `dev` when `production || development`, then calls `read(options.start, opts, cb)` (`lib/index.js:267-276`, `298`).
- **SRC-CODE-ARGS**: `lib/args.js` normalizes `direct` to `0` (direct-only) or `Infinity` (recursive), which is consumed as loader depth (`lib/args.js:79-83`; `lib/index.js:270`).
- **SRC-README-LOAD**: README states dependency walking is performed through `read-installed` from the module tree (`README.md:177-182`).
- **SRC-TEST-ERR**: Error path for invalid/unresolvable start location returns an error via callback (`tests/test.js:379-385`).

## Feature Boundary & I/O Contract

**In scope (F-002 only)**
- Building dependency-loader options for traversal behavior (`depth`, `dev`, logger linkage).
- Invoking dependency graph loading from a configured start path.
- Returning loader success/error outcome to the next pipeline boundary.

**Out of scope (handled by other functionality IDs)**
- Graph flattening semantics (F-003).
- License detection (F-004).
- Filtering/policy exits (F-005).
- Output formatting/export (F-006).

**Inputs**
- `start` (string path): root location to scan. **[Sources: SRC-AN-002A, SRC-CODE-INIT]**
- `direct` (normalized numeric depth): `0` for direct-only, recursive sentinel otherwise. **[Source: SRC-CODE-ARGS]**
- `production` / `development` (boolean flags): drive loader dev toggle behavior. **[Sources: SRC-AN-002A, SRC-CODE-INIT]**

**Outputs**
- On success: loaded dependency tree passed into downstream flattening stage. **[Sources: SRC-AN-002B, SRC-CODE-INIT]**
- On failure: non-null loader error returned through callback/error channel. **[Source: SRC-TEST-ERR]**

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Load dependency graph from a selected start path (Priority: P1)

As a user, I can start a scan from a specific project path so the system builds the dependency graph for that project before any license processing.

**Why this priority**: All downstream behaviors (flattening, filtering, policy checks, output) depend on a successfully loaded dependency graph. Without this, no scan can proceed. (Sources: SRC-AN-002B, SRC-OV-R002)

**Independent Test**: Can be fully tested by running a scan with a valid start path and verifying the loader is invoked with that path and returns data to the next stage.

**Acceptance Scenarios**:

1. **Given** a valid project directory, **When** scan initialization runs, **Then** dependency loading MUST start from that exact path.
2. **Given** dependency loading succeeds, **When** callback returns, **Then** graph data is passed forward to flattening/post-processing.

---

### User Story 2 - Limit traversal to direct dependencies only (Priority: P2)

As a user, I can request direct-only traversal so transitive dependencies are not traversed during graph loading.

**Why this priority**: Direct-only scope is a common audit need and materially changes scan scope while keeping the rest of the workflow unchanged. (Sources: SRC-CODE-ARGS, SRC-CODE-INIT)

**Independent Test**: Can be tested by enabling direct mode and confirming loader depth is set to direct-only semantics (`0`) rather than recursive traversal.

**Acceptance Scenarios**:

1. **Given** direct mode is enabled, **When** loader options are created, **Then** traversal depth MUST be set to direct-only.
2. **Given** direct mode is not enabled, **When** loader options are created, **Then** traversal depth MUST allow recursive traversal.

---

### User Story 3 - Control dev-dependency loading behavior for scoped scans (Priority: P3)

As a user, I can run production/development-scoped scans that adjust dependency loading behavior before downstream filtering.

**Why this priority**: This preserves compatibility with existing scan scope behavior tied to production/development intent. (Sources: SRC-AN-002A, SRC-CODE-INIT)

**Independent Test**: Can be tested by toggling production/development flags and verifying the loader `dev` option changes as specified.

**Acceptance Scenarios**:

1. **Given** neither production nor development scope is requested, **When** loader options are built, **Then** dev-loading remains enabled.
2. **Given** production or development scope is requested, **When** loader options are built, **Then** dev-loading is disabled in loader options.

### Edge Cases

- What happens when the start path is invalid or cannot resolve installed packages? The load step MUST surface an error to callback and stop normal graph handoff. (Source: SRC-TEST-ERR)
- How does system handle explicit production/development flags combined with depth settings? Depth behavior and dev-loading behavior MUST be applied independently (depth from normalized direct; dev toggle from scope flags). (Sources: SRC-CODE-ARGS, SRC-CODE-INIT)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST begin dependency-graph loading from the caller-provided start path. **[Sources: SRC-AN-002A, SRC-CODE-INIT]**
- **FR-002**: System MUST construct loader options that include traversal depth derived from normalized direct mode. **[Sources: SRC-CODE-ARGS, SRC-CODE-INIT]**
- **FR-003**: System MUST use direct-only depth semantics when direct mode is enabled. **[Sources: SRC-CODE-ARGS, SRC-CODE-INIT]**
- **FR-004**: System MUST use recursive-depth semantics when direct mode is not enabled. **[Sources: SRC-CODE-ARGS, SRC-CODE-INIT]**
- **FR-005**: System MUST set dev-loading enabled by default for graph loading unless production/development scope is requested. **[Sources: SRC-AN-002A, SRC-CODE-INIT]**
- **FR-006**: System MUST disable loader dev-loading when either production or development scope flag is set. **[Sources: SRC-AN-002A, SRC-CODE-INIT]**
- **FR-007**: System MUST surface dependency-loader errors through its callback/error channel and MUST NOT continue as a successful load. **[Sources: SRC-TEST-ERR, SRC-AN-002B]**
- **FR-008**: On successful load, System MUST pass the loaded dependency tree to the next pipeline stage (flattening). **[Sources: SRC-AN-002B, SRC-CODE-INIT]**

### Key Entities *(include if feature involves data)*

- **Scan Options**: Runtime scan inputs relevant to this feature (`start`, `direct`-normalized depth, `production`, `development`). **[Sources: SRC-AN-002A, SRC-CODE-ARGS]**
- **Loader Options**: Options passed to dependency-tree loader (`depth`, `dev`, logging hook). **[Source: SRC-CODE-INIT]**
- **Loaded Dependency Tree**: Loader output representing discovered dependency graph to be consumed by flattening. **[Sources: SRC-AN-002B, SRC-CODE-INIT]**

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For 100% of valid scan invocations, dependency loading is initiated with the exact provided start path.
- **SC-002**: For 100% of runs with direct mode enabled, effective loader depth is direct-only; for 100% without direct mode, effective loader depth is recursive.
- **SC-003**: For 100% of runs where `production` or `development` is set, loader dev-loading is disabled; for runs with neither set, it remains enabled.
- **SC-004**: For invalid/unresolvable start locations, loader failure is surfaced via error callback in 100% of cases and no successful graph handoff is reported.

## Assumptions

- F-001-compatible option normalization is already in place before this feature executes, especially direct-depth normalization (`0` vs recursive). **[Source: SRC-CODE-ARGS]**
- This feature covers only graph-loading and traversal-control behavior; flattening/filtering/output are handled by other functionality slices. **[Sources: SRC-AN-002B, SRC-OV-R002]**
- Existing dependency-loader behavior (read-installed style traversal) is preserved for compatibility rather than replaced in this requirement slice. **[Sources: SRC-FI-002, SRC-README-LOAD]**

## Open Questions

- None blocking for F-002 at this stage. **[Source: SRC-FI-002]**

## Relevant Files for This Requirement

- `target/specs/ANALYSIS.md`
- `target/specs/OVERVIEW.md`
- `target/specs/FUNCTIONALITY-INDEX.md`
- Upstream evidence anchors: `lib/index.js`, `lib/args.js`, `README.md`, `tests/test.js` (as cited above)
