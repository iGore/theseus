---
name: traceability-check
description: Verify that specification or implementation claims are backed by source evidence. Use when a verifier subagent must check whether statements can be traced to concrete source artifacts.
---

# Traceability Check

## Purpose
Verify that relevant claims in the specification or implementation can be traced to concrete source evidence.

## When to use
- A verifier must check whether the specification is grounded.
- A release gate depends on evidence quality.

## Inputs
- Source artifacts
- Specification statements
- Optional implementation evidence

## Workflow
1. Extract relevant claims.
2. Link each claim to a source artifact.
3. Mark unsupported claims as blockers or open questions.

## Outputs
- Traceability matrix
- Unsupported-claim list
- Blocker list

## Guardrails
- Evidence must be specific, not generic.
- Unsupported claims must never be silently accepted.
