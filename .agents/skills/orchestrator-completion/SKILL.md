---
name: orchestrator-completion
description: >
  Use this skill before declaring any phase or the overall workflow complete.
  It defines explicit exit criteria for each phase gate and for the full
  pipeline. Run through the checklist before stopping — do not rely on a
  subjective sense of "done".
---

# Orchestrator Completion

## Purpose

Use this skill to verify that a phase or the full pipeline is genuinely
complete before stopping. It replaces subjective "done" judgments with
explicit, artifact-based exit criteria.

## When to use

- Before closing Phase 1 and handing off to Architect
- Before closing Phase 2 and handing off to Task Decomposer
- Before closing Phase 3 and declaring the workflow complete
- Any time you are about to stop and believe the workflow is finished

## Phase 1 exit criteria — Rekonstruktion

Before proceeding to Phase 2, verify all of the following:

- [ ] `ANALYSIS.md` exists in `target/specs/`
- [ ] `OVERVIEW.md` exists in `target/specs/`
- [ ] `FUNCTIONALITY-INDEX.md` exists in `target/specs/` and every item has a status of `ready` or an explicit `blocked` reason — no item is left blank
- [ ] One `SPEC.md` exists for every item marked `ready` in `FUNCTIONALITY-INDEX.md`, each in its own `target/specs/<use-case-slug>/` folder
- [ ] `SPEC-INDEX.md` exists in `target/specs/` and lists every generated requirement folder
- [ ] No Spec-Writer run is still in progress or unconfirmed

## Phase 2 exit criteria — Transformationsplanung

Before proceeding to Phase 3, verify all of the following:

- [ ] `ARCHITECTURE` exists in `target/specs/`
- [ ] One `PLAN.md` exists for every use-case folder listed in `SPEC-INDEX.md`
- [ ] No Planner run is still in progress or unconfirmed

## Phase 3 exit criteria — Neuimplementierung

Before declaring the workflow complete, verify all of the following:

- [ ] One `TASKS.md` exists for every use-case folder that has a `PLAN.md`
- [ ] One `BUILD.md` exists for every use-case folder that has a `TASKS.md`
- [ ] All checkboxes in every `TASKS.md` are checked — no open items remain
- [ ] No Task Decomposer or Builder run is still in progress or unconfirmed

## Guardrails

- Do not declare a phase complete if any checklist item above is unverified
- Do not assume a background run finished — confirm the artifact exists on disk
- If a required artifact is missing, route the work back to the responsible role rather than skipping the gate
- Do not substitute a summary or status note for an actual artifact check
