---
name: cli-analyzer
description: >
  Use this skill when the source module under analysis is a command-line tool.
  It extends the Sourcebot skill with CLI-specific reconstruction focus:
  entry points, flag and argument parsing, subcommand dispatch, environment
  variables, exit codes, and dependency-ordered requirements. Use it to ensure
  the USE-CASES.md captures the CLI surface in the order downstream
  Spec-Writers need.
---

# CLI Analyzer

## Purpose

Extend the standard Analyzer reconstruction with CLI-specific depth.
A CLI tool's surface contract (entry points, flags, args, exit codes) is
often the most important behavioral specification and must be captured
before any internal logic can be meaningfully decomposed.

## When to use

- Source module is a CLI tool (has a `bin/`, `cmd/`, or executable entry point)
- Reconstruction must capture flag and argument contracts, not just internal logic
- Requirements ordering matters for downstream Spec-Writer fan-out

## CLI-Specific Reconstruction Focus

### 1. Entry Point Mapping

- Locate the executable entry point (e.g. `bin/`, `main.js`, `cmd/root.go`)
- Identify the command dispatch pattern: single command, subcommand tree, or plugin-based
- Record the full invocation signature: `tool [global-flags] <command> [command-flags] [args]`

### 2. Flag and Argument Inventory

For every flag and positional argument, record:
- Name and aliases (short/long form)
- Type (string, boolean, integer, enum)
- Default value and whether it is required
- Mutual exclusion or dependency with other flags
- Description as documented in help text or source comments

Group flags by scope: global flags vs. subcommand-local flags.

### 3. Subcommand Structure

- Map every subcommand and its own flag/arg surface
- Note which subcommands share flags via parent inheritance
- Record aliases and hidden/deprecated commands

### 4. Environment Variables and Config Files

- List every environment variable the tool reads, with fallback precedence
- Note any config file paths and their loading order relative to flags

### 5. Exit Codes and Stderr Contract

- Document known exit codes and their semantics
- Note what goes to stdout vs. stderr (output vs. diagnostic)
- Record any structured output formats (JSON, CSV, plain text) activated by flags

### 6. Golden Output Snippets from the Source System

For every observable output mode the source tool supports — whether
selected by a flag, a subcommand, an environment variable, or a config
toggle — ANALYSIS.md MUST capture a verbatim sample of real
source-system output. Do not paraphrase or reconstruct these from source
code; execute the source system on a minimal fixture and record the
actual result.

For each output mode:

- Construct a minimal fixture that is just large enough to exercise the
  mode meaningfully (one or two representative inputs is usually enough).
  Keep the fixture small enough that the captured output stays readable.
- Run the source system once with exactly the flag / subcommand / env
  configuration that selects the mode.
- Capture verbatim stdout (or a short representative excerpt when the
  full output is large), verbatim stderr, and the observed exit code
  into a dedicated subsection of ANALYSIS.md, clearly labelled per
  output mode.

These snippets become the canonical reference for the Architect and the
Spec-Writer downstream: any output-format requirement in SPEC.md MUST
reference one of these snippets as its acceptance baseline, not a
narrative description. This closes the class of defects where the
target system silently invents its own output schema because the spec
only described the format in prose.

**Byte-exactness and edge coverage (required).** Snippets are a parity
baseline, so capture them byte-for-byte — preserve field order, exact
sentinel strings, indentation, and trailing newlines, and note where
stdout (often newline-terminated by a print wrapper) differs from file
output (often written raw). "Semantic" or "selected-snippet" comparison
is insufficient and is the documented cause of past divergence. The
fixture set MUST exercise, in addition to the happy path, at least:

- a record with **missing optional fields** (no author, no repository)
- a record using each **sentinel/default** path (e.g. private → `UNLICENSED`,
  nothing found → `UNKNOWN`)
- inputs that force **value-derivation from secondary sources** (e.g. a
  license determinable only from file text, not metadata)
