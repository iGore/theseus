---
name: go-architect
description: >
  Use this skill when deriving a Go target architecture from specification
  artifacts. It captures Go-specific architectural best practices: package
  layout, interface design, error handling strategy, stdlib-first selection,
  and testability. Use it to anchor architecture decisions in documented
  Go conventions before writing ARCHITECTURE.md.
---

# Go Architect

## Purpose

Ensure that architectural decisions for a Go reimplementation are grounded
in idiomatic Go conventions and verifiable documentation rather than
framework habits from other ecosystems.

## When to use

- Producing the ARCHITECTURE.md artifact for a Go target
- Selecting packages and frameworks with Context7 backing
- Deciding on package layout, interface boundaries, and error strategy

## Package Layout

Follow the standard Go project layout:

- `cmd/<tool>/main.go` — entry point(s); keep thin, delegate to internal packages
- `internal/` — private packages not importable by external modules
- `internal/<domain>/` — one package per bounded domain concept
- `pkg/` — only for packages explicitly intended for external reuse
- Avoid deep nesting; prefer flat package hierarchies within `internal/`

Do not import across domain boundaries without an explicit interface contract.

## Composition Root (mandatory for CLI tools)

For any target system that is a CLI tool, ARCHITECTURE.md MUST contain a
section titled `## Composition Root` that specifies the end-to-end wiring
of the target binary. This section is the bridge between the per-slice
specifications and the runnable artifact, and is the single source of
truth the Builder follows when implementing `cmd/<tool>/main.go`.

The Composition Root section MUST include:

- **Entry-point path**: explicit file path, normally `cmd/<tool>/main.go`,
  declaring `package main` and a `func main()` that delegates to internal
  packages without containing business logic itself
- **Stage sequence**: the ordered list of internal stages the binary
  executes between CLI invocation and output emission. Each stage MUST
  be referenced back to the functionality slice that owns it. The
  sequence mirrors the `## End-to-End Pipeline` section produced by the
  Analyzer in ANALYSIS.md, expressed in the target language's package
  vocabulary
- **Construction of each stage**: which concrete type implements each
  stage interface, where it is constructed, and which option fields
  configure it. Interfaces remain at the consumer side as required by
  the Interface Design section above; the Composition Root names the
  concrete implementations chosen for the binary
- **Flag-to-stage propagation table**: for every flag, subcommand,
  environment variable, or config-file option declared in the
  CLI-options slice, name the stage that consumes it and the option
  field that carries it from parsed input to that stage. Inputs without
  an entry here are an architectural defect — they will be silently
  dropped at runtime
- **Acceptance reference**: a link to the Golden Output Snippets in
  ANALYSIS.md that the end-to-end binary must reproduce semantically.
  These snippets are the architectural contract; the Builder is
  expected to smoke-test against them before marking BUILD.md complete

A target system whose ARCHITECTURE.md lacks a Composition Root section
is not architecturally complete, regardless of how well the per-package
decisions are documented. The Architect MUST produce this section
explicitly before Phase 3 may begin.

## Interface Design

- Define interfaces at the consumer, not the producer
- Keep interfaces small: one or two methods is better than one large interface
- Name interfaces after their behaviour, not their implementation (`Reader`, not `FileReader`)
- Accept interfaces, return concrete types — unless the return type must be abstracted
- Use `io.Reader`, `io.Writer`, `io.Closer` from stdlib where applicable

## Error Handling

- Always return `error` as the last return value; never panic for recoverable conditions
- Wrap errors with context at each layer: `fmt.Errorf("scanning %s: %w", path, err)`
- Use `errors.Is` and `errors.As` for programmatic error inspection
- Define sentinel errors (`var ErrNotFound = errors.New(...)`) only when callers need to match them
- Do not swallow errors; if an error cannot be handled, propagate it

## Stdlib-First Selection

Before reaching for an external library, verify that stdlib covers the need:
- File I/O: `os`, `io`, `bufio`, `path/filepath`
- JSON: `encoding/json`
- CLI flags: `flag` (for simple cases), or justify a framework like `cobra` with Context7 docs
- Testing: `testing` + `testify/assert` for assertions only when stdlib is insufficient

Use Context7 to verify that any chosen external package is actively maintained
and has stable API surface before committing to it in ARCHITECTURE.md.

**Stdlib-first is not "reimplement-by-hand-first".** Stdlib-first applies to
*generic plumbing* (I/O, JSON, flags). It does **not** license re-deriving
*standardized or non-trivial domain behavior* that the source delegated to a
library (identifier/expression validation, dependency-graph resolution,
canonical formatting/escaping). Hand-rolling such behavior under a stdlib-only
goal is the dominant cause of output divergence in rewrites. For those cases,
run the Source-Dependency Replacement Evaluation below before defaulting to a
hand-written approximation.

## Source-Dependency Replacement Evaluation (mandatory for migrations)

For every entry in the Analyzer's `## Source-Dependency Contracts`, ARCHITECTURE.md
MUST record an explicit decision. Do not silently drop a source dependency into
a hand-written helper.

For each source dependency:

1. **Search for a target-language equivalent via Context7.** Look up candidate
   Go libraries that cover the same behavior and read their documented input
   and output contract.
2. **Compare contracts side by side.** Put the source library's documented I/O
   and worked examples (from ANALYSIS.md) next to the candidate Go library's
   documented I/O and examples. Note where they agree and, critically, where
   they differ (ordering, escaping, sentinels, validation strictness,
   whitespace).
3. **Decide and justify**, choosing exactly one:
   - **Adopt a Go library** when it reproduces the documented contract (or the
     gap is small and explicitly specified away). Pin name + version and cite
     the Context7 evidence.
   - **Faithful reimplementation** when no library matches and the behavior is
     standardized/deep: the spec MUST carry the full rule table and example
     outputs, and the architecture MUST budget for reproducing them exactly —
     not a "common cases" subset.
   - **Approximate** only when the behavior is genuinely trivial and the
     approximation is proven equivalent on the captured examples.
4. **Carry example outputs into the contract.** Whichever option is chosen, the
   decision MUST reference the worked `input → output` examples that the chosen
   implementation has to satisfy, and link them to the acceptance baseline in
   the Composition Root.

Record this as a `## Dependency Decisions` table in ARCHITECTURE.md with
columns: source dependency, behavior relied on, decision, Go library (if any),
contract differences, and the example(s) that pin the expected output.

## Concurrency

- Prefer sequential code unless parallelism is explicitly required by the spec
- Use goroutines + channels for producer-consumer patterns; use `sync.WaitGroup` for fan-out
- Always document goroutine ownership and lifetime
- Use `context.Context` for cancellation propagation; pass it as the first argument

## Testability

- Design packages so that external dependencies are injected as interfaces
- Avoid global state; if unavoidable, document it explicitly
- Table-driven tests (`[]struct{ name, input, want }`) for all boundary conditions
- Integration tests in a separate `test/integration/` package

## Guardrails

- Do not import packages without a documented reason in ARCHITECTURE.md
- Do not choose a framework because it is familiar from another ecosystem
- Do not defer error handling decisions to the Builder — decide at architecture time
- Keep ARCHITECTURE.md concrete enough for Planner and Builder to act on without reopening the spec
- Do not invoke "stdlib-first" to justify re-deriving standardized or non-trivial
  behavior a source library provided, without first running the Source-Dependency
  Replacement Evaluation and recording the decision
- Do not finalize a dependency decision without comparing the source and
  candidate I/O contracts on concrete example outputs
