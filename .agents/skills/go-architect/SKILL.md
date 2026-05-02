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
