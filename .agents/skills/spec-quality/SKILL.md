---
name: spec-quality
description: >
  Use this skill when writing or finalizing a SPEC.md artifact. It defines
  what "implementation-ready" means for a specification and provides a
  self-check before handoff to the Architect phase. Use it to decide whether
  a functionality item is fully captured or still has blocking open questions.
---

# Spec-Quality

## Purpose

Use this skill to determine whether a SPEC.md artifact is complete and
precise enough to pass the gate from Phase 1 (Rekonstruktion) into Phase 2
(Transformationsplanung).

This skill is a quality guide for the specification artifact. It does not
replace the SPEC.md template in the Spec-Writer agent profile — it
supplements it with explicit acceptance criteria.

## When to use

- Before finalizing any SPEC.md artifact
- When deciding whether a functionality item is fully captured
- When handling missing or ambiguous information from source evidence

## What makes a spec implementation-ready

### Mandatory elements

- **Bounded scope**: the spec covers exactly one functionality item from `USE-CASES.md`
- **Traceable source**: every core claim links to evidence from `ANALYSIS.md` or a Sourcebot finding
- **Explicit I/O**: inputs and outputs are named with types, formats, and constraints — not left implicit
- **Acceptance scenarios**: at least one scenario per user story in Gherkin format (`Given / When / Then`), independently testable, tagged with a source reference (e.g. `[SA-001]`). Inspired by Spec-Kit.
- **Explicit open questions**: any unclear point is listed as an open question — never silently assumed or guessed

### Quality indicators

- Preconditions and postconditions are stated where relevant
- Edge cases and error paths are named, even if only as open questions
- Functional requirements use unambiguous obligation language (MUST / MUST NOT)
- Success criteria are measurable, not subjective

### Output-identity indicators (for output-producing or value-transforming slices)

- The SPEC has an **Output Contract** pinned to a verbatim golden baseline:
  field order, conditional fields, exact sentinels/literals, sort order, and
  trailing-newline / stdout-vs-file behavior — no prose-only format description
- Any classification/normalization rule is captured as the **full ordered rule
  table** including the **fallback branch**, with verbatim output literals
- Source dependencies whose behavior must be replicated (identifier/expression
  validator, tree/table/markdown formatter, dependency-graph resolver) are named,
  with an FR requiring the reimplementation to match them
- Each such dependency references the Analyzer's `## Source-Dependency Contracts`
  entry (documented input/output) and the Architect's dependency decision
  (adopt Go lib / faithful reimplementation / approximate)
- **Concrete example outputs are always present**: every output/transformation
  requirement carries at least one worked `input → output` example (from a golden
  snippet or the dependency contract), never a prose-only description
- At least one success criterion is a **byte-level differential-parity** baseline
  over an edge-inclusive fixture set, not just the happy path

## Handling missing information

- If source evidence is absent, record an open question — do not fill the gap with a guess
- If behavior cannot be determined from `ANALYSIS.md`, mark the requirement as `NEEDS CLARIFICATION`
- Never write a spec that requires Architect or Builder to resolve ambiguity that belongs here

## Self-check before handoff

Before submitting the final `SPEC.md`, verify:

- [ ] Scope is bounded to exactly one functionality item
- [ ] Every functional requirement has a source reference or an explicit open-question marker
- [ ] At least one acceptance scenario is independently testable without re-opening the source
- [ ] All open questions are explicitly listed — none are silently omitted
- [ ] I/O formats, types, and constraints are explicit
- [ ] (Output slices) Output Contract cites a byte-exact golden baseline and lists
      field order, sentinels/literals, sort order, and newline behavior
- [ ] (Transform slices) The full ordered rule table incl. fallback branch is
      transcribed verbatim, not summarized
- [ ] Replicated-dependency behaviors (SPDX, formatters, graph resolver) are named
      and covered by an FR + acceptance scenario
- [ ] A byte-level differential-parity success criterion covers the edge set, not
      just the happy path
- [ ] Every output/transform requirement carries a concrete worked example
      (`input → output`); none rely on prose alone
- [ ] Replicated source dependencies link to their `Source-Dependency Contracts`
      entry and the Architect's dependency decision

## Guardrails

- Do not mark a spec as complete when it still contains critical unresolved open questions
- Do not expand scope beyond the assigned functionality item
- Do not use Context7 for spec writing — the spec describes source behavior, not target implementation
- Do not conflate source reconstruction with target design decisions
- Do not pass an output slice whose format is described only in prose, or whose
  parity criterion is "semantic"/happy-path only — require a byte-exact baseline
- Do not pass a transform slice whose rule table is a "common cases" subset
- Do not pass any output/transform requirement that lacks a concrete worked
  example output
- Do not pass a replicated-dependency requirement that lacks a documented
  input/output contract reference and an explicit architecture decision
