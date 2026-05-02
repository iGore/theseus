# ANALYSIS.md — `github.com/davglass/license-checker` (master)

## Scope

- **Bounded module:** Node.js CLI/package scanner in `github.com/davglass/license-checker` (default branch `master`), focused on:
  - CLI entrypoint and argument model
  - dependency traversal and flattening
  - license detection/normalization
  - filtering/policy enforcement
  - output formatting/export behavior
- **Out of scope:** full ecosystem architecture beyond this repo; downstream rewrite implementation.

## Entry Points

- **CLI entrypoint:** `bin/license-checker` (shebang + orchestration) (`bin/license-checker:1-110`)
- **Programmatic API entrypoint:** `checker.init(options, callback)` (`lib/index.js:261-469`)
- **Argument parser entrypoint:** `parse(args)` (`lib/args.js:88-91`)

## Inputs and Outputs

### Inputs

- **CLI flags / args** via `nopt` (`lib/args.js:7-35`), including:
  - dependency scope: `production`, `development`, `direct`, `start`
  - output mode: `json`, `csv`, `markdown`, `summary`, `out`, `files`, `color`
  - license filtering/policy: `unknown`, `onlyunknown`, `exclude`, `failOn`, `onlyAllow`
  - package filtering: `packages`, `excludePackages`, `excludePrivatePackages`
  - formatting config: `customPath`, `customFormat`, `relativeLicensePath`
  - meta: `help`, `version`
- **Filesystem inputs:** project dependency tree under `start`; package metadata (`package.json`), `README.md`, license/notice files (`lib/index.js:114-160`, `162-232`).
- **Config JSON file:** custom format file path (`--customPath`) parsed by `parseJson` (`lib/index.js:264-266`, `593-608`).
- **Environment variable:** `DEBUG=license-checker*` enables debug logging (README `162-170`; logger setup `lib/index.js:23-26`).

### Outputs

- **In-memory result object:** `restricted` map keyed as `name@version` with module metadata (`lib/index.js:31`, `404-467`).
- **CLI textual output modes:**
  - tree (`asTree`) (`lib/index.js:475-477`)
  - JSON (stringified in CLI) (`bin/license-checker:84-86`)
  - CSV (`asCSV`) (`lib/index.js:508-562`)
  - Markdown (`asMarkDown`) (`lib/index.js:571-591`)
  - Summary (`asSummary`) (`lib/index.js:479-506`)
- **Filesystem side effects:**
  - write formatted output to `--out` path (`bin/license-checker:98-102`)
  - write per-module license file copies to `--files` dir (`lib/index.js:610-627`)
- **Process behavior:**
  - `--help` exits `0` (`bin/license-checker:17-47`)
  - `--version` prints version and exits `1` (`bin/license-checker:49-52`)
  - `--failOn`/`--onlyAllow` policy violations exit `1` (`lib/index.js:437-457`)
  - invalid combo `--failOn` + `--onlyAllow` exits `1` (`bin/license-checker:54-57`)

## Function Inventory

| Function | File | Role | Functionality ID |
|---|---|---|---|
| `parse(args)` | `lib/args.js:88-91` | Parse CLI/options and apply defaults | F-001 |
| `setDefaults(parsed)` | `lib/args.js:65-86` | Default `start`, `color`, `direct`, `relativeLicensePath` | F-001 |
| `raw(args)` / `clean(args)` / `has(a)` | `lib/args.js:41-63` | Lower-level nopt handling and option detection | F-001 |
| CLI main callback + `shouldColorizeOutput(args)` | `bin/license-checker:65-109` | Invoke scan and choose output/write path | F-006 |
| `flatten(options)` | `lib/index.js:27-259` | Recursive dependency flatten + metadata/license collection | F-003, F-004 |
| internal `include(property)` | `lib/index.js:56-58` | Respect custom format property exclusion | F-007 |
| `init(options, callback)` | `lib/index.js:261-469` | Main orchestration: read tree, filter, enforce policy | F-002, F-005 |
| `asTree(sorted)` / `print(sorted)` | `lib/index.js:471-477` | Tree output | F-006 |
| `asSummary(sorted)` | `lib/index.js:479-506` | Count grouped by license + tree output | F-006 |
| `asCSV(sorted, customFormat, csvComponentPrefix)` | `lib/index.js:508-562` | CSV export (default/custom schemas) | F-006 |
| `asMarkDown(sorted, customFormat)` | `lib/index.js:571-591` | Markdown export | F-006 |
| `parseJson(jsonPath)` | `lib/index.js:593-608` | Safe JSON parse returning object or Error | F-007 |
| `asFiles(json, outDir)` | `lib/index.js:610-627` | Export detected license files into output dir | F-006 |
| `module.exports(str)` | `lib/license.js:22-83` | Normalize/deduce license strings and SPDX expressions | F-004 |
| `module.exports(dirFiles)` | `lib/license-files.js:14-29` | License filename detection in precedence order | F-004 |
| `Stack` (`add/test/done`) | `lib/stack.js:6-44` | Async stack utility; no observed usage in repo | F-008 |

