# AGENTS.md

## Purpose

This repository is a starter for a subagent-based migration workflow.
It defines a small, spec-first pipeline and the supporting assets needed to run it in Codex-style and OpenCode environments.

## Repository structure

- `.codex/agents/`: TOML agent profiles
- `.opencode/agents/`: Markdown agent profiles with YAML frontmatter
- `opencode.jsonc`: OpenCode configuration for MCP servers and default agent selection
- `.agents/skills/`: Agent Skills directories with `SKILL.md`
- `docker-compose.yml`: Sourcebot bootstrap
- `config.json`: Sourcebot example configuration
- `README.md`: quickstart and MCP setup notes

## Canonical role order

Use the following flow unless there is a strong reason to deviate:

1. `Analyzer`
2. `Spec-Writer`
3. `Verifier`
4. `Architect`
5. `Planner`
6. `Task Decomposer`
7. `Builder`

The `Orchestrator` coordinates handoffs, artifact gates, and retry decisions across the full flow.

This role set now matches the full conceptual model from the report: `Orchestrator`, `Analyzer`, `Spec-Writer`, `Verifier`, `Architect`, `Planner`, `Task Decomposer`, and `Builder`.

## Phase model

- Phase 1 `Rekonstruktion`: `Analyzer -> Spec-Writer -> Verifier`
- Phase 2 `Transformationsplanung`: `Architect -> Planner`
- Phase 3 `Neuimplementierung`: `Task Decomposer -> Builder`

The Orchestrator spans all three phases, coordinates handoffs, and enforces the transition gates between them.

## Role intent

- `Analyzer`: reconstruct the bounded source module from repository structure, symbols, dependencies, build metadata, and supporting docs
- `Spec-Writer`: condense Analyzer findings into the central textual specification artifact in Markdown
- `Verifier`: check whether the specification is traceable back to source-context evidence before later phases are allowed
- `Architect`: derive target architecture, framework/library choices, and implementation structure from the verified specification
- `Planner`: translate verified specification and architecture decisions into an ordered migration plan with work packages
- `Task Decomposer`: break the migration plan into small, implementation-ready coding tasks for the builder
- `Builder`: implement target code from the verified specification and architecture decisions, then record technical verification evidence

## MCP usage

- `Sourcebot`: primary MCP for reconstruction work - repository structure, symbol lookup, references, and code-intelligence tasks
- `Context7`: primary MCP for target-stack planning and implementation - framework, package, library, and best-practice documentation

Do not model these MCPs as skills. They are external context/access layers.

Keep secrets such as Sourcebot API keys out of tracked repo files. Use local-only files such as `.codex/config.local.toml`, `.codex/secrets.toml`, or `opencode.local.jsonc`. A tracked template is available as `opencode.local.jsonc.example`. Use the shared env key `SOURCEBOT_BEARER_TOKEN` when wiring Sourcebot auth across Codex and OpenCode.

## Skill usage

Skills live under `.agents/skills/<skill-name>/SKILL.md`.

They are optional helpers for execution, not a canonical layer of the report's pipeline model. If you use them, prefer only the skills that are actually reflected in the report and starter flow:

- `reverse-engineering`
- `code-mapping`
- `dependency-inspection`
- `spec-driven-design`
- `specification-drafting`
- `contract-writing`
- `edge-case-structuring`
- `traceability-check`
- `consistency-review`
- `acceptance-validation`
- `evaluation-rubric`
- `ecosystem-mapping`
- `package-selection`
- `spec-to-code-translation`
- `test-scaffolding`
- `build-verification`

## Working rules

- Keep the workflow spec-first: do not jump from analysis straight to implementation.
- Keep the migration slice bounded and explicit.
- Have all agents work inside `target/` for generated artifacts, plans, and implementation outputs unless the user explicitly requests a different location.
- Treat the textual specification as the main handoff artifact between reconstruction and implementation.
- The Verifier is a release gate: unresolved source/spec mismatches block later phases.
- Record architecture and package decisions before build work starts.
- Record a migration plan before detailed implementation work starts.
- Use task decomposition to keep the builder on bounded, implementation-ready work packages.
- Give each agent only the context and MCP access needed for its current phase.
- Prefer documented evidence over intuition.
- Keep agent names aligned with the report: `Analyzer`, `Spec-Writer`, `Verifier`, `Architect`, `Planner`, `Task Decomposer`, `Builder`, `Orchestrator`.

## Sourcebot bootstrap

To fetch the Sourcebot compose file:

```bash
curl -o docker-compose.yml https://raw.githubusercontent.com/sourcebot-dev/sourcebot/main/docker-compose.yml
```

`Sourcebot` is configured for both Codex and OpenCode against `http://localhost:3000/api/mcp`.

## Context7 bootstrap

To register Context7 for Codex:

```bash
codex mcp add context7 -- npx -y @upstash/context7-mcp
```

To register Context7 for OpenCode:

```bash
opencode mcp add context7 -- npx -y @upstash/context7-mcp
```
