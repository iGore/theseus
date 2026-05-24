---
name: cli-analyzer
description: >
  Use this skill when the source module under analysis is a command-line tool.
  It extends the Sourcebot skill with CLI-specific reconstruction focus:
  entry points, flag and argument parsing, subcommand dispatch, environment
  variables, exit codes, and dependency-ordered requirements. Use it to ensure
  the FUNCTIONALITY-INDEX.md captures the CLI surface in the order downstream
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

### 7. End-to-End Behaviour as an Architectural Anchor

The end-to-end behaviour of the tool (invoke → parse → load → process →
render → exit) is not itself a functionality and MUST NOT be added as a
synthetic entry in FUNCTIONALITY-INDEX.md. It is an architectural
concern and is propagated to ARCHITECTURE.md by the Architect.

ANALYSIS.md MUST therefore contain a dedicated `## End-to-End Pipeline`
subsection that names, for the source system, the ordered sequence of
internal stages between CLI invocation and output rendering, and links
each stage to the functionality slice that implements it. This ordered
sequence is the input the Architect needs to define the Composition
Root of the target system.

## Dependency-Ordered Requirements

When producing FUNCTIONALITY-INDEX.md for a CLI tool, order functionality
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

Mark each item in FUNCTIONALITY-INDEX.md with its upstream dependency IDs
so the Orchestrator can enforce ordering in Spec-Writer fan-out.

## Guardrails

- Do not reconstruct internal logic before the CLI surface is fully mapped
- Do not invent flag semantics — derive from help text, source, or tests only
- If a flag's behavior is ambiguous, record it as an open question rather than guessing
- Keep CLI surface reconstruction separate from internal domain reconstruction