- input values **outside the common set** that the source still classifies
  via its full rule table (e.g. GPL/LGPL/ISC-from-text, SPDX expressions),
  to expose any fallback branch that dumps raw input
- every **output mode crossed with the flags that alter it** (custom format,
  component prefix, relative paths, summary), not just the default mode

**Value-transformation tables (required).** When an output field is the
result of a classification/normalization function, ANALYSIS.md MUST
include the function's *complete, ordered* rule table with the verbatim
output literal of every branch and the no-match/fallback branch — copied
from source, not summarized. The first matching rule wins, so an omitted
earlier rule changes results. Record any third-party transformation the
source delegates to (SPDX validator, tree/CSV/markdown formatter) and the
exact behavior a stdlib-only target must replicate.

### 7. Input-Field Tracing (Data-Flow Reconstruction)

The CLI surface (entry points, flags, exit codes) captures the
*control-flow contract* of the tool, but it does not capture the
*data-flow contract*: how each externally-supplied input field is read,
transformed, and possibly overridden before reaching output. A flag-
centric reconstruction systematically misses behaviour that lives as a
flag-less side effect on a single input field — for instance an input
marker that silently rewrites a sibling field before rendering, an
auxiliary descriptor that suppresses a diagnostic, or a configuration
key that overrides a CLI flag under one undocumented condition. These
defects are invisible to a CLI-first reading and only surface in
differential testing against the source system.

To close this class, ANALYSIS.md MUST contain a section titled
`## Input-Field Inventory` that enumerates, for every input field the
source tool reads from external data sources (configuration files,
manifest or descriptor files, environment variables, stdin, network
responses, or any other input outside the CLI argument vector):

- **Field name and source**: the exact key and where it is read from,
  using the source-system's own naming for both the file or stream and
  the key path within it
- **Default path**: how the field flows through the pipeline when no
  special condition applies — which stage reads it, which fields it
  populates in the in-memory model, which output property it surfaces
  through
- **Conditional overrides**: every code location that assigns this
  field a value different from the one read, regardless of whether a
  CLI flag triggers it. For each override, record: source file and line
  number, the triggering condition expressed in domain terms, the new
  value, and whether any CLI flag, env var, or configuration key
  controls it. An override with no controlling input is an
  *unconditional override* and MUST be called out explicitly as such
- **Test coverage reference**: the test in the source repository that
  exercises each override, if one exists. Tests whose grouping label
  (e.g. `describe` block, test class, or fixture name) suggests a flag
  while the test body actually exercises a flag-less override MUST be
  cross-referenced here, not only under the flag they appear to
  belong to

For every unconditional override identified, the Analyzer MUST also
construct a probe fixture that triggers the override path, execute the
source system on it, and capture the resulting output as an additional
Golden Output Snippet labelled with the override it exercises. The
Spec-Writer downstream MUST surface each override as an explicit FR
in the slice that owns the affected output field, not in the slice
that owns the surface flag whose name the field happens to share.

Skipping this section is permitted only when the source tool consumes
no external data beyond its CLI arguments — a rare case in practice.

### 8. End-to-End Behaviour as an Architectural Anchor

The end-to-end behaviour of the tool (invoke → parse → load → process →
render → exit) is not itself a functionality and MUST NOT be added as a
synthetic entry in USE-CASES.md. It is an architectural
concern and is propagated to ARCHITECTURE.md by the Architect.

ANALYSIS.md MUST therefore contain a dedicated `## End-to-End Pipeline`
subsection that names, for the source system, the ordered sequence of
internal stages between CLI invocation and output rendering, and links
each stage to the functionality slice that implements it. This ordered
sequence is the input the Architect needs to define the Composition
Root of the target system.

### 9. Source-Dependency Behavior Contracts (documented I/O + examples)

When the source delegates non-trivial behavior to a third-party library, the
output identity of the whole tool depends on that library's exact semantics.
Reconstructing the *calling code* is not enough — the Analyzer MUST also
reconstruct the *library's contract*, because a stdlib-only reimplementation
that re-derives it by hand is the dominant cause of divergence.

