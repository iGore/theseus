---
name: edge-case-structuring
description: Capture edge cases in a reusable form for specification, verification, and testing. Use when failure conditions, boundary values, or exception behavior must be made explicit.
---

# Edge-Case Structuring

## Purpose
Make edge cases explicit so they survive the transition from analysis to specification, tests, and implementation.

## When to use
- Inputs, states, or dependencies can fail in non-happy-path scenarios.
- Verification and tests need structured exception coverage.

## Inputs
- Specification draft
- Risk notes
- Failure modes

## Workflow
1. Group edge cases by category.
2. Describe expected behavior for each case.
3. Link the case to validation or tests.

## Outputs
- Edge-case list
- Expected behavior notes
- Test or validation hooks

## Guardrails
- Do not mix edge cases with normal flow.
- Keep wording concrete and observable.
