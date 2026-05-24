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

## Guardrails

- Do not reconstruct internal logic before the CLI surface is fully mapped
- Do not invent flag semantics — derive from help text, source, or tests only
- If a flag's behavior is ambiguous, record it as an open question rather than guessing
- Keep CLI surface reconstruction separate from internal domain reconstruction
