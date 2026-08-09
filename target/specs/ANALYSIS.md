# ANALYSIS — davglass/license-checker CLI rewrite to Go

## Scope

- **Bounded source module:** the Node.js CLI/package `davglass/license-checker` only.
- **Pinned source version:** commit `de6e9a42513aa38a58efc6b202ee5281ed61f486` (`de6e9a4 changelog`), `package.json` version `25.0.1`.
- **Evidence source:** Sourcebot skill was used as the reconstruction workflow. The configured local Sourcebot endpoint was unavailable, so concrete line evidence below comes from a read-only clone of the pinned GitHub commit. The live source system was **not installed or executed**.
- **Target constraint:** downstream implementation is Go.
- **Out of scope:** whole-system architecture and any `SPEC.md`, `ARCHITECTURE.md`, `PLAN.md`, `TASKS.md`, `BUILD.md`, or Go implementation files.

## Entry Points

| Surface | Source reference | Contract |
|---|---|---|
| CLI binary | `package.json:91-94`, `bin/license-checker:1-15` | Executable is `license-checker`; script parses args, calls `checker.init`, then renders/writes output. |
| Programmatic API | `README.md:139-155`, `lib/index.js:261-469` | `checker.init(options, callback)` scans installed npm packages and returns a sorted object. |
| Argument parser | `lib/args.js:7-43`, `lib/args.js:65-90` | Single command, no subcommands. Signature: `license-checker [flags]`. Short aliases: `-v` -> `--version`, `-h` -> `--help` (`lib/args.js:36-39`). |

All known flags are declared in `lib/args.js:9-35`; defaults are in `lib/args.js:65-85`.

| Flag | Type/default | Behavior references |
|---|---|---|
| `--production` | Boolean | Filters `json.extraneous`; passes `dev:false` (`lib/index.js:49`, `273-275`). |
| `--development` | Boolean | Filters non-extraneous non-root packages; passes `dev:false` (`lib/index.js:49`, `273-275`). |
| `--json` | Boolean | `JSON.stringify(json, null, 2) + '\n'`; disables color (`bin/license-checker:84-85`, `lib/args.js:73-75`). |
| `--csv` | Boolean | `checker.asCSV`; disables color (`bin/license-checker:86-87`, `lib/index.js:508-562`). |
| `--csvComponentPrefix` | String | Adds `"component"` column and repeats supplied value (`lib/index.js:510-525`, `538-550`). |
| `--markdown` | Boolean | `checker.asMarkDown(...) + '\n'`; disables color (`bin/license-checker:88-89`). |
| `--summary` | Boolean | Count licenses, output as tree (`bin/license-checker:90-91`, `lib/index.js:479-505`). |
| `--out` | path | Create parent dir; write raw formatted output to file (`bin/license-checker:96-104`). |
| `--files` | path | Copy license files, bypass normal renderer (`bin/license-checker:96-98`, `lib/index.js:610-627`). |
| `--start` | String, default `process.cwd()` | Scan root passed to `read-installed` (`lib/args.js:76`, `lib/index.js:298`). |
| `--unknown` | Boolean | Adds `dependencyPath`; rewrites guessed `*` licenses to `UNKNOWN` (`lib/index.js:89-92`, `331-339`). |
| `--onlyunknown` | Boolean | Keeps only licenses containing `*` or `UNKNOWN` (`lib/index.js:342-346`). |
| `--exclude` | comma string | Excludes by escaped-comma split, literal invalid SPDX match, or SPDX satisfaction (`lib/index.js:313-315`, `357-399`). |
| `--failOn` | semicolon string | Mutually exclusive with `--onlyAllow`; exits `1` on exact license match (`bin/license-checker:54-63`, `lib/index.js:277-295`, `437-442`). |
| `--onlyAllow` | semicolon string | Exits `1` when license string contains none of the allowed tokens (`bin/license-checker:54-63`, `lib/index.js:444-456`). |
| `--packages` | semicolon string | Whitelist exact `name@version` keys (`lib/index.js:407-415`). |
| `--excludePackages` | semicolon string | Blacklist exact `name@version` keys (`lib/index.js:418-426`). |
| `--excludePrivatePackages` | Boolean | Deletes records with `private` flag (`lib/index.js:428-435`). |
| `--relativeLicensePath` | Boolean | Makes `licenseFile`/`noticeFile` relative to root (`lib/args.js:77`, `lib/index.js:307`, `178`, `230`). |
| `--customPath` | path | Loads custom JSON format (`lib/index.js:264-266`, `593-608`). |
| `--customFormat` | Object API option | Includes/excludes/defaults output fields (`lib/index.js:56-58`, `95-105`). |
| `--direct` | Boolean | Present -> depth `0`; absent -> `Infinity` (`lib/args.js:79-83`, `lib/index.js:270`). |
| `--color` | Boolean/default `chalk.supportsColor` | Colorizes tree keys only without `--out`, `--csv`, `--json`, or `--markdown` (`lib/args.js:70-75`, `bin/license-checker:74-82`, `107-109`). |
| `--help` | Boolean | Prints version and usage to stderr, exits `0` (`bin/license-checker:17-47`). |
| `--version` | Boolean | Prints version to stderr, exits `1` (`bin/license-checker:49-52`). |

