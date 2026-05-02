# Context7

## Purpose

Use this skill when a role needs external framework, library, package, or best-practice documentation from Context7.

This skill is a usage guide for the Context7 MCP. It does not replace Context7 as an external documentation layer.

## When to use

- select packages or frameworks for the target stack
- validate implementation patterns against current documentation
- support architecture and planning decisions with package evidence
- resolve API usage questions during implementation
- gather documentation support for `ARCHITECTURE.md`, `PLAN.md`, or `BUILD.md`

## Required inputs

- target framework, package, or library name
- concrete documentation question to answer
- bounded migration slice or implementation context
- target artifact that will consume the documentation evidence

## Workflow

### 1. Bound the documentation need

- confirm the exact library, framework, or package in scope
- avoid broad ecosystem research unless the task explicitly requires it

### 2. Resolve the library

- identify the correct Context7 library ID before querying docs
- prefer the most authoritative and relevant documentation source

### 3. Query focused documentation

- ask specific questions about APIs, setup, constraints, or best practices
- gather only the sections needed for the current bounded task

### 4. Record documentation evidence

- capture package names, APIs, version-relevant notes, and short factual summaries
- separate documented guidance from local design decisions
- keep notes tied to the current migration slice

### 5. Handoff cleanly

- summarize only the documentation needed by the downstream artifact
- include unresolved gaps or version questions as open questions

## Expected outputs

Provide compact, traceable documentation support with:

- bounded documentation scope
- selected package or library
- relevant APIs or conventions
- documentation references
- open questions or version notes

## Guardrails

- prefer official or high-authority documentation surfaced through Context7
- do not use Context7 for repository reconstruction when Sourcebot is the better tool
- do not expand scope beyond the agreed migration slice
- keep documented facts separate from architecture or implementation decisions
