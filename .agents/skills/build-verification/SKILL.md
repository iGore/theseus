---
name: build-verification
description: Confirm that the current implementation builds and is supported by executable evidence. Use when a builder subagent must report build, test, and residual-risk status.
---

# Build Verification

## Purpose
Confirm that the current implementation is buildable and supported by executable evidence.

## When to use
- The builder finishes an implementation step.
- The workflow needs evidence for release or review.

## Inputs
- Build command
- Test command
- Current code state

## Workflow
1. Run the available build and test checks.
2. Record outcomes and failures.
3. Summarize remaining known risks.

## Outputs
- Build result
- Test result
- Known-risk note

## Guardrails
- Record failures explicitly.
- Do not claim readiness without evidence.