Environment variables: `DEBUG=license-checker*` enables debug output per `README.md:158-170` and `lib/index.js:21-25`. Chalk may also derive color support from terminal/environment, but the project references only `chalk.supportsColor`.

## Inputs and Outputs

### Expected inputs

- CLI flags above.
- A start directory containing an installed npm package tree. `read-installed` receives `options.start`, `{ dev, log, depth }` (`lib/index.js:267-298`).
- Package metadata nodes returned by `read-installed`: `name`, `version`, `private`, `repository`, `url`, `author`, `license`, `licenses`, `readme`, `path`, `dependencies`, `extraneous`, and `root` (`lib/index.js:27-258`).
- Files under package dirs: `README.md`, license-like files, and `NOTICE` files (`lib/index.js:114-120`, `151-230`; `lib/license-files.js:3-29`).
- Optional custom format JSON via `--customPath` (`lib/index.js:593-608`).

### Expected outputs

- Default tree: `treeify.asTree(sorted, true)` (`lib/index.js:475-477`).
- JSON, CSV, Markdown, summary tree, or copied files depending on flags (`bin/license-checker:84-104`, `lib/index.js:479-627`).
- Diagnostics to stderr for help/version, mutual exclusion, comma warning, scan errors, fail policies, and missing copied license files (`bin/license-checker:17-72`, `lib/index.js:437-456`, `624`).
- Exit code `0` for normal CLI (`tests/bin-test.js:7-14`), `0` for help, `1` for version/fail conditions (`bin/license-checker:46-56`, `49-52`, `lib/index.js:440-455`, `tests/failOn-test.js:7-42`).

### Golden Output Snippets

No live CLI execution was performed. These are committed examples/tests from the pinned source; byte-exact stdout for some modes remains limited because tests assert substrings or parsed JSON rather than full stdout.

Default tree example (`README.md:24-57`, excerpt):

```text
├─ cli@0.4.3
│  ├─ repository: http://github.com/chriso/cli
│  └─ licenses: MIT
├─ glob@3.1.14
│  ├─ repository: https://github.com/isaacs/node-glob
│  └─ licenses: UNKNOWN
└─ yui-lint@0.1.1
   ├─ licenses: BSD
      └─ repository: http://github.com/yui/yui-lint
```

Guessed license marker (`README.md:59-67`):

```text
└─ debug@2.0.0
   ├─ repository: https://github.com/visionmedia/debug
   └─ licenses: MIT*
```

CSV snippets (`tests/test.js:37-39`, `77-91`):

```text
"module name","license","repository"
"abbrev@1.0.9","ISC","https://github.com/isaacs/abbrev-js"
```

```text
"module name","name","description","pewpew"
"abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"
```

```text
"component","module name","name","description","pewpew"
"main-module","abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"
```

Markdown snippets (`tests/test.js:41-43`, `95-103`):

```text
[abbrev@1.0.9](https://github.com/isaacs/abbrev-js) - ISC
```

```text
 - **[abbrev@1.0.9](https://github.com/isaacs/abbrev-js)**
```

