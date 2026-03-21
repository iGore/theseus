---
name: package-selection
description: Select suitable packages and libraries for the target ecosystem based on documented best practices, compatibility, and scope. Use when an architect subagent must justify dependency choices.
---

# Package Selection

## Purpose
Choose packages that fit the target ecosystem and migration scope.

## When to use
- The target stack needs concrete library choices.
- There are multiple candidate packages for a required capability.

## Inputs
- Target ecosystem map
- Functional requirements
- Documentation and best-practice references

## Workflow
1. List the capability that needs package support.
2. Compare candidate packages using documentation and best practices.
3. Record the chosen package and the reason for the decision.

## Outputs
- Package shortlist
- Selected package set
- Selection rationale

## Guardrails
- Avoid adding packages without a clear requirement.
- Favor stable, well-documented packages when possible.
