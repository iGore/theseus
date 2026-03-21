---
name: contract-writing
description: Write explicit interface and behavior contracts for specifications. Use when a spec-writer subagent must define inputs, outputs, preconditions, postconditions, and failure cases precisely.
---

# Contract Writing

## Purpose
Write explicit contracts so interfaces and behaviors can be implemented and verified consistently.

## When to use
- Interfaces or data structures must be fixed precisely.
- Specifications need implementation-facing contract sections.

## Inputs
- Specification draft
- Source evidence
- Target constraints

## Workflow
1. Name the contract and its scope.
2. Define inputs, outputs, and invariants.
3. Add preconditions, postconditions, and failure cases.

## Outputs
- Contract name
- Inputs and outputs
- Preconditions and postconditions
- Failure cases

## Guardrails
- Avoid ambiguous field names.
- Keep contracts aligned with source evidence.
