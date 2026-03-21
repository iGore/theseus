---
name: test-scaffolding
description: Create a minimal but meaningful test scaffold for a bounded implementation slice. Use when a builder subagent needs executable evidence tied to contracts or acceptance criteria.
---

# Test Scaffolding

## Purpose
Create the smallest useful set of tests that can validate the implementation slice.

## When to use
- A new implementation slice needs verification.
- The builder must show evidence beyond code alone.

## Inputs
- Specification
- Contracts
- Edge cases

## Workflow
1. Select the most relevant test types.
2. Add a minimal but meaningful scaffold.
3. Ensure tests map back to contracts or acceptance criteria.

## Outputs
- Test file skeletons
- Test intent map
- Coverage notes

## Guardrails
- Favor useful coverage over volume.
- Keep tests within the bounded slice.