Diagnostic snippet (`bin/license-checker:59-62`, asserted in `tests/failOn-test.js:27-40`):

```text
Warning: As of v17 the --failOn argument takes semicolons as delimeters instead of commas (some license names can contain commas)
```

Other exact stderr strings: mutual exclusion `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` (`bin/license-checker:54-56`); failOn `Found license defined by the --failOn flag: "<license>". Exiting.` (`lib/index.js:439-441`); onlyAllow `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` (`lib/index.js:453-455`).

## Function Inventory

| Function/symbol | File | Role | Slice |
|---|---|---|---|
| top-level CLI | `bin/license-checker` | Validate CLI, call init, choose renderer/write path. | F-001/F-006/F-007 |
| `shouldColorizeOutput(args)` | `bin/license-checker:107-109` | Color eligibility. | F-006 |
| `raw(args)` | `lib/args.js:41-43` | Calls `nopt`. | F-002 |
| `has(a)` | `lib/args.js:45-57` | Detects cooked option presence. | F-002 |
| `clean(args)` | `lib/args.js:59-63` | Removes `argv`. | F-002 |
| `setDefaults(parsed)` | `lib/args.js:65-86` | Applies defaults. | F-002 |
| `parse(args)` | `lib/args.js:88-91` | Parse plus defaults. | F-002 |
| `flatten(options)` | `lib/index.js:27-258` | Converts dependency tree into package records and recurses children. | F-003/F-004 |
| `include(property)` | `lib/index.js:56-58` | Custom-format inclusion gate. | F-003/F-006 |
| `exports.init(options, callback)` | `lib/index.js:261-469` | Main scan/sort/filter/restrict/fail pipeline. | F-003/F-005/F-007 |
| `colorizeString(string)` | `lib/index.js:318-321` | Chalk red/bold sentinel formatting. | F-006 |
| `transformBSD(spdx)` | `lib/index.js:357-360` | `BSD` exclusion alias expansion. | F-005 |
| `invert(fn)` | `lib/index.js:361` | Predicate inversion utility. | F-005 |
| `spdxIsValid(spdx)` | `lib/index.js:362` | `spdxCorrect(spdx) === spdx`. | F-005 |
| `exports.print(sorted)` | `lib/index.js:471-473` | Prints tree. | F-006 |
| `exports.asTree(sorted)` | `lib/index.js:475-477` | Tree rendering via `treeify`. | F-006 |
| `exports.asSummary(sorted)` | `lib/index.js:479-505` | License counts sorted by descending count. | F-006 |
| `exports.asCSV(...)` | `lib/index.js:508-562` | CSV rendering. | F-006 |
| `exports.asMarkDown(...)` | `lib/index.js:571-590` | Markdown rendering. | F-006 |
| `exports.parseJson(jsonPath)` | `lib/index.js:593-608` | Custom JSON loading. | F-002/F-006 |
| `exports.asFiles(json, outDir)` | `lib/index.js:610-627` | License-file copy output. | F-006 |
| license classifier | `lib/license.js:22-83` | SPDX pass-through, regex classification, fallback. | F-004 |
| license file detector | `lib/license-files.js:14-29` | License filename precedence. | F-004 |
| `Stack` | `lib/stack.js:6-43` | Async utility; not referenced by current CLI path. | low-priority |

## Input-Field Inventory

