---
name: dependency-inspection
description: Identify internal and external dependencies that constrain a migration slice. Use when a subagent must detect coupling, side effects, frameworks, services, or storage dependencies before specification or build work.
---

# Dependency Inspection

## Purpose
Expose dependencies that constrain or complicate the migration slice.

## When to use
- A module interacts with frameworks, services, storage, or shared state.
- The migration scope must be risk-assessed before specification.

## Inputs
- Source module
- Build and runtime context

## Workflow
1. List internal and external dependencies.
2. Identify side effects, persistence, and network interactions.
3. Flag hidden coupling and migration risk.

## Outputs
- Dependency name
- Dependency type
- Migration risk note

## Guardrails
- Distinguish direct from indirect dependencies.
- Highlight unknowns instead of guessing.
