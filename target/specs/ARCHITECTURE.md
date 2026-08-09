# ARCHITECTURE — davglass/license-checker Go Rewrite

## Summary

The target is a bounded Go reimplementation of the pinned `davglass/license-checker` CLI/package at commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`, package version `25.0.1`. The implementation shape is a single Go command named `license-checker` with a thin `package main` entry point, internal domain packages for parsing, npm installed-tree scanning, record construction, license classification, policy filtering, rendering, diagnostics, and ordered-output support.

The migration goal is source-compatible behavior for the Phase-1 feature slices F-001 through F-007. Byte identity is required where Analyzer/SPEC artifacts provide exact snippets or source-derived exact strings. Where Analyzer recorded missing full-byte live fixtures, this architecture requires target-side golden tests from committed snippets and source-derived contracts, while reserving live source-vs-target differential execution as the maintainers' out-of-flow manual step.

Context7 was used for Go architecture/package guidance:

- `/golang/go`: supports a `package main` entry point, `internal/` visibility boundaries, `flag.FlagSet`, `encoding/json.MarshalIndent`, standard `os`/`path/filepath` filesystem primitives, and standard `testing`/`go test` workflows.
- `/spdx/tools-golang` and `/git-pkgs/spdx`: evaluated for SPDX replacement behavior. `spdx/tools-golang` is document/SBOM-oriented and not sufficient for `spdx-satisfies`; `git-pkgs/spdx` documents parse/validate/satisfy behavior but has module availability risk, so the implementation should use a real Go module with equivalent expression APIs behind an adapter, selected below.

## Technical Context

| Dimension | Decision |
|---|---|
| Language | Go, target module rooted under `target/`. |
| CLI framework | No full CLI framework. Use a source-compatible internal parser, with `flag.FlagSet` considered only as a helper where it does not conflict with `nopt` behavior. |
| Composition style | `target/cmd/license-checker/main.go` delegates to `internal/app.Run`; business logic lives under `internal/`. |
| Storage | Filesystem only: npm project tree under `--start`; optional `--customPath`; optional `--out`; optional `--files`. No database. |
| Ordered data | Do not use plain Go maps for observable JSON/tree/CSV traversal. Use explicit ordered record/map types with custom JSON/tree rendering. |
| SPDX | Use `internal/spdxcompat` facade backed by `github.com/github/go-spdx/v2/spdxexp` pinned to `v2.7.0`, plus compatibility code for `spdx-correct`-style correction decisions not covered by the library. |
| Testing | Standard `testing` package, table-driven unit tests, CLI integration tests using built binary/`os/exec`, golden fixtures under target testdata derived from `ANALYSIS.md` snippets and each `SPEC.md`. Run with `go test ./...`. |
| Target platform | POSIX-like CLI on darwin/linux first; Windows path behavior must be explicit in tests where `filepath` changes separators. |
| Constraints | No live source-system execution in-flow. Preserve exact source literals, field order, renderer newline distinctions, and stderr/stdout routing from Phase-1 artifacts. |

## Structure Decision

Use a standard Go project layout with a thin command and private internal packages:

```text
target/
  go.mod
  cmd/
    license-checker/
      main.go
  internal/
    app/
    cli/
    color/
    customformat/
    debuglog/
    diagnostics/
    files/
    license/
    npmgraph/
    ordered/
    policy/
    records/
    render/
    spdxcompat/
```

Boundaries and interfaces:

- `cmd/license-checker`: owns only `func main()`, process args/env/stdout/stderr wiring, and `os.Exit(code)`.
- `internal/app`: composition and orchestration. Exposes `Run(ctx context.Context, inv Invocation) int`.
- `internal/cli`: source-compatible argument parsing/defaults and preflight. It returns a concrete `Options` value plus a `PreflightResult`.
- `internal/npmgraph`: defines the consumer-side `Scanner` interface and a filesystem implementation that reconstructs `read-installed`-compatible installed npm tree nodes.
- `internal/records`: flattens `npmgraph.Node` values into ordered package records, preserving source insertion order and sentinels.
- `internal/license`: exact classifier, license-file precedence, README/license/NOTICE field extraction.
- `internal/spdxcompat`: small adapter for valid SPDX expression parsing/correction/satisfaction semantics used by `license` and `policy`.
- `internal/policy`: applies `--unknown`, `--onlyunknown`, `--exclude`, package filters, private exclusion, `--failOn`, and `--onlyAllow` over ordered records.
- `internal/render`: tree, JSON, CSV, Markdown, summary, `--out`, and `--files` behavior. It depends on `ordered` and consumes interfaces for filesystem/writers where useful.
- `internal/diagnostics`: exact messages, stderr channel selection, and exit-code decisions.

Interfaces are defined at the consuming package. Example contracts:

```go
// internal/app
type Scanner interface {
    Scan(ctx context.Context, root string, opts npmgraph.Options) (*npmgraph.Node, error)
}