| Field/source | Default path | Overrides and trigger conditions | Probe / golden evidence |
|---|---|---|---|
| CLI `start` | `parsed.start || process.cwd()` -> `read(options.start, ...)` | None (`lib/args.js:76`, `lib/index.js:298`). | Defaults asserted `tests/test.js:392-413`. |
| CLI `direct` | Absent -> `Infinity`; present -> `0`; passed as `depth` | Conditional CLI override (`lib/args.js:79-83`, `lib/index.js:270`). | `tests/test.js:404-413`. |
| Package `name`, `version` | Builds key `name@version` | If missing, `delete data[key]` (`lib/index.js:31`, `255-257`). | No dedicated fixture. |
| Package `private` | Sets `moduleInfo.private = true` | **Unconditional override:** during sorted pass, `licenses` becomes `UNLICENSED` whenever `data[item].private` is true, regardless of package license (`lib/index.js:37-39`, `323-326`). `--excludePrivatePackages` can later delete record (`428-435`). | Probe fixture `tests/fixtures/privateModule/package.json:1-5`; `UNLICENSED` asserted `tests/test.js:291-310`; empty JSON key set asserted `tests/packages-test.js:37-44`. No live byte stdout captured. |
| Package `repository.url` | Copied to `repository` when object URL string | **Unconditional normalization:** `git+ssh://git@` -> `git://`; `git+https://github.com` -> `https://github.com`; `git://github.com` -> `https://github.com`; `git@github.com:` -> `https://github.com/`; trailing `.git` removed (`lib/index.js:60-68`). | Normalized abbrev URL in CSV/Markdown tests (`tests/test.js:37-43`). |
| Package `url.web`, `author.url` | `json.url.web` populates `url`; author fields populate `publisher`, `email`, `url` | **Unconditional override:** `json.author.url` overwrites earlier `json.url.web` when included (`lib/index.js:70-86`). | No dedicated fixture; open golden gap. |
| Package `license`/`licenses` | `json.license || json.licenses`; objects/strings classified, arrays mapped (`lib/index.js:112-138`) | README can supply license when no package license (`139-140`). License file can override missing/UNKNOWN/Custom licenses (`168-172`). | Classifier tests `tests/license.js:10-181`; custom fixtures under `tests/fixtures/custom-license-*`. |
| Package `readme` | Used for classification if no license fields | **Unconditional override:** if missing or contains `no readme data found`, local `README.md` is read and assigned (`lib/index.js:114-120`). | Fixture README exists for private module; no exact output assertion. |
| Package `path` | Output `path`; filesystem root for file scan | `dependencyPath` added only with `--unknown`; license/notice paths relative only with `--relativeLicensePath` (`lib/index.js:89-92`, `108-110`, `177-179`, `226-230`, `307`). | `tests/test.js:480-530`. |
| Directory filenames | `licenseFiles(dirFiles)` | Precedence list selects first match per basename class (`lib/license-files.js:3-29`). | Exact expectations `tests/license-files-test.js:10-87`. |
| License file content | Highest-precedence file can classify and fill `licenseText`/`copyright` | **Unconditional override:** if license is absent, includes `UNKNOWN`, or starts `Custom:`, file content reclassifies `licenses` (`lib/index.js:162-172`). `licenseText` newline/quote behavior changes with CSV (`181-190`). | “license file over custom urls” expects `MIT*` (`tests/test.js:313-323`). |
| NOTICE files | Adds `noticeFile` | Relative only with `--relativeLicensePath` (`lib/index.js:155-159`, `226-230`). | No dedicated test. |
| Custom format JSON | Parsed object controls fields | `false` suppresses a field; otherwise missing field receives configured default (`lib/index.js:56-58`, `95-105`, `593-608`). | `customFormatExample.json:1-10`; `tests/test.js:429-476`, `690-716`. |

## Value-Transformation Tables

### License classifier exact ordered rule set (`lib/license.js:22-83`)

1. Try `spdxExpressionParse(str || '')`; on success return raw `str` unchanged (`24-27`).
2. If `str` truthy, replace the first newline only: `str = str.replace('\n', '')` (`31`).
3. Undefined/falsey -> `Undefined` (`33-34`).
4. `/The ISC License/` -> `ISC*` (`35-36`).
5. `/ermission is hereby granted, free of charge, to any/` -> `MIT*` (`37-38`).
6. `/edistribution and use in source and binary forms, with or withou/` -> `BSD*` (`39-40`).
7. `/edistribution and use of this software in source and binary forms, with or withou/` -> `BSD-Source-Code*` (`41-43`). This appears after the broader BSD rule and can be shadowed.
8. `/DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE/` -> `WTFPL*` (`44-45`).
9. `/\bISC\b/` -> `ISC*` (`46-47`).
10. `/\bMIT\b/` -> `MIT*` (`48-49`).
11. `/\bBSD\b/` -> `BSD*` (`50-51`).
12. `/\bWTFPL\b/` -> `WTFPL*` (`52-53`).
13. `/\bApache License\b/` -> `Apache*` (`54-55`).
14. CC0 deed regex -> `CC0-1.0*` (`56-57`).
15. GPL regex `/\bGNU GENERAL PUBLIC LICENSE\s*Version ([^,]*)/i`: capture version, append `.0` when length 1, return `GPL-<version>*` (`58-65`).
16. LGPL regex `/(?:LESSER|LIBRARY) GENERAL PUBLIC LICENSE\s*Version ([^,]*)/i`: same version padding, return `LGPL-<version>*` (`66-72`).
17. `/[Pp]ublic [Dd]omain/` -> `Public Domain` (`73-74`).
18. URL regex or `SEE LICENSE IN (.*)` -> `Custom: <capture>` (`76-79`).
19. Fallback/no match -> `null` (`79-81`).