## Feature Slices

### F-001 — CLI options model and argument normalization
- **What it does:** Parses known flags, sets defaults, enforces early CLI guardrails.
- **Main evidence:** `lib/args.js:7-99`, `bin/license-checker:17-63`.
- **Behavior:**
  - default scan path = current working dir (`lib/args.js:76`)
  - `direct` toggles traversal depth (`true => 0`, else `Infinity`) (`lib/args.js:79-83`)
  - disables color when json/csv/markdown selected (`lib/args.js:73-75`)
  - warns when `failOn/onlyAllow` use commas, expects semicolons (`bin/license-checker:58-63`)

### F-002 — Dependency graph loading and traversal control
- **What it does:** Calls `read-installed` to load dependency graph from start path with depth/dev rules.
- **Main evidence:** `lib/index.js:267-276`, `298-309`; README explanation (`180-182`).
- **Behavior:**
  - `opts.depth` driven by normalized `direct`
  - `opts.dev` false when production/development filtering requested

### F-003 — Flattening dependency graph into package map
- **What it does:** Converts nested tree to a unique `name@version` map with metadata fields.
- **Main evidence:** `lib/index.js:27-259`.
- **Behavior:**
  - circular guard via `if (data[key]) return data` (`41-47`)
  - applies prod/dev gate (`49-51`)
  - extracts repository/author/url/path/private markers (`60-110`, `37-39`)
  - recursively traverses `json.dependencies` (`235-254`)

### F-004 — License detection and normalization pipeline
- **What it does:** Resolves license from package fields, README, then license files, using SPDX-aware and heuristic parser.
- **Main evidence:** `lib/index.js:112-172`, `150-160`; `lib/license.js:22-83`; `lib/license-files.js:3-29`.
- **Behavior:**
  - supports `license` / `licenses` array/object/string (`112-137`)
  - fallback to README text (`139-141`)
  - scans files in precedence `LICENSE*`, `LICENCE*`, `COPYING`, `README` (`license-files.js`)
  - marks inferred results with `*` (heuristic parser returns e.g., `MIT*`)
  - handles custom URL/file refs (`Custom: ...`) (`license.js:76-79`)

### F-005 — Filtering, restriction, and policy failure behavior
- **What it does:** Applies license/package/private filters and policy exits (`failOn`, `onlyAllow`).
- **Main evidence:** `lib/index.js:313-457`; tests `tests/test.js:137-288`, `tests/packages-test.js:8-45`, `tests/failOn-test.js:7-41`.
- **Behavior:**
  - `unknown` converts guessed (`*`) to `UNKNOWN` (`331-339`)
  - `onlyunknown` keeps only guessed/unknown licenses (`342-347`)
  - `exclude` supports escaped commas and SPDX matching (`313-399`)
  - package include/exclude via semicolon-separated keys (`407-426`)
  - `excludePrivatePackages` removes `private` modules (`428-435`)
  - exits process on `failOn` or disallowed `onlyAllow` match (`438-457`)

### F-006 — Output rendering and persistence
- **What it does:** Produces tree/json/csv/markdown/summary output and writes to stdout/file(s).
- **Main evidence:** `bin/license-checker:84-104`; `lib/index.js:475-562`, `571-591`, `610-627`.
- **Behavior:**
  - mode priority in CLI: json > csv > markdown > summary > tree
  - optional colorized package key when interactive and non-structured output (`74-83`, `107-109`)
  - `--out` creates parent dirs before writing (`99-102`)
  - `--files` copies detected license files per module (`610-627`)

### F-007 — Custom format/config model
- **What it does:** Adds user-defined metadata schema from inline object or JSON file.
- **Main evidence:** `lib/index.js:95-106`, `181-191`, `193-221`, `264-266`, `593-608`; tests `tests/test.js:429-478`, `690-716`.
- **Behavior:**
  - missing fields can be filled from default values in `customFormat`
  - supports keys such as `licenseText`, `licenseFile`, and hidden `copyright`
  - JSON parse failures return `Error` from `parseJson`

### F-008 — Ancillary utility module (`Stack`)
- **What it does:** Generic async aggregation utility.
- **Main evidence:** `lib/stack.js:6-44`; no references found via search.
- **Behavior:** Present but not wired into current scanner flow.

## Domain Map

- **Primary entities**
  - `moduleInfo`: per-package output object (`licenses`, repository, author-derived fields, path, etc.)
  - `sorted/filtered/restricted`: successive filtered views of package map (`lib/index.js:311-467`)
  - `customFormat`: schema/default-value contract for extra output fields (`95-106`, `513-520`)

