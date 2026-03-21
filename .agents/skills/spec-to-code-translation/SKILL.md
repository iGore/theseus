---
name: spec-to-code-translation
description: Translate a validated specification into an implementation plan and a bounded code slice. Use when a builder subagent starts implementation from approved contracts and acceptance criteria.
---

# Spec-to-Code Translation

## Purpose
Translate validated specification content into an implementation plan and a minimal code slice.

## When to use
- The builder receives an approved specification.
- The workflow moves from specification to implementation.

## Inputs
- Technical specification
- Contracts
- Acceptance criteria

## Workflow
1. Derive implementation tasks from the contracts.
2. Implement only what the approved scope requires.
3. Preserve traceability back to specification and criteria.

## Outputs
- Implementation plan
- Code slice
- Traceability notes

## Guardrails
- Treat the specification as the source of truth.
- Do not add features outside the approved scope.
