# Theseus OpenCode Starter

Minimal starter for the spec-first migration workflow described in the report, adapted to OpenCode's agent format.

The OpenCode agent prompts are now structured in a `spec-kit`-inspired template style: each role has explicit required inputs, a single named output artifact, and a repeatable output format.

## Structure

- `.opencode/agents/`: agent profiles as Markdown files with YAML frontmatter
- `.agents/skills/`: reusable cross-agent skill directories
- `opencode.jsonc`: project-local OpenCode configuration with MCP setup

## Included Roles

- `orchestrator.md`: primary agent that coordinates stage order, artifact gates, and retry paths without taking over specialist stage work itself
- `analyzer.md`: subagent for reconstruction of the bounded source module
- `spec-writer.md`: subagent for the central textual specification artifact
- `verifier.md`: subagent for source/spec traceability and release-gate checks
- `architect.md`: subagent for target architecture and package planning from the verified spec
- `planner.md`: subagent for the ordered migration plan and work-package sequence
- `task-decomposer.md`: subagent for turning the plan into implementation-ready coding tasks
- `builder.md`: subagent for target-code implementation from verified artifacts

## Prompt Template Style

The agent files intentionally mirror patterns from `github/spec-kit/templates`:

- explicit `Required Input`
- one named artifact per stage
- template-shaped output sections
- explicit gates and retry rules

This keeps the prompts reusable and makes stage handoffs easier to audit.

## Quick Start

1. Describe the source system and target system in one task.
2. Choose a bounded demonstration module for the PoC.
3. Start `orchestrator` as the primary agent and keep the order Analyzer -> Spec-Writer -> Verifier -> Architect -> Planner -> Task Decomposer -> Builder.
4. Keep `orchestrator` coordination-only: it should route, gate, and hand off work, not author stage artifacts for the specialist roles.
5. Persist each stage result as an artifact (Markdown or JSON), with the textual specification as the main handoff between reconstruction and implementation.

## Three phases from the report

- Phase 1 `Rekonstruktion`: `Analyzer -> Spec-Writer -> Verifier`
- Phase 2 `Transformationsplanung`: `Architect -> Planner`
- Phase 3 `Neuimplementierung`: `Task Decomposer -> Builder`

## Configured MCPs

- `Context7`: documentation and library reference context for authoring, orchestration, and implementation support
- `Sourcebot`: code-intelligence and repository context for reconstruction, specification, and verification

Add `Context7` with:

```bash
opencode mcp add context7 -- npx -y @upstash/context7-mcp
```

`Sourcebot` is also preconfigured in `opencode.jsonc` with the local MCP endpoint:

```text
http://localhost:3000/api/mcp
```
