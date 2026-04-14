# Sourcebot

## Purpose

Use this skill when a role needs repository reconstruction, symbol lookup, cross-reference tracing, or code-intelligence support from Sourcebot.

This skill is a usage guide for the Sourcebot MCP. It does not replace Sourcebot as an external access layer.

## When to use

- reconstruct a bounded source module
- inspect repository structure
- locate definitions, references, and ownership
- gather traceable repository evidence for `ANALYSIS.md`, `SPEC.md`, or `VERIFY.md`
- trace requirements back to concrete code locations

## Required inputs

- bounded module or migration slice
- repository or path scope when known
- concrete question to answer
- target artifact that will consume the evidence

## Workflow

### 1. Bound the search

- confirm the slice, package, directory, or symbol to inspect
- avoid whole-repository exploration unless the task explicitly requires it

### 2. Collect structure evidence

- use Sourcebot tree and repository search features to map relevant files
- identify likely entry points, interfaces, models, and tests

### 3. Collect symbol evidence

- locate definitions for key functions, classes, types, routes, or commands
- trace references to understand dependencies and call paths

### 4. Record evidence

- capture file paths, symbols, and short factual notes
- distinguish direct evidence from inferred conclusions
- keep evidence tied to the current bounded slice

### 5. Handoff cleanly

- summarize only the evidence needed by the downstream artifact
- include unresolved ambiguities as open questions

## Expected outputs

Provide compact, traceable evidence with:

- bounded scope
- key files and symbols
- dependency or call-flow notes
- source references
- open questions or ambiguity notes

## Guardrails

- prefer Sourcebot evidence over intuition
- do not invent behavior that is not visible in the repository context
- do not expand scope beyond the agreed migration slice
- use Context7 separately for external package or framework documentation
