---
name: go-craftsman
description: >
  Use this skill when implementing Go target code from TASKS.md, SPEC.md,
  and ARCHITECTURE.md. It provides idiomatic Go implementation patterns,
  naming conventions, and efficiency guidelines. Use it to ensure generated
  code is not just correct but maintainable, readable, and idiomatic — not
  a mechanical translation from the source language.
---

# Go Craftsman

## Purpose

Ensure that every Go file produced by the Builder is idiomatic, efficient,
and maintainable — not a syntactic translation from the source ecosystem
but a proper Go reimplementation grounded in Go conventions.

## When to use

- Implementing any task from TASKS.md that produces Go source files
- Reviewing generated code before marking a task complete in TASKS.md
- Deciding between equivalent implementation approaches

## Naming Conventions

- Package names: short, lowercase, no underscores (`scanner`, not `license_scanner`)
- Exported names: clear and self-documenting without the package prefix
  (`scanner.Detect`, not `scanner.ScannerDetect`)
- Unexported names: short but clear within the package context
- Receiver names: one or two letters matching the type (`s` for `Scanner`)
- Error variables: `ErrXxx` for sentinel errors
- Interface names: `-er` suffix for single-method interfaces (`Resolver`, `Detector`)

## Idiomatic Patterns

### Error handling
```go
// Do this
result, err := scan(path)
if err != nil {
    return fmt.Errorf("scan %s: %w", path, err)
}

// Not this
result, _ := scan(path)
```

### Early return over nesting
```go
// Do this
if err != nil {
    return err
}
// continue happy path

// Not this
if err == nil {
    // deeply nested happy path
}
```

### Defer for cleanup
```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close()
```

### Table-driven tests
```go
tests := []struct {
    name  string
    input string
    want  string
}{
    {"mit", "MIT", "MIT"},
    {"unknown", "", "UNKNOWN"},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := normalize(tt.input)
        assert.Equal(t, tt.want, got)
    })
}
```

## Efficiency Guidelines

- Preallocate slices when length is known: `make([]T, 0, n)`
- Use `strings.Builder` for string concatenation in loops
- Prefer `[]byte` over `string` for I/O-heavy code paths
- Use `bufio.Scanner` or `bufio.Reader` for line-by-line file reading
- Avoid unnecessary allocations in hot paths; profile before optimizing

## CLI-Specific Patterns

- Use `os.Exit` only in `main()`; return errors from all other functions
- Write errors to `os.Stderr`, results to `os.Stdout`
- Respect exit codes: `0` for success, `1` for usage errors, `2` for runtime errors
- Parse flags before any I/O; validate all required flags upfront

## Completeness Check Before Marking Done

Before checking off a task in TASKS.md, verify:

- [ ] Error from every function call is checked and wrapped with context
- [ ] No global mutable state introduced without explicit justification
- [ ] Test file exists alongside every new package file
- [ ] Exported types and functions have doc comments
- [ ] Code compiles and tests pass (`go build ./...` and `go test ./...`)
- [ ] No linter warnings (`go vet ./...`)

## Guardrails

- Do not transliterate source language patterns into Go — reimpliment idiomatically
- Do not skip error handling to keep code shorter
- Do not add dependencies not listed in ARCHITECTURE.md without flagging it as a blocker
- Do not mark a task complete if the corresponding acceptance scenario from SPEC.md
  is not demonstrably satisfied