- **Core transformations**
  1. dependency tree (`read-installed`) → flattened map
  2. raw license signals (SPDX/text/file) → normalized license label
  3. normalized map → constrained map (`exclude`, packages, private, fail/allow gates)
  4. constrained map → output serializer (tree/json/csv/md/summary/files)

## Flow Map

1. CLI parses args and applies defaults (`lib/args.js:88-91`, `65-86`).
2. CLI validates incompatible flags and semicolon warning (`bin/license-checker:54-63`).
3. `init` builds `read-installed` opts and policy lists (`lib/index.js:267-297`).
4. `read-installed` returns dependency tree (`298`).
5. `flatten` recursively extracts metadata and license evidence (`299-309`, `27-259`).
6. Post-processing applies unknown/onlyunknown/exclude/package/private filters (`323-435`).
7. Policy checks may hard exit (`438-457`).
8. Callback returns results; CLI formats and prints/writes (`467`, `84-104`).

## Dependency Notes

Important runtime dependencies for this bounded module:

- `read-installed` (`package.json:64`, `lib/index.js:11`): loads node dependency tree from `start` path.
- `nopt` (`package.json:63`, `lib/args.js:7`): CLI option parsing and typing.
- `spdx-expression-parse` (`package.json:67`, `lib/license.js:1`): validates SPDX expressions directly.
- `spdx-correct` + `spdx-satisfies` (`package.json:66,68`, `lib/index.js:18-19`, `362-391`): normalize and evaluate SPDX license filters (`--exclude`).
- `treeify` (`package.json:69`, `lib/index.js:13`, `475-506`): renders tree/summary output.
- `chalk` (`package.json:60`, `lib/args.js:8`, `lib/index.js:12`, `bin/license-checker:14`): terminal color support.
- `mkdirp` (`package.json:62`, `lib/index.js:17`, `bin/license-checker:12`): ensure output dirs exist before writes.
- `debug` (`package.json:61`, `lib/index.js:16`, README `163-170`): namespaced debug logging.

External systems/side effects:

- local filesystem reads/writes (`fs` in CLI/core parser)
- process control (`process.exit`, stdout/stderr)
- environment-controlled debug channel (`DEBUG`)

## Risk Map

### Invariants (observed)
- Output key format is `name@version` (`lib/index.js:31`).
- One module entry per key due to duplicate guard (`41-47`).
- License-file precedence order is deterministic (`lib/license-files.js:3-29`).
- Structured output modes disable color by default (`lib/args.js:73-75`).

### Critical failure paths
- Invalid start path / dependency load failure bubbles via callback error (`tests/test.js:379-386`).
- Policy violation on `failOn`/`onlyAllow` exits process (`lib/index.js:438-457`).
- Invalid JSON in custom format file yields `Error` result from `parseJson` (`593-608`, `tests/test.js:701-716`).

### Edge cases (observed)
- Exclude list supports escaped commas in license strings (`lib/index.js:313-315`, `tests/test.js:152-165`).
- `BSD` exclusion expands to multiple SPDX BSD variants (`359`, `384-386`, `tests/test.js:167-181`).
- Private modules are marked `UNLICENSED` and can be excluded (`324-326`, `428-435`, `tests/packages-test.js:37-44`).
- URL/file license references become `Custom: ...` (`lib/license.js:76-79`, `tests/license.js:133-147`).

### Direct evidence vs inference
- **Direct:** all behavior above tied to code/tests line references.
- **Inference:** `lib/stack.js` appears vestigial because no references found in repository search.

## Open Questions

1. **`--version` exits with code `1`** (`bin/license-checker:49-52`). Intentional convention or legacy quirk?
2. **`--files` and `--markdown` are implemented but under-documented** (present in code, not in README option list). Should rewrite preserve hidden behavior or align docs/options?
3. **`parseJson` error propagation in `init`**: if `--customPath` is invalid, `customFormat` becomes an `Error` object (`264-266`). Should this hard-fail early in rewrite?
4. **`onlyAllow` matching uses substring logic** (`447-451`) rather than exact SPDX semantic matching; confirm whether compatibility requires preserving this behavior.

## Evidence Pointers

- Repo structure: `sourcebot_list_tree` on `github.com/davglass/license-checker` root and `tests/fixtures`
- CLI behavior: `bin/license-checker:17-109`
- Option model: `lib/args.js:7-99`, README `69-90`
- Core pipeline: `lib/index.js:27-469`
- Serializers/export: `lib/index.js:471-627`
- License parsing heuristics: `lib/license.js:1-83`
- License file precedence: `lib/license-files.js:3-29`
- Tests validating edge conditions: `tests/test.js`, `tests/license.js`, `tests/failOn-test.js`, `tests/packages-test.js`, `tests/license-files-test.js`
