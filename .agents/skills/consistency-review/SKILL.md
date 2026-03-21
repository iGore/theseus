---
name: consistency-review
description: Check whether analysis, specification, and implementation artifacts contradict each other. Use when a verifier subagent must detect naming, behavior, or scope inconsistencies across the workflow.
---

# Consistency Review

## Purpose
Check that artifacts remain aligned across analysis, specification, verification, and implementation outputs.

## When to use
- Multiple artifacts exist for the same migration slice.
- The verifier must detect contradictions early.

## Inputs
- Analysis artifacts
- Specification artifacts
- Build notes or code outputs

## Workflow
1. Compare naming and concept usage.
2. Compare behavioral and error-handling statements.
3. Flag scope drift or contradictions.

## Outputs
- Consistency findings
- Contradiction list
- Required follow-up actions

## Guardrails
- Focus on meaningful contradictions, not wording style.
- Record where each inconsistency appears.
