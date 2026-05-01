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

- **Bounded scope**: the spec covers exactly one functionality item from `FUNCTIONALITY-INDEX.md`
- **Traceable source**: every core claim links to evidence from `ANALYSIS.md` or a Sourcebot finding
- **Explicit I/O**: inputs and outputs are named with types, formats, and constraints — not left implicit
- **Acceptance scenarios**: at least one Given/When/Then scenario per user story that is independently testable
- **Explicit open questions**: any unclear point is listed as an open question — never silently assumed or guessed

### Quality indicators

- Preconditions and postconditions are stated where relevant
- Edge cases and error paths are named, even if only as open questions
- Functional requirements use unambiguous obligation language (MUST / MUST NOT)
- Success criteria are measurable, not subjective

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

## Guardrails

- Do not mark a spec as complete when it still contains critical unresolved open questions
- Do not expand scope beyond the assigned functionality item
- Do not use Context7 for spec writing — the spec describes source behavior, not target implementation
- Do not conflate source reconstruction with target design decisions
