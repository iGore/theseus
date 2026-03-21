---
name: reverse-engineering
description: Analyze one bounded source module to recover responsibilities, behavior, data flow, and failure modes. Use when a subagent must understand existing code before writing specs or implementation plans.
---

# Reverse Engineering

## Purpose
Build a reliable understanding of a bounded source module before specification or implementation work begins.

## When to use
- A source module must be understood before writing specifications.
- Behavior is implicit in code and not documented elsewhere.

## Inputs
- Module or package path
- Runtime or framework context
- Available code, tests, and configuration

## Workflow
1. Identify entry points, public surfaces, and main execution paths.
2. Trace core data flow, state changes, and side effects.
3. Record assumptions, failure modes, and open questions.

## Outputs
- Responsibility summary
- Inputs and outputs map
- Dependency list
- Failure-mode notes

## Guardrails
- Do not redesign the module.
- Separate observed behavior from assumptions.
