---
name: ecosystem-mapping
description: Map source-side findings to a suitable target ecosystem, including frameworks, architectural conventions, and migration-relevant best practices. Use when an architect subagent must translate analysis into target-stack choices.
---

# Ecosystem Mapping

## Purpose
Turn analysis findings into a coherent target-ecosystem view.

## When to use
- A migration slice must be aligned to a target stack.
- Framework and architecture options must be compared before specification.

## Inputs
- Analyzer artifacts
- Target-stack constraints
- External documentation context

## Workflow
1. Read the recovered source responsibilities and dependencies.
2. Map them to fitting target-stack concepts and architecture patterns.
3. Document why the selected ecosystem direction is appropriate.

## Outputs
- Target ecosystem map
- Architecture option notes
- Decision rationale

## Guardrails
- Keep the mapping tied to the bounded migration slice.
- Prefer documented conventions over personal style.