// internal/render
type FileSystem interface {
    ReadFile(name string) ([]byte, error)
    WriteFile(name string, data []byte, perm fs.FileMode) error
    MkdirAll(path string, perm fs.FileMode) error
    Stat(name string) (fs.FileInfo, error)
}
```

Errors return as the final value, are wrapped with context using `fmt.Errorf("...: %w", err)`, and are converted to source-compatible diagnostics only at `internal/app`/`internal/diagnostics` boundaries.

## Composition Root

### Entry point

- Path: `target/cmd/license-checker/main.go`
- Package: `package main`
- Responsibility: build an `app.Invocation` from `os.Args[1:]`, `os.Environ()`, `os.Stdout`, `os.Stderr`, `os.Getwd`, and `os.Exit(app.Run(...))`. No business logic belongs in `main.go`.

### Ordered stage sequence

1. `cli.Parse(argv, cwd, color.Supports)` — F-001/F-002; source-compatible known flags, aliases, defaults, `--no-color` presence tracking, `direct` normalization.
2. `diagnostics.Preflight(options, stderr)` — F-001/F-007; help, version, `failOn`/`onlyAllow` conflict, comma warning.
3. `customformat.Load(options.CustomPath)` — F-002; attach `CustomFormat` before scanning.
4. `npmgraph.Scanner.Scan(ctx, options.Start, npmgraph.Options{Dev, Depth, Log})` — F-003; reproduce `read-installed`-compatible tree semantics.
5. `records.Flatten(tree, records.Options{...})` — F-003; create ordered `name@version` records, metadata, repository normalization, path fields.
6. `license.Apply(record, node, license.Options{...})` during flatten — F-004; package/readme/license-file classifier, file precedence, NOTICE fields.
7. `policy.Apply(records, policy.Options{...}, stderr)` — F-005/F-007; sorted pass, filters, restrictions, fail policies, exact diagnostics/exits.
8. `render.Format(records, render.Options{...})` — F-006; renderer precedence: JSON, CSV, Markdown, summary, else tree.
9. `render.Emit(records, formatted, render.SideEffectOptions{Files, Out}, stdout, stderr, fs)` — F-006/F-007; `--files`, `--out`, stdout, missing-file warnings.
10. Return `diagnostics.ExitCode` — F-007; `0` normal/help, `1` version/policy fatal paths, no invented non-zero exit solely for `checker.init` callback-equivalent errors.

### Concrete construction

`app.NewDefault()` constructs:

- `cli.Parser{}`
- `color.ChalkCompatDetector{Env: os.Environ, StdoutFD: os.Stdout.Fd}`
- `debuglog.Namespaced{Env: os.Environ, Namespaces: []string{"license-checker:log", "license-checker:error"}}`
- `npmgraph.FileSystemScanner{FS: os.DirFS("/"), ReadFile: os.ReadFile, ReadDir: os.ReadDir}`
- `records.Flattener{License: license.Detector{SPDX: spdxcompat.Adapter{}}}`
- `policy.Engine{SPDX: spdxcompat.Adapter{}}`
- `render.Renderer{FS: files.OS{}}`
- `diagnostics.Writer{Stderr: invocation.Stderr}`

### Flag-to-stage propagation table

| Input | Parsed option field | Consuming stage/package | Propagated option/value |
|---|---|---|---|
| `--production` | `Options.Production` | `npmgraph`, `records` | provider `dev=false`; omit `extraneous` nodes. |
| `--development` | `Options.Development` | `npmgraph`, `records` | provider `dev=false`; keep only `extraneous` or `root`. |
| `--json` | `Options.JSON` | `cli`, `render` | disables color default; selects JSON formatter. |
| `--csv` | `Options.CSV` | `cli`, `render` | disables color default; selects CSV formatter. |
| `--csvComponentPrefix` | `Options.CSVComponentPrefix` | `render` | optional leading CSV `component` column/value. |
| `--markdown` | `Options.Markdown` | `cli`, `render` | disables color default; selects Markdown formatter. |
| `--summary` | `Options.Summary` | `render` | selects summary after JSON/CSV/Markdown. |
| `--out` | `Options.Out` | `render.Emit`, `files` | mkdir parent; write raw formatted string; disables color eligibility. |
| `--files` | `Options.Files` | `render.Emit`, `files` | bypass formatted output; copy `<module>-LICENSE.txt`. |
| `--start` | `Options.Start` | `npmgraph` | scan root; default current working directory. |
| `--unknown` | `Options.Unknown` | `records`, `policy` | add `dependencyPath`; rewrite guessed `*` licenses to `UNKNOWN`. |
| `--onlyunknown` | `Options.OnlyUnknown` | `policy` | keep only current license containing `*` or `UNKNOWN`. |
| `--exclude` | `Options.Exclude` | `policy`, `spdxcompat` | escaped-comma split, BSD alias, SPDX/literal exclusion. |
| `--failOn` | `Options.FailOn` | `diagnostics.Preflight`, `policy` | conflict/warning; semicolon exact-match fail policy. |
| `--onlyAllow` | `Options.OnlyAllow` | `diagnostics.Preflight`, `policy` | conflict/warning; semicolon substring allow policy. |
| `--packages` | `Options.Packages` | `policy` | semicolon exact package whitelist. |
| `--excludePackages` | `Options.ExcludePackages` | `policy` | semicolon exact package blacklist built from filtered map. |
| `--excludePrivatePackages` | `Options.ExcludePrivatePackages` | `policy` | delete truthy-private records after restrictions. |
| `--relativeLicensePath` | `Options.RelativeLicensePath` | `license`, `records`, `render` | convert license/notice paths relative to scan root. |
| `--customPath` | `Options.CustomPath` | `customformat` | parse UTF-8 JSON and assign `CustomFormat`. |
| programmatic `customFormat` | `Options.CustomFormat` | `records`, `render` | custom fields/defaults/order. |
| `--direct` | `Options.DirectDepth` | `npmgraph` | `0` when present/truthy; unbounded sentinel otherwise. |
| `--color` / `--no-color` | `Options.Color`, `Options.ColorExplicit` | `cli`, `color`, `render` | default from color detector unless explicit; tree key/sentinel ANSI only when eligible. |
| `--help`, `-h` | `Options.Help` | `diagnostics.Preflight` | help stderr and exit `0`. |
| `--version`, `-v` | `Options.Version` | `diagnostics.Preflight` | `25.0.1` stderr and exit `1`. |
| `DEBUG=license-checker*` | `Invocation.Env` | `debuglog`, `npmgraph`, `app` | enable `license-checker:log` and `license-checker:error`; no exit impact. |

### Acceptance reference

End-to-end acceptance is anchored to `target/specs/ANALYSIS.md` **Golden Output Snippets** lines 68-127: default tree excerpt, guessed license marker, CSV snippets, Markdown snippets, and diagnostics. Builder must create committed target golden tests from these snippets and from exact strings in the individual `SPEC.md` files before claiming completion.

## Package Decisions

| Package/library | Decision | Reason | Alternatives considered | Supporting documentation/reference |
|---|---|---|---|---|
| `cmd/license-checker` / `package main` | Use | Source exposes one executable named `license-checker`; Go command entry point maps directly. | Multiple subcommands rejected: source has none. | Context7 `/golang/go` package-main examples. |
| `internal/` layout | Use | Prevents leaking migration internals as public API and matches Go visibility boundaries. | `pkg/` rejected because no reusable public API is in scope. | Context7 `/golang/go` internal package visibility examples. |
| `context` | Use | First parameter for scanner/app cancellation boundaries and future test control. | Global cancellation rejected. | Go conventions from `/golang/go`; go-architect skill. |
| `flag` | Do not use as the primary parser | `flag.FlagSet` is documented and stable, but source uses `nopt` behavior: aliases, path type, `--no-color`, metadata removal, unknown/malformed edge risks. | Cobra/urfave rejected as unnecessary and likely divergent; custom parser selected. | Context7 `/golang/go` documents `flag.FlagSet.Parse`, `Bool`, `String`, `Visit`; F-002 requires `nopt` compatibility. |
| `encoding/json` | Use with custom ordered marshaler | `MarshalIndent` gives two-space formatting, but map key ordering/struct schema alone cannot preserve JavaScript insertion order. Use `ordered.Object.MarshalJSON` and then indent. | Plain `map[string]any` rejected for observable order. | Context7 `/golang/go` `MarshalIndent` docs; F-006 JSON contract. |
| `os`, `io/fs`, `path/filepath` | Use | Covers read/write/stat/mkdir and platform path joins; `os.MkdirAll` replaces `mkdirp`. | External mkdir package rejected: stdlib is sufficient. | Context7 `/golang/go`; F-006 mkdirp behavior is recursive directory creation only. |
| `regexp`, `strings`, `sort`, `strconv` | Use | Needed for exact classifier regex order, list splitting/trimming, lexical key sort, summary counts. | External text libraries rejected. | Stdlib-first rule; F-004/F-005 transformation tables. |
| `errors`, `fmt` | Use | Explicit wrapping and source-compatible conversion at diagnostics boundary. | Panics rejected for recoverable errors. | Go error conventions from skill and `/golang/go` docs. |
| `testing`, `os/exec` | Use | Table-driven unit tests and CLI process tests are sufficient; no assertion framework required. | Testify rejected unless later tests become unreadable. | Context7 `/golang/go` testing examples and `go test` workflow. |
| `github.com/github/go-spdx/v2/spdxexp v2.7.0` | Use behind `internal/spdxcompat` | Provides real Go SPDX expression validation/satisfaction APIs, closer to `spdx-expression-parse`/`spdx-satisfies` than document-oriented SPDX tools. Must be wrapped because source also uses `spdx-correct` semantics. | `/spdx/tools-golang` rejected: SBOM/document-oriented; `/git-pkgs/spdx` not selected because module availability/version risk despite Context7 docs. Hand-rolled SPDX rejected as standardized/deep behavior. | Context7 SPDX evaluation plus pkg.go.dev/web lookup for `github.com/github/go-spdx/v2` v2.7.0. |

## Dependency Decisions

| Source dependency | Behavior relied on | Decision | Go library/package | Contract differences | Examples pinning expected output |
|---|---|---|---|---|---|
| `nopt ^4.0.1` | Known option typing, aliases `-v`/`-h`, path-typed options, `--no-` cooked presence, `argv` metadata deletion. | Faithful internal reimplementation for declared flags. | `internal/cli`; optional `flag.FlagSet` helpers only if invisible. | Go `flag` does not match `nopt` edge behavior; parser edge gaps remain conformance risks. | F-002 defaults/direct/customPath tests; help/version preflight. |
| `read-installed ~4.0.3` | Installed npm dependency graph with `extraneous`, `root`, `path`, `readme`, nested `dependencies`, dedupe behavior. | Faithful reimplementation, not naive directory walk. | `internal/npmgraph` using `os`/`filepath`. | Exact npm tree semantics are high risk; package-lock behavior is not source-pinned. | `abbrev@1.0.9` CSV/Markdown snippets; production/development/direct fixtures. |
| `spdx-expression-parse ^3.0.0` | Valid SPDX expressions pass through unchanged; invalid inputs fall through to regex classifier. | Adopt adapter over Go SPDX library plus pass-through rule. | `internal/spdxcompat` over `github.com/github/go-spdx/v2/spdxexp v2.7.0`. | Go library may normalize expressions; classifier must return the original raw string on success. | F-004 examples: `MIT`, `LGPL-2.0`, `(GPL-2.0+ WITH Bison-exception-2.2)`, `Apache-2.0 OR ISC OR MIT`. |
| `spdx-correct ^3.0.0` | Determines valid exclusions when `spdxCorrect(x) === x`; corrects candidate licenses before satisfaction. | Compatibility wrapper with explicit table, using SPDX library for validation plus source examples for correction behavior. | `internal/spdxcompat`. | Go library validation is not identical to JS correction; unsupported corrections must be fixture-driven, not guessed. | F-005 examples: MIT/ISC, BSD alias, Public Domain, escaped comma Apache literal. |
| `spdx-satisfies ^4.0.0` | Checks candidate license satisfaction against constructed OR expression. | Adapter over Go SPDX satisfaction API. | `internal/spdxcompat` with tests mapping source argument order. | API direction may differ (`expression` vs allowed licenses); wrapper owns translation. | Exclude MIT/ISC, BSD-3-Clause under `BSD` alias. |
| `treeify ^1.1.0` | `asTree(obj, true)` glyphs, indentation, traversal order for default and summary trees. | Faithful internal renderer. | `internal/render/tree.go`. | No external Go treeify chosen; byte identity must be snippet-tested. | README tree excerpt and guessed `MIT*` snippet in `ANALYSIS.md`. |
| `chalk ^2.4.1` | `supportsColor`; blue/dim/green key styling; bold red sentinels. | Small internal color adapter; no-color baseline exact, forced-color tests for ANSI roles. | `internal/color`. | Terminal auto-detection may not exactly match chalk; color rows remain compatible unless test pins env. | F-006 color eligibility and no-color snippets. |
| `debug ^3.1.0` | `DEBUG=license-checker*` opt-in namespaces. | Internal namespace gate only; do not emulate byte debug formatting. | `internal/debuglog`. | Debug timestamp/color formatting is not captured; no exit impact. | F-007 debug namespace requirements. |
| `mkdirp ^0.5.1` | Recursive directory creation before `--out`/`--files`. | Replace with stdlib. | `os.MkdirAll`. | None for required behavior: existing dirs succeed. | F-006 `--out`, `--files` side-effect tests. |

## Implementation Order

1. Create `target/go.mod`, `cmd/license-checker/main.go`, `internal/app`, and shared `Options`, `ExitCode`, ordered data types.
2. Implement `internal/diagnostics` exact literals and `internal/cli` parser/defaults/preflight tests for F-001/F-002/F-007.
3. Implement `internal/customformat` and ordered object/record JSON support.
4. Implement `internal/spdxcompat` wrapper and lock its behavior with F-004/F-005 SPDX examples before integrating classifier/policy.
5. Implement `internal/license` classifier, license-file precedence, license/notice extraction tests.
6. Implement `internal/npmgraph` scanner with fixtures for installed npm tree semantics, then `internal/records` flattening/order/repository/path/custom field behavior.
7. Implement `internal/policy` filters/restrictions/fail policies and exact stderr/exit tests.
8. Implement `internal/render` tree/JSON/CSV/Markdown/summary and side-effect outputs, including newline and no-general-CSV-escaping behavior.
9. Wire `internal/app.Run` composition and CLI integration tests for each conformance-matrix row.
10. Run `go test ./...`; record any remaining non-identical rows as explicit conformance limitations before Phase 3 completion.

## Conformance Matrix

Status values use the repository contract: `identical`, `compatible`, or `out-of-scope`. `identical` means byte/function parity is required for the recorded/source-derived contract. It does not imply live full-byte parity where Analyzer explicitly recorded missing full-byte fixtures; those rows prescribe golden tests from committed snippets/source-derived strings.

| Recorded mode/flag row | Status | Required acceptance/golden coverage | Justification / limitation |
|---|---|---|---|
| Default invocation / tree | identical | Golden tests from `ANALYSIS.md` default tree and guessed `MIT*` snippets; source-derived tree renderer fixtures. | Full live stdout absent, but tree glyph/field order contract is source/snippet-backed. |
| `--production` | identical | Scanner/flatten fixtures for `extraneous` pruning and rendered key absence. | Functional parity required; no standalone byte output fixture. |
| `--development` | identical | Scanner/flatten fixtures for non-`extraneous`, non-`root` pruning. | Functional parity required; no standalone byte output fixture. |
| `--json` | identical | Ordered-record JSON golden using two-space indentation, field order, and CLI/file newline rules. | Full live JSON stdout absent; builder must derive fixture from SPEC contracts. |
| `--csv` | identical | `"module name","license","repository"` and `abbrev@1.0.9` CSV snippets. | Byte snippet exists. |
| `--csvComponentPrefix` | identical | Component-prefix CSV snippet with `main-module`. | Byte snippet exists. |
| `--markdown` | identical | Default Markdown and custom package-row snippets. | Byte snippets exist. |
| `--summary` | identical | Source-derived summary count/sort/tree fixtures; avoid deterministic equal-count ties. | Full live summary bytes absent; tie ordering not claimed. |
| `--out` | identical | File-content tests for raw formatted string and no stdout-added newline. | Full live file fixture absent; source algorithm explicit. |
| `--files` | identical | `foo-LICENSE.txt` naming and missing-file warning golden. | Full copied contents fixture limited; behavior source/snippet-backed. |
| `--start` | identical | Scanner root/path propagation tests. | Functional contract; renderer bytes depend on fixture. |
| `--unknown` | identical | `dependencyPath` presence and guessed `*` -> `UNKNOWN` tests. | Byte output through renderers must use existing snippets/source-derived fixtures. |
| `--onlyunknown` | identical | Keeps only `*` or `UNKNOWN`; package map/key tests. | Functional parity required. |
| `--exclude` | identical | MIT/ISC, escaped comma Apache, BSD alias, Public Domain, Custom no-match fixtures. | SPDX adapter must match examples. |
| `--failOn` | identical | Exact warning and fail stderr snippets; exit `1`. | Byte diagnostic strings captured/source-derived. |
| `--onlyAllow` | identical | Exact warning and violation stderr snippets; exit `1`. | Byte diagnostic strings captured/source-derived. |
| `--packages` | identical | Exact semicolon split/no-trim package key-set fixtures. | Functional parity required. |
| `--excludePackages` | identical | Exact blacklist rebuild from filtered map fixtures. | Functional parity required. |
| `--excludePrivatePackages` | identical | Private fixture: `UNLICENSED` then deletion/empty key set. | Functional parity required. |
| `--relativeLicensePath` | identical | License/notice path relative-to-root fixtures. | No dedicated NOTICE source test; target must use source-derived constructed fixture. |
| `--customPath` | identical | Valid JSON, invalid/missing/non-string path tests; custom CSV snippets. | Parse-error downstream bytes not fully captured; parser behavior required. |
| Programmatic `customFormat` | identical | Field order/default/suppression tests and CSV/Markdown snippets. | Byte snippets exist for common custom format. |
| `--direct` | identical | Provider depth `0` vs unbounded sentinel tests. | Functional parity required. |
| `--color` / implicit color support | compatible | No-color byte baseline exact; forced-color unit tests for ANSI role placement only. | Chalk terminal detection and exact ANSI bytes are terminal-dependent and not live-captured. |
| `--help`, `-h` | identical | Exact help lines and exit `0`; final Node two-argument trailing whitespace remains documented limitation. | Source-derived help text exists; full trailing byte capture absent. |
| `--version`, `-v` | identical | Exact `25.0.1\n` on stderr, no stdout, exit `1`. | Exact string captured/source-derived. |
| `DEBUG=license-checker*` | compatible | Namespace enablement tests; no exit/output contract changes. | Debug dependency formatting not byte-captured. |
| Init callback error object formatting | compatible | Must emit `Found error`; second error-object line requires maintainers-approved baseline before byte identity. | Analyzer explicitly lacks Node `console.error(err)` bytes. |
| Programmatic API beyond `init` and renderer helpers | compatible | Expose internal/app or library facade only if Planner keeps it in bounded scope. | Rewrite is CLI-oriented; full npm package API surface is lower-priority UC-5 and not requested as public Go API. |

## Risks

| Risk | Severity | Builder constraint / mitigation |
|---|---:|---|
| `read-installed` semantics differ from a normal filesystem walk. | High | Implement `npmgraph` with compatibility fixtures for `extraneous`, `root`, nested dependencies, path/readme propagation, duplicate key guard. |
| SPDX behavior is standardized/deep and JS dependencies are not identical to Go APIs. | High | Keep all SPDX calls behind `spdxcompat`; pass F-004/F-005 examples before use; do not hand-roll broad SPDX parsing. |
| Go map ordering would break JSON/tree/CSV parity. | High | Use `ordered.Map`/`ordered.Record` everywhere observable; forbid plain maps at renderer boundaries. |
| Missing full live byte fixtures for several modes. | High | Golden-test committed snippets and source-derived contracts; record limitations; leave live differential run manual/out-of-flow. |
| Color/debug output is terminal/environment dependent. | Medium | No-color is the exact baseline; forced-color/debug tests only assert documented role/namespace behavior. |
| Node error-object formatting is not captured. | Medium | Emit `Found error` exactly; keep second line compatible until maintainers provide a baseline. |
| Summary equal-count order is unspecified. | Medium | Sort by descending count only; do not add lexical tie sorting unless declared non-identical. |
| Scoped package key color split on `@` is source-bug-compatible. | Low/Medium | Preserve first-two-parts behavior; do not “fix” scoped package formatting silently. |
| Windows paths may differ from Node/POSIX snippets. | Low/Medium | Treat darwin/linux as primary parity platform; isolate path separator tests and document Windows compatibility if added. |