ANALYSIS.md MUST contain a section titled `## Source-Dependency Contracts`
that lists every dependency whose behavior shapes observable output
(examples of such categories, abstractly: dependency-graph/resolution
libraries, standardized-identifier validators or expression parsers,
output/tree/table formatters, escaping or serialization helpers, path or
URL normalizers). For each one, record:

- **Name, version, and role**: the package as pinned in the source manifest,
  and the one behavior the tool relies on.
- **Documented input contract**: accepted input shapes/types and relevant
  options, taken from the library's own documentation/README — not guessed
  from the call site. Cite the doc source.
- **Documented output contract**: return shape, ordering guarantees, escaping
  rules, sentinel/error values, and any whitespace/formatting it controls.
- **Worked example (mandatory)**: at least one concrete `input → output`
  example for the relied-upon behavior, copied from the library's docs or
  produced by exercising the library on a minimal input. Edge inputs
  (empty, malformed, multi-value) should each get an example when they change
  the result. These examples are the contract the target must reproduce.
- **Re-derivation risk**: a one-line note on how hard the behavior is to
  reproduce faithfully (trivial / moderate / standardized-and-deep). This
  feeds the Architect's replacement decision (see the `go-architect` skill).

This section gives the Architect the evidence to decide whether to adopt an
equivalent target-language library or reimplement, and gives the Spec-Writer
concrete example I/O to attach to every transformation requirement.

## Dependency-Ordered Requirements

When producing USE-CASES.md for a CLI tool, order functionality
items so that downstream Spec-Writers can work in dependency sequence:

1. **Entry point and invocation model** — must be specified first; all other
   slices depend on knowing how the tool is invoked
2. **Flag and argument parsing** — must come before any behavioral slice
   because every behavioral requirement references specific flags
3. **Core domain logic slices** — in dependency order; a slice that produces
   data another slice consumes must be specified first
4. **Output and rendering** — depends on core logic producing its input
5. **Error handling and exit codes** — cross-cutting, but specify after
   the happy path is clear
6. **Configuration loading** — specify after the flag surface is clear,
   since config often mirrors or overrides flags

Mark each item in USE-CASES.md with its upstream dependency IDs
so the Orchestrator can enforce ordering in Spec-Writer fan-out.

## Parity strategy: in-flow golden tests vs. out-of-flow live run

When the goal is a behavior-identical rewrite, capturing snippets is not
enough — but note the split:

- **In-flow (test-driven):** the Golden Output Snippets captured here (from the
  pinned source version) are committed as **fixture files**. The Builder writes
  byte-level golden tests against them; the Verifier confirms those tests cover
  the matrix and pass. The agent flow does **not** install or run the live
  source system — the captured examples are the oracle.
- **Out-of-flow (manual, by maintainers):** a live differential run (real source
  vs. target across the full matrix, via a maintainer-run differential script) is performed once at
  the end, by hand, to catch anything the captured examples missed. Keep the
  example set broad (edge set above) so this final run is cheap and likely green.

For both, comparison is a **byte-level diff**, normalizing only volatile values
that legitimately differ (e.g. absolute paths under a temp root) and documenting
each normalization explicitly. Any remaining diff is a defect to be specified
away, not waved through as "semantically equivalent".

## Guardrails

- Do not reconstruct internal logic before the CLI surface is fully mapped
- Do not invent flag semantics — derive from help text, source, or tests only
- If a flag's behavior is ambiguous, record it as an open question rather than guessing
- Keep CLI surface reconstruction separate from internal domain reconstruction
- Do not summarize classification/transformation logic — transcribe the full
  ordered rule table including the fallback branch and verbatim output literals
- Do not accept prose output descriptions or "semantic"/partial snippet
  comparison as a parity baseline; require byte-exact golden snippets across
  the full edge set and flag matrix
