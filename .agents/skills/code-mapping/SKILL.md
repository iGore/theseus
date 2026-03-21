---
name: code-mapping
description: Map source-level constructs to migration-relevant concepts and responsibilities. Use when analysis findings must be translated into target-facing concepts for specs, contracts, or design decisions.
---

# Code Mapping

## Purpose
Translate source-level constructs into migration-relevant concepts that can be reused in specifications and design decisions.

## When to use
- Source code must be abstracted into business or technical concepts.
- A later specification should not depend on raw code fragments alone.

## Inputs
- Source files or symbols
- Reverse-engineering notes

## Workflow
1. Group source elements by responsibility.
2. Map code elements to target-facing concepts.
3. Record why each mapping is valid.

## Outputs
- Source element
- Mapped concept
- Mapping rationale

## Guardrails
- Do not invent concepts without evidence.
- Keep mappings stable across artifacts.