Examples are asserted in `tests/license.js:10-181`, including SPDX pass-through for `MIT`, `LGPL-2.0`, `Apache-2.0`, `BSD-2-Clause`, `(GPL-2.0+ WITH Bison-exception-2.2)`, `LGPL-2.0 OR (ISC AND BSD-3-Clause+)`, and `Apache-2.0 OR ISC OR MIT`.

### License file precedence (`lib/license-files.js:3-29`)

Evaluation order: `^LICENSE$`, `^LICENSE\-\w+$`, `^LICENCE$`, `^LICENCE\-\w+$`, `^COPYING$`, `^README$`. Basename is uppercased and extension is ignored; only first matching filename per pattern is pushed. Tests cover empty/no-match, extension-insensitive, case-insensitive, multiple candidates, and first `LICENSE-*` only (`tests/license-files-test.js:10-87`).

### Exclusion transformation (`lib/index.js:313-399`)

- Split with `/([^\\\][^,]|\\,)+/g`, unescape `\,`, trim whitespace (`313-315`).
- Transform exact exclusion `BSD` to `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)` (`357-360`).
- Valid exclusions satisfy `spdxCorrect(spdx) === spdx`; invalid exclusions match literally (`362-366`, `388-390`).
- Candidate package license arrays are normalized by removing trailing `*` and transforming exact `BSD`; `UNKNOWN` is never excluded (`374-397`).

## Output Field Contract

### Package record field order and sentinels

JavaScript insertion order controls `JSON.stringify` and `treeify` traversal. Per record:

1. `licenses` is inserted first by `var moduleInfo = { licenses: UNKNOWN }` (`lib/index.js:27-29`). Initial sentinel is `UNKNOWN` (`line 7`).
2. `private` is inserted next only when `json.private` is truthy (`37-39`). Later, `licenses` is overwritten to `UNLICENSED` (`323-326`) but remains first.
3. Subsequent conditional insertions occur in source order: `repository`, `url`, `publisher`, `email`, author `url`, `dependencyPath`, custom-format fields, `path`, `licenseFile`, `licenseText`, `copyright`, `noticeFile` (`60-230`).
4. `licenses` can be overwritten by package/readme/file classification and by `--unknown`, but not re-ordered (`112-146`, `168-172`, `331-339`). Sentinels/literals include `UNKNOWN`, `UNLICENSED`, `Undefined`, guessed `*` values, `Custom: <value>`, `Public Domain`, and raw valid SPDX expressions.
5. Top-level package keys are sorted lexically before filtering by `Object.keys(data).sort().forEach` (`323`). Whitelist/blacklist rebuild objects in filtered key order (`407-426`).

### Renderer contracts

- **Tree:** `treeify.asTree(sorted, true)` (`475-477`). Stdout uses `console.log`, adding an extra newline (`bin/license-checker:102-104`).
- **JSON:** `JSON.stringify(json, null, 2) + '\n'` (`bin/license-checker:84-85`). Stdout gets an additional `console.log` newline; `--out` writes only the formatted string.
- **CSV:** no trailing newline from `asCSV`; default header is `"module name","license","repository"`; custom header is `"module name"` plus `Object.keys(customFormat)` with optional leading `"component"` (`lib/index.js:508-562`). General CSV escaping is not implemented; values are wrapped in double quotes.
- **Markdown:** default row `[key](repository) - licenses`; custom row ` - **[key](repository)**` followed by four-space-indented custom fields (`571-590`). CLI appends one newline, then stdout adds another.
- **Summary:** count exact `licenses` string, sort by descending count, render through treeify (`479-505`). Equal-count ordering is not stabilized by source code.
- **Files:** file name `<moduleName>-LICENSE.txt`; missing `licenseFile` warns `no license file found for: <moduleName>` (`610-627`).

