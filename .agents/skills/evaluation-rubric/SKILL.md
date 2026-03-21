---
name: evaluation-rubric
description: Apply a compact quality rubric for release and review decisions. Use when a verifier or orchestrator needs a reproducible gate result with blocker classification.
---

# Evaluation Rubric

## Purpose
Apply a compact quality rubric so release or review decisions stay explicit and reproducible.

## When to use
- A verifier or orchestrator needs a clear gate decision.
- Different artifacts must be judged against the same criteria.

## Inputs
- Verification findings
- Acceptance results
- Build and test evidence

## Workflow
1. Review evidence by criterion.
2. Assign PASS, CONDITIONAL, or FAIL.
3. Record blockers and remediation.

## Outputs
- Rubric result
- Blocker summary
- Remediation summary

## Guardrails
- Use the same criteria across iterations.
- Explain any non-pass result concisely.
