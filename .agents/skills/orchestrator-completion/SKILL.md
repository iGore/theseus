---
name: orchestrator-completion
description: >
  Use this skill before declaring any phase or the overall workflow complete.
  It defines explicit exit criteria for each phase gate and for the full
  pipeline. The Orchestrator runs in a continuous verification loop: emit a
  completion promise only after all checklist items pass, then run Oracle
  verification before stopping. Do not rely on a subjective sense of "done".
---

# Orchestrator Completion

## Purpose

The Orchestrator operates in a self-referential completion loop. Believing
the work is done is not the same as the work being done. This skill provides
the Oracle verification step that must pass before the Orchestrator may stop.

A subjective "I think this phase is complete" is a **completion promise** —
not a completion fact. The Oracle (this checklist) must confirm the promise
before the Orchestrator proceeds or stops.

## Completion Loop Protocol

1. When you believe a phase or the workflow is complete, do **not** stop.
2. Emit an internal completion promise: note that you believe the phase is done.
3. Immediately run the Oracle verification for that phase (checklist below).
4. Only if **every item passes** may you treat the phase as complete.
5. If any item fails, continue working rather than stopping or skipping.
6. Repeat until Oracle verification passes for all phases.

The loop only ends after Phase 3 Oracle verification confirms every item.
The `steps` parameter in the agent profile ensures enough iterations to
complete this loop without premature termination.

## Phase 1 Oracle — Rekonstruktion

Before proceeding to Phase 2, verify all of the following:

- [ ] `ANALYSIS.md` exists in `target/specs/`
- [ ] `OVERVIEW.md` exists in `target/specs/`
- [ ] `FUNCTIONALITY-INDEX.md` exists in `target/specs/` and every item has a status of `ready` or an explicit `blocked` reason — no item is left blank
- [ ] One `SPEC.md` exists for every item marked `ready` in `FUNCTIONALITY-INDEX.md`, each in its own `target/specs/<use-case-slug>/` folder
- [ ] `SPEC-INDEX.md` exists in `target/specs/` and lists every generated requirement folder
- [ ] No Spec-Writer run is still in progress or unconfirmed

**Oracle result**: if all pass → Phase 1 complete. If any fail → continue Phase 1.

## Phase 2 Oracle — Transformationsplanung

Before proceeding to Phase 3, verify all of the following:

- [ ] `ARCHITECTURE` exists in `target/specs/`
- [ ] One `PLAN.md` exists for every use-case folder listed in `SPEC-INDEX.md`
- [ ] No Planner run is still in progress or unconfirmed

**Oracle result**: if all pass → Phase 2 complete. If any fail → continue Phase 2.

## Phase 3 Oracle — Neuimplementierung

Before declaring the workflow complete, verify all of the following:

- [ ] One `TASKS.md` exists for every use-case folder that has a `PLAN.md`
- [ ] One `BUILD.md` exists for every use-case folder that has a `TASKS.md`
- [ ] All checkboxes in every `TASKS.md` are checked — no open items remain
- [ ] No Task Decomposer or Builder run is still in progress or unconfirmed

**Oracle result**: if all pass → workflow complete. If any fail → continue Phase 3.

## Guardrails

- Do not declare a phase complete if any Oracle item is unverified
- Do not assume a background run finished — confirm the artifact exists on disk
- If a required artifact is missing, route the work back to the responsible role rather than skipping the gate
- Do not substitute a summary or status note for an actual artifact check
- A completion promise without a passing Oracle result has no effect — keep working