## Feature Slices

| ID | Title | Status | Main source | Notes |
|---|---|---|---|---|
| F-001 | CLI entry point and invocation model | ready | `bin/license-checker`, `package.json` | Single command with stderr help/version. |
| F-002 | Argument parsing, defaults, and custom format loading | ready | `lib/args.js`, `lib/index.js:264-270`, `593-608` | Includes `nopt`, `direct`, `color`, `start`, `customPath`. |
| F-003 | Dependency scan and metadata flattening | ready | `lib/index.js:27-258`, `261-309` | Depends on `read-installed` tree semantics. |
| F-004 | License detection/classification and license file precedence | ready | `lib/license.js`, `lib/license-files.js`, `lib/index.js:112-230` | Exact ordered rules captured. |
| F-005 | Filtering/restriction/compliance decisions | ready | `lib/index.js:313-458` | Includes exclude, packages, private, failOn/onlyAllow. |
| F-006 | Output rendering and side-effect outputs | ready | `bin/license-checker:84-104`, `lib/index.js:471-627` | Full JSON byte fixture is an open parity gap, but renderer source is explicit. |
| F-007 | Diagnostics, error handling, and exit codes | ready | `bin/license-checker:17-72`, `lib/index.js:437-467` | Includes stderr contracts and process exits. |

## Domain Map

- **Package node:** installed package metadata returned by `read-installed`.
- **Output record:** mutable object keyed by `name@version`, initially `{ licenses: 'UNKNOWN' }`.
- **License value:** raw SPDX expression, classifier literal, guessed `*` literal, `Custom: ...`, `Public Domain`, `UNKNOWN`, `UNLICENSED`, or `Undefined`.
- **Custom format:** JSON object whose keys define additional/excluded output fields and defaults.
- **Restriction set:** package whitelist/blacklist and license-policy filters applied after flatten/sort.

## Flow Map

1. CLI parses args with nopt and defaults (F-001/F-002).
2. `checker.init` loads custom format and calls `read-installed` (F-002/F-003).
3. `flatten` recursively builds records, normalizes metadata, reads files, classifies licenses (F-003/F-004).
4. Keys sort lexically; private/missing/unknown sentinels are applied (F-003/F-004).
5. License exclusions, package restrictions, private deletion, fail policies run (F-005/F-007).
6. CLI chooses renderer or file-copy/out-file side effect (F-006).
7. Diagnostics/errors emit to stderr; process exits where source explicitly calls `process.exit` (F-007).

## End-to-End Pipeline

This is an architectural anchor, not a separate functionality item: `bin/license-checker` -> `args.parse()`/defaults (F-001/F-002) -> help/version/fail flag preflight (F-001/F-007) -> `checker.init(args, cb)` (F-003/F-005/F-007) -> optional custom JSON parse (F-002) -> `read-installed` dependency tree load (F-003) -> recursive `flatten` record extraction and license detection (F-003/F-004) -> lexical sort and sentinel overrides (F-003/F-004) -> exclude/package/private/fail filters (F-005/F-007) -> renderer selection or `--files`/`--out` side effect (F-006) -> stdout/stderr and exit behavior (F-007).

## Dependency Notes

Manifest dependencies (`package.json:59-70`): `chalk`, `debug`, `mkdirp`, `nopt`, `read-installed`, `semver`, `spdx-correct`, `spdx-expression-parse`, `spdx-satisfies`, `treeify`. There is no committed lockfile, so exact transitive versions are not pinned by the repository artifact.

Important output-shaping dependencies:

