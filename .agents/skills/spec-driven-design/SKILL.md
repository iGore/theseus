---
name: spec-driven-design
description: Create implementation-ready technical specifications that act as the source of truth in a spec-first workflow. Use when analysis findings must become a structured spec, plan, and task sequence for downstream subagents.
---

# Spec-Driven Design

## Purpose
Create technical specifications that directly drive planning, implementation, and verification.

## When to use
- Analysis findings must be turned into an implementation-ready specification.
- The workflow expects the specification to act as the source of truth.
- A team wants a Spec-Kit-style flow: specification first, then planning, then task breakdown, then implementation.

## Inputs
- Reverse-engineering artifacts
- Target-stack constraints
- Migration scope
- Any explicit non-goals or out-of-scope constraints

## Workflow
1. Write the core specification first.
   - Define scope, interfaces, behavioral rules, data contracts, error handling, and acceptance criteria.
   - Treat this file as the source of truth for downstream work.
2. Derive an implementation plan from the approved specification.
   - Identify architectural decisions, sequencing, constraints, and validation strategy.
   - Capture what must happen before coding begins.
3. Break the plan into tasks.
   - Produce bounded implementation tasks that are traceable to the specification.
   - Keep tasks small enough to execute and verify independently.
4. Hand the plan and tasks to downstream implementers.
   - Implementation should refine the spec only through explicit feedback, not silent drift.

## Spec-Kit-inspired flow

Use this order whenever possible:

1. `spec.md` - what must be true
2. `plan.md` - how the specification will be realized
3. `tasks.md` - what concrete implementation steps follow
4. implementation and verification artifacts

This flow is useful because it prevents direct prompt-to-code jumps and keeps planning tied to a stable specification.

## Outputs
- `spec.md` or equivalent specification artifact
- `plan.md` or equivalent implementation plan
- `tasks.md` or equivalent execution breakdown
- Traceable links between requirements, tasks, and later verification

## Prompt templates

For reusable prompt patterns, see:

- `references/PROMPTS.md`
- `references/FLOW.md`

## Guardrails
- Avoid vague requirements.
- Every important rule must be testable.
- Do not skip directly from analysis to code if a specification is still missing.
- Keep planning and tasks traceable to the specification.