- `read-installed ~4.0.3`: builds the dependency graph and provides `extraneous`, `root`, `path`, `readme`, dependencies, etc. High parity risk.
- `nopt ^4.0.1`: CLI type coercion and `--no-` cooked parsing. Moderate risk.
- `spdx-expression-parse ^3.0.0`: validates SPDX expressions; valid input returns raw string. Standardized-and-deep risk.
- `spdx-correct ^3.0.0` and `spdx-satisfies ^4.0.0`: exclusion matching. Standardized-and-deep risk.
- `treeify ^1.1.0`: box-drawing tree output. Moderate byte-output risk.
- `chalk ^2.4.1`: ANSI colors for tree keys/sentinel strings. Moderate terminal-dependent risk.
- `debug ^3.1.0`: opt-in diagnostics controlled by `DEBUG`. Low core-output risk.
- `mkdirp ^0.5.1`: creates output dirs. Low formatting risk.

## Source-Dependency Contracts

- **`nopt` (`^4.0.1`)**: caller supplies known option types and short aliases (`lib/args.js:7-43`). It parses `process.argv`/provided args and returns an object plus `argv` metadata, which this source deletes. Worked example from source: `-v` maps to `version`; path-typed `out`, `files`, `customPath` are declared using `require('path')`. Re-derivation risk: moderate.
- **`read-installed` (`~4.0.3`)**: called as `read(options.start, {dev, log, depth}, cb)` (`lib/index.js:267-298`). Output contract relied upon: tree nodes with metadata fields and nested `dependencies`; `extraneous`/`root` flags drive production/development filters. Worked example in tests uses the repository itself and expects `abbrev@1.0.9` with `ISC` (`tests/test.js:19-44`). Re-derivation risk: high.
- **`spdx-expression-parse` (`^3.0.0`)**: parser accepts valid SPDX identifiers/expressions and throws on invalid. Source only uses success/failure and returns the original string on success (`lib/license.js:24-27`). Worked examples in `tests/license.js:155-181`. Re-derivation risk: standardized-and-deep.
- **`spdx-correct`/`spdx-satisfies` (`^3.0.0`/`^4.0.0`)**: source partitions valid exclusions using `spdxCorrect(spdx) === spdx`, corrects package license strings, and checks satisfaction against an OR expression (`lib/index.js:357-391`). Worked fixtures cover excluding MIT/ISC, escaped-comma Apache string, BSD alias, and Public Domain (`tests/test.js:137-197`). Re-derivation risk: standardized-and-deep.
- **`treeify` (`^1.1.0`)**: source calls `asTree(obj, true)` for full tree/summary (`lib/index.js:475-477`, `505`). README snippets show `├─`, `│`, `└─` output. Re-derivation risk: moderate for byte identity.
- **`chalk` (`^2.4.1`)**: source uses `chalk.supportsColor`, `chalk.blue`, `chalk.dim`, `chalk.green`, `chalk.bold.red` (`lib/args.js:70-75`, `bin/license-checker:74-82`, `lib/index.js:318-321`). Re-derivation risk: terminal-dependent; likely compatible ANSI library needed if color supported.

## Risk Map

| Risk | Severity | Evidence | Mitigation for downstream specs |
|---|---|---|---|
| `read-installed` graph semantics are not a filesystem walk | high | `lib/index.js:298`, fields like `extraneous`/`root` used in filters | Specify logical npm installed-tree behavior, not naive directory traversal. |
| SPDX semantics delegated to libraries | high | `lib/license.js:24-27`, `lib/index.js:357-391` | Use equivalent Go SPDX packages or explicitly reproduce examples/edge cases. |
| Tree/JSON/CSV byte details are output-sensitive | high | renderers `lib/index.js:475-590`; newline split in `bin/license-checker:84-104` | Preserve field order, newlines, no CSV escaping bug, and tree glyphs. |
| Unconditional private/license-file overrides can hide source metadata | high | `lib/index.js:323-326`, `168-172` | Spec as explicit FRs in field extraction/license slices. |
| Color output is terminal/environment dependent | medium | `chalk.supportsColor`; `shouldColorizeOutput` | Treat no-color as baseline unless color flag/environment is in test matrix. |
| No lockfile pins dependency versions | medium | no `package-lock.json`; ranges in `package.json` | Pin behavior to commit and recorded examples, not latest dependency docs. |
| Summary equal-count ordering unspecified | low/medium | numeric sort only `lib/index.js:497-499` | Do not promise deterministic ties unless Go implementation mimics JS engine used by baseline. |

## Open Questions

1. Sourcebot was unavailable; if later accessible, re-check line evidence against indexed commit.
2. Full byte-exact stdout fixtures for JSON, full tree, summary, `--files`, `--out`, color, and error paths cannot be captured without executing the source; this run intentionally did not execute it.
3. There is no dedicated source test for `NOTICE` files or author URL overriding `url.web`.
4. Dependency documentation was summarized from package roles and public npm documentation; no lockfile fixes exact dependency implementation versions.

## Evidence Pointers

- `package.json:59-70`, `91-100`, `106-110` — dependencies, bin/main, scripts, license/repository.
- `README.md:69-90`, `115-137`, `158-182` — CLI docs, custom format docs, debug/env, license-finding prose.
- `bin/license-checker:17-104` — CLI control flow and render/write selection.
- `lib/args.js:9-98` — flags, aliases, defaults.
- `lib/index.js:27-258` — flatten/metadata/license-file extraction.
- `lib/index.js:261-469` — init pipeline, filters, fail policies.
- `lib/index.js:471-627` — renderers and side-effect file output.
- `lib/license.js:1-83` — exact license classifier.
- `lib/license-files.js:3-29` — license filename precedence.
- Tests: `tests/test.js`, `tests/license.js`, `tests/license-files-test.js`, `tests/failOn-test.js`, `tests/packages-test.js`, `tests/bin-test.js`.

## Use Cases

- UC-1: User invokes `license-checker` with flags and receives documented help/version or scan output (F-001/F-002/F-007).
- UC-2: User scans an npm project and obtains package records containing license and metadata (F-003/F-004).
- UC-3: User constrains results by dependency scope, package list, private packages, unknown licenses, or license policy (F-005).
- UC-4: User emits tree, JSON, CSV, Markdown, summary, output file, or copied license-file artifacts (F-006).
- UC-5: Programmatic caller imports `license-checker` and uses `init` plus renderer helpers (F-002 through F-006).

## App Tasks

1. Parse command line and defaults.
2. Resolve dependency graph from start path.
3. Extract and normalize package metadata.
4. Determine license from package fields, readme, license file, and SPDX/regex rules.
5. Apply filters/restrictions/fail policies.
6. Render or write output in selected mode.
7. Emit diagnostics and exit with source-compatible codes.

## Requirements Task List

- RT-001: Implement single-command CLI with all flags, aliases, defaults, mutual exclusions, and stderr help/version behavior (evidence: `lib/args.js`, `bin/license-checker`).
- RT-002: Implement dependency scan abstraction with `production`, `development`, and `direct` semantics (evidence: `lib/index.js:267-309`).
- RT-003: Build package records with source field order, sentinels, repository normalization, author/url handling, custom format defaults, path fields, and private override (evidence: `lib/index.js:27-110`, `323-339`).
- RT-004: Implement exact license classifier and file precedence rules (evidence: `lib/license.js`, `lib/license-files.js`, `tests/license*.js`).
- RT-005: Implement all filters/restrictions, including escaped-comma excludes, BSD alias, SPDX matching, semicolon package lists, and fail exits (evidence: `lib/index.js:313-458`).
- RT-006: Implement renderers and side effects preserving field order, line endings, stdout vs file newline behavior, and copied file naming (evidence: `bin/license-checker:84-104`, `lib/index.js:471-627`).
- RT-007: Implement debug/error/stderr contracts and exit code behavior (evidence: `README.md:158-170`, `bin/license-checker`, `tests/failOn-test.js`).

## Feature-to-Task Mapping

| Feature | Requirements |
|---|---|
| F-001 | RT-001, RT-007 |
| F-002 | RT-001, RT-006 |
| F-003 | RT-002, RT-003 |
| F-004 | RT-004 |
| F-005 | RT-005, RT-007 |
| F-006 | RT-006 |
| F-007 | RT-007 |
